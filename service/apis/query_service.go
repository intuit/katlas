package apis

import (
	"github.com/intuit/katlas/service/db"
)

//IQueryService ...define interfaces to query data
type IQueryService interface {
	GetQueryResult(queryMap map[string][]string) (map[string]interface{}, error)
}

//QueryService ...
type QueryService struct {
	dbclient db.IDGClient
}

//NewQueryService ...Create NewQueryService
func NewQueryService(dc *db.DGClient) *QueryService {
	return &QueryService{dc}
}

//GetQueryResult ...Api to get Query Results
func (s QueryService) GetQueryResult(queryMap map[string][]string) (map[string]interface{}, error) {
	return s.dbclient.GetQueryResult(queryMap)
}
