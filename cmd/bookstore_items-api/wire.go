package main

//+build wireinject

import (
	"net/http"

	"github.com/chernyshev-alex/bookstore-oauth-go/oauth"
	"github.com/chernyshev-alex/bookstore_items-api/app"
	"github.com/chernyshev-alex/bookstore_items-api/controllers"
	"github.com/chernyshev-alex/bookstore_items-api/services"
	"github.com/google/wire"
)

func inject() *app.Application {
	panic(
		wire.Build(

			wire.Value(http.DefaultClient),
			oauth.ProvideOAuthClient,
			wire.Bind(new(oauth.HTTPClientInterface), new(*http.Client)),
			wire.Bind(new(oauth.OAuthInterface), new(*oauth.OAuthClient)),

			services.NewItemsService,
			controllers.NewItemController,
			app.NewApp,
		))
}
