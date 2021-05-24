package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/chernyshev-alex/bookstore-oauth-go/oauth"
	"github.com/chernyshev-alex/bookstore_items-api/domain/items"
	"github.com/chernyshev-alex/bookstore_items-api/mocks"
	"github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// type OAuthIntefaceMock interface {
// 	AuthenticateRequest(*http.Request) *rest_errors.RestErr
// 	GetCallerId(*http.Request) int64
// }

// type ItemsServiceInterfaceMock interface {
// 	Create(items.Item) (*items.Item, rest_errors.RestErr)
// 	Get(string) (*items.Item, rest_errors.RestErr)
// 	Search(queries.EsQuery) ([]items.Item, rest_errors.RestErr)
// }

// https://github.com/kyleconroy/sqlc

// type oauthServiceMock struct {
// 	isPublicFn            func(*http.Request) bool
// 	getCallerIdFn         func(*http.Request) int64
// 	getClientIdFn         func(*http.Request) int64
// 	authenticateRequestFn func(*http.Request) rest_errors.RestErr
// }

var (
	mockedItemsService = new(mocks.ItemsServiceInterface)
	itemsController    = NewItemController(
		oauth.ProvideOAuthClient(http.DefaultClient),
		mockedItemsService,
	)
)

func TestCreateOk(t *testing.T) {
	var (
		callerId int64 = 100
		item           = items.Item{Id: "", Seller: callerId}
	)

	bytes, _ := json.Marshal(item)
	req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(string(bytes)))

	req.Header.Add("X-Caller-Id", strconv.FormatInt(callerId, 10))

	resp := httptest.NewRecorder()

	mockedItemsService.On("Create", mock.IsType(item)).Return(func(item items.Item) *items.Item {
		item.Id = "assgined"
		return &item
	}, nil)

	itemsController.Create(resp, req)

	var itemResult items.Item
	if err := json.Unmarshal(resp.Body.Bytes(), &itemResult); err != nil {
		t.Error("bad body response", resp.Body)
	}

	mockedItemsService.AssertExpectations(t)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Equal(t, "assgined", itemResult.Id)
	assert.Equal(t, callerId, itemResult.Seller)
}

type ARestErr struct {
	Fmessage string        `json:"message"`
	Fstatus  int           `json:"status"`
	Ferror   string        `json:"error"`
	Fcauses  []interface{} `json:"causes"`
}

func TestCreateSaveFailed(t *testing.T) {
	var (
		callerId int64 = 100
		item           = items.Item{Id: "", Seller: callerId}
	)

	bytes, _ := json.Marshal(item)
	req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(string(bytes)))

	req.Header.Add("X-Caller-Id", strconv.FormatInt(callerId, 10))

	resp := httptest.NewRecorder()

	mockedItemsService.On("Create", mock.IsType(item)).Return(nil,
		func(item items.Item) rest_errors.RestErr {
			return rest_errors.NewInternalServerError("item service", errors.New("failed"))
		})

	itemsController.Create(resp, req)

	restError, err := rest_errors.NewRestErrorFromBytes(resp.Body.Bytes())
	if err != nil {
		t.Error("bad body response", resp.Body)
	}

	mockedItemsService.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, restError.Status())
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}
