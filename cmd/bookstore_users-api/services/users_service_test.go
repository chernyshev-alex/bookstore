package services

import (
	"errors"
	"net/http"
	"testing"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/domain/users"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/rest_errors"
	"github.com/stretchr/testify/assert"
)

type userDaoMock struct {
	getFn     func(int64) (*users.User, rest_errors.RestErr)
	modifyFn  func(*users.User) rest_errors.RestErr
	findFn    func(string) ([]users.User, rest_errors.RestErr)
	findGetFn func(*users.User) rest_errors.RestErr
}

func (m userDaoMock) Get(userId int64) (*users.User, rest_errors.RestErr) {
	return m.getFn(userId)
}
func (m userDaoMock) Save(u *users.User) rest_errors.RestErr {
	return m.modifyFn(u)
}
func (m userDaoMock) Update(u *users.User) rest_errors.RestErr {
	return m.modifyFn(u)
}
func (m userDaoMock) Delete(u *users.User) rest_errors.RestErr {
	return m.modifyFn(u)
}
func (m userDaoMock) FindByStatus(status string) ([]users.User, rest_errors.RestErr) {
	return m.findFn(status)
}
func (m userDaoMock) FindByEmailAndPsw(u *users.User) rest_errors.RestErr {
	return m.findGetFn(u)
}

// helpers

func withMock(configFn func(*userDaoMock)) *UsersService {
	userDaoMock := new(userDaoMock)
	configFn(userDaoMock)
	return ProvideUserService(userDaoMock)
}

// tests

func TestGetUserOk(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.getFn = func(userId int64) (*users.User, rest_errors.RestErr) {
			return &users.User{Id: 1, FirstName: "fname", LastName: "lname", Email: "email"}, nil
		}
	})
	u, err := usersService.GetUser(1)
	assert.Nil(t, err)
	assert.EqualValues(t, u.Id, 1)
}

func TestGetUserNotFound(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.getFn = func(userId int64) (*users.User, rest_errors.RestErr) {
			return nil, rest_errors.NewNotFoundError("user not found")
		}
	})
	u, err := usersService.GetUser(0)
	assert.Nil(t, u)
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestGetUserFailed(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.getFn = func(userId int64) (*users.User, rest_errors.RestErr) {
			return nil, rest_errors.NewInternalServerError("failed", errors.New("db error"))
		}
	})
	u, err := usersService.GetUser(-1)
	assert.Nil(t, u)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestCreateUserOk(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.modifyFn = func(u *users.User) rest_errors.RestErr {
			u.Id = 1
			return nil
		}
	})

	original := users.User{Email: "xxx@xxx.com", Password: "111"}
	created, err := usersService.CreateUser(original)

	assert.Nil(t, err)
	assert.EqualValues(t, created.Id, 1)
}

func TestCreateUserFailedMandatoryFieldsRequired(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.modifyFn = func(u *users.User) rest_errors.RestErr {
			u.Id = 1
			return nil
		}
	})

	original := users.User{}
	created, err := usersService.CreateUser(original)

	assert.Nil(t, created)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
}

func TestUpdateUserOk(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.getFn = func(userId int64) (*users.User, rest_errors.RestErr) {
			return &users.User{Id: 1, FirstName: "old_value"}, nil
		}
		mock.modifyFn = func(u *users.User) rest_errors.RestErr {
			u.FirstName = "new_value"
			return nil
		}
	})

	original := users.User{Id: 1, FirstName: "new_value", Email: "xxx@xxx.com", Password: "111"}
	updated, err := usersService.UpdateUser(true, original)

	assert.Nil(t, err)
	assert.EqualValues(t, "new_value", updated.FirstName)
}

func TestUpdateUserNotFound(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.getFn = func(userId int64) (*users.User, rest_errors.RestErr) {
			return nil, rest_errors.NewNotFoundError("user not found")
		}
		mock.modifyFn = func(u *users.User) rest_errors.RestErr {
			u.FirstName = "new_value"
			return nil
		}
	})

	original := users.User{Id: 1, FirstName: "new_value", Email: "xxx@xxx.com", Password: "111"}
	updated, err := usersService.UpdateUser(true, original)

	assert.Nil(t, updated)
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestDeleteUserOk(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.modifyFn = func(u *users.User) rest_errors.RestErr {
			return nil
		}
	})
	err := usersService.DeleteUser(1)
	assert.Nil(t, err)
}

func TestDeleteFailed(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.modifyFn = func(u *users.User) rest_errors.RestErr {
			return rest_errors.NewNotFoundError("not found")
		}
	})
	err := usersService.DeleteUser(-1)
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestSearchUserOk(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.findFn = func(status string) ([]users.User, rest_errors.RestErr) {
			return []users.User{{Id: 1}}, nil
		}
	})
	ls, err := usersService.SearchUser(users.StatusActive)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, len(ls))
}

func TestLoginUserOk(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.findGetFn = func(u *users.User) rest_errors.RestErr {
			u.Id = 1
			return nil
		}
	})
	lrq := users.LoginRequest{Email: "xxx@xxx.com", Password: "111"}
	u, err := usersService.LoginUser(lrq)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, u.Id)
}

func TestLoginNotFound(t *testing.T) {
	usersService := withMock(func(mock *userDaoMock) {
		mock.findGetFn = func(u *users.User) rest_errors.RestErr {
			return rest_errors.NewNotFoundError("not found")
		}
	})
	lrq := users.LoginRequest{Email: "xxx@xxx.com", Password: "111"}
	u, err := usersService.LoginUser(lrq)

	assert.Nil(t, u)
	assert.EqualValues(t, err.Status(), http.StatusNotFound)
}
