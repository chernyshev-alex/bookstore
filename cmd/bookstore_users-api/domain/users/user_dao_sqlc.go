package users

import (
	"context"
	"database/sql"
	"time"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/sqlc/user_dao"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/logger"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/rest_errors"
)

type UserDao struct {
	SqlClient *sql.DB
	dbq       *user_dao.Queries
}

func NewUsersDao(client *sql.DB) *UserDao {
	return &UserDao{SqlClient: client,
		dbq: user_dao.New(client),
	}
}

func asNullStr(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}

func (d *UserDao) Get(id int64) (*User, rest_errors.RestErr) {
	result, err := d.dbq.FindUser(context.Background(), int32(id))
	if err != nil {
		logger.Error("get user", err)
		return nil, rest_errors.NewInternalServerError("db error", err)
	}

	u := User{
		Id:          int64(result.ID),
		FirstName:   result.FirstName.String,
		LastName:    result.LastName.String,
		Email:       result.Email,
		DateCreated: result.DateCreated.String(),
		Status:      result.Status.String,
	}
	return &u, nil
}

func (d *UserDao) Save(u *User) rest_errors.RestErr {
	args := user_dao.InsertUserParams{
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

func (d *UserDao) Update(u *User) rest_errors.RestErr {
	args := user_dao.UpdateUserParams{
		FirstName: asNullStr(u.FirstName),
		LastName:  asNullStr(u.LastName),
		Email:     u.Email,
	}
	_, err := d.dbq.UpdateUser(context.Background(), args)
	if err != nil {
		logger.Error("get user", err)
		return rest_errors.NewInternalServerError("db error", err)
	}
	return nil
}

func (d *UserDao) Delete(u *User) rest_errors.RestErr {
	result, err := d.dbq.DeleteUser(context.Background(), int32(u.Id))
	if err != nil {
		logger.Error("delete user", err)
		return rest_errors.NewInternalServerError("db error", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		logger.Error("RowsAffected", err)
	}
	if rows == 0 {
		logger.Info("no row to delete")
	}
	return nil
}

func (d *UserDao) FindByStatus(status string) ([]User, rest_errors.RestErr) {
	result, err := d.dbq.FindByStatus(context.Background(), asNullStr(status))
	if err != nil {
		logger.Error("FindByStatus", err)
		return nil, rest_errors.NewInternalServerError("db error", err)
	}

	ls := make([]User, 0, len(result))
	for _, rec := range result {
		u := User{
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

func (d *UserDao) FindByEmailAndPsw(u *User) rest_errors.RestErr {
	if u == nil {
		logger.Error("FindByEmailAndPsw not expected null", nil)
		return rest_errors.NewBadRequestError("bad input")
	}
	args := user_dao.FindByEMailAndPswParams{
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
