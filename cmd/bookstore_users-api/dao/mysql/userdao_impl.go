package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/dao/intf"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/dao/mysql/gen"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/logger"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/rest_errors"
)

type userDao struct {
	SqlClient *sql.DB
	dbq       *gen.Queries
}

func NewUserDao(client *sql.DB) intf.UserDao {
	return &userDao{SqlClient: client,
		dbq: gen.New(client),
	}
}

func nillableStr(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}

func (d *userDao) Get(id int64) (*gen.FindUserRow, rest_errors.RestErr) {
	fmt.Println("called Get", id)
	result, err := d.dbq.FindUser(context.Background(), int32(id))
	if err != nil {
		logger.Error("get user", err)
		return nil, rest_errors.NewInternalServerError("db error", err)
	}
	return &result, nil
}

func (d *userDao) Save(u gen.InsertUserParams) (int64, rest_errors.RestErr) {
	result, err := d.dbq.InsertUser(context.Background(), u)
	if err != nil {
		logger.Error("save user", err)
		return -1, rest_errors.NewInternalServerError("db error", err)
	}

	userId, err := result.LastInsertId()
	if err != nil {
		logger.Error("LastInsertId failed", err)
		return -1, rest_errors.NewInternalServerError("db error", err)
	}
	return userId, nil
}

func (d *userDao) Update(u gen.UpdateUserParams) rest_errors.RestErr {
	_, err := d.dbq.UpdateUser(context.Background(), u)
	if err != nil {
		logger.Error("update user", err)
		return rest_errors.NewInternalServerError("db error", err)
	}
	return nil
}

func (d *userDao) Delete(userId int64) rest_errors.RestErr {
	result, err := d.dbq.DeleteUser(context.Background(), int32(userId))
	if err != nil {
		logger.Error("delete user", err)
		return rest_errors.NewInternalServerError("db error", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		logger.Info("RowsAffected failed")
		return nil
	}
	if rows == 0 {
		msg := fmt.Sprintf("no row to delete %d", userId)
		logger.Info(msg)
		return rest_errors.NewNotFoundError(msg)
	}
	return nil
}

func (d *userDao) FindByStatus(status string) ([]gen.FindByStatusRow, rest_errors.RestErr) {
	result, err := d.dbq.FindByStatus(context.Background(), nillableStr(status))
	if err != nil {
		logger.Error("FindByStatus", err)
		return nil, rest_errors.NewInternalServerError("db error", err)
	}
	return result, nil
}

func (d *userDao) FindByEmailAndPsw(arg gen.FindByEMailAndPswParams) (gen.FindByEMailAndPswRow, rest_errors.RestErr) {
	result, err := d.dbq.FindByEMailAndPsw(context.Background(), arg)
	if err != nil {
		logger.Error("find user", err)
		return result, rest_errors.NewInternalServerError("find user error", err)
	}
	return result, nil
}
