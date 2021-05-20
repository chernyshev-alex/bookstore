package rest

import (
	"encoding/json"
	"time"

	"github.com/chernyshev-alex/bookstore_users-api/domain/users"
	"github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
	"github.com/mercadolibre/golang-restclient/rest"
)

var (
	restClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8082",
		Timeout: 100 * time.Millisecond,
	}
)

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, rest_errors.RestErr)
}

type usersRepository struct {
}

func NewRestUsersRepository() RestUsersRepository {
	return &usersRepository{}
}

func (r *usersRepository) LoginUser(email string, psw string) (*users.User, rest_errors.RestErr) {
	req := users.LoginRequest{
		Email:    email,
		Password: psw,
	}
	resp := restClient.Post("/users/login", req)
	if resp == nil || resp.Response == nil {
		return nil, rest_errors.NewInternalServerError("user login timeout")
	}
	if resp.StatusCode > 299 {
		var jsonErr rest_errors.RestErr
		if err := json.Unmarshal(resp.Bytes(), &jsonErr); err != nil {
			return nil, rest_errors.NewInternalServerError("unmarshal response error")
		}
		return nil, &jsonErr
	}

	var user users.User
	if err := json.Unmarshal(resp.Bytes(), &user); err != nil {
		return nil, rest_errors.NewInternalServerError("bad response from user api (login)")
	}
	return &user, nil
}
