package apis

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/intuit/katlas/service/db"
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
	var q string
	var err error

	val, ok := queryMap[QueryParamKeyword]
	if ok {
		if val[0] == "" {
			err := fmt.Errorf("Value not specified for Query Param [%s]", QueryParamKeyword)
			return nil, err
		}
		q, err = s.getQueryResultByKeyword(val[0])
		if err != nil {
			log.Debug(err)
			return nil, err
		}
	} else {
		if len(queryMap) == 0 {
			err := fmt.Errorf("Query Params not specified")
			return nil, err
		}
		q = getQueryResultByKeyValue(queryMap)
	}

	return s.dbclient.GetQueryResult(q)
}

// Keyword query http://<dgraph ip:port>/v1/query?keyword=pod
func (s QueryService) getQueryResultByKeyword(keyword string) (string, error) {
	smds, err := s.dbclient.GetSchema()
	if err != nil {
		log.Debug(err)
		return "", err
	}
	cnt := 0
	var qr string
	qr = "{"

	for _, schemanode := range smds {
<<<<<<< HEAD
		log.Debugf("Predicate: %v Type: %v tokenizer: %v\n", schemanode.Predicate, schemanode.Type, schemanode.Tokenizer)
=======
		fmt.Printf("Predicate: %v Type: %v tokenizer: %v\n", schemanode.Predicate, schemanode.Type, schemanode.Tokenizer)
>>>>>>> cb3cb8587ca3500ecf39dc3806791f1784aaa65a

		if schemanode.Type == "string" && schemanode.Index == true && len(schemanode.Tokenizer) > 0 {
			for _, tokenizer := range schemanode.Tokenizer {
				tk := tokenizer
				if tk == "trigram" {
<<<<<<< HEAD
					log.Debugf("Found ***** Predicate: %v Type: %v tokenizer: %v\n", schemanode.Predicate, schemanode.Type, schemanode.Tokenizer)
=======
					fmt.Printf("Found ***** Predicate: %v Type: %v tokenizer: %v\n", schemanode.Predicate, schemanode.Type, schemanode.Tokenizer)
>>>>>>> cb3cb8587ca3500ecf39dc3806791f1784aaa65a
					filter := "obj" + strconv.Itoa(cnt) + "(func:regexp(" + schemanode.Predicate + ",/" + keyword + "/i)) {"
					qr = qr + filter + `
						uid
	            		expand(_all_) {
							uid
							expand(_all_)
						}
					}
					`
				}

			}
		}
		cnt++
	}
	qr = qr + "}"
	log.Debugf("Query string is =%v\n", qr)
	return qr, nil
}

// Key-Value query http://<dgraph ip:port>/v1/query?name=pod01&objtype=Pod
func getQueryResultByKeyValue(queryMap map[string][]string) string {

	//Only indexed fields can be filtered on
	//Time must be in correct format "2018-10-18 14:36:32 -0700 PDT"
	qps := []string{}
	var funcStr, filterStr string

	for k, v := range queryMap {
		qp := "eq(" + k + ",\"" + v[0] + "\")"
		qps = append(qps, qp)
	}

	funcStr = "(func:" + qps[0] + ") "
	filters := qps[1:]
	if len(filters) > 0 {
		filterStr = "@filter(" + strings.Join(filters, " AND ") + ")"
	}

	q := `
	{
		objects` + funcStr + filterStr + ` {
			uid
			expand(_all_) {
				uid
				expand(_all_)
			}
		}
	}
	`
	return q
}
