//+build wireinject

package main

import (
	"net/http"

	"github.com/chernyshev-alex/bookstore-oauth-go/oauth"
	"github.com/chernyshev-alex/bookstore_items-api/app"
	"github.com/chernyshev-alex/bookstore_items-api/config"
	"github.com/chernyshev-alex/bookstore_items-api/controllers"
	"github.com/chernyshev-alex/bookstore_items-api/services"
	"github.com/google/wire"
)

func inject(appConfig *config.Config) *app.Application {
	panic(
		wire.Build(

			app.NewApp,
			controllers.NewItemController,
			services.NewItemsService,

			oauth.ProvideOAuthClient,
			wire.Bind(new(oauth.OAuthInterface), new(*oauth.OAuthClient)),

			wire.Bind(new(oauth.HTTPClientInterface), new(*http.Client)),
			wire.Value(http.DefaultClient),
		))
}
