package apis

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

var block_regex = `([a-zA-Z0-9]+)\[(?:(\@[\"\,\@\=\>\<a-zA-Z0-9\-\.\|\:_]*|\*))\]\{([\*|[\,\@\"\=a-zA-Z0-9\-]*)`
var filter_regex = `\@([a-zA-Z0-9]*)([\<\>\=]*)\"([a-zA-Z0-9\-\.\|\:_]*)\"`

type QSLService struct {
	host    string
	metaSvc *MetaService
}

func NewQSLService(host string, m *MetaService) *QSLService {
	return &QSLService{host, m}
}

// filterfunc
// @name="cluster1" -> eq(name,cluster1)
// @name="paas-preprod-west2.cluster.k8s.local",@k8sobj="K8sObj",@resourceid="paas-preprod-west2.cluster.k8s.local"
// -> eq(name,paas-preprod-west2.cluster.k8s.local) and eq(k8sobj,K8sObj) and eq(resourceid,paas-preprod-west2.cluster.k8s.local)
// filterdeclaraction
// @name="paas-preprod-west2.cluster.k8s.local",@k8sobj="K8sObj",@resourceid="paas-preprod-west2.cluster.k8s.local"
// -> , $name: string, $k8sobj: string, $resourceid: string
func CreateFiltersQuery(filterlist string) (string, string, error) {
	// split the whole string by the | "or" symbol because of higher priority for ands
	// e.g. a&b&c|d&e == (a&b&c) | (d&e)
	splitlist := strings.Split(filterlist, "|")
	// the variable definitions e.g. $name: string,
	filterdeclaration := ""
	// the eq functions eq(name,paas-preprod-west2.cluster.k8s.local)
	filterfunc := []string{}

	operator_map := map[string]string{
		">":  "gt",
		">=": "ge",
		"<=": "le",
		"<":  "lt",
		"=":  "eq",
	}

	for i, item := range splitlist {
		_ = i
		splitstring := strings.Split(item, ",")

		interfilterfunc := []string{}

		for j, item2 := range splitstring {
			_ = j
			// kvsplit := strings.Split(item2, "=")
			r := regexp.MustCompile(filter_regex)
			matches := r.FindStringSubmatch(item2)
			log.Debugf("filtermatches %s %#v\n", item2, matches)

			if len(matches) < 4 {
				return "", "", errors.New("Invalid filters in " + filterlist)
			}

			keyname := matches[1]
			operator := matches[2]
			value := matches[3]

			//filterdeclaration = filterdeclaration + ", $" + kvsplit[0][1:] + ": string"
			filterdeclaration = filterdeclaration + ", $" + keyname + ": string"
			// interfilterfunc = append(interfilterfunc, " eq("+kvsplit[0][1:]+",\""+kvsplit[1][1:len(kvsplit[1])-1]+"\") ")
			interfilterfunc = append(interfilterfunc, " "+operator_map[operator]+"("+keyname+",\""+value+"\") ")
		}

		filterfunc = append(filterfunc, strings.Join(interfilterfunc, "and"))

	}

	return filterdeclaration, " (" + strings.Join(filterfunc, "or") + ")", nil

}

// creates a list of the fields of an object we want to return
// will be joined with newlines for the resulting query
// e.g. @name,@resourceversion -> [name, resourceversion]
func CreateFieldsQuery(fieldlist string, metafieldslist []MetadataField, tabs int) ([]string, error) {
	// TODO: potentially add some character that we can use to get all related object types
	// kind of like expand all
	// get a list of fields since we can't use expand all due to repeating subgroups
	if fieldlist == "*" {
		returnlist := []string{}
		for i, item := range metafieldslist {
			_ = i
			if item.FieldType != "relationship" {
				returnlist = append(returnlist, strings.Repeat("\t", tabs+1)+item.FieldName)

			}

		}
		returnlist = append(returnlist, strings.Repeat("\t", tabs+1)+"uid")
		return returnlist, nil
	} else if string(fieldlist[0]) == "*" && len(fieldlist) > 1 {
		// TODO: actually check if all of the characters are * instead of just assuming
		returnlist := []string{}
		for i := 0; i < len(fieldlist); i++ {
			returnlist = append([]string{strings.Repeat("\t", len(fieldlist)-i) + "expand(_all_){"}, returnlist...)
			returnlist = append(returnlist, strings.Repeat("\t", len(fieldlist)-i)+"}")

		}

		return returnlist, nil
	}

	splitlist := strings.Split(fieldlist, ",")
	returnlist := []string{}

	for i, item := range splitlist {
		_ = i
		if strings.HasPrefix(item, "@") {
			if len(item) > 1 && !strings.Contains(item[1:], "@") {
				returnlist = append(returnlist, strings.Repeat("\t", tabs+1)+item[1:])
			} else {
				return nil, errors.New("Field names must not contain @ sign other than in prefix [" + item[1:] + "]")
			}

		} else {
			return nil, errors.New("Field names must be prefixed with @ sign [" + item + "]")
		}

	}

	returnlist = append(returnlist, strings.Repeat("\t", tabs+1)+"uid")
	log.Debugf("fields %#v\n", returnlist)
	return returnlist, nil

}

