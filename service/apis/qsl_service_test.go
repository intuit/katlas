package apis

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/intuit/katlas/service/db"
)

type FResult struct {
	values []string
	err    error
}

func TestCreateFiltersQuery(t *testing.T) {

	tests := map[string]FResult{
		`@name="paas-preprod-west2.cluster.k8s.local"|@k8sobj="K8sObj"|@resourceid="paas-preprod-west2.cluster.k8s.local"`: FResult{
			[]string{
				", $name: string, $k8sobj: string, $resourceid: string",
				` ( eq(name,"paas-preprod-west2.cluster.k8s.local") or eq(k8sobj,"K8sObj") or eq(resourceid,"paas-preprod-west2.cluster.k8s.local") )`,
			},
			nil,
		},
		`@name="paas-preprod-west2.cluster.k8s.local",@k8sobj="K8sObj",@resourceid="paas-preprod-west2.cluster.k8s.local"`: FResult{
			[]string{
				", $name: string, $k8sobj: string, $resourceid: string",
				` ( eq(name,"paas-preprod-west2.cluster.k8s.local") and eq(k8sobj,"K8sObj") and eq(resourceid,"paas-preprod-west2.cluster.k8s.local") )`,
			},
			nil,
		},
		`@name="paas-preprod-west2.cluster.k8s.local"|@k8sobj="K8sObj",@resourceid="paas-preprod-west2.cluster.k8s.local"`: FResult{
			[]string{
				", $name: string, $k8sobj: string, $resourceid: string",
				` ( eq(name,"paas-preprod-west2.cluster.k8s.local") or eq(k8sobj,"K8sObj") and eq(resourceid,"paas-preprod-west2.cluster.k8s.local") )`,
			},
			nil,
		},
		`@name="paas-preprod-west2.cluster.k8s.local"|@k8sobj="K8sObj""`: FResult{
			[]string{
				", $name: string, $k8sobj: string",
				` ( eq(name,"paas-preprod-west2.cluster.k8s.local") or eq(k8sobj,"K8sObj") )`,
			},
			nil,
		},
		`@k8sobj="K8sObj",@resourceid="paas-preprod-west2.cluster.k8s.local"`: FResult{
			[]string{
				", $k8sobj: string, $resourceid: string",
				` ( eq(k8sobj,"K8sObj") and eq(resourceid,"paas-preprod-west2.cluster.k8s.local") )`,
			},
			nil,
		},
		`@name="paas-preprod-west2.cluster.k8s.local""`: FResult{
			[]string{
				", $name: string",
				` ( eq(name,"paas-preprod-west2.cluster.k8s.local") )`,
			},
			nil,
		},
		`@name="paas-preprod-west2.cluster.k8s.local?"`: FResult{
			[]string{
				", $name: string",
				` ( eq(name,"paas-preprod-west2.cluster.k8s.local") )`,
			},
			errors.New("Invalid filters in @name=\"paas-preprod-west2.cluster.k8s.local?\""),
		},
	}

	for k, v := range tests {
		realOut1, realOut2, err := CreateFiltersQuery(k)
		if err == nil {
			if !(v.values[0] == realOut1) {
				t.Errorf("filter declaration incorrect\n input: %s\n testdec: %s\n realdec: %s", k, v.values[0], realOut1)
			}

			if !(v.values[1] == realOut2) {
				t.Errorf("filter function incorrect\n input: %s\n testfunc: %s\n realfunc: %s", k, v.values[1], realOut2)
			}
		} else {
			if err.Error() != v.err.Error() {
				t.Errorf("filter error incorrect\n input: %s\n testfiltererr: %s\n realfiltererr: %s", k, v.err.Error(), err.Error())
			}
		}

	}

}

