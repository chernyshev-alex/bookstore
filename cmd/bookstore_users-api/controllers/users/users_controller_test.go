package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chernyshev-alex/bookstore_users-api/domain/users"
	"github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type usersServiceMock struct {
	createUserFn func(users.User) (*users.User, rest_errors.RestErr)
	getUserFn    func(int64) (*users.User, rest_errors.RestErr)
	updateUserFn func(bool, users.User) (*users.User, rest_errors.RestErr)
	deleteUserFn func(int64) rest_errors.RestErr
	searchUserFn func(string) (users.Users, rest_errors.RestErr)
	loginUserFn  func(users.LoginRequest) (*users.User, rest_errors.RestErr)
}

func (m usersServiceMock) CreateUser(u users.User) (*users.User, rest_errors.RestErr) {
	return m.createUserFn(u)
}
func (m usersServiceMock) GetUser(userId int64) (*users.User, rest_errors.RestErr) {
	return m.getUserFn(userId)
}
func (m usersServiceMock) UpdateUser(isPartial bool, u users.User) (*users.User, rest_errors.RestErr) {
	return m.updateUserFn(isPartial, u)
}
func (m usersServiceMock) DeleteUser(userId int64) rest_errors.RestErr {
	return m.deleteUserFn(userId)
}
func (m usersServiceMock) SearchUser(status string) (users.Users, rest_errors.RestErr) {
	return m.searchUserFn(status)
}
func (m usersServiceMock) LoginUser(req users.LoginRequest) (*users.User, rest_errors.RestErr) {
	return m.loginUserFn(req)
}

type oauthServiceMock struct {
	isPublicFn            func(*http.Request) bool
	getCallerIdFn         func(*http.Request) int64
	getClientIdFn         func(*http.Request) int64
	authenticateRequestFn func(*http.Request) rest_errors.RestErr
}

func (m oauthServiceMock) IsPublic(rq *http.Request) bool {
	return m.isPublicFn(rq)
}
func (m oauthServiceMock) GetCallerId(rq *http.Request) int64 {
	return m.getCallerIdFn(rq)
}
func (m oauthServiceMock) GetClientId(rq *http.Request) int64 {
	return m.getClientIdFn(rq)
}
func (m oauthServiceMock) AuthenticateRequest(rq *http.Request) rest_errors.RestErr {
	return m.authenticateRequestFn(rq)
}

// tests

func TestCreateUserOk(t *testing.T) {
	userController := withMock(
		func(mock *usersServiceMock) {
			mock.createUserFn = func(u users.User) (*users.User, rest_errors.RestErr) {
				return &u, nil
			}
		}, nil, nil, nil)

	ctx, _ := withRequest(func(c *gin.Context) *http.Request {
		jsonUser, _ := json.Marshal(users.User{Id: 1, FirstName: "fname"})
		return httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(jsonUser)))
	})

	userController.Create(ctx)
	assertStatus(t, ctx, http.StatusCreated)
}

func TestCreateUserBadJson(t *testing.T) {
	userController := withMock(func(mock *usersServiceMock) {
		mock.createUserFn = func(u users.User) (*users.User, rest_errors.RestErr) {
			return &u, nil
		}
	}, nil, nil, nil)

	ctx, _ := withRequest(func(c *gin.Context) *http.Request {
		return httptest.NewRequest(http.MethodPost, "/", strings.NewReader("bad json"))
	})

	userController.Create(ctx)
	assertStatus(t, ctx, http.StatusBadRequest)
}

func TestCreateUserFailed(t *testing.T) {
	userController := withMock(func(mock *usersServiceMock) {
		mock.createUserFn = func(u users.User) (*users.User, rest_errors.RestErr) {
			return nil, rest_errors.NewInternalServerError("err", errors.New("db error"))
		}
	}, nil, nil, nil)

	ctx, _ := withRequest(func(c *gin.Context) *http.Request {
		jsonUser, _ := json.Marshal(users.User{Id: 1, FirstName: "fname"})
		return httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(jsonUser)))
	})

	userController.Create(ctx)
	assertStatus(t, ctx, http.StatusInternalServerError)
}