// creates the inner nested clauses searching for relationships
func (qa *QSLService) CreateDgraphQueryHelper(query []string, tabs int, parent string) ([]string, error) {

	// string that we're going to replace
	basequery := []string{
		strings.Repeat("\t", tabs) + "$RELATION @filter(eq(objtype, $OBJTYPE) and$FILTERSFUNC){",
	}
	_ = basequery

	// regex to match the string pattern
	//r := regexp.MustCompile(`([a-zA-Z]+)\[(\@[\,\@\"\=a-zA-Z0-9\-\.\|]*)\]\{([\*|[\,\@\"\=a-zA-Z0-9\-]*)`)
	//r := regexp.MustCompile(`([a-zA-Z]+)\[(?:(\@[\,\@\"\=a-zA-Z0-9\-\.\|\:_]*|\*))\]\{([\*|[\,\@\"\=a-zA-Z0-9\-]*)`)
	r := regexp.MustCompile(block_regex)
	matches := r.FindStringSubmatch(query[0])
	// log.Infof("%#v\n", splitquery)
	log.Debugf("helpermatches %#v\n", matches)

	// extract the values of the form objtype[filters]fields and assign to individual variables
	objtype := strings.Title(matches[1])
	filters := matches[2]
	fields := matches[3]
	// // replicaset
	// if strings.HasSuffix(objtype, "set") {
	// 	objtype = objtype[0:len(objtype)-3] + "Set"
	// }
	// log.Debugf("helperobjtype %#v\n", objtype)

	// get a list of the metadata fields for this object type
	metafieldslist, err := qa.metaSvc.GetMetadataFields(objtype)
	if err != nil {
		log.Error("err in getting metadata fields for " + objtype)
		log.Error(err)
		return nil, errors.New("Failed to connect to dgraph to get metadata")
	}

	if len(metafieldslist) == 0 {
		log.Error("metadata for " + objtype + " not found in db. Will not be able to use * or find relationships")
	}

	log.Debugf("metadata fields for %s: %#v", objtype, metafieldslist)

	// declare relation variable
	relation := ""
	// see if we can find the reverse relation from this object to its parent
	found := false
	// look in the list of fields for the metadata and
	// find if there's a relationship between the parent's and this object's type
	// e.g. if we had cluster[...]{...}.pod[...]{...} parent=cluster
	// and we will find the pods relation to cluster is called ~cluster
	for i, item := range metafieldslist {
		_ = i

		if item.FieldType == "relationship" {
			log.Debugf("1 found relationship for %s-%s->%s", objtype, item.FieldName, item.RefDataType)
			if item.RefDataType == strings.Title(parent) {
				relation = "~" + strings.ToLower(item.FieldName)
				found = true
			}
		}
	}

	if !found {
		// if not, see if we can find the relation from the parent to this object
		metafieldslist2, err := qa.metaSvc.GetMetadataFields(parent)
		if err != nil {
			log.Error("err in getting metadata fields")
			log.Error(err)
			return nil, errors.New("Failed to connect to dgraph to get metadata")
		}
		log.Debugf("couldn't find relation for %s->%s,", parent, objtype)
		log.Debugf("metadata fields for %s: %#v", parent, metafieldslist)

		for i, item := range metafieldslist2 {
			_ = i
			if item.FieldType == "relationship" {
				log.Debugf("2 found relationship for %s-%s->%s", parent, item.FieldName, item.RefDataType)
				if item.RefDataType == strings.Title(objtype) {
					relation = strings.ToLower(item.FieldName)
					found = true
				}
			}
		}
	}

	// no relation found between the two objects
	if !found {
		return nil, errors.New("no relation found between " + objtype + " and " + parent)
	}

	_, ff, err := CreateFiltersQuery(filters)
	if err != nil {
		return nil, err
	}
	fl, err := CreateFieldsQuery(fields, metafieldslist, tabs)
	if err != nil {
		return nil, err
	}

	// replace filters and object type accordingly
	basequery[0] = strings.Replace(basequery[0], "$FILTERSFUNC", ff, -1)
	basequery[0] = strings.Replace(basequery[0], "$OBJTYPE", strings.Title(objtype), -1)

	// add the tilde because we are adding an inverse relationship
	basequery[0] = strings.Replace(basequery[0], "$RELATION", relation, -1)

	// append the fields to be returned
	basequery = append(basequery, fl...)

	// if there's more, recursively create the relationship clauses and add them
	if len(query) > 1 && !(string(fields[0]) == "*" && len(fields) > 1) {
		middlefilter, err := qa.CreateDgraphQueryHelper(query[1:], tabs+1, objtype)
		if err != nil {
			return nil, err
		}
		basequery = append(basequery, middlefilter...)
	}

	// add closing braces and return
	basequery = append(basequery, []string{strings.Repeat("\t", tabs) + "}"}...)

	log.Debugf("helper query for %s: %#v\n", objtype, basequery)
	return basequery, nil
}

