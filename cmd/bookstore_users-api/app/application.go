package app

import (
	"github.com/chernyshev-alex/bookstore_users-api/controllers/ping"
	"github.com/chernyshev-alex/bookstore_users-api/controllers/users"
	"github.com/chernyshev-alex/bookstore_utils_go/logger"
	"github.com/gin-gonic/gin"
)

type Application struct {
	router         *gin.Engine
	pingController *ping.PingController
	userController *users.UserController
}

func ProvideApp(pingController *ping.PingController, userController *users.UserController) Application {
	return Application{
		router:         gin.Default(),
		pingController: pingController,
		userController: userController,
	}
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

func (app *Application) StartApplication() {
	app.mapUrls()
	logger.Info("starting ..")
	app.router.Run(":8081")
}
