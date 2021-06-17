package app

import (
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/conf"
	c "github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/controllers"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore-oauth-go/oauth"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/logger"
	"github.com/gin-gonic/gin"
)

type Application struct {
	router         *gin.Engine
	pingController *c.PingController
	userController *c.UserController
	appConfig      *conf.Config
}

func ProvideApp(appConfig *conf.Config, pingController *c.PingController,
	userController *c.UserController) Application {
	return Application{
		router:         gin.Default(),
		pingController: pingController,
		userController: userController,
		appConfig:      appConfig,
	}
}

func NewOAuthClient(httpClient oauth.HttpClientInterface, appConf *conf.Config) *oauth.OAuthClient {
	return oauth.NewAuthClient(httpClient, appConf.OAuth.URL)
}

func (app *Application) mapUrls() {
	app.router.GET("/ping", app.pingController.Ping)
	app.router.POST("/users", app.userController.Create)
	app.router.GET("/users/:user_id", app.userController.Get)
	app.router.PUT("/users/:user_id", app.userController.Update)
	app.router.PATCH("/users/:user_id", app.userController.Update)
	app.router.DELETE("/users/:user_id", app.userController.Delete)
	app.router.GET("/internal/users/search", app.userController.Search)
	app.router.POST("/users/login", app.userController.Login)
}

func (app *Application) StartApp() {
	app.mapUrls()
	logger.Info("listening on port  " + app.appConfig.Server.Port)
	app.router.Run(":" + app.appConfig.Server.Port)
}
