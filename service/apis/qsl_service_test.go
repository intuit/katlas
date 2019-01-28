package apis

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/intuit/katlas/service/db"
)

// FResult values is the expected output, err is the expected error
type FResult struct {
	values []string
	err    error
}

func TestCreateFiltersQuery(t *testing.T) {

	tests := map[string]FResult{
		`@name="paas-preprod-west2.cluster.k8s.local"||@k8sobj="K8sObj"||@resourceid="paas-preprod-west2.cluster.k8s.local"`: FResult{
			[]string{
				", $name: string, $k8sobj: string, $resourceid: string",
				`@filter( eq(name,"paas-preprod-west2.cluster.k8s.local") or eq(k8sobj,"K8sObj") or eq(resourceid,"paas-preprod-west2.cluster.k8s.local") )`,
				"",
			},
			nil,
		},
		`@name="paas-preprod-west2.cluster.k8s.local"&&@k8sobj="K8sObj"&&@resourceid="paas-preprod-west2.cluster.k8s.local"`: FResult{
			[]string{
				", $name: string, $k8sobj: string, $resourceid: string",
				`@filter( eq(name,"paas-preprod-west2.cluster.k8s.local") and eq(k8sobj,"K8sObj") and eq(resourceid,"paas-preprod-west2.cluster.k8s.local") )`,
				"",
			},
			nil,
		},
		`@name="paas-preprod-west2.cluster.k8s.local"||@k8sobj="K8sObj"&&@resourceid="paas-preprod-west2.cluster.k8s.local"`: FResult{
			[]string{
				", $name: string, $k8sobj: string, $resourceid: string",
				`@filter( eq(name,"paas-preprod-west2.cluster.k8s.local") or eq(k8sobj,"K8sObj") and eq(resourceid,"paas-preprod-west2.cluster.k8s.local") )`,
				"",
			},
			nil,
		},
		`@name="paas-preprod-west2.cluster.k8s.local"||@k8sobj="K8sObj"`: FResult{
			[]string{
				", $name: string, $k8sobj: string",
				`@filter( eq(name,"paas-preprod-west2.cluster.k8s.local") or eq(k8sobj,"K8sObj") )`,
				"",
			},
			nil,
		},
		`@k8sobj="K8sObj"&&@resourceid="paas-preprod-west2.cluster.k8s.local"`: FResult{
			[]string{
				", $k8sobj: string, $resourceid: string",
				`@filter( eq(k8sobj,"K8sObj") and eq(resourceid,"paas-preprod-west2.cluster.k8s.local") )`,
				"",
			},
			nil,
		},
		`@name="paas-preprod-west2.cluster.k8s.local"`: FResult{
			[]string{
				", $name: string",
				`@filter( eq(name,"paas-preprod-west2.cluster.k8s.local") )`,
				"",
			},
			nil,
		},
		`@name="paas-preprod-west2.cluster.k8s.local?"`: FResult{
			[]string{
				", $name: string",
				`@filter( eq(name,"paas-preprod-west2.cluster.k8s.local") )`,
				"",
			},
			errors.New("Invalid filters in @name=\"paas-preprod-west2.cluster.k8s.local?\""),
		},
		`@numreplicas>=1`: FResult{
			[]string{
				", $numreplicas: int",
				`@filter( ge(numreplicas,1) )`,
				"",
			},
			nil,
		},
		`@name="paas-preprod-west2.cluster.k8s.local"$$first=2`: FResult{
			[]string{
				", $name: string",
				`@filter( eq(name,"paas-preprod-west2.cluster.k8s.local") )`,
				",first: 2",
			},
			nil,
		},
		`@name="paas-preprod-west2.cluster.k8s.local"$$first=2,offset=4`: FResult{
			[]string{
				", $name: string",
				`@filter( eq(name,"paas-preprod-west2.cluster.k8s.local") )`,
				",first: 2,offset: 4",
			},
			nil,
		},
	}

	for k, v := range tests {
		realOut1, realOut2, realPag, err := CreateFiltersQuery(k)
		if err == nil {
			if !(v.values[0] == realOut1) {
				t.Errorf("filter declaration incorrect\n input: %s\n testdec: %s\n realdec: %s", k, v.values[0], realOut1)
			}

			if !(v.values[1] == realOut2) {
				t.Errorf("filter function incorrect\n input: %s\n testfunc: %s\n realfunc: %s", k, v.values[1], realOut2)
			}

			if !(v.values[2] == realPag) {
				t.Errorf("filter pagination incorrect\n input: %s\n testpag: %s\n realpag: %s", k, v.values[2], realPag)
			}
		} else {
			if v.err != nil {
				if err.Error() != v.err.Error() {
					t.Errorf("filter error incorrect\n input: %s\n testfiltererr: %s\n realfiltererr: %s", k, v.err.Error(), err.Error())
				}
			} else {
				t.Errorf("filter error incorrect\n input: %s\n testfiltererr: %s\n realfiltererr: %s", k, "no error expected", err.Error())
			}

		}

	}

}