func (qa *QSLService) CreateDgraphQuery(query string) (string, error) {
	// e.g. cluster[@name="cluster1"]{@name,@region}.pod[@name="pod1"]{@phase,@image}
	// split by }. to get each individual block
	log.Info("Received Query: ", strings.Split(query, "}."))

	splitquery := strings.Split(query, "}.")
	basequery := []string{
		"query objects($objtype: string$FILTERSDEC){",
		"objects(func: eq(objtype, $OBJTYPE)) @filter($FILTERSFUNC){",
	}
	_ = basequery
	// extract the objtype, filters and fields to return from the query string
	// r := regexp.MustCompile(`([a-zA-Z]+)\[(\@[\,\@\"\=a-zA-Z0-9\-\.]*)\]\{([\*|[\,\@\"\=a-zA-Z0-9\-]*)`)
	// r := regexp.MustCompile(`([a-zA-Z]+)\[(?:(\@[\,\@\"\=a-zA-Z0-9\-\.\|\:_]*|\*))\]\{([\*|[\,\@\"\=a-zA-Z0-9\-]*)`)
	r := regexp.MustCompile(block_regex)
	matches := r.FindStringSubmatch(splitquery[0])
	log.Debugf("matches %#v\n", matches)
	log.Debugf("%#v\n", r.SubexpNames())

	if len(matches) < 2 {
		log.Error("Malformed Query received: " + query)
		return "", errors.New("Malformed Query: " + query)
	}

	// extract the values of the form objtype[filters]fields and assign to individual variables
	objtype := strings.Title(matches[1])
	filters := matches[2]
	fields := matches[3]
	_ = fields

	log.Debugf("helperobjtype %#v\n", objtype)

	log.Debugf("objecttype %#v\n", objtype)
	// get a list of the fields supoorted for this object type
	metafieldslist, err := qa.metaSvc.GetMetadataFields(objtype)
	if err != nil {
		log.Error("err in getting metadata fields")
		log.Error(err)
		return "", errors.New("Failed to connect to dgraph to get metadata")
	}

	if len(metafieldslist) == 0 {
		log.Error("metadata for " + objtype + " not found in db. Will not be able to use * or find relationships")
	}

	log.Debugf("metafields for %s: %#v\n", objtype, metafieldslist)

	// convert the filters and fields to corresponding dgraph lines for the query
	fd, ff, err := CreateFiltersQuery(filters)
	if err != nil {
		return "", err
	}
	fl, err := CreateFieldsQuery(fields, metafieldslist, 0)
	if err != nil {
		return "", err
	}

	log.Debugf("Objtype: %#v\n", objtype)
	log.Debugf("Filterfunc: %#v\nFilterdec: %#v\n", ff, fd)
	log.Debugf("Fieldlist: %#v\n", fl)

	// replace the filters and object type and add the list of fields
	basequery[0] = strings.Replace(basequery[0], "$FILTERSDEC", fd, -1)
	basequery[1] = strings.Replace(basequery[1], "$FILTERSFUNC", ff, -1)
	basequery[1] = strings.Replace(basequery[1], "$OBJTYPE", objtype, -1)
	basequery = append(basequery[0:2], fl...)

	// if there's relations, this length will be greater than 1
	if len(splitquery) > 1 && !(string(fields[0]) == "*" && len(fields) > 1) {
		// recursively create the intermediate relations
		middlefilter, err := qa.CreateDgraphQueryHelper(splitquery[1:], 1, objtype)
		if err != nil {
			return "", err
		}
		// and add them to the end of this list with the closing braces
		basequery = append(basequery, middlefilter...)
		basequery = append(basequery, []string{"}", "}"}...)
	} else {
		// just add closing braces
		basequery = append(basequery, []string{"}", "}"}...)
	}

	// create the string joining all the lines with newlines
	finalquery := strings.Join(basequery, "\n")
	// log.Infof("New Query: %s\n", finalquery)
	return finalquery, nil

}

// Run query on a dgraph instance
func (qa *QSLService) ExecuteDgraphQuery(query string) (map[string]interface{}, error) {
	conn, err := grpc.Dial(qa.host, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	dgclient := dgo.NewDgraphClient(
		api.NewDgraphClient(conn),
	)

	txn := dgclient.NewTxn()
	defer txn.Discard(context.Background())

	resp, err := txn.Query(context.Background(), query)
	if err != nil {
		log.Errorf("query err: %#v\n", err)
		return nil, errors.New("Could not successfully execute query. Please try again later.\n" + err.Error())
	}

	respjson := map[string]interface{}{}

	err = json.Unmarshal(resp.GetJson(), &respjson)
	if err != nil {
		log.Errorf("unmarshal err: %#v\n", err)
		return nil, errors.New("Could not successfully handle data from query. Please try again later.")
	}

	// log.Infof("response from executing dgraph query: %#v\n", respjson)
	return respjson, nil

}
