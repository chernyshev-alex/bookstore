package http

import (
	"net/http"
	"strings"

	"github.com/chernyshev-alex/bookstore_oauth-api/domain/access_token"
	"github.com/chernyshev-alex/bookstore_oauth-api/services"
	"github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
	"github.com/gin-gonic/gin"
)

type AccessTokenHandler interface {
	GetById(*gin.Context)
	Create(ctx *gin.Context)
}

type accessTokenHandler struct {
	service services.Service
}

func NewHandler(srv services.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: srv,
	}
}

func (h *accessTokenHandler) GetById(ctx *gin.Context) {
	tokenId := strings.TrimSpace(ctx.Param("access_token_id"))
	at, err := h.service.GetById(tokenId)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	ctx.JSON(http.StatusOK, at)
}

func (h *accessTokenHandler) Create(ctx *gin.Context) {

	var rq access_token.AccessTokenRequest

	if err := ctx.ShouldBindJSON(rq); err != nil {
		restErr := rest_errors.NewBadRequestError("wrong json body")
		ctx.JSON(restErr.Status, restErr)
		return
	}

	accessToken, err := h.service.Create(rq)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	ctx.JSON(http.StatusCreated, accessToken)
}
