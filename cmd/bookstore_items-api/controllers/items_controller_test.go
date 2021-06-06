package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_items-api/mocks"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_items_api/domain/items"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_items_api/domain/queries"
	oaumocks "github.com/chernyshev-alex/bookstore/pkg/bookstore-oauth-go/mocks"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/rest_errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ItemControllerSuite struct {
	suite.Suite
	mockedItemsService *mocks.ItemsServiceInterface
	mockedOAuthService *oaumocks.OAuthInterface
	itemsController    ItemControllerInterface
}

func TestItemControllerSuite(t *testing.T) {
	suite.Run(t, new(ItemControllerSuite))
}

func (s *ItemControllerSuite) SetupTest() {
	s.mockedOAuthService = new(oaumocks.OAuthInterface)
	s.mockedItemsService = new(mocks.ItemsServiceInterface)
	s.itemsController = NewItemController(
		s.mockedOAuthService,
		s.mockedItemsService,
	)
}

func (s *ItemControllerSuite) TestCreateOk() {
	var (
		callerId int64 = 100
		item           = items.Item{Id: "", Seller: callerId}
	)

	req := requestForBodyItem(http.MethodPost, "/items", &item)

	req.Header.Add("X-Caller-Id", strconv.FormatInt(callerId, 10))

	resp := httptest.NewRecorder()

	s.mockedOAuthService.On("AuthenticateRequest", mock.IsType(req)).Return(nil)
	s.mockedOAuthService.On("GetCallerId", mock.IsType(req)).Return(int(callerId))

	s.mockedItemsService.On("Create", mock.IsType(item)).Return(func(item items.Item) *items.Item {
		item.Id = "assigned"
		return &item
	}, nil)

	s.itemsController.Create(resp, req)

	var itemResult items.Item
	err := json.Unmarshal(resp.Body.Bytes(), &itemResult)

	s.mockedItemsService.AssertExpectations(s.T())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusCreated, resp.Code)
	assert.Equal(s.T(), "assigned", itemResult.Id)
	assert.Equal(s.T(), callerId, itemResult.Seller)
}

func (s *ItemControllerSuite) TestCreateFailedNotAuthenticated() {
	req := requestForBodyItem(http.MethodPost, "/items", &items.Item{})
	resp := httptest.NewRecorder()

	s.mockedOAuthService.On("AuthenticateRequest", mock.IsType(req)).Return(
		rest_errors.NewAuthorizationError("oauth error"))

	s.itemsController.Create(resp, req)
	restError, err := rest_errors.NewRestErrorFromBytes(resp.Body.Bytes())

	s.mockedItemsService.AssertExpectations(s.T())
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusUnauthorized, restError.Status())
	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}

func (s *ItemControllerSuite) TestCreateFailedCallerIdIsZero() {
	req := requestForBodyItem(http.MethodPost, "/items", &items.Item{})
	resp := httptest.NewRecorder()

	s.mockedOAuthService.On("AuthenticateRequest", mock.IsType(req)).Return(nil)
	s.mockedOAuthService.On("GetCallerId", mock.IsType(req)).Return(0)

	s.itemsController.Create(resp, req)
	restError, err := rest_errors.NewRestErrorFromBytes(resp.Body.Bytes())

	s.mockedItemsService.AssertExpectations(s.T())
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusUnauthorized, restError.Status())
	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}

func (s *ItemControllerSuite) TestCreateFailedOnSave() {
	var (
		callerId int64 = 100
		item           = items.Item{Id: "", Seller: callerId}
	)

	req := requestForBodyItem(http.MethodPost, "/items", &item)

	req.Header.Add("X-Caller-Id", strconv.FormatInt(callerId, 10))

	resp := httptest.NewRecorder()

	s.mockedOAuthService.On("AuthenticateRequest", mock.IsType(req)).Return(nil)
	s.mockedOAuthService.On("GetCallerId", mock.IsType(req)).Return(1)

	s.mockedItemsService.On("Create", mock.IsType(item)).Return(nil,
		func(item items.Item) rest_errors.RestErr {
			return rest_errors.NewInternalServerError("item service", errors.New("failed"))
		})

	s.itemsController.Create(resp, req)

	restError, err := rest_errors.NewRestErrorFromBytes(resp.Body.Bytes())

	s.mockedItemsService.AssertExpectations(s.T())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusInternalServerError, restError.Status())
	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
}

