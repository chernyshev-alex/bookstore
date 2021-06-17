package intf

import (
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/dao/mysql/gen"
	"github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
)

type UserDao interface {
	Get(int64) (*gen.FindUserRow, rest_errors.RestErr)
	Save(gen.InsertUserParams) (uint64, rest_errors.RestErr)
	Update(gen.UpdateUserParams) rest_errors.RestErr
	Delete(userId int64) rest_errors.RestErr
	FindByStatus(status string) ([]gen.FindByStatusRow, rest_errors.RestErr)
	FindByEmailAndPsw(gen.FindByEMailAndPswParams) (gen.FindByEMailAndPswRow, rest_errors.RestErr)
}
