package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/hashicorp/golang-lru"
	"github.com/intuit/katlas/service/apis"
	"github.com/intuit/katlas/service/cfg"
	"github.com/intuit/katlas/service/db"
	"github.com/intuit/katlas/service/resources"
	"github.com/intuit/katlas/service/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
)

const cacheSize = 10

//Health checks service health
func Health(w http.ResponseWriter, r *http.Request) {
	log.Info("RestService is still running")
	w.Write([]byte("RestService is still running"))

}

//Up checks service status
func Up(w http.ResponseWriter, r *http.Request) {
	log.Info("Up")
	w.Write([]byte("Up"))
}

func serve() {
	router := mux.NewRouter()

	dc := db.NewDGClient(cfg.ServerCfg.DgraphHost)
	defer dc.Close()
	metaSvc := apis.NewMetaService(dc)
	entitySvc := apis.NewEntityService(dc)
	querySvc := apis.NewQueryService(dc)
	qslSvc := apis.NewQSLService(dc)
	res := resources.ServerResource{EntitySvc: entitySvc, QuerySvc: querySvc, MetaSvc: metaSvc, QSLSvc: qslSvc}
	// Entity APIs v1
	router.HandleFunc("/v1/entity/{metadata}/{uid}", res.EntityGetHandler).Methods("GET")
	router.HandleFunc("/v1/entity/{metadata}", res.EntityCreateHandler).Methods("POST")
	router.HandleFunc("/v1/entity/{metadata}/{uid}", res.EntityUpdateHandler).Methods("POST")
	router.HandleFunc("/v1/entity/{metadata}/{resourceid}", res.EntityDeleteHandler).Methods("DELETE")
	router.HandleFunc("/v1/sync/{metadata}", res.EntitySyncHandler).Methods("POST")
	// Query APIs
	router.HandleFunc("/v1/query", res.QueryHandler).Methods("GET")
	router.HandleFunc("/v1/qsl/{query}", res.QSLHandler).Methods("GET")
	//Metadata
	router.HandleFunc("/v1/metadata/{name}", res.MetaGetHandler).Methods("GET")
	router.HandleFunc("/v1/metadata/{name}", res.MetaDeleteHandler).Methods("DELETE")
	router.HandleFunc("/v1/metadata", res.MetaCreateHandler).Methods("POST")
	router.HandleFunc("/v1/metadata/{name}", res.MetaUpdateHandler).Methods("POST")
	router.HandleFunc("/v1/schema", res.SchemaUpsertHandler).Methods("POST")
	router.HandleFunc("/v1/schema/{name}", res.SchemaDropHandler).Methods("DELETE")

	// Entity APIs v1.1
	router.HandleFunc("/v1.1/entity/{uid}", res.EntityGetHandlerV1_1).Methods("GET")
	router.HandleFunc("/v1.1/entity", res.EntityCreateHandlerV1_1).Methods("POST")
	router.HandleFunc("/v1.1/entity/{uid}", res.EntityUpdateHandlerV1_1).Methods("POST")
	router.HandleFunc("/v1.1/entity/{uid}", res.EntityDeleteHandlerV1_1).Methods("DELETE")
	router.HandleFunc("/v1.1/sync/{metadata}", res.EntitySyncHandlerV1_1).Methods("POST")
	// Query APIs v1.1
	router.HandleFunc("/v1.1/query", res.QueryHandlerV1_1).Methods("GET")
	router.HandleFunc("/v1.1/qsl/{query}", res.QSLHandlerV1_1).Methods("GET")
	//Metadata v1.1
	router.HandleFunc("/v1.1/metadata/{name}", res.MetaGetHandlerV1_1).Methods("GET")
	router.HandleFunc("/v1.1/metadata/{name}", res.MetaDeleteHandlerV1_1).Methods("DELETE")
	router.HandleFunc("/v1.1/metadata", res.MetaCreateHandlerV1_1).Methods("POST")
	router.HandleFunc("/v1.1/metadata/{name}", res.MetaUpdateHandlerV1_1).Methods("POST")
	router.HandleFunc("/v1.1/schema", res.SchemaUpsertHandlerV1_1).Methods("POST")
	router.HandleFunc("/v1.1/schema/{name}", res.SchemaDropHandlerV1_1).Methods("DELETE")

	// Status
	router.HandleFunc("/health", Health).Methods("GET")
	router.HandleFunc("/", Up).Methods("GET", "POST")
	// Monitoring
	router.Handle("/prometheus_metrics", promhttp.Handler()).Methods("GET")
	util.RegisterHistogramMetrics()

	//Creates an LRU cache of the given size
	var err error
	db.LruCache, err = lru.New(cacheSize)
	if err != nil {
		log.Errorf("err: %v", err)
	}
	log.Infoln("LRU cache created with given size")
	log.Infoln("Starting initialize schema and metadata... ")
	// Create dgraph schema
	data, err := ioutil.ReadFile("data/dbschema.json")
	if err != nil {
		log.Fatalf("Schema file error: %v\n", err)
	}
	var predicates []db.Schema
	json.Unmarshal(data, &predicates)
	for _, p := range predicates {
		metaSvc.CreateSchema(p)
	}
	// Initialize metadata
	meta, err := ioutil.ReadFile("data/meta.json")
	if err != nil {
		log.Fatalf("Metadata file error: %v\n", err)
	}
	var jsonData []map[string]interface{}
	json.Unmarshal(meta, &jsonData)
	for _, data := range jsonData {
		qm := map[string][]string{util.Name: {data[util.Name].(string)}, util.ObjType: {util.Metadata}}
		metas, _ := querySvc.GetQueryResult(qm)
		if len(metas[util.Objects].([]interface{})) > 0 {
			metaSvc.UpdateMetadata(data[util.Name].(string), data)
		} else {
			metaSvc.CreateMetadata(data)
		}
	}
	log.Infof("Service started on port:8011, mode:%s", cfg.ServerCfg.ServerType)
	if strings.EqualFold(cfg.ServerCfg.ServerType, "https") {
		log.Fatal(http.ListenAndServeTLS(":8011", "server.crt", "server.key", router))
	} else {
		log.Fatal(http.ListenAndServe(":8011", router))
	}
}

func main() {
	log.SetLevel(log.DebugLevel)
	// parse and print command line flags
	flag.Parse()
	log.Infof("EnvNamespace=%s", cfg.ServerCfg.EnvNamespace)
	log.Infof("ServerType=%s", cfg.ServerCfg.ServerType)
	log.Infof("DgraphHost=%s", cfg.ServerCfg.DgraphHost)

	if cfg.ServerCfg.DgraphHost == "" {
		flag.PrintDefaults()
		log.Fatal("Invalid input params")
	}
	serve()
}
