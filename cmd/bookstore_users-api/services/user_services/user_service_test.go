package user_services

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/dao/mysql/gen"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/models"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/rest_errors"
	"github.com/stretchr/testify/assert"
)

type userDaoMock struct {
	getFn     func(int64) (*gen.FindUserRow, rest_errors.RestErr)
	saveFn    func(gen.InsertUserParams) (int64, rest_errors.RestErr)
	updateFn  func(gen.UpdateUserParams) rest_errors.RestErr
	deleteFn  func(int64) rest_errors.RestErr
	findFn    func(string) ([]gen.FindByStatusRow, rest_errors.RestErr)
	findGetFn func(gen.FindByEMailAndPswParams) (gen.FindByEMailAndPswRow, rest_errors.RestErr)
}

func (m userDaoMock) Get(userId int64) (*gen.FindUserRow, rest_errors.RestErr) {
	return m.getFn(userId)
}
func (m userDaoMock) Save(p gen.InsertUserParams) (int64, rest_errors.RestErr) {
	return m.saveFn(p)
}
func (m userDaoMock) Update(p gen.UpdateUserParams) rest_errors.RestErr {
	return m.updateFn(p)
}
func (m userDaoMock) Delete(userId int64) rest_errors.RestErr {
	return m.deleteFn(userId)
}
func (m userDaoMock) FindByStatus(status string) ([]gen.FindByStatusRow, rest_errors.RestErr) {
	return m.findFn(status)
}
func (m userDaoMock) FindByEmailAndPsw(p gen.FindByEMailAndPswParams) (gen.FindByEMailAndPswRow, rest_errors.RestErr) {
	return m.findGetFn(p)
}

// helpers

func withMock(configFn func(*userDaoMock)) UserService {
	userDaoMock := new(userDaoMock)
	configFn(userDaoMock)
	return NewService(userDaoMock)
}

// tests

func TestGetUserOk(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.getFn = func(userId int64) (*gen.FindUserRow, rest_errors.RestErr) {
			return &gen.FindUserRow{
				ID: 1, FirstName: nillableStr("fname"),
			}, nil
		}
	})
	u, err := usersService.GetUser(1)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, u.Id)
}

func TestGetUserNotFound(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.getFn = func(userId int64) (*gen.FindUserRow, rest_errors.RestErr) {
			return nil, rest_errors.NewNotFoundError("user not found")
		}
	})
	u, err := usersService.GetUser(0)
	assert.Nil(t, u)
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestGetUserFailed(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.getFn = func(userId int64) (*gen.FindUserRow, rest_errors.RestErr) {
			return nil, rest_errors.NewInternalServerError("failed", errors.New("db error"))
		}
	})
	u, err := usersService.GetUser(-1)
	assert.Nil(t, u)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestCreateUserOk(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.saveFn = func(p gen.InsertUserParams) (int64, rest_errors.RestErr) {
			return int64(1), nil
		}
	})

	original := models.User{Email: "xxx@xxx.com", Password: "111"}
	created, err := usersService.CreateUser(original)

	assert.Nil(t, err)
	assert.EqualValues(t, created.Id, 1)
}

func TestCreateUserFailedMandatoryFieldsRequired(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.saveFn = func(p gen.InsertUserParams) (int64, rest_errors.RestErr) {
			return int64(1), nil
		}
	})

	original := models.User{}
	created, err := usersService.CreateUser(original)

	assert.Nil(t, created)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
}

func TestUpdateUserOk(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.getFn = func(userId int64) (*gen.FindUserRow, rest_errors.RestErr) {
			return &gen.FindUserRow{
				ID:          1,
				FirstName:   nillableStr("fname"),
				LastName:    nillableStr("lname"),
				Email:       "email",
				DateCreated: time.Now(),
				Status:      sql.NullString{},
			}, nil
		}
		mock.updateFn = func(p gen.UpdateUserParams) rest_errors.RestErr {
			p.FirstName = nillableStr("new_value")
			return nil
		}
	})

	original := models.User{Id: 1, FirstName: "new_value", Email: "xxx@xxx.com", Password: "111"}
	updated, err := usersService.UpdateUser(true, original)

	assert.Nil(t, err)
	assert.EqualValues(t, "new_value", updated.FirstName)
}

func TestUpdateUserNotFound(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.getFn = func(userId int64) (*gen.FindUserRow, rest_errors.RestErr) {
			return nil, rest_errors.NewNotFoundError("user not found")
		}
		mock.updateFn = func(p gen.UpdateUserParams) rest_errors.RestErr {
			p.FirstName = nillableStr("new_value")
			return nil
		}
	})

	original := models.User{Id: 1, FirstName: "new_value", Email: "xxx@xxx.com", Password: "111"}
	updated, err := usersService.UpdateUser(true, original)

	assert.Nil(t, updated)
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestDeleteUserOk(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.deleteFn = func(userId int64) rest_errors.RestErr {
			return nil
		}
	})
	err := usersService.DeleteUser(1)
	assert.Nil(t, err)
}

func TestDeleteFailed(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.deleteFn = func(userId int64) rest_errors.RestErr {
			return rest_errors.NewNotFoundError("not found")
		}
	})
	err := usersService.DeleteUser(-1)
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestSearchUserOk(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.findFn = func(status string) ([]gen.FindByStatusRow, rest_errors.RestErr) {
			return []gen.FindByStatusRow{{
				ID: 1, FirstName: nillableStr("fname"),
			}}, nil
		}
	})
	ls, err := usersService.SearchUsersByStatus(models.STATUS_ACTIVE)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, len(ls))
}

func TestLoginUserOk(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.findGetFn = func(p gen.FindByEMailAndPswParams) (gen.FindByEMailAndPswRow, rest_errors.RestErr) {
			return gen.FindByEMailAndPswRow{
				ID:        1,
				FirstName: nillableStr("fname"),
			}, nil
		}
	})
	lrq := models.LoginRequest{Email: "xxx@xxx.com", Password: "111"}
	u, err := usersService.LoginUser(lrq)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, u.Id)
}

func TestLoginNotFound(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.findGetFn = func(p gen.FindByEMailAndPswParams) (gen.FindByEMailAndPswRow, rest_errors.RestErr) {
			return gen.FindByEMailAndPswRow{
				ID: 0,
			}, rest_errors.NewNotFoundError("not found")
		}
	})
	lrq := models.LoginRequest{Email: "xxx@xxx.com", Password: "111"}
	u, err := usersService.LoginUser(lrq)

	assert.Nil(t, u)
	assert.EqualValues(t, err.Status(), http.StatusNotFound)
}
