//+build wireinject

package main

import (
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/google/wire"

	"github.com/chernyshev-alex/bookstore-oauth-go/oauth"
	"github.com/chernyshev-alex/bookstore_users-api/app"
	"github.com/chernyshev-alex/bookstore_users-api/controllers/ping"
	"github.com/chernyshev-alex/bookstore_users-api/controllers/users"
	"github.com/chernyshev-alex/bookstore_users-api/datasources/mysql/users_db"
	dao_users "github.com/chernyshev-alex/bookstore_users-api/domain/users"
	"github.com/chernyshev-alex/bookstore_users-api/services"
)

func inject(mysqlConf mysql.Config) app.Application {
	panic(wire.Build(

		wire.Value(http.DefaultClient),

		users_db.ProvideSqlClient,

		ping.ProvidePingController,

		services.ProvideUserService,
		wire.Bind(new(services.UsersServiceInterface), new(*services.UsersService)),

		dao_users.ProvideUserDao,
		wire.Bind(new(dao_users.UsersDAOInterface), new(*dao_users.UserDAO)),

		oauth.ProvideOAuthClient,
		wire.Bind(new(oauth.HTTPClientInterface), new(*http.Client)),
		wire.Bind(new(oauth.OAuthInterface), new(*oauth.OAuthClient)),

		users.ProvideUserController,
		app.ProvideApp,
	))
}
