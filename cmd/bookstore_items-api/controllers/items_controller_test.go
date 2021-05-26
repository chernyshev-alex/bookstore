package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/chernyshev-alex/bookstore_items-api/domain/items"
	"github.com/chernyshev-alex/bookstore_items-api/mocks"
	"github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockOAuth struct {
	mock.Mock
}

func (m *MockOAuth) AuthenticateRequest(rq *http.Request) rest_errors.RestErr {
	ret := m.Called(rq)
	var r1 rest_errors.RestErr
	if rf, ok := ret.Get(0).(func(*http.Request) rest_errors.RestErr); ok {
		r1 = rf(rq)
	} else {
		if ret.Get(0) != nil {
			r1 = ret.Get(0).(rest_errors.RestErr)
		}
	}
	return r1
}

func (m *MockOAuth) IsPublic(rq *http.Request) bool {
	args := m.Called(rq)
	return args.Bool(0)
}

func (m *MockOAuth) GetCallerId(rq *http.Request) int64 {
	args := m.Called(rq)
	return int64(args.Int(0))
}
func (m *MockOAuth) GetClientId(rq *http.Request) int64 {
	args := m.Called(rq)
	return int64(args.Int(0))
}

type ItemControllerSuite struct {
	suite.Suite
	mockedItemsService *mocks.ItemsServiceInterface
	mockedOAuthService *MockOAuth
	itemsController    ItemControllerInterface
}

func TestItemControllerSuite(t *testing.T) {
	suite.Run(t, new(ItemControllerSuite))
}

func (s *ItemControllerSuite) SetupTest() {
	s.mockedOAuthService = new(MockOAuth)
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

	req := requestForBody(http.MethodPost, "/items", &item)

	req.Header.Add("X-Caller-Id", strconv.FormatInt(callerId, 10))

	resp := httptest.NewRecorder()

	s.mockedOAuthService.On("AuthenticateRequest", mock.IsType(req)).Return(nil)

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
	req := requestForBody(http.MethodPost, "/items", &items.Item{})
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

func (s *ItemControllerSuite) TestCreateFailedOnSave() {
	var (
		callerId int64 = 100
		item           = items.Item{Id: "", Seller: callerId}
	)

	req := requestForBody(http.MethodPost, "/items", &item)

	req.Header.Add("X-Caller-Id", strconv.FormatInt(callerId, 10))

	resp := httptest.NewRecorder()

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

// helpers

func requestForBody(httpMethod, urlPath string, item *items.Item) *http.Request {
	bytes, _ := json.Marshal(item)
	return httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(string(bytes)))
}
