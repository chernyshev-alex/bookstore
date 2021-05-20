package users

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/chernyshev-alex/bookstore_utils_go/logger"
	"github.com/chernyshev-alex/bookstore_utils_go/mysql_utils"
	"github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
)

// data access layer
const (
	queryInsertUser        = "insert into users(first_name,last_name,email,date_created, status, password) values(?,?,?,?,?,?);"
	queryUpdateUser        = "update users set first_name=?,last_name=?,email=? where id=?;"
	queryFindUser          = "select id, first_name,last_name,email,date_created, status from users where id=?;"
	queryDeleteUser        = "delete from users where id=?;"
	queryFindByStatus      = "select id, first_name,last_name,email,date_created, status from users where status=?;"
	queryFindByEMailAndPsw = "select id, first_name,last_name,email,date_created, status from users where email=? and password=? and status=?;"
)

type UserDAO struct {
	SqlClient *sql.DB
}

type UsersDAOInterface interface {
	Get(id int64) (*User, rest_errors.RestErr)
	Save(u *User) rest_errors.RestErr
	Update(u *User) rest_errors.RestErr
	Delete(u *User) rest_errors.RestErr
	FindByStatus(status string) ([]User, rest_errors.RestErr)
	FindByEmailAndPsw(u *User) rest_errors.RestErr
}

func ProvideUserDao(client *sql.DB) *UserDAO {
	return &UserDAO{
		SqlClient: client,
	}
}

func (d UserDAO) Get(id int64) (*User, rest_errors.RestErr) {
	stmt, err := d.SqlClient.Prepare(queryFindUser)
	if err != nil {
		logger.Error("get prepare", err)
		return nil, rest_errors.NewInternalServerError("db error", err)
	}

	defer stmt.Close()

	u := User{Id: id}
	if getErr := stmt.QueryRow(u.Id).Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email,
		&u.DateCreated, &u.Status); getErr != nil {

		logger.Error("get user", getErr)
		if errors.Is(getErr, sql.ErrNoRows) {
			return nil, rest_errors.NewNotFoundError("user not found")
		}
		return nil, rest_errors.NewInternalServerError("db error", getErr)
	}
	return &u, nil
}

func (d UserDAO) Save(u *User) rest_errors.RestErr {
	stmt, err := d.SqlClient.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("prepare failed", err)
		return rest_errors.NewInternalServerError(err.Error(), err)
	}

	defer stmt.Close()

	result, saveErr := stmt.Exec(u.FirstName, u.LastName, u.Email, u.DateCreated, u.Status, u.Password)
	if saveErr != nil {
		logger.Error("exec failed", err)
		return mysql_utils.ParseErrors(saveErr)
	}

	userId, err := result.LastInsertId()
	if err != nil {
		logger.Error("get last id failed", err)
		return mysql_utils.ParseErrors(err)
	}

	u.Id = userId
	return nil
}

func (d UserDAO) Update(u *User) rest_errors.RestErr {
	stmt, err := d.SqlClient.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("prepare failed", err)
		return rest_errors.NewInternalServerError(err.Error(), err)
	}

	defer stmt.Close()

	if _, err := stmt.Exec(u.FirstName, u.LastName, u.Email, u.Id); err != nil {
		logger.Error("update failed", err)
		return mysql_utils.ParseErrors(err)

	}
	return nil
}

func (d UserDAO) Delete(u *User) rest_errors.RestErr {
	stmt, err := d.SqlClient.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("prepare failed", err)
		return rest_errors.NewInternalServerError(err.Error(), err)
	}

	defer stmt.Close()

	if _, err := stmt.Exec(u.Id); err != nil {
		logger.Error("delete failed", err)
		return mysql_utils.ParseErrors(err)
	}
	return nil
}

func (d UserDAO) FindByStatus(status string) ([]User, rest_errors.RestErr) {
	stmt, err := d.SqlClient.Prepare(queryFindByStatus)
	if err != nil {
		logger.Error("prepare failed", err)
		return nil, rest_errors.NewInternalServerError(err.Error(), err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("query faailed", err)
		return nil, rest_errors.NewInternalServerError(err.Error(), err)
	}
	defer rows.Close()

	result := make([]User, 0)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.DateCreated, &u.Status); err != nil {
			logger.Error("scan failed", err)
			return nil, mysql_utils.ParseErrors(err)
		}
		result = append(result, u)
	}
	if len(result) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("not found users with status %s", status))
	}
	return result, nil
}

func (d UserDAO) FindByEmailAndPsw(u *User) rest_errors.RestErr {
	stmt, err := d.SqlClient.Prepare(queryFindByEMailAndPsw)
	if err != nil {
		logger.Error("prepare failed", err)
		return rest_errors.NewInternalServerError("db error", err)
	}

	defer stmt.Close()

	err = stmt.QueryRow(u.Email, u.Password, StatusActive).Scan(&u.Id, &u.FirstName,
		&u.LastName, &u.Email, &u.DateCreated, &u.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			return rest_errors.NewNotFoundError("user not found")
		}
		logger.Error("query failed", err)
		return rest_errors.NewInternalServerError("db error", err)
	}
	return nil
}
