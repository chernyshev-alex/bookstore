package user_services

import (
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/models"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/rest_errors"
)

type UserService interface {
	GetUser(int64) (*models.User, rest_errors.RestErr)
	CreateUser(models.User) (*models.User, rest_errors.RestErr)
	UpdateUser(bool, models.User) (*models.User, rest_errors.RestErr)
	DeleteUser(int64) rest_errors.RestErr
	SearchUsersByStatus(string) ([]models.User, rest_errors.RestErr)
	LoginUser(models.LoginRequest) (*models.User, rest_errors.RestErr)
}
