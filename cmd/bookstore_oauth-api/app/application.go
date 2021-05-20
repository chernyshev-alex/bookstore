package app

import (
	"github.com/chernyshev-alex/bookstore_oauth-api/http"
	"github.com/chernyshev-alex/bookstore_oauth-api/repository/db"
	"github.com/chernyshev-alex/bookstore_oauth-api/repository/rest"
	"github.com/chernyshev-alex/bookstore_oauth-api/services"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func StartApplication() {

	//	_ = cassandra.GetSession()

	atService := services.NewService(rest.NewRestUsersRepository(), db.NewRepository())

	atHandler := http.NewHandler(atService)

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)

	router.Run(":8080")
}
