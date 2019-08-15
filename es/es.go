package es

import (
	"github.com/olivere/elastic/v7"
)

func NewClient() *elastic.Client {
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}
	return client
}
