package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/intuit/katlas/service/apis"
	"github.com/intuit/katlas/service/db"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	dbHost := flag.String("dbhost", "localhost", "dgraph server name")
	port := flag.Int("port", 9080, "dgraph server port")
	flag.Parse()

	dc := db.NewDGClient(*dbHost + ":" + strconv.Itoa(*port))
	// create dgraph schema
	data, err := ioutil.ReadFile("./dbschema.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	var predicates []db.Schema
	json.Unmarshal(data, &predicates)
	for _, p := range predicates {
		dc.CreateSchema(p)
		fmt.Printf("Schema %s was created\n", p.Predicate)
	}

	m := apis.NewMetaService(dc)
	e := apis.NewEntityService(dc, m)

	// create k8s metadata
	meta, err := ioutil.ReadFile("./meta.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	var jsonData []map[string]interface{}
	json.Unmarshal(meta, &jsonData)
	for _, data := range jsonData {
		e.CreateEntity("metadata", data)
	}
}