func TestCreateFieldsQuery(t *testing.T) {
	tests := map[string]FResult{
		"*":                                  FResult{[]string{"k8sobj", "objtype", "name", "resourceid", "resourceversion", "uid"}, nil},
		"**":                                 FResult{[]string{"expand(_all_){", "\texpand(_all_){", "\t}", "}"}, nil},
		"***":                                FResult{[]string{"expand(_all_){", "\texpand(_all_){", "\t\texpand(_all_){", "\t\t}", "\t}", "}"}, nil},
		"?":                                  FResult{nil, errors.New("Field names must be prefixed with @ sign and followed by an alphanumeric field name [?]")},
		"name":                               FResult{nil, errors.New("Field names must be prefixed with @ sign and followed by an alphanumeric field name [name]")},
		"@n@me":                              FResult{nil, errors.New("Field names must be composed of only alphanumeric characters [n@me]")},
		"@*":                                 FResult{nil, errors.New("Field names must be composed of only alphanumeric characters [*]")},
		"*@":                                 FResult{nil, errors.New("Fields may be a string of * indicating how many levels, or a list of fields @field1,@field2,... not both [*@]")},
		"*,@name":                            FResult{nil, errors.New("Fields may be a string of * indicating how many levels, or a list of fields @field1,@field2,... not both [*,@name]")},
		"@name,**":                           FResult{nil, errors.New("Field names must be prefixed with @ sign and followed by an alphanumeric field name [**]")},
		"**,*":                               FResult{nil, errors.New("Fields may be a string of * indicating how many levels, or a list of fields @field1,@field2,... not both [**,*]")},
		"@name":                              FResult{[]string{"name", "uid"}, nil},
		"@name,@resourceversion":             FResult{[]string{"name", "resourceversion", "uid"}, nil},
		"@name,@resourceversion,@resourceid": FResult{[]string{"name", "resourceversion", "resourceid", "uid"}, nil},
		"@name,@resourceversion,@k8sobj":     FResult{[]string{"name", "resourceversion", "k8sobj", "uid"}, nil},
	}

	namespacemetafieldslist := []MetadataField{MetadataField{FieldName: "k8sobj", FieldType: "string", Mandatory: true, RefDataType: "", Cardinality: "one"}, MetadataField{FieldName: "objtype", FieldType: "string", Mandatory: true, RefDataType: "", Cardinality: "one"}, MetadataField{FieldName: "name", FieldType: "string", Mandatory: true, RefDataType: "", Cardinality: "one"}, MetadataField{FieldName: "resourceid", FieldType: "string", Mandatory: false, RefDataType: "", Cardinality: "one"}, MetadataField{FieldName: "resourceversion", FieldType: "string", Mandatory: true, RefDataType: "", Cardinality: "one"}}

	for k, v := range tests {
		output, err := CreateFieldsQuery(k, namespacemetafieldslist, -1)
		if err == nil {
			if !(reflect.DeepEqual(v.values, output)) {
				t.Errorf("fields incorrect\n input: %s\n testfields: %s\n realfields: %s", k, v.values, output)
			}
		} else {
			if v.err != nil {
				if err.Error() != v.err.Error() {
					t.Errorf("field error incorrect\n input: %s\n testfielderr: %s\n realfielderr: %s", k, v.err.Error(), err.Error())
				}
			} else {
				t.Errorf("field error incorrect\n input: %s\n testfielderr: %s\n realfielderr: %s", k, "no error expected", err.Error())
			}

		}

	}
}

