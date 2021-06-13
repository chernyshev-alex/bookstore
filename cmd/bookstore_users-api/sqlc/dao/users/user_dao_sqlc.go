package users

import (
	"context"
	"database/sql"
	"time"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/domain/users"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/sqlc/dao/users/gen"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/logger"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/rest_errors"
)

type IUserDao interface {
	Get(id int32) (*users.User, rest_errors.RestErr)
	Save(u *users.User) rest_errors.RestErr
	Update(u *users.User) rest_errors.RestErr
	Delete(u *users.User) rest_errors.RestErr
	FindByStatus(status string) ([]users.User, rest_errors.RestErr)
	FindByEmailAndPsw(u *users.User) rest_errors.RestErr
}
type userDao struct {
	SqlClient *sql.DB
	dbq       *gen.Queries
}

func NewUserDao(client *sql.DB) IUserDao {
	return &userDao{SqlClient: client,
		dbq: gen.New(client),
	}
}

func asNullStr(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}

func (d *userDao) Get(id int32) (*users.User, rest_errors.RestErr) {
	result, err := d.dbq.FindUser(context.Background(), id)
	if err != nil {
		logger.Error("get user", err)
		return nil, rest_errors.NewInternalServerError("db error", err)
	}

	u := users.User{
		Id:          int64(result.ID),
		FirstName:   result.FirstName.String,
		LastName:    result.LastName.String,
		Email:       result.Email,
		DateCreated: result.DateCreated.String(),
		Status:      result.Status.String,
	}
	return &u, nil
}

func (d *userDao) Save(u *users.User) rest_errors.RestErr {
	args := gen.InsertUserParams{
		FirstName:   asNullStr(u.FirstName),
		LastName:    asNullStr(u.LastName),
		Email:       u.Email,
		DateCreated: time.Now(),
		Status:      asNullStr(u.Status),
		Password:    asNullStr(u.Password),
	}

	result, err := d.dbq.InsertUser(context.Background(), args)
	if err != nil {
		logger.Error("get user", err)
		return rest_errors.NewInternalServerError("db error", err)
	}

	userId, err := result.LastInsertId()
	if err != nil {
		logger.Error("LastInsertId failed", err)
		return rest_errors.NewInternalServerError("db error", err)
	}
	u.Id = userId
	return nil
}

func (d *userDao) Update(u *users.User) rest_errors.RestErr {
	args := gen.UpdateUserParams{
		ID:        int32(u.Id),
		FirstName: asNullStr(u.FirstName),
		LastName:  asNullStr(u.LastName),
		Email:     u.Email,
	}
	_, err := d.dbq.UpdateUser(context.Background(), args)
	if err != nil {
		logger.Error("update user", err)
		return rest_errors.NewInternalServerError("db error", err)
	}
	return nil
}

func (d *userDao) Delete(u *users.User) rest_errors.RestErr {
	result, err := d.dbq.DeleteUser(context.Background(), int32(u.Id))
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
		logger.Info("no row to delete")
		return rest_errors.NewNotFoundError(u.Email)
	}
	return nil
}

func (d *userDao) FindByStatus(status string) ([]users.User, rest_errors.RestErr) {
	result, err := d.dbq.FindByStatus(context.Background(), asNullStr(status))
	if err != nil {
		logger.Error("FindByStatus", err)
		return nil, rest_errors.NewInternalServerError("db error", err)
	}

	ls := make([]users.User, 0, len(result))
	for _, rec := range result {
		u := users.User{
			Id:        int64(rec.ID),
			FirstName: rec.FirstName.String,
			LastName:  rec.LastName.String,
			Email:     rec.Email,
			//	DateCreated: date_utils.rec.DateCreated,  TODO
			Status: rec.Status.String,
		}
		ls = append(ls, u)
	}
	return ls, nil
}

func (d *userDao) FindByEmailAndPsw(u *users.User) rest_errors.RestErr {
	if u == nil {
		logger.Error("FindByEmailAndPsw not expected null", nil)
		return rest_errors.NewBadRequestError("bad input")
	}
	args := gen.FindByEMailAndPswParams{
		Email:    u.Email,
		Password: asNullStr(u.Password),
		Status:   asNullStr(u.Status),
	}
	result, err := d.dbq.FindByEMailAndPsw(context.Background(), args)
	if err != nil {
		logger.Error("get user", err)
		return rest_errors.NewInternalServerError("db error", err)
	}

	u.Id = int64(result.ID)
	u.FirstName = result.FirstName.String
	u.LastName = result.LastName.String
	u.Email = result.Email
	//	DateCreated: date_utils.rec.DateCreated,  TODO
	u.Status = result.Status.String

	return nil
}
