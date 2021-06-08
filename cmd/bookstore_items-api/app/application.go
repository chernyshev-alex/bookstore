package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_items-api/client/es"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_items-api/config"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_items-api/controllers"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore-oauth-go/oauth"
	"github.com/gorilla/mux"
)

type Application struct {
	config *config.Config
	router *mux.Router
	items  controllers.ItemControllerInterface
}

func NewApp(appConf *config.Config,
	itemsController controllers.ItemControllerInterface) *Application {
	return &Application{
		config: appConf,
		router: mux.NewRouter(),
		items:  itemsController,
	}
}

// wired OAuthClient constructor
func ProvideOAuthClient(httpClient oauth.HttpClientInterface, appConf *config.Config) *oauth.OAuthClient {
	return oauth.NewAuthClient(httpClient, appConf.OAuth.URL)
}

func (app *Application) StartApp() {

	es.Init()

	app.router.HandleFunc("/ping", app.items.Ping).Methods(http.MethodGet)
	app.router.HandleFunc("/items", app.items.Create).Methods(http.MethodPost)
	app.router.HandleFunc("/items/{id}", app.items.Get).Methods(http.MethodGet)
	app.router.HandleFunc("/items/search", app.items.Search).Methods(http.MethodPost)

	listenAddress := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port)
	srv := http.Server{
		Addr:         listenAddress,
		Handler:      app.router,
		WriteTimeout: 200 * time.Millisecond,
		ReadTimeout:  20 * time.Millisecond,
		IdleTimeout:  10 * time.Millisecond,
	}

	fmt.Println("listening on ", listenAddress)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
