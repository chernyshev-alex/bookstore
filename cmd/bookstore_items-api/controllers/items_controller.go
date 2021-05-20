package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/chernyshev-alex/bookstore-oauth-go/oauth"
	"github.com/chernyshev-alex/bookstore_items-api/domain/items"
	"github.com/chernyshev-alex/bookstore_items-api/domain/queries"
	"github.com/chernyshev-alex/bookstore_items-api/services"
	"github.com/chernyshev-alex/bookstore_items-api/utils/http_utils"
	"github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
	"github.com/gorilla/mux"
)

type ItemControllerInterface interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	Ping(http.ResponseWriter, *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
}

type itemController struct{}

var (
	ItemController ItemControllerInterface = &itemController{}
)

func (c *itemController) Ping(w http.ResponseWriter, rq *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func (c *itemController) Create(w http.ResponseWriter, rq *http.Request) {

	if err := oauth.AuthenticateRequest(rq); err != nil {
		http_utils.ResponseError(w, err)
		return
	}

	callerId := oauth.GetCallerId(rq)
	if callerId == 0 {
		http_utils.ResponseError(w, rest_errors.NewAuthorizationError("no user info in the token"))
		return
	}

	buf, err := ioutil.ReadAll(rq.Body)
	if err != nil {
		http_utils.ResponseError(w, rest_errors.NewBadRequestError(err.Error()))
	}

	defer rq.Body.Close()

	var item items.Item
	err = json.Unmarshal(buf, &item)
	if err != nil {
		http_utils.ResponseError(w, rest_errors.NewBadRequestError(err.Error()))
	}

	item.Seller = callerId
	result, createErr := services.ItemService.Create(item)
	if createErr != nil {
		http_utils.ResponseError(w, createErr)
	}

	http_utils.ResponseJson(w, http.StatusCreated, result)
}

func (c *itemController) Get(w http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	itemId := strings.TrimSpace(vars["id"])
	item, err := services.ItemService.Get(itemId)
	if err != nil {
		http_utils.ResponseError(w, err)
	}
	http_utils.ResponseJson(w, http.StatusOK, item)
}

func (c *itemController) Search(w http.ResponseWriter, rq *http.Request) {
	b, err := ioutil.ReadAll(rq.Body)
	if err != nil {
		http_utils.ResponseError(w, rest_errors.NewBadRequestError("bad json body"))
		return
	}
	defer rq.Body.Close()

	var q queries.EsQuery
	if err := json.Unmarshal(b, &q); err != nil {
		apiErr := rest_errors.NewBadRequestError("bad json body")
		http_utils.ResponseError(w, apiErr)
		return
	}

	items, searchErr := services.ItemService.Search(q)
	if searchErr != nil {
		http_utils.ResponseError(w, searchErr)
		return
	}

	http_utils.ResponseJson(w, http.StatusOK, items)
}