func TestGetUserOk(t *testing.T) {
	userController := withMock(
		func(mock *usersServiceMock) {
			mock.getUserFn = func(userId int64) (*users.User, rest_errors.RestErr) {
				return &users.User{Id: 1, FirstName: "fname"}, nil
			}
		},
		func(mock *oauthServiceMock) {
			mock.authenticateRequestFn = func(rq *http.Request) rest_errors.RestErr {
				return nil
			}
		},
		func(mock *oauthServiceMock) {
			mock.getCallerIdFn = func(r *http.Request) int64 {
				return 1
			}
		},
		func(mock *oauthServiceMock) {
			mock.isPublicFn = func(r *http.Request) bool {
				return true
			}
		})

	ctx, _ := withRequest(func(c *gin.Context) *http.Request {
		c.Params = gin.Params{gin.Param{Key: "user_id", Value: "1"}}
		return nil
	})

	userController.Get(ctx)
	assertStatus(t, ctx, http.StatusOK)
}

func TestGetUserNotAuthenticated(t *testing.T) {
	userController := withMock(
		nil,
		func(mock *oauthServiceMock) {
			mock.authenticateRequestFn = func(rq *http.Request) rest_errors.RestErr {
				return rest_errors.NewAuthorizationError("not authenticated")
			}
		}, nil, nil)

	ctx, _ := withRequest(func(c *gin.Context) *http.Request {
		c.Params = gin.Params{gin.Param{Key: "user_id", Value: "1"}}
		return nil
	})

	userController.Get(ctx)
	assertStatus(t, ctx, http.StatusUnauthorized)
}

func TestGetUserNotFound(t *testing.T) {
	userController := withMock(
		func(mock *usersServiceMock) {
			mock.getUserFn = func(userId int64) (*users.User, rest_errors.RestErr) {
				return nil, rest_errors.NewNotFoundError("not found")
			}
		},
		func(mock *oauthServiceMock) {
			mock.authenticateRequestFn = func(rq *http.Request) rest_errors.RestErr {
				return nil
			}
		}, nil, nil)

	ctx, _ := withRequest(func(c *gin.Context) *http.Request {
		c.Params = gin.Params{gin.Param{Key: "user_id", Value: "-1"}}
		return nil
	})

	userController.Get(ctx)
	assertStatus(t, ctx, http.StatusNotFound)
}

func TestUpdateUserOk(t *testing.T) {
	userController := withMock(func(mock *usersServiceMock) {
		mock.updateUserFn = func(isPublic bool, u users.User) (*users.User, rest_errors.RestErr) {
			u.FirstName = "updated"
			return &u, nil
		}
		mock.getUserFn = func(userId int64) (*users.User, rest_errors.RestErr) {
			return &users.User{Id: userId, FirstName: "change_me"}, nil
		}
	}, nil, nil, nil)

	ctx, response := withRequest(func(c *gin.Context) *http.Request {
		c.Params = gin.Params{gin.Param{Key: "user_id", Value: "1"}}
		bytes, _ := json.Marshal(users.User{Id: 1, FirstName: "update_me"})
		return httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(bytes)))
	})

	userController.Update(ctx)

	var u users.User
	if err := json.Unmarshal(response.Body.Bytes(), &u); err != nil {
		t.Error("invalid body response", response.Body)
	}

	fmt.Println("sdsds", response.Body)
	if u.FirstName != "updated" {
		t.Error("not expected", u.FirstName)
	}

	assertStatus(t, ctx, http.StatusOK+1)
}

func TestUpdateUserFailed(t *testing.T) {
	userController := withMock(func(mock *usersServiceMock) {
		mock.updateUserFn = func(isPublic bool, u users.User) (*users.User, rest_errors.RestErr) {
			return nil, rest_errors.NewInternalServerError("err", errors.New("db error"))
		}

		mock.getUserFn = func(userId int64) (*users.User, rest_errors.RestErr) {
			return &users.User{Id: userId, FirstName: "change_me"}, nil
		}
	}, nil, nil, nil)

	ctx, _ := withRequest(func(c *gin.Context) *http.Request {
		c.Params = gin.Params{gin.Param{Key: "user_id", Value: "1"}}
		bytes, _ := json.Marshal(users.User{Id: 1, FirstName: "update_me"})
		return httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(bytes)))
	})

	userController.Update(ctx)
	assertStatus(t, ctx, http.StatusInternalServerError)
}

