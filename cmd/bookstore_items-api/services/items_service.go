package services

import (
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_items-api/domain/items"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_items-api/domain/queries"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/rest_errors"
)

type ItemsServiceInterface interface {
	Create(items.Item) (*items.Item, rest_errors.RestErr)
	Get(string) (*items.Item, rest_errors.RestErr)
	Search(queries.EsQuery) ([]items.Item, rest_errors.RestErr)
}

type itemsService struct {
	persist items.ItemsPersistInterface
}

func NewItemsService(persist items.ItemsPersistInterface) ItemsServiceInterface {
	return &itemsService{
		persist: persist,
	}
}

func (s *itemsService) Create(items items.Item) (*items.Item, rest_errors.RestErr) {
	if err := s.persist.Save(&items); err != nil {
		return nil, err
	}
	return &items, nil
}

func (s *itemsService) Get(Id string) (*items.Item, rest_errors.RestErr) {
	item := items.Item{Id: Id}
	res, err := s.persist.Get(item)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *itemsService) Search(q queries.EsQuery) ([]items.Item, rest_errors.RestErr) {
	return s.persist.Search(q)
}
