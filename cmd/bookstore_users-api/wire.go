//+build wireinject

package main

import (
	"net/http"

	"github.com/google/wire"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/app"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/conf"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/controllers"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/dao/mysql"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/services/user_services"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore-oauth-go/oauth"
)

func inject(conf *conf.Config) app.Application {
	panic(wire.Build(

		app.ProvideApp,
		controllers.ProvidePingController,
		controllers.ProvideUserController,

		app.NewOAuthClient,
		wire.Bind(new(oauth.OAuthInterface), new(*oauth.OAuthClient)),

		user_services.NewService,

		mysql.NewUserDao,

		wire.Bind(new(oauth.HttpClientInterface), new(*http.Client)),
		mysql.MakeConfig,
		mysql.NewSqlClient,

		wire.Value(http.DefaultClient),
	))

}
