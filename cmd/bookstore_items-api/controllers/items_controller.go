package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_items_api/domain/items"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_items_api/domain/queries"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_items_api/services"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore-oauth-go/oauth"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/rest_errors"
	"github.com/gorilla/mux"
)

type ItemControllerInterface interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	Ping(http.ResponseWriter, *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
}

type itemController struct {
	oauthService oauth.OAuthInterface
	itemsService services.ItemsServiceInterface
}

func NewItemController(oauthService oauth.OAuthInterface,
	itemService services.ItemsServiceInterface) ItemControllerInterface {
	return &itemController{
		oauthService: oauthService,
		itemsService: itemService,
	}
}

func (c *itemController) Ping(w http.ResponseWriter, rq *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func (c *itemController) Create(w http.ResponseWriter, rq *http.Request) {
	if err := c.oauthService.AuthenticateRequest(rq); err != nil {
		rest_errors.ResponseError(w, err)
		return
	}

	callerId := c.oauthService.GetCallerId(rq)
	if callerId == 0 {
		rest_errors.ResponseError(w, rest_errors.NewAuthorizationError("no user info in the token"))
		return
	}

	buf, err := ioutil.ReadAll(rq.Body)
	if err != nil {
		rest_errors.ResponseError(w, rest_errors.NewBadRequestError(err.Error()))
	}

	defer rq.Body.Close()

	var item items.Item
	err = json.Unmarshal(buf, &item)
	if err != nil {
		rest_errors.ResponseError(w, rest_errors.NewBadRequestError(err.Error()))
	}

	item.Seller = callerId
	result, createErr := c.itemsService.Create(item)
	if createErr != nil {
		rest_errors.ResponseError(w, createErr)
		return
	}
	rest_errors.ResponseJson(w, http.StatusCreated, result)
}

func (c *itemController) Get(w http.ResponseWriter, rq *http.Request) {
	itemId := rq.Context().Value("id")
	if itemId == nil {
		itemId = strings.TrimSpace(mux.Vars(rq)["id"])
	}

	item, err := c.itemsService.Get(itemId.(string))
	if err != nil {
		rest_errors.ResponseError(w, err)
		return
	}
	rest_errors.ResponseJson(w, http.StatusOK, item)
}

func (c *itemController) Search(w http.ResponseWriter, rq *http.Request) {
	b, err := ioutil.ReadAll(rq.Body)
	if err != nil {
		rest_errors.ResponseError(w, rest_errors.NewBadRequestError("bad json body"))
		return
	}
	defer rq.Body.Close()

	var q queries.EsQuery
	if err := json.Unmarshal(b, &q); err != nil {
		apiErr := rest_errors.NewBadRequestError("bad query")
		rest_errors.ResponseError(w, apiErr)
		return
	}

	items, searchErr := c.itemsService.Search(q)
	if searchErr != nil {
		rest_errors.ResponseError(w, searchErr)
		return
	}

	rest_errors.ResponseJson(w, http.StatusOK, items)
}