func (s *ItemControllerSuite) TestCreateFailedBadRequest() {
	var (
		callerId int64 = 100
		item           = items.Item{Id: "", Seller: callerId}
	)

	req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader("invalid item json"))

	req.Header.Add("X-Caller-Id", strconv.FormatInt(1, 10))

	resp := httptest.NewRecorder()

	s.mockedOAuthService.On("AuthenticateRequest", mock.IsType(req)).Return(nil)
	s.mockedOAuthService.On("GetCallerId", mock.IsType(req)).Return(1)

	s.mockedItemsService.On("Create", mock.IsType(item)).Return(func(item items.Item) *items.Item {
		item.Id = "assigned"
		return &item
	}, nil)

	s.itemsController.Create(resp, req)

	var itemResult items.Item
	err := json.Unmarshal(resp.Body.Bytes(), &itemResult)

	s.mockedItemsService.AssertExpectations(s.T())

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), http.StatusBadRequest, resp.Code)
}

func (s *ItemControllerSuite) TestGetItemOk() {
	var (
		itemId = "100"
		item   = items.Item{Id: itemId}
	)

	ctx := context.WithValue(context.Background(), "id", itemId)
	req := httptest.NewRequest(http.MethodGet, "/items/{id}", nil).WithContext(ctx)

	resp := httptest.NewRecorder()

	s.mockedItemsService.On("Get", mock.IsType("")).Return(&item, nil)

	s.itemsController.Get(resp, req)

	var itemResult items.Item
	err := json.Unmarshal(resp.Body.Bytes(), &itemResult)

	s.mockedItemsService.AssertExpectations(s.T())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), itemResult.Id, itemId)
	assert.Equal(s.T(), http.StatusOK, resp.Code)
}

func (s *ItemControllerSuite) TestGetItemFailed() {
	req := httptest.NewRequest(http.MethodGet, "/items/{id}", nil)

	resp := httptest.NewRecorder()

	s.mockedItemsService.On("Get", mock.IsType("100")).Return(nil,
		func(itemId string) rest_errors.RestErr {
			return rest_errors.NewInternalServerError("item service", errors.New("failed"))
		})

	s.itemsController.Get(resp, req)

	restError, _ := rest_errors.NewRestErrorFromBytes(resp.Body.Bytes())

	s.mockedItemsService.AssertExpectations(s.T())

	assert.NotNil(s.T(), restError)
	assert.Equal(s.T(), http.StatusInternalServerError, restError.Status())
	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
}

func (s *ItemControllerSuite) TestSearchOk() {
	var q = queries.EsQuery{}

	req := requestForBodyQuery(http.MethodPost, "/items/search", &q)
	resp := httptest.NewRecorder()

	s.mockedItemsService.On("Search", mock.IsType(q)).Return(
		func(q queries.EsQuery) []items.Item { return []items.Item{{Id: "1"}} }, nil)

	s.itemsController.Search(resp, req)

	var items []items.Item
	err := json.Unmarshal(resp.Body.Bytes(), &items)

	s.mockedItemsService.AssertExpectations(s.T())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(items))
	assert.Equal(s.T(), http.StatusOK, resp.Code)
}

func (s *ItemControllerSuite) TestSearchFailedService() {
	var q = queries.EsQuery{}

	req := requestForBodyQuery(http.MethodPost, "/items/search", &q)
	resp := httptest.NewRecorder()

	s.mockedItemsService.On("Search", mock.IsType(q)).Return(nil,
		func(q queries.EsQuery) rest_errors.RestErr {
			return rest_errors.NewInternalServerError("search service", errors.New("failed"))
		})

	s.itemsController.Search(resp, req)

	var items []items.Item
	err := json.Unmarshal(resp.Body.Bytes(), &items)

	s.mockedItemsService.AssertExpectations(s.T())

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
}

func (s *ItemControllerSuite) TestSearchBadRequest() {
	req := httptest.NewRequest(http.MethodPost, "/items/search", strings.NewReader("bad query"))
	resp := httptest.NewRecorder()

	s.itemsController.Search(resp, req)

	re, err := rest_errors.NewRestErrorFromBytes(resp.Body.Bytes())

	s.mockedItemsService.AssertExpectations(s.T())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusBadRequest, re.Status())
}

// helpers

func requestForBodyItem(httpMethod, urlPath string, item *items.Item) *http.Request {
	bytes, _ := json.Marshal(item)
	return httptest.NewRequest(httpMethod, urlPath, strings.NewReader(string(bytes)))
}

func requestForBodyQuery(httpMethod, urlPath string, q *queries.EsQuery) *http.Request {
	bytes, _ := json.Marshal(q)
	return httptest.NewRequest(httpMethod, urlPath, strings.NewReader(string(bytes)))
}