func TestCreateDgraphQuery(t *testing.T) {
	tests := map[string]FResult{
		`namespace[@name="default"]{*}`: FResult{[]string{
			"{ A as var(func: eq(objtype, namespace)) @filter( eq(name,\"default\") ) @cascade {",
			"\tcount(uid)",
			"}",
			"}",
			"{ objects(func: uid(A),first:1000,offset:0) {",
			"\tresourceversion",
			"\tcreationtime",
			"\tk8sobj",
			"\tobjtype",
			"\tname",
			"\tresourceid",
			"\tlabels",
			"\tuid",
			"}",
			"}",
		}, nil},
		`cluster[@name="paas-preprod-west2.cluster.k8s.local"]{@name}`: FResult{[]string{
			"{ A as var(func: eq(objtype, cluster)) @filter( eq(name,\"paas-preprod-west2.cluster.k8s.local\") ) @cascade {",
			"\tcount(uid)",
			"}",
			"}",
			"{ objects(func: uid(A),first:1000,offset:0) {",
			"\tname",
			"\tuid",
			"}",
			"}",
		}, nil},
		`cluster[@name="paas-preprod-west2.cluster.k8s.local"]{*}.namespace[@name="opa"||@name="default"]{*}`: FResult{[]string{
			"{ A as var(func: eq(objtype, cluster)) @filter( eq(name,\"paas-preprod-west2.cluster.k8s.local\") ) @cascade {",
			"\tcount(uid)",
			"\t~cluster @filter(eq(objtype, namespace) and eq(name,\"opa\") or eq(name,\"default\") ){",
			"\tcount(uid)",
			"}",
			"}",
			"}",
			"{ objects(func: uid(A),first:1000,offset:0) {",
			"\tresourceid",
			"\tresourceversion",
			"\tcreationtime",
			"\tk8sobj",
			"\tobjtype",
			"\tname",
			"\tuid",
			"\t~cluster @filter(eq(objtype, namespace) and eq(name,\"opa\") or eq(name,\"default\") )(first:1000,offset:0){",
			"\t	objtype",
			"\t	name",
			"\t	resourceid",
			"\t	labels",
			"\t	resourceversion",
			"\t	creationtime",
			"\t	k8sobj",
			"\t	uid",
			"}",
			"}",
			"}",
		}, nil},
		`cluster[@name="paas-preprod-west2.cluster.k8s.local"]{*}.namespace[@name="opa"&&@k8sobj="k8sobj"]{*}`: FResult{[]string{
			"{ A as var(func: eq(objtype, cluster)) @filter( eq(name,\"paas-preprod-west2.cluster.k8s.local\") ) @cascade {",
			"\tcount(uid)",
			"\t~cluster @filter(eq(objtype, namespace) and eq(name,\"opa\") and eq(k8sobj,\"k8sobj\") ){",
			"\tcount(uid)",
			"}",
			"}",
			"}",
			"{ objects(func: uid(A),first:1000,offset:0) {",
			"\tname",
			"\tresourceid",
			"\tresourceversion",
			"\tcreationtime",
			"\tk8sobj",
			"\tobjtype",
			"\tuid",
			"\t~cluster @filter(eq(objtype, namespace) and eq(name,\"opa\") and eq(k8sobj,\"k8sobj\") )(first:1000,offset:0){",
			"\tk8sobj",
			"\tobjtype",
			"\tresourceid",
			"\tlabels",
			"\tresourceversion",
			"\tcreationtime",
			"\tname",
			"\tuid",
			"}",
			"}",
			"}",
		}, nil},
		`cluster[@name="paas-preprod-west2.cluster.k8s.local"]{*}.namespace[@name="default"]{**}`: FResult{[]string{
			"{ A as var(func: eq(objtype, cluster)) @filter( eq(name,\"paas-preprod-west2.cluster.k8s.local\") ) @cascade {",
			"\tcount(uid)",
			"\t~cluster @filter(eq(objtype, namespace) and eq(name,\"default\") ){",
			"\tcount(uid)",
			"}",
			"}",
			"}",
			"{ objects(func: uid(A),first:1000,offset:0) {",
			"\tresourceversion",
			"\tcreationtime",
			"\tk8sobj",
			"\tobjtype",
			"\tname",
			"\tresourceid",
			"\tuid",
			"\t~cluster @filter(eq(objtype, namespace) and eq(name,\"default\") )(first:1000,offset:0){",
			"\texpand(_all_){",
			"\texpand(_all_){",
			"}",
			"}",
			"}",
			"}",
			"}",
		}, nil},
		`cluster[@name="paas-preprod-west2.cluster.k8s.local"]{**}`: FResult{[]string{
			"{ A as var(func: eq(objtype, cluster)) @filter( eq(name,\"paas-preprod-west2.cluster.k8s.local\") ) @cascade {",
			"\tcount(uid)",
			"}",
			"}",
			"{ objects(func: uid(A),first:1000,offset:0) {",
			"\texpand(_all_){",
			"\texpand(_all_){",
			"}",
			"}",
			"}",
			"}",
		}, nil},
		`cluster[]{*}`: FResult{[]string{
			"{ A as var(func: eq(objtype, cluster))  @cascade {",
			"\tcount(uid)",
			"}",
			"}",
			"{ objects(func: uid(A),first:1000,offset:0) {",
			"\tresourceversion",
			"\tcreationtime",
			"\tk8sobj",
			"\tobjtype",
			"\tname",
			"\tresourceid",
			"\tuid",
			"}",
			"}",
		}, nil},
		`cluster[ ]{ }`: FResult{[]string{
			"{ A as var(func: eq(objtype, cluster))  @cascade {",
			"\tcount(uid)",
			"}",
			"}",
			"{ objects(func: uid(A),first:1000,offset:0) {",
			"}",
			"}",
		}, nil},
		`cluster[@name =  "paas-preprod-west2.cluster.k8s.local"]{ * }`: FResult{[]string{
			"{ A as var(func: eq(objtype, cluster)) @filter( eq(name,\"paas-preprod-west2.cluster.k8s.local\") ) @cascade {",
			"\tcount(uid)",
			"}",
			"}",
			"{ objects(func: uid(A),first:1000,offset:0) {",
			"\tobjtype",
			"\tname",
			"\tresourceid",
			"\tresourceversion",
			"\tcreationtime",
			"\tk8sobj",
			"\tuid",
			"}",
			"}",
		}, nil},
		`cluster[]{@name}`: FResult{[]string{
			"{ A as var(func: eq(objtype, cluster))  @cascade {",
			"\tcount(uid)",
			"}",
			"}",
			"{ objects(func: uid(A),first:1000,offset:0) {",
			"\tname",
			"\tuid",
			"}",
			"}",
		}, nil},
		`cluster[]{}.namespace[@name="default"]{*}`: FResult{[]string{
			"{ A as var(func: eq(objtype, cluster))  @cascade {",
			"\tcount(uid)",
			"\t~cluster @filter(eq(objtype, namespace) and eq(name,\"default\") ){",
			"\tcount(uid)",
			"}",
			"}",
			"}",
			"{ objects(func: uid(A),first:1000,offset:0) {",
			"\t~cluster @filter(eq(objtype, namespace) and eq(name,\"default\") )(first:1000,offset:0){",
			"\tk8sobj",
			"\tobjtype",
			"\tname",
			"\tresourceid",
			"\tlabels",
			"\tresourceversion",
			"\tcreationtime",
			"\tuid",
			"}",
			"}",
			"}",
		}, nil},
		`cluster[$$first=2]{}.namespace[@name="default"]{*}`: FResult{[]string{
			"{ A as var(func: eq(objtype, cluster))  @cascade {",
			"\tcount(uid)",
			"\t~cluster @filter(eq(objtype, namespace) and eq(name,\"default\") ){",
			"\tcount(uid)",
			"}",
			"}",
			"}",
			"{ objects(func: uid(A),first: 2) {",
			"\t~cluster @filter(eq(objtype, namespace) and eq(name,\"default\") )(first:1000,offset:0){",
			"\tlabels",
			"\tresourceversion",
			"\tcreationtime",
			"\tk8sobj",
			"\tresourceid",
			"\tname",
			"\tobjtype",
			"\tuid",
			"}",
			"}",
			"}",
		}, nil},
		`cluster[$$first=2,offset=2]{}.namespace[@name="default"$$first=2]{*}`: FResult{[]string{
			"{ A as var(func: eq(objtype, cluster))  @cascade {",
			"\tcount(uid)",
			"\t~cluster @filter(eq(objtype, namespace) and eq(name,\"default\") ){",
			"\tcount(uid)",
			"}",
			"}",
			"}",
			"{ objects(func: uid(A),first: 2,offset: 2) {",
			"\t~cluster @filter(eq(objtype, namespace) and eq(name,\"default\") )(first: 2){",
			"\tlabels",
			"\tresourceversion",
			"\tcreationtime",
			"\tk8sobj",
			"\tresourceid",
			"\tname",
			"\tobjtype",
			"\tuid",
			"}",
			"}",
			"}",
		}, nil},
		`cluster[@objtype="cluster"$$first=2,offset=2]{}.namespace[@name="default"$$first=2,offset=2]{*}`: FResult{[]string{
			"{ A as var(func: eq(objtype, cluster)) @filter( eq(objtype,\"cluster\") ) @cascade {",
			"\tcount(uid)",
			"\t~cluster @filter(eq(objtype, namespace) and eq(name,\"default\") ){",
			"\tcount(uid)",
			"}",
			"}",
			"}",
			"{ objects(func: uid(A),first: 2,offset: 2) {",
			"\t~cluster @filter(eq(objtype, namespace) and eq(name,\"default\") )(first: 2,offset: 2){",
			"\tlabels",
			"\tresourceversion",
			"\tcreationtime",
			"\tk8sobj",
			"\tresourceid",
			"\tname",
			"\tobjtype",
			"\tuid",
			"}",
			"}",
			"}",
		}, nil},
	}

	dgraphHost := "127.0.0.1:9080"
	dc := db.NewDGClient(dgraphHost)
	defer dc.Close()
	metaSvc := NewMetaService(dc)
	qslSvc := NewQSLService(dc, metaSvc)

	// Initialize metadata
	meta, err := ioutil.ReadFile("../data/meta.json")
	if err != nil {
		log.Fatalf("Metadata file error: %v\n", err)
	}
	var jsonData []map[string]interface{}
	json.Unmarshal(meta, &jsonData)
	for _, data := range jsonData {
		metaSvc.CreateMetadata(data)
	}

	for k, v := range tests {
		output, err := qslSvc.CreateDgraphQuery(k, false)
		if err != nil {
			if v.err != nil {
				if err.Error() != v.err.Error() {
					t.Errorf("query error incorrect\n input: %s\n testqueryerr: %s\n realqueryerr: %s", k, v.err.Error(), err.Error())
				}
			} else {
				t.Errorf("query error incorrect\n input: %s\n testqueryerr: %s\n realqueryerr: %s", k, "no error expected", err.Error())
			}
		} else {
			// check to see that the output has the same lines as the test output
			// because we can't assure that the metadata api will return fields
			// in the same order every time
			testmap := make(map[string]bool)

			for _, line := range v.values {
				testmap[strings.TrimSpace(line)] = true
			}

			for _, line := range strings.Split(output, "\n") {
				if !(testmap[strings.TrimSpace(line)]) {
					t.Errorf("query incorrect\n input: %s\n testquery: \n%s\n realquery: \n%s", k, strings.Join(v.values, "\n"), output)
					break
				}
			}

		}
	}

}
