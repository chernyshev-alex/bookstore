package queries

import "github.com/olivere/elastic"

type FieldValue struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

type EsQuery struct {
	Equals []FieldValue `json:"equals"`
}

func (q EsQuery) Build() elastic.Query {
	query := elastic.NewBoolQuery()
	queries := make([]elastic.Query, 0)
	for _, q := range q.Equals {
		queries = append(queries, elastic.NewMatchQuery(q.Field, q.Value))
	}
	return query.Must(queries...)
}
