package es

import (
	"context"
	"time"

	"github.com/chernyshev-alex/bookstore_utils_go/logger"
	"github.com/olivere/elastic"
)

var (
	Client esClientInterface = &esClient{}
)

type esClientInterface interface {
	SetClient(*elastic.Client)
	Index(string, string, interface{}) (*elastic.IndexResponse, error)
	Get(string, string, string) (*elastic.GetResult, error)
	Search(string, elastic.Query) (*elastic.SearchResult, error)
}

type esClient struct {
	client *elastic.Client
}

func Init() {
	c, err := elastic.NewClient(
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetErrorLog(logger.GetLogger()),
		elastic.SetInfoLog(logger.GetLogger()))

	if err != nil {
		panic(err)
	}

	Client.SetClient(c)
}

func (es *esClient) SetClient(c *elastic.Client) {
	es.client = c
}

func (es *esClient) Index(index string, docType string, doc interface{}) (*elastic.IndexResponse, error) {
	result, err := es.client.Index().
		Index(index).
		Type(docType).
		BodyJson(doc).
		Do(context.Background())
	if err != nil {
		logger.Error("idexing document error", err)
		return nil, err
	}
	return result, nil
}

func (es esClient) Get(index string, docType string, id string) (*elastic.GetResult, error) {
	ctx := context.Background()
	result, err := es.client.Get().
		Index(index).
		Type(docType).
		Id(id).
		Do(ctx)
	if err != nil {
		logger.Error("get document error", err)
		return nil, err
	}
	return result, nil
}

func (es esClient) Search(index string, q elastic.Query) (*elastic.SearchResult, error) {
	ctx := context.Background()
	result, err := es.client.Search(index).
		Query(q).
		RestTotalHitsAsInt(true).
		Do(ctx)
	if err != nil {
		logger.Error("search documents error", err)
		return nil, err
	}
	return result, nil
}
