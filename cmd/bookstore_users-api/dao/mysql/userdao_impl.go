package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/dao/intf"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/dao/mysql/gen"
	"github.com/chernyshev-alex/bookstore_utils_go/logger"
	"github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
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
	return sql.NullString{s, true}
}

func (d *userDao) Get(id int32) (*gen.FindUserRow, rest_errors.RestErr) {
	result, err := d.dbq.FindUser(context.Background(), id)
	if err != nil {
		logger.Error("get user", err)
		return nil, rest_errors.NewInternalServerError("db error", err)
	}
	return &result, nil
}

func (d *userDao) Save(u gen.InsertUserParams) (uint64, rest_errors.RestErr) {
	result, err := d.dbq.InsertUser(context.Background(), u)
	if err != nil {
		logger.Error("get user", err)
		return uint64(0), rest_errors.NewInternalServerError("db error", err)
	}

	userId, err := result.LastInsertId()
	if err != nil {
		logger.Error("LastInsertId failed", err)
		return uint64(0), rest_errors.NewInternalServerError("db error", err)
	}
	return userId, nil
}

func (d *userDao) Update(u gen.UpdateUserParams) rest_errors.RestErr {
	result, err := d.dbq.UpdateUser(context.Background(), u)
	if err != nil {
		logger.Error("update user", err)
		return rest_errors.NewInternalServerError("db error", err)
	}
	nrows, err := result.RowsAffected()
	if nrows != 1 {
		logger.Error("affected", nrows)
		return rest_errors.NewInternalServerError("update error", err)
	}
	return nil
}

func (d *userDao) Delete(userId int64) rest_errors.RestErr {
	result, err := d.dbq.DeleteUser(context.Background(), userId)
	if err != nil {
		logger.Error("delete user", err)
		return rest_errors.NewInternalServerError("db error", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		logger.Error("RowsAffected", err)
		return rest_errors.NewInternalServerError("get rows affected failed", err)
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
		return rest_errors.NewInternalServerError("find user error", err)
	}
	return result, err
}
