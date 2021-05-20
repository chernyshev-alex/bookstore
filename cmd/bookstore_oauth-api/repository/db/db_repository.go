package db

import (
	"github.com/chernyshev-alex/bookstore_oauth-api/clients/cassandra"
	"github.com/chernyshev-alex/bookstore_oauth-api/domain/access_token"
	"github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
	"github.com/gocql/gocql"
)

const (
	getAccessToken    = "select access_token, user_id, cclient_id, expires from access_tokens where sccess_token=?;"
	createAccessToken = "insert into access_tokens(access_token, user_id, cclient_id, expires) values(?,?,?,?);"
	updateExpires     = "updte access_tokens set expires=? where  access_token=?;"
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, rest_errors.RestErr)
	Create(access_token.AccessToken) rest_errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) rest_errors.RestErr
}

type dbRepository struct {
}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, rest_errors.RestErr) {
	var result access_token.AccessToken

	ss := cassandra.GetSession()
	if err := ss.Query(getAccessToken, id).Scan(&result.AccessToken,
		&result.UserId, &result.ClientId, &result.Expires); err != nil {

		if err == gocql.ErrNotFound {
			return nil, rest_errors.NewNotFoundError("token not found")
		}
		return nil, rest_errors.NewInternalServerError(err.Error())
	}

	return &result, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) rest_errors.RestErr {
	ss := cassandra.GetSession()
	defer ss.Close()
	if err := ss.Query(createAccessToken,
		&at.AccessToken,
		&at.UserId,
		&at.ClientId,
		&at.Expires).Exec(); err != nil {
		return rest_errors.NewInternalServerError(err.Error())
	}
	return nil
}

func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken) rest_errors.RestErr {
	ss := cassandra.GetSession()
	defer ss.Close()
	if err := ss.Query(updateExpires, &at.Expires, &at.AccessToken).Exec(); err != nil {
		return rest_errors.NewInternalServerError(err.Error())
	}
	return nil
}
