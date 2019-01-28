package apis

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode"

	"crypto/rand"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/intuit/katlas/service/db"
	"strconv"
)

// regex to get objtype[filters]{fields}
var blockRegex = `([a-zA-Z0-9]+)\[?(?:(\@[\(\)\"\,\@\$\=\>\<\!a-zA-Z0-9\-\.\|\&\:_]*|\**|\$\$[a-zA-Z0-9\,\=]+))\]?\{([\*|[\,\@\"\=a-zA-Z0-9\-]*)`

// regex to get KeyOperatorValue from something like numreplicas>=2
var filterRegex = `\@([a-zA-Z0-9\(\)]*)([\!\<\>\=]*)(\"?[a-zA-Z0-9\-\.\|\&\:_]*\"?)`

// QSLService service for QSL
type QSLService struct {
	DBclient db.IDGClient
	metaSvc  *MetaService
}

// MaximumLimit define pagination limit
const MaximumLimit = 10000

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

// GetMetadata - get a list of the fields supoorted for this object type
func (qa *QSLService) GetMetadata(objtype string) ([]MetadataField, error) {
	metafieldslist, err := qa.metaSvc.GetMetadataFields(objtype)
	if err != nil {
		log.Error(err)
		return []MetadataField{}, errors.New("Failed to connect to dgraph to get metadata")
	}

	if len(metafieldslist) == 0 {
		log.Error("metadata for " + objtype + " not found in db. Will not be able to use * or find relationships")
	}
	return metafieldslist, err
}

// NewQSLService creates an instance of a QSLService
func NewQSLService(host db.IDGClient, m *MetaService) *QSLService {
	return &QSLService{host, m}
}

// CreateFiltersQuery translates the filters part of the qsl string to dgraph
// input @name="name",@objtype="objtype"$$first=2,offset=2
// filterfunc
// @name="cluster1" -> eq(name,cluster1)
// @name="paas-preprod-west2.cluster.k8s.local",@k8sobj="K8sObj",@resourceid="paas-preprod-west2.cluster.k8s.local"
// -> @filter( eq(name,paas-preprod-west2.cluster.k8s.local) and eq(k8sobj,K8sObj) and eq(resourceid,paas-preprod-west2.cluster.k8s.local) )
// filterdeclaraction
// @name="paas-preprod-west2.cluster.k8s.local",@k8sobj="K8sObj",@resourceid="paas-preprod-west2.cluster.k8s.local"
// -> , $name: string, $k8sobj: string, $resourceid: string
// pagination
// $$first=2,offset=2
// -> first: 2,offset: 2
func CreateFiltersQuery(filterlist string) (string, string, string, error) {
	// default for empty filters is assume no filters
	if len(filterlist) == 0 {
		return "", "", "", nil
	}

	// for the pagination
	paginate := ""

	filtersAndPages := strings.Split(filterlist, "$$")
	if len(filtersAndPages) > 1 {
		splitlist := strings.Split(filtersAndPages[1], ",")
		for _, item := range splitlist {
			splitval := strings.Split(item, "=")

			if splitval[0] == "first" || splitval[0] == "offset" {
				paginate += "," + splitval[0] + ": " + splitval[1]
				val, err := strconv.Atoi(splitval[1])
				if err != nil || val > MaximumLimit {
					return "", "", "", fmt.Errorf("pagination format error or exceeding maxiumum limit %d", MaximumLimit)
				}
			} else {
				return "", "", "", errors.New("Invalid pagination filters in " + filterlist)
			}

		}
		// get rid of the first comma
		paginate = paginate[0:]
		if strings.HasPrefix(filterlist, "$$") {
			return "", "", paginate, nil
		}
	}

	// split the whole string by the | "or" symbol because of higher priority for ands
	// e.g. a&b&c|d&e == (a&b&c) | (d&e)
	splitlist := strings.Split(filtersAndPages[0], "||")
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
		"!=": "not eq",
	}

	for _, item := range splitlist {
		splitstring := strings.Split(item, "&&")

		interfilterfunc := []string{}

		for _, item2 := range splitstring {
			// use regex to get the key, operator and value
			r := regexp.MustCompile(filterRegex)
			matches := r.FindStringSubmatch(item2)
			// should be 4 elements in matches
			// the whole string, key, operator, value
			if len(matches) < 4 {
				return "", "", "", errors.New("Invalid filters in " + filterlist)
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
				return "", "", "", errors.New("Invalid filters in " + filterlist)
			}

			// if keyname is count filter, add prefix cnt_ to the var name
			if strings.Contains(keyname, "count(") {
				keyname = "val(cnt_" + keyname[6:]
			}

			filterdeclaration = filterdeclaration + ", $" + keyname + dectype
			interfilterfunc = append(interfilterfunc, " "+operatorMap[operator]+"("+keyname+","+value+") ")
		}

		filterfunc = append(filterfunc, strings.Join(interfilterfunc, "and"))

	}
	return filterdeclaration, "@filter(" + strings.Join(filterfunc, "or") + ")", paginate, nil

}