func TestUpdateUserNotExist(t *testing.T) {
	userController := withMock(func(mock *usersServiceMock) {
		mock.updateUserFn = func(isPublic bool, u users.User) (*users.User, rest_errors.RestErr) {
			return nil, rest_errors.NewNotFoundError("err")
		}
		mock.getUserFn = func(userId int64) (*users.User, rest_errors.RestErr) {
			return nil, rest_errors.NewNotFoundError("err")
		}
	}, nil, nil, nil)

	ctx, _ := withRequest(func(c *gin.Context) *http.Request {
		c.Params = gin.Params{gin.Param{Key: "user_id", Value: "1"}}
		bytes, _ := json.Marshal(users.User{Id: 1, FirstName: "update_me"})
		return httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(string(bytes)))
	})

	userController.Update(ctx)
	assertStatus(t, ctx, http.StatusNotFound)
}

func TestDelUserOk(t *testing.T) {
	userController := withMock(func(mock *usersServiceMock) {
		mock.deleteUserFn = func(userId int64) rest_errors.RestErr {
			return nil
		}
	}, nil, nil, nil)

	ctx, _ := withRequest(func(c *gin.Context) *http.Request {
		c.Params = gin.Params{gin.Param{Key: "user_id", Value: "1"}}
		return nil
	})

	userController.Delete(ctx)
	assertStatus(t, ctx, http.StatusOK)
}

func TestDelUserNoIdParam(t *testing.T) {
	userController := withMock(func(mock *usersServiceMock) {
		mock.deleteUserFn = func(userId int64) rest_errors.RestErr {
			return nil
		}
	}, nil, nil, nil)

	ctx, _ := withRequest(func(c *gin.Context) *http.Request {
		return nil
	})

	userController.Delete(ctx)
	assertStatus(t, ctx, http.StatusBadRequest)
}

func TestDelUserServiceError(t *testing.T) {
	userController := withMock(func(mock *usersServiceMock) {
		mock.deleteUserFn = func(userId int64) rest_errors.RestErr {
			return rest_errors.NewInternalServerError("err", errors.New("db error"))
		}
	}, nil, nil, nil)

	ctx, _ := withRequest(func(c *gin.Context) *http.Request {
		c.Params = gin.Params{gin.Param{Key: "user_id", Value: "1"}}
		return nil
	})

	userController.Delete(ctx)
	assertStatus(t, ctx, http.StatusInternalServerError)
}

func TestSearchUserOk(t *testing.T) {
	userController := withMock(func(mock *usersServiceMock) {
		mock.searchUserFn = func(searchQuery string) (users.Users, rest_errors.RestErr) {
			return []users.User{}, nil
		}
	}, nil, nil, nil)

	ctx, _ := withRequest(func(c *gin.Context) *http.Request {
		return httptest.NewRequest(http.MethodPatch, "/status=active", nil)
	})

	userController.Search(ctx)
	assertStatus(t, ctx, http.StatusOK)
}

func TestLoginOk(t *testing.T) {
	lr := users.LoginRequest{Email: "email", Password: "pws"}

	userController := withMock(func(mock *usersServiceMock) {
		mock.loginUserFn = func(lr users.LoginRequest) (*users.User, rest_errors.RestErr) {
			return &users.User{Id: 1, Email: lr.Email}, nil
		}
	}, nil, nil, nil)

	ctx, _ := withRequest(func(c *gin.Context) *http.Request {
		bytes, _ := json.Marshal(lr)
		return httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(bytes)))
	})
	userController.Login(ctx)
	assertStatus(t, ctx, http.StatusOK)
}

// helpers

func withMock(fns ...interface{}) *UserController {

	userMock := new(usersServiceMock)
	oauthMock := new(oauthServiceMock)

	for _, fn := range fns {
		if fn != nil {
			switch fncall := fn.(type) {
			case func(*usersServiceMock):
				fncall(userMock)
			case func(*oauthServiceMock):
				fncall(oauthMock)
			default:
				panic("unknown mock type")
			}
		}
	}
	return ProvideUserController(userMock, oauthMock)
}

func withRequest(makeRequestFn func(*gin.Context) *http.Request) (*gin.Context, *httptest.ResponseRecorder) {
	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)
	ctx.Request = makeRequestFn(ctx)
	return ctx, response
}

func assertStatus(t *testing.T, c *gin.Context, expectedStatus int) {
	assert.EqualValues(t, c.Writer.Status(), expectedStatus)
}