func TestCreateFieldsQuery(t *testing.T) {
	tests := map[string]FResult{
		"@name":                              FResult{[]string{"name", "uid"}, nil},
		"@name,@resourceversion":             FResult{[]string{"name", "resourceversion", "uid"}, nil},
		"@name,@resourceversion,@resourceid": FResult{[]string{"name", "resourceversion", "resourceid", "uid"}, nil},
		"@name,@resourceversion,@k8sobj":     FResult{[]string{"name", "resourceversion", "k8sobj", "uid"}, nil},
		"*":     FResult{[]string{"k8sobj", "objtype", "name", "resourceid", "resourceversion", "uid"}, nil},
		"**":    FResult{[]string{"\texpand(_all_){", "\t\texpand(_all_){", "\t\t}", "\t}"}, nil},
		"***":   FResult{[]string{"\texpand(_all_){", "\t\texpand(_all_){", "\t\t\texpand(_all_){", "\t\t\t}", "\t\t}", "\t}"}, nil},
		"?":     FResult{nil, errors.New("Field names must be prefixed with @ sign and followed by an alphanumeric field name [?]")},
		"@n@me": FResult{nil, errors.New("Field names must be composed of only alphanumeric characters [n@me]")},
		"@*":    FResult{nil, errors.New("Field names must be composed of only alphanumeric characters [*]")},
		"*@":    FResult{nil, errors.New("Fields may be a string of * indicating how many levels, or a list of fields @field1,@field2,... [*@]")},
	}

	namespacemetafieldslist := []MetadataField{MetadataField{FieldName: "k8sobj", FieldType: "string", Mandatory: true, Index: true, RefDataType: "", Cardinality: "One"}, MetadataField{FieldName: "objtype", FieldType: "string", Mandatory: true, Index: true, RefDataType: "", Cardinality: "One"}, MetadataField{FieldName: "name", FieldType: "string", Mandatory: true, Index: true, RefDataType: "", Cardinality: "One"}, MetadataField{FieldName: "resourceid", FieldType: "string", Mandatory: false, Index: true, RefDataType: "", Cardinality: "One"}, MetadataField{FieldName: "resourceversion", FieldType: "string", Mandatory: true, Index: false, RefDataType: "", Cardinality: "One"}}

	for k, v := range tests {
		output, err := CreateFieldsQuery(k, namespacemetafieldslist, -1)
		if err == nil {
			if !(reflect.DeepEqual(v.values, output)) {
				t.Errorf("fields incorrect\n input: %s\n testfields: %s\n realfields: %s", k, v.values, output)
			}
		} else {
			if v.err != nil {
				if err.Error() != v.err.Error() {
					t.Errorf("filter error incorrect\n input: %s\n testfielderr: %s\n realfielderr: %s", k, v.err.Error(), err.Error())
				}
			} else {
				t.Errorf("filter error incorrect\n input: %s\n testfielderr: %s\n realfielderr: %s", k, "no error expected", err.Error())
			}

		}

	}
}

func TestCreateDgraphQuery(t *testing.T) {
	tests := map[string]string{
		`namespace[@name="default"]{*}`: strings.Join([]string{
			"query objects($objtype: string, $name: string){",
			"objects(func: eq(objtype, Namespace)) @filter( ( eq(name,\"default\") )){",
			"\tresourceversion",
			"\tcreationtime",
			"\tk8sobj",
			"\tlabels",
			"\tname",
			"\tresourceid",
			"\tobjtype",
			"\tuid",
			"}",
			"}",
		}, "\n"),
		`cluster[@name="paas-preprod-west2.cluster.k8s.local"]{@name}`: strings.Join([]string{
			"query objects($objtype: string, $name: string){",
			"objects(func: eq(objtype, Cluster)) @filter( ( eq(name,\"paas-preprod-west2.cluster.k8s.local\") )){",
			"\tname",
			"\tuid",
			"}",
			"}",
		}, "\n"),

		`cluster[@name="paas-preprod-west2.cluster.k8s.local"]{*}.namespace[@name="opa"|@name="default"]{*}`: strings.Join([]string{
			"query objects($objtype: string, $name: string){",
			"objects(func: eq(objtype, Cluster)) @filter( ( eq(name,\"paas-preprod-west2.cluster.k8s.local\") )){",
			"\tcreationtime",
			"\tk8sobj",
			"\tobjtype",
			"\tname",
			"\tresourceid",
			"\tresourceversion",
			"\tuid",
			"\t~cluster @filter(eq(objtype, Namespace) and ( eq(name,\"opa\") or eq(name,\"default\") )){",
			"\t\tresourceversion",
			"\t\tcreationtime",
			"\t\tk8sobj",
			"\t\tlabels",
			"\t\tname",
			"\t\tresourceid",
			"\t\tobjtype",
			"\t\tuid",
			"\t}",
			"}",
			"}",
		}, "\n"),
		`cluster[@name="paas-preprod-west2.cluster.k8s.local"]{**}`: strings.Join([]string{
			"query objects($objtype: string, $name: string){",
			"objects(func: eq(objtype, Cluster)) @filter( ( eq(name,\"paas-preprod-west2.cluster.k8s.local\") )){",
			"\texpand(_all_){",
			"\t\texpand(_all_){",
			"\t\t}",
			"\t}",
			"}",
			"}",
		}, "\n"),
	}

	dgraphHost := "127.0.0.1:9080"
	dc := db.NewDGClient(dgraphHost)
	defer dc.Close()
	metaSvc := NewMetaService(dc)
	qslSvc := NewQSLService(dgraphHost, metaSvc)

	for k, v := range tests {
		output, err := qslSvc.CreateDgraphQuery(k)
		if err != nil {
			t.Error(err)
		}
		if !(output == v) {
			t.Errorf("query incorrect\n input: %s\n testquery: \n%s\n realquery: \n%s", k, v, output)
		}
	}

}
