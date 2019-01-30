package apis

import (
	"fmt"
	"strconv"
	"strings"

	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/intuit/katlas/service/db"
	"github.com/intuit/katlas/service/util"
)

//QueryParamKeyword ...param for keyword query
const QueryParamKeyword = "keyword"

//IQueryService ...define interfaces to query data
type IQueryService interface {
	GetQueryResult(queryMap map[string][]string) (map[string]interface{}, error)
}

//QueryService ...
type QueryService struct {
	dbclient db.IDGClient
}

//NewQueryService ...Create NewQueryService
func NewQueryService(dc db.IDGClient) *QueryService {
	return &QueryService{dc}
}

//GetQueryResult ...Api to get Query Results
func (s QueryService) GetQueryResult(queryMap map[string][]string) (map[string]interface{}, error) {
	var err error
	// default limit is 1000
	first, offset := 1000, 0
	// limit should be number and less than 10000
	if val, ok := queryMap[util.First]; ok {
		first, err = strconv.Atoi(val[0])
		if err != nil || first > MaximumLimit {
			return nil, fmt.Errorf("pagination format error or exceeding maxiumum limit %d", MaximumLimit)
		}
	}
	// offset should be number
	if val, ok := queryMap[util.Offset]; ok {
		offset, err = strconv.Atoi(val[0])
		if err != nil {
			return nil, err
		}
	}
	// keyword search
	if val, ok := queryMap[QueryParamKeyword]; ok {
		if val[0] == "" {
			err := fmt.Errorf("Value not specified for Query Param [%s]", QueryParamKeyword)
			return nil, err
		}
		// generate queries include count query
		q, cntQry, err := s.getQueryResultByKeyword(val[0], first, offset)
		if err != nil {
			log.Debug(err)
			return nil, err
		}
		// execute query to get result
		ret, err := s.dbclient.GetQueryResult(cntQry)
		if err != nil {
			log.Debug(err)
			return nil, err
		}
		total := GetTotalCnt(ret)
		ret, err = s.dbclient.GetQueryResult(q)
		if err != nil {
			log.Debug(err)
			return nil, err
		}
		ret[util.Count] = total
		return ret, nil
	}
	// key value query
	if len(queryMap) == 0 {
		err := fmt.Errorf("Query Params not specified")
		return nil, err
	}
	q, cntQry := getQueryResultByKeyValue(queryMap, first, offset)
	ret, err := s.dbclient.GetQueryResult(cntQry)
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	total := GetTotalCnt(ret)
	ret, err = s.dbclient.GetQueryResult(q)
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	ret[util.Count] = total
	return ret, nil
}

// GetTotalCnt find count from returned value
func GetTotalCnt(data map[string]interface{}) float64 {
	var total float64
	for _, obj := range data[util.Objects].([]interface{}) {
		val, ok := obj.(map[string]interface{})[util.Count]
		if ok {
			total = val.(float64)
		}
	}
	return total
}

// Keyword query http://<dgraph ip:port>/v1/query?keyword=pod
func (s QueryService) getQueryResultByKeyword(keyword string, first, offset int) (string, string, error) {
	smds, err := s.dbclient.GetSchemaFromCache(db.LruCache)
	if err != nil {
		log.Debug(err)
		return "", "", err
	}
	// generate query as following
	//{
	//	A as var(func: regexp(name, /test/i)) {}
	//	B as var(func: regexp(labels, /test/i)) {}
	//	me(func: uid(A,B), first:1000,offset:0) {
	//	expand(_all_) {
	//	expand(_all_)
	//}}}
	cnt := 0
	statements := []string{"{"}
	for _, schemanode := range smds {
		if schemanode.Type == "string" && schemanode.Index == true && len(schemanode.Tokenizer) > 0 {
			for _, tokenizer := range schemanode.Tokenizer {
				tk := tokenizer
				if tk == "trigram" {
					filter := "obj" + strconv.Itoa(cnt) + " as var(func:regexp(" + schemanode.Predicate + ",/" + keyword + "/i)) {}"
					statements = append(statements, filter)
					cnt++
				}
			}
		}
	}
	var buf bytes.Buffer
	for i := 0; i < cnt; i++ {
		buf.WriteString("obj")
		buf.WriteString(strconv.Itoa(i))
		if i < cnt-1 {
			buf.WriteString(",")
		}
	}

	cntOnlyStatements := make([]string, len(statements))
	copy(cntOnlyStatements, statements)
	cntTemplate := `objects(func: uid(%s)) { %s }`
	cntQuery := fmt.Sprintf(cntTemplate, buf.String(), "count(uid)")
	cntOnlyStatements = append(cntOnlyStatements, cntQuery)
	cntOnlyStatements = append(cntOnlyStatements, "}")
	template := `objects(func: uid(%s), first:%d,offset:%d) { %s }`
	query := fmt.Sprintf(template, buf.String(), first, offset, "uid expand(_all_) { uid expand(_all_) }")
	statements = append(statements, query)
	statements = append(statements, "}")
	return strings.Join(statements, "\n"), strings.Join(cntOnlyStatements, "\n"), nil
}

// Key-Value query http://<dgraph ip:port>/v1/query?name=pod01&objtype=Pod
func getQueryResultByKeyValue(queryMap map[string][]string, first, offset int) (string, string) {
	//Only indexed fields can be filtered on
	//Time must be in correct format "2018-10-18 14:36:32 -0700 PDT"
	qps := []string{}
	var funcStr, filterStr string
	for k, v := range queryMap {
		if k != util.First && k != util.Offset && k != util.Print {
			qp := "eq(" + k + ",\"" + v[0] + "\")"
			qps = append(qps, qp)
		}
	}
	print := "expand(_all_)"
	if p, ok := queryMap[util.Print]; ok {
		if "*" != p[0] {
			print = p[0]
		}
	}
	funcStr = fmt.Sprintf("(func:%s, first:%d, offset:%d) ", qps[0], first, offset)
	cntStr := fmt.Sprintf("(func:%s)", qps[0])
	filters := qps[1:]
	if len(filters) > 0 {
		filterStr = "@filter(" + strings.Join(filters, " AND ") + ")"
	}
	q := fmt.Sprintf(`{objects %s %s {uid %s { uid %s }}}`, funcStr, filterStr, print, print)
	cntQry := fmt.Sprintf(`{objects %s %s {count(uid)}}`, cntStr, filterStr)
	return q, cntQry
}
