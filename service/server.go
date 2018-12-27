package main

import (
	"flag"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	lru "github.com/hashicorp/golang-lru"
	"github.com/intuit/katlas/service/apis"
	"github.com/intuit/katlas/service/cfg"
	"github.com/intuit/katlas/service/db"
	"github.com/intuit/katlas/service/resources"
)

const CacheSize = 10

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
	entitySvc := apis.NewEntityService(dc, metaSvc)
	querySvc := apis.NewQueryService(dc)
	res := resources.ServerResource{EntitySvc: entitySvc, QuerySvc: querySvc, MetaSvc: metaSvc}
	router.HandleFunc("/v1/entity/{metadata}/{uid}", res.EntityGetHandler).Methods("GET")
	// TODO: wire up more resource APIs here
	router.HandleFunc("/v1/entity/{metadata}", res.EntityCreateHandler).Methods("POST")
	router.HandleFunc("/v1/sync/{metadata}", res.EntitySyncHandler).Methods("POST")
	router.HandleFunc("/v1/entity/{metadata}/{resourceid}", res.EntityDeleteHandler).Methods("DELETE")
	router.HandleFunc("/v1/query", res.QueryHandler).Methods("GET")

	//Metadata
	router.HandleFunc("/v1/meta/{name}", res.MetaGetHandler).Methods("GET")

	router.HandleFunc("/health", Health).Methods("GET")
	router.HandleFunc("/", Up).Methods("GET", "POST")

	//Creates an LRU cache of the given size
	var err error
	db.LruCache, err = lru.New(CacheSize)
	if err != nil {
		log.Errorf("err: %v", err)
	}
	log.Infoln("LRU cache created with given size")

	log.Infof("Service started on port:8011, mode:%s", cfg.ServerCfg.ServerType)

	if strings.EqualFold(cfg.ServerCfg.ServerType, "https") {
		log.Fatal(http.ListenAndServeTLS(":8011", "server.crt", "server.key", router))
	} else {
		log.Fatal(http.ListenAndServe(":8011", router))
	}
}

func main() {
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
