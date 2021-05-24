package services

import (
	"github.com/chernyshev-alex/bookstore_items-api/domain/items"
	"github.com/chernyshev-alex/bookstore_items-api/domain/queries"
	"github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
)

var (
	ItemService ItemsServiceInterface = &itemsService{}
)

type ItemsServiceInterface interface {
	Create(items.Item) (*items.Item, rest_errors.RestErr)
	Get(string) (*items.Item, rest_errors.RestErr)
	Search(queries.EsQuery) ([]items.Item, rest_errors.RestErr)
}

type itemsService struct{}

func NewItemsService() ItemsServiceInterface {
	return &itemsService{}
}

func (s *itemsService) Create(item items.Item) (*items.Item, rest_errors.RestErr) {
	if err := item.Save(); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *itemsService) Get(id string) (*items.Item, rest_errors.RestErr) {
	item := items.Item{Id: id}
	if err := item.Get(); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *itemsService) Search(q queries.EsQuery) ([]items.Item, rest_errors.RestErr) {
	return new(items.Item).Search(q)
}
