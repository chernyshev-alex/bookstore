//+build wireinject

package main

import (
	"net/http"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_items-api/app"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_items-api/config"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_items-api/controllers"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_items-api/domain/items"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_items-api/services"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore-oauth-go/oauth"
	"github.com/google/wire"
)

func inject(appConfig *config.Config) *app.Application {
	panic(
		wire.Build(

			app.NewApp,

			app.ProvideOAuthClient,
			wire.Bind(new(oauth.OAuthInterface), new(*oauth.OAuthClient)),

			controllers.NewItemController,
			services.NewItemsService,
			items.NewItemPersister,

			wire.Bind(new(oauth.HttpClientInterface), new(*http.Client)),
			wire.Value(http.DefaultClient),
		))
}
