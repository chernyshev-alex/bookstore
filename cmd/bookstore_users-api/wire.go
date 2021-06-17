package main

//+build wireinject

import (
	"net/http"

	"github.com/google/wire"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/app"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/conf"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/controllers"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/dao/mysql"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/services"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore-oauth-go/oauth"
)

func inject(conf conf.Config) app.Application {
	panic(wire.Build(

		app.ProvideApp,
		controllers.ProvideUserController,
		controllers.ProvidePingController,

		app.NewOAuthClient,
		wire.Bind(new(oauth.OAuthInterface), new(*oauth.OAuthClient)),

		services.NewService,
		//	wire.Bind(new(services.UsersServiceInterface), new(*services.UsersService)),

		mysql.NewUserDao,
		wire.Bind(new(mysql.UsersDAOInterface), new(*mysql.UserDAO)),
		//	users_db.ProvideSqlClient,

		wire.Value(http.DefaultClient),
	))
	return nil
}
