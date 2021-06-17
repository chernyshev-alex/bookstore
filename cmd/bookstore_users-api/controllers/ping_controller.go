package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	Controller = PingController{}
)

type PingController struct{}

func ProvidePingController() *PingController {
	return &PingController{}
}

func (pc PingController) Ping(c *gin.Context) {
	c.String(http.StatusOK, "**** pong *****")
}
