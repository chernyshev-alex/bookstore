package rest

import (
	"net/http"
	"os"
	"testing"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"xxx@gmail.com", "password":"xxxx}`,
		RespHTTPCode: -1,
		RespBody:     `{}`,
	})

	repository := usersRepository{}
	u, err := repository.LoginUser("xxx@gmail.com", "pwd")
	assert.Nil(t, u)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"xxx@gmail.com", "password":"xxxx}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message":"unknown response from user api (login)", "status":"404", "error":"not_found"}`,
	})

	repository := usersRepository{}
	u, err := repository.LoginUser("xxx@gmail.com", "pwd")
	assert.Nil(t, u)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)

}
func TestLoginUserInvalidLoginCredetials(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"xxx@gmail.com", "password":"xxxx}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message":"unknown response from user api (login)", "status":"404", "error":"not_found"}`,
	})

	repository := usersRepository{}
	u, err := repository.LoginUser("xxx@gmail.com", "pwd")
	assert.Nil(t, u)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
}
func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"xxx@gmail.com", "password":"xxxx}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id":1}`,
	})

	repository := usersRepository{}
	u, err := repository.LoginUser("xxx@gmail.com", "pwd")
	assert.Nil(t, u)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
}

func TestLoginNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"xxx@gmail.com", "password":"xxxx}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id":1}`,
	})

	repository := usersRepository{}
	u, err := repository.LoginUser("xxx@gmail.com", "pwd")
	assert.Nil(t, err)
	assert.NotNil(t, u)
	assert.EqualValues(t, u.Id, 1)
}
