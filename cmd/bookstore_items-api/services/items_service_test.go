package services

import (
	"errors"
	"testing"

	"github.com/chernyshev-alex/bookstore_items-api/domain/items"
	"github.com/chernyshev-alex/bookstore_items-api/domain/queries"
	"github.com/chernyshev-alex/bookstore_items-api/mocks"
	"github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ItemServiceSuite struct {
	suite.Suite
	itemsService ItemsServiceInterface
	daoItemsMock mocks.ItemsPersistInterface
}

func TestItemServiceSuite(t *testing.T) {
	suite.Run(t, new(ItemServiceSuite))
}

func (s *ItemServiceSuite) SetupTest() {
	s.daoItemsMock = mocks.ItemsPersistInterface{}
	s.itemsService = NewItemsService(&s.daoItemsMock)
}

func (s *ItemServiceSuite) TestCreateOk() {
	var item = items.Item{}
	s.daoItemsMock.On("Save", mock.IsType(&item)).Return(func(item *items.Item) rest_errors.RestErr {
		item.Id = "assigned"
		return nil
	})

	result, err := s.itemsService.Create(item)

	s.daoItemsMock.AssertExpectations(s.T())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "assigned", result.Id)
	assert.True(s.T(), len(result.Id) > 0)
}

func (s *ItemServiceSuite) TestCreateFailedPersist() {
	var item = items.Item{}
	s.daoItemsMock.On("Save", mock.IsType(&item)).Return(func(item *items.Item) rest_errors.RestErr {
		return rest_errors.NewInternalServerError("failed", errors.New("save"))
	})

	result, err := s.itemsService.Create(item)

	s.daoItemsMock.AssertExpectations(s.T())
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *ItemServiceSuite) TestGetOk() {
	const objId = "111"
	it := items.Item{Id: objId}
	s.daoItemsMock.On("Get", mock.IsType(it)).Return(&items.Item{Id: objId}, nil)

	result, err := s.itemsService.Get(objId)

	s.daoItemsMock.AssertExpectations(s.T())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), objId, result.Id)
}

func (s *ItemServiceSuite) TestGetNotFound() {
	const objId = "111"
	it := items.Item{Id: objId}
	s.daoItemsMock.On("Get", mock.IsType(it)).Return(nil, rest_errors.NewNotFoundError(objId))

	_, err := s.itemsService.Get(objId)

	s.daoItemsMock.AssertExpectations(s.T())
	assert.NotNil(s.T(), err)
}

func (s *ItemServiceSuite) TestSearchOk() {
	q := queries.EsQuery{}

	s.daoItemsMock.On("Search", mock.IsType(q)).Return([]items.Item{}, nil)

	result, err := s.itemsService.Search(q)

	s.daoItemsMock.AssertExpectations(s.T())
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)
}