// CreateFieldsQuery translates the fields part of the qsl string to dgraph
// input @name,@resourceversion, metadata fields for this block's object type,
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
	return returnlist, nil

}

// CreateDgraphQuery translates the querystring to a dgraph query
func (qa *QSLService) CreateDgraphQuery(query string, cntOnly bool) (string, error) {
	log.Info("Received Query: ", strings.Split(query, "}."))

	// remove all whitespace
	whitespace := regexp.MustCompile("\\s*")
	querys := whitespace.ReplaceAllString(query, "")

	// e.g. cluster[@name="cluster1.k8s.local"]{@name,@region}.pod[@name="pod1"]{@phase,@image}
	// split by }. to get each individual block
	// cannot split by . because . may be present in some object names
	splitQuery := strings.Split(querys, "}.")
	rootTemplate := "{ A as var(func: eq(objtype, $OBJTYPE)) $FILTERSFUNC @cascade {"
	edgeTemplate := "\t$RELATION @filter(eq(objtype, $OBJTYPE) $FILTERSFUNC)"
	pageTemplate := "{ objects(func: uid(A)$PAGINATE) {"
	brakets := []string{"}", "}"}

	root, objType, rootCntFilter, err := qa.buildRootQuery(splitQuery[0], rootTemplate, true)
	if err != nil {
		return "", err
	}
	parentType := objType
	edgeCntFilters := []string{}
	for i := 1; i < len(splitQuery); i++ {
		edges, ptype, edgeCntFilter, err := qa.buildEdgeQuery(splitQuery[i], edgeTemplate, parentType, true)
		if err != nil {
			return "", err
		}
		parentType = ptype
		root = append(root, edges...)
		brakets = append(brakets, "}")
		if edgeCntFilter != "" {
			edgeCntFilters = append(edgeCntFilters, edgeCntFilter)
		}
	}
	root = append(root, brakets...)
	pages, _, _, err := qa.buildRootQuery(splitQuery[0], pageTemplate, cntOnly)
	if err != nil {
		return "", err
	}
	root = append(root, pages...)
	parentType = objType
	for i := 1; i < len(splitQuery); i++ {
		edges, ptype, _, err := qa.buildEdgeQuery(splitQuery[i], edgeTemplate, parentType, cntOnly)
		if err != nil {
			return "", err
		}
		parentType = ptype
		root = append(root, edges...)
	}
	root = append(root, brakets...)
	if rootCntFilter != "" {
		root = append(root, rootCntFilter)
	}
	if len(edgeCntFilters) > 0 {
		root = append(root, edgeCntFilters...)
	}

	return strings.Join(root, "\n"), nil
}

func (qa *QSLService) buildRootQuery(qry string, template string, cntOnly bool) ([]string, string, string, error) {
	ret := []string{template}

	objType, filters, fields, err := parseQuery(qry)
	if err != nil {
		return nil, "", "", err
	}
	// create var for dgraph if filter has count
	cntFilterQry, err := qa.getCntFilter(filters, objType)
	if err != nil {
		return nil, "", "", err
	}

	_, ff, pag, err := CreateFiltersQuery(filters)
	if err != nil {
		return nil, "", "", err
	}
	// replace the filters and object type and add the list of fields
	ret[0] = strings.Replace(ret[0], "$FILTERSFUNC", ff, -1)
	ret[0] = strings.Replace(ret[0], "$OBJTYPE", objType, -1)
	if cntOnly {
		ret = append(ret, "\tcount(uid)")
		if strings.Contains(template, "$PAGINATE") {
			ret[0] = strings.Replace(ret[0], "$PAGINATE", "", -1)
		}
	} else {
		// get metadata fields for projection
		metafieldslist, err := qa.GetMetadata(objType)
		if err != nil {
			return nil, "", "", err
		}
		fl, err := CreateFieldsQuery(fields, metafieldslist, 0)
		if err != nil {
			return nil, "", "", err
		}
		ret = append(ret, fl...)
		if strings.Contains(template, "$PAGINATE") {
			if pag != "" {
				ret[0] = strings.Replace(ret[0], "$PAGINATE", pag, -1)
			} else {
				ret[0] = strings.Replace(ret[0], "$PAGINATE", ",first:1000,offset:0", -1)
			}
		}
	}
	return ret, objType, cntFilterQry, nil
}

