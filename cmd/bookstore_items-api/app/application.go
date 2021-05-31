package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chernyshev-alex/bookstore_items-api/client/es"
	"github.com/chernyshev-alex/bookstore_items-api/controllers"

	"github.com/gorilla/mux"
)

type Application struct {
	router *mux.Router
	items  controllers.ItemControllerInterface
}

func NewApp(itemsController controllers.ItemControllerInterface) *Application {
	return &Application{
		router: mux.NewRouter(),
		items:  itemsController,
	}
}

func (app *Application) StartApp() {

	es.Init()

	app.router.HandleFunc("/ping", app.items.Ping).Methods(http.MethodGet)
	app.router.HandleFunc("/items", app.items.Create).Methods(http.MethodPost)
	app.router.HandleFunc("/items/{id}", app.items.Get).Methods(http.MethodGet)
	app.router.HandleFunc("/items/search", app.items.Search).Methods(http.MethodPost)

	srv := http.Server{
		Addr:         "127.0.0.1:8080",
		Handler:      app.router,
		WriteTimeout: 200 * time.Millisecond,
		ReadTimeout:  20 * time.Millisecond,
		IdleTimeout:  10 * time.Millisecond,
	}

	fmt.Println("bookstore items started on ", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
