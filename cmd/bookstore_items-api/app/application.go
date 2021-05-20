package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chernyshev-alex/bookstore_items-api/client/es"
	"github.com/chernyshev-alex/bookstore_items-api/controllers"

	"github.com/gorilla/mux"
)

var (
	router = mux.NewRouter()
)

func StartApp() {

	es.Init()

	router.HandleFunc("/ping", controllers.ItemController.Ping).Methods(http.MethodGet)
	router.HandleFunc("/items", controllers.ItemController.Create).Methods(http.MethodPost)
	router.HandleFunc("/items{id}", controllers.ItemController.Get).Methods(http.MethodGet)
	router.HandleFunc("/items/search", controllers.ItemController.Search).Methods(http.MethodPost)

	srv := http.Server{
		Addr:         "127.0.0.1:8080",
		Handler:      router,
		WriteTimeout: 200 * time.Millisecond,
		ReadTimeout:  20 * time.Millisecond,
		IdleTimeout:  10 * time.Millisecond,
	}

	fmt.Println("bookstore items started on ", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