func (qa *QSLService) buildEdgeQuery(qry string, template string, parent string, cntOnly bool) ([]string, string, string, error) {
	ret := []string{template}
	objType, filters, fields, err := parseQuery(qry)
	if err != nil {
		return nil, "", "", err
	}
	// get a list of the metadata fields for this object type
	metafieldslist, err := qa.GetMetadata(objType)
	if err != nil {
		return nil, "", "", errors.New("Failed to connect to dgraph to get metadata")
	}

	// declare relation variable
	relation, err := qa.getRelationName(objType, parent)
	if err != nil {
		return nil, "", "", err
	}

	// no relation found between the two objects
	if relation == "" {
		return nil, "", "", errors.New("no relation found between " + objType + " and " + parent)
	}

	// create var for dgraph if filter has count
	cntFilterQry, err := qa.getCntFilter(filters, objType)
	if err != nil {
		return nil, "", "", err
	}
	_, ff, pag, err := CreateFiltersQuery(filters)
	if err != nil {
		return nil, "", "", err
	}
	if len(ff) > 0 {
		ff = "and" + ff[8:len(ff)-1]
	}
	// replace filters and object type accordingly
	ret[0] = strings.Replace(ret[0], "$FILTERSFUNC", ff, -1)
	ret[0] = strings.Replace(ret[0], "$OBJTYPE", objType, -1)
	// add the tilde because we are adding an inverse relationship
	ret[0] = strings.Replace(ret[0], "$RELATION", relation, -1)

	if cntOnly {
		ret[0] += "{"
		ret = append(ret, "\tcount(uid)")
	} else {
		// if pagination values were supplied, add and get rid of the comma at the beginning
		if len(pag) > 1 {
			ret[0] += "(" + pag[1:] + ")"
		} else {
			ret[0] += "(first:1000,offset:0)"
		}
		ret[0] += "{"
		fl, err := CreateFieldsQuery(fields, metafieldslist, 1)
		if err != nil {
			return nil, "", "", err
		}
		// append the fields to be returned
		ret = append(ret, fl...)
	}
	return ret, objType, cntFilterQry, nil
}

func parseQuery(qry string) (string, string, string, error) {
	r := regexp.MustCompile(blockRegex)
	matches := r.FindStringSubmatch(qry)
	if len(matches) < 2 {
		log.Error("Malformed Query received: " + qry)
		return "", "", "", errors.New("Malformed Query: " + qry)
	}
	// extract the values of the form objtype[filters]fields and assign to individual variables
	objType := strings.ToLower(matches[1])
	filters := matches[2]
	fields := strings.ToLower(matches[3])
	return objType, filters, fields, nil
}

func (qa *QSLService) getCntFilter(query, objType string) (string, error) {
	ret := ""
	// create var for dgraph if filter has count
	if strings.Contains(query, "count(") {
		unix32bits := uint32(time.Now().UTC().Unix())
		buff := make([]byte, 4)
		rand.Read(buff)
		seq := fmt.Sprintf("%x%x", unix32bits, buff)
		tmp := query[strings.Index(query, "count(")+6:]
		relType := strings.ToLower(tmp[:strings.Index(tmp, ")")])
		relation, err := qa.getRelationName(relType, objType)
		if err != nil {
			return "", err
		}
		cntTemplate := `{ var (func: eq(objtype, %s)) { %s @filter (eq (objtype, %s)) { cnt%s as count(name) } cnt_%s as sum(val(cnt%s)) }}`
		ret = fmt.Sprintf(cntTemplate, objType, relation, relType, seq, relType, seq)
	}
	return ret, nil
}

func (qa *QSLService) getRelationName(objType string, parent string) (string, error) {
	// get a list of the metadata fields for this object type
	metafieldslist, err := qa.GetMetadata(objType)
	if err != nil {
		return "", errors.New("Failed to connect to dgraph to get metadata")
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
			log.Debugf("1 found relationship for %s-%s->%s", objType, item.FieldName, item.RefDataType)
			for _, dtype := range strings.Split(item.RefDataType, ",") {
				if dtype == parent {
					relation = "~" + strings.ToLower(item.FieldName)
					found = true
					break
				}
			}
		}
	}

	if !found {
		// if not, see if we can find the relation from the parent to this object
		metafieldslist2, err := qa.metaSvc.GetMetadataFields(parent)
		if err != nil {
			log.Error(err)
			return "", errors.New("Failed to connect to dgraph to get metadata")
		}
		log.Debugf("couldn't find relation for %s->%s,", parent, objType)
		log.Debugf("metadata fields for %s: %#v", parent, metafieldslist)

		for _, item := range metafieldslist2 {
			if item.FieldType == "relationship" {
				log.Debugf("2 found relationship for %s-%s->%s", parent, item.FieldName, item.RefDataType)
				for _, dtype := range strings.Split(item.RefDataType, ",") {
					if dtype == objType {
						relation = strings.ToLower(item.FieldName)
						found = true
						break
					}
				}
			}
		}
	}
	return relation, nil
}
