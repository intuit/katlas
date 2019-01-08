package apis

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode"

	log "github.com/Sirupsen/logrus"
	"github.com/intuit/katlas/service/db"
)

// regex to get objtype[filters]{fields}
var blockRegex = `([a-zA-Z0-9]+)\[(?:(\@[\"\,\@\=\>\<a-zA-Z0-9\-\.\|\&\:_]*|\**))\]\{([\*|[\,\@\"\=a-zA-Z0-9\-]*)`

// regex to get KeyOperatorValue from something like numreplicas>=2
var filterRegex = `\@([a-zA-Z0-9]*)([\<\>\=]*)(\"?[a-zA-Z0-9\-\.\|\&\:_]*\"?)`

// QSLService service for QSL
type QSLService struct {
	DBclient db.IDGClient
	metaSvc  *MetaService
}

// IsAlphaNum determine if string is made up of only alphanumeric characters
func IsAlphaNum(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

// IsStar determine if string is made up of only *
func IsStar(s string) bool {
	for _, r := range s {
		if string(r) != "*" {
			return false
		}
	}
	return true
}

func (qa *QSLService) GetMetadata(objtype string) ([]MetadataField, error) {
	// get a list of the fields supoorted for this object type
	start := time.Now()
	metafieldslist, err := qa.metaSvc.GetMetadataFields(objtype)
	if err != nil {
		log.Error("err in getting metadata fields")
		log.Error(err)
		return []MetadataField{}, errors.New("Failed to connect to dgraph to get metadata")
	}

	if len(metafieldslist) == 0 {
		log.Error("metadata for " + objtype + " not found in db. Will not be able to use * or find relationships")
	}
	end := time.Now()
	elapsed := end.Sub(start)
	log.Debugf("[elapsedtime: %s] metafields for %s: %#v\n", elapsed, objtype, metafieldslist)
	return metafieldslist, err
}

// NewQSLService creates an instance of a QSLService
func NewQSLService(host db.IDGClient, m *MetaService) *QSLService {
	return &QSLService{host, m}
}

// CreateFiltersQuery translates the filters part of the qsl string to dgraph
// filterfunc
// @name="cluster1" -> eq(name,cluster1)
// @name="paas-preprod-west2.cluster.k8s.local",@k8sobj="K8sObj",@resourceid="paas-preprod-west2.cluster.k8s.local"
// -> eq(name,paas-preprod-west2.cluster.k8s.local) and eq(k8sobj,K8sObj) and eq(resourceid,paas-preprod-west2.cluster.k8s.local)
// filterdeclaraction
// @name="paas-preprod-west2.cluster.k8s.local",@k8sobj="K8sObj",@resourceid="paas-preprod-west2.cluster.k8s.local"
// -> , $name: string, $k8sobj: string, $resourceid: string
func CreateFiltersQuery(filterlist string) (string, string, error) {
	// default for empty filters is assume no filters
	if len(filterlist) == 0 {
		return "", "", nil
	}

	// split the whole string by the | "or" symbol because of higher priority for ands
	// e.g. a&b&c|d&e == (a&b&c) | (d&e)
	splitlist := strings.Split(filterlist, "||")
	// the variable definitions e.g. $name: string,
	filterdeclaration := ""
	// the eq functions eq(name,paas-preprod-west2.cluster.k8s.local)
	filterfunc := []string{}

	operatorMap := map[string]string{
		">":  "gt",
		">=": "ge",
		"<=": "le",
		"<":  "lt",
		"=":  "eq",
	}

	for _, item := range splitlist {
		splitstring := strings.Split(item, "&&")

		interfilterfunc := []string{}

		for _, item2 := range splitstring {
			// use regex to get the key, operator and value
			r := regexp.MustCompile(filterRegex)
			matches := r.FindStringSubmatch(item2)
			log.Debugf("filtermatches %s %#v\n", item2, matches)

			// should be 4 elements in matches
			// the whole string, key, operator, value
			if len(matches) < 4 {
				return "", "", errors.New("Invalid filters in " + filterlist)
			}

			keyname := matches[1]
			operator := matches[2]
			value := matches[3]

			// if the value is meant to be an int, it won't have quotes around it
			dectype := ": string"
			if string(value[0]) != "\"" {
				dectype = ": int"
			}

			// if the value is a string make sure it has quotes on both sides
			if string(value[0]) == "\"" && !(string(value[len(value)-1]) == "\"") {
				return "", "", errors.New("Invalid filters in " + filterlist)
			}

			filterdeclaration = filterdeclaration + ", $" + keyname + dectype
			interfilterfunc = append(interfilterfunc, " "+operatorMap[operator]+"("+keyname+","+value+") ")
		}

		filterfunc = append(filterfunc, strings.Join(interfilterfunc, "and"))

	}

	log.Debugf("filter:[ %s ]\n filterdec: %s \n filterfunc: %s", filterlist, filterdeclaration, strings.Join(filterfunc, "or"))

	return filterdeclaration, "@filter(" + strings.Join(filterfunc, "or") + ")", nil

}

// CreateFieldsQuery translates the fields part of the qsl string to dgraph
// creates a list of the fields of an object we want to return
// will be joined with newlines for the resulting query
// e.g. @name,@resourceversion -> [name, resourceversion]
func CreateFieldsQuery(fieldlist string, metafieldslist []MetadataField, tabs int) ([]string, error) {
	// default case for empty fields is to display nothing
	if len(fieldlist) == 0 {
		return []string{}, nil
	}
	// if one star, show all fields
	if fieldlist == "*" {
		returnlist := []string{}
		// use list of metadatafields from metadata api to get names of all the fields
		// for this object type
		for _, item := range metafieldslist {
			if item.FieldType != "relationship" {
				returnlist = append(returnlist, strings.Repeat("\t", tabs+1)+item.FieldName)

			}

		}
		returnlist = append(returnlist, strings.Repeat("\t", tabs+1)+"uid")
		return returnlist, nil

		// if n stars show the direct relationships n levels deep
	} else if string(fieldlist[0]) == "*" && len(fieldlist) > 1 {
		returnlist := []string{}
		if IsStar(fieldlist) {
			for i := 0; i < len(fieldlist); i++ {
				returnlist = append([]string{strings.Repeat("\t", len(fieldlist)-i+tabs) + "expand(_all_){"}, returnlist...)
				returnlist = append(returnlist, strings.Repeat("\t", len(fieldlist)-i+tabs)+"}")
			}
			return returnlist, nil
		}
		return nil, errors.New("Fields may be a string of * indicating how many levels, or a list of fields @field1,@field2,... not both [" + fieldlist + "]")

	}

	splitlist := strings.Split(fieldlist, ",")
	returnlist := []string{}
	// if we have a list of fields e.g. @name,@resourceversion,@creationtime
	for _, item := range splitlist {
		// each item must begin with @ followed by an alphanumeric string
		if strings.HasPrefix(item, "@") && len(item) > 1 {
			if IsAlphaNum(item[1:]) {
				returnlist = append(returnlist, strings.Repeat("\t", tabs+1)+item[1:])
			} else {
				return nil, errors.New("Field names must be composed of only alphanumeric characters [" + item[1:] + "]")
			}

		} else {
			return nil, errors.New("Field names must be prefixed with @ sign and followed by an alphanumeric field name [" + item + "]")
		}

	}

	returnlist = append(returnlist, strings.Repeat("\t", tabs+1)+"uid")
	log.Debugf("fields %#v\n", returnlist)
	return returnlist, nil

}

// CreateDgraphQueryHelper creates the inner nested clauses searching for relationships
func (qa *QSLService) CreateDgraphQueryHelper(query []string, tabs int, parent string) ([]string, error) {

	// string that we're going to replace
	basequery := []string{
		strings.Repeat("\t", tabs) + "$RELATION @filter(eq(objtype, $OBJTYPE) $FILTERSFUNC){",
	}

	// regex to match the string pattern
	r := regexp.MustCompile(blockRegex)
	matches := r.FindStringSubmatch(query[0])
	log.Debugf("helpermatches %#v\n", matches)

	// extract the values of the form objtype[filters]fields and assign to individual variables
	objtype := strings.Title(matches[1])
	filters := matches[2]
	fields := matches[3]

	// get a list of the metadata fields for this object type
	metafieldslist, err := qa.GetMetadata(objtype)
	if err != nil {
		return nil, errors.New("Failed to connect to dgraph to get metadata")
	}

	// declare relation variable
	relation := ""
	// see if we can find the reverse relation from this object to its parent
	found := false
	// look in the list of fields for the metadata and
	// find if there's a relationship between the parent's and this object's type
	// e.g. if we had cluster[...]{...}.pod[...]{...} parent=cluster
	// and we will find the pods relation to cluster is called ~cluster
	for _, item := range metafieldslist {

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

		for _, item := range metafieldslist2 {
			if item.FieldType == "relationship" {
				log.Debugf("2 found relationship for %s-%s->%s", parent, item.FieldName, item.RefDataType)
				if item.RefDataType == strings.Title(objtype) {
					relation = strings.ToLower(item.FieldName)
					found = true
				}
			}
		}
	}

	// still no relation found between the two objects
	if !found {
		return nil, errors.New("no relation found between " + objtype + " and " + parent)
	}

	_, ff, err := CreateFiltersQuery(filters)
	if err != nil {
		return nil, err
	}
	if len(ff) > 0 {
		ff = "and" + ff[8:len(ff)-1]
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

// CreateDgraphQuery translates the querystring to a dgraph query
func (qa *QSLService) CreateDgraphQuery(query string) (string, error) {
	log.Info("Received Query: ", strings.Split(query, "}."))

	// remove all whitespace
	whitespace := regexp.MustCompile("\\s*")
	querys := whitespace.ReplaceAllString(query, "")

	// e.g. cluster[@name="cluster1.k8s.local"]{@name,@region}.pod[@name="pod1"]{@phase,@image}
	// split by }. to get each individual block
	// cannot split by . because . may be present in some object names
	splitquery := strings.Split(querys, "}.")
	basequery := []string{
		"query objects($objtype: string$FILTERSDEC){",
		"objects(func: eq(objtype, $OBJTYPE)) $FILTERSFUNC{",
	}
	// extract the objtype, filters and fields to return from the query string
	r := regexp.MustCompile(blockRegex)
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
	fd, ff, err := CreateFiltersQuery(filters)
	if err != nil {
		return "", err
	}

	log.Debugf("helperobjtype %#v\n", objtype)

	log.Debugf("objecttype %#v\n", objtype)

	metafieldslist, err := qa.GetMetadata(objtype)
	if err != nil {
		return "", err
	}

	// convert the filters and fields to corresponding dgraph lines for the query

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
	if len(splitquery) > 1 && !(len(fields) > 1 && string(fields[0]) == "*") {
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
	return finalquery, nil

}
