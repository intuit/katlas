package db

//QueryPredicates...Query params to be used for the Keyword query
//Each field must be a string, have an index with tokenizer term

type predicate struct {
	Type      string
	IsIndex   bool
	IndexType []string
}

//QueryPredicates ...Predicates to filter query for Keyword Query
var QueryPredicates = map[string]predicate{
	"name":              predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"objtype":           predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"k8sobj":            predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"resourceversion":   predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"resourceid":        predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"availablereplicas": predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"clusterip":         predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"containers":        predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"creationtime":      predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"defaultbackend": predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"ip": predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"labels": predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"numreplicas": predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"phase": predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"podspec": predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"ports": predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"rules": predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"selector": predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"servicetype": predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"strategy": predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"tls": predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
	"volumes": predicate{Type: "string", IsIndex: true, IndexType: []string{"term", "trigram"}},
}

// IsExpectedIndexType to check index type
func IsExpectedIndexType(indexTypes []string, indexType string) bool {

	for _, index := range indexTypes {
		if index == indexType {
			return true
		}
	}
	return false
}
