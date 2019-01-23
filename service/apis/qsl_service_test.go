package apis

import (
		"errors"
			"reflect"
		"testing"

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
