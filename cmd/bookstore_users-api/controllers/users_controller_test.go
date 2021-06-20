package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	mock_srv "github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/mocks"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/models"
	mocks_oauth "github.com/chernyshev-alex/bookstore/pkg/bookstore-oauth-go/mocks"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/rest_errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UCServiceSuite struct {
	suite.Suite
	mockedUserService  *mock_srv.UserService
	mockedOAuthService *mocks_oauth.OAuthInterface
	userController     *UserController // TODO intf
	ctx                *gin.Context
	response           *httptest.ResponseRecorder
}

func TestUCServiceSuite(t *testing.T) {
	suite.Run(t, new(UCServiceSuite))
}

func (s *UCServiceSuite) SetupTest() {
	s.mockedUserService = new(mock_srv.UserService)
	s.mockedOAuthService = new(mocks_oauth.OAuthInterface)
	s.userController = ProvideUserController(s.mockedUserService, s.mockedOAuthService)

	s.response = httptest.NewRecorder()
	s.ctx, _ = gin.CreateTestContext(s.response)
}

// tests
//go:generate mockery  --name=UsersService --dir=../../services  --output ../../mocks

func (s *UCServiceSuite) TestCreateUserOk() {
	u := models.User{Id: 1, FirstName: "fname"}

	s.requestWithUserAndParams(http.MethodPost, &u, nil)

	s.mockedUserService.On("CreateUser", mock.IsType(u)).Return(&u, nil)
	s.userController.Create(s.ctx)

	s.mockedUserService.AssertExpectations(s.T())
	assert.EqualValues(s.T(), http.StatusCreated, s.ctx.Writer.Status())
}

func (s *UCServiceSuite) TestCreateUserBadJson() {
	s.requestWithJson(http.MethodPost, "bad json")

	s.userController.Create(s.ctx)

	s.mockedUserService.AssertExpectations(s.T())
	assert.EqualValues(s.T(), http.StatusBadRequest, s.ctx.Writer.Status())
}

func (s *UCServiceSuite) TestCreateUserServiceFailed() {
	u := models.User{Id: 1, FirstName: "fname"}

	s.requestWithUserAndParams(http.MethodPost, &u, nil)

	s.mockedUserService.On("CreateUser", mock.IsType(u)).Return(nil,
		rest_errors.NewInternalServerError("err", errors.New("db error")))

	s.userController.Create(s.ctx)

	s.mockedUserService.AssertExpectations(s.T())
	assert.EqualValues(s.T(), http.StatusInternalServerError, s.ctx.Writer.Status())
}

func (s *UCServiceSuite) TestGetUserOk() {
	userId := int64(1)

	params := gin.Param{Key: "user_id", Value: strconv.FormatInt(userId, 10)}
	s.requestWithUserAndParams(http.MethodGet, nil, gin.Params{params})

	u := models.User{Id: userId}

	rq := s.ctx.Request
	s.mockedUserService.On("GetUser", userId).Return(&u, nil)
	s.mockedOAuthService.On("AuthenticateRequest", mock.IsType(rq)).Return(nil)
	s.mockedOAuthService.On("GetCallerId", mock.IsType(rq)).Return(userId)
	s.mockedOAuthService.On("IsPublic", mock.IsType(rq)).Return(true)

	s.userController.Get(s.ctx)
	s.mockedOAuthService.AssertExpectations(s.T())
	s.mockedUserService.AssertExpectations(s.T())
	assert.EqualValues(s.T(), http.StatusOK, s.ctx.Writer.Status())
}

func (s *UCServiceSuite) TestGetUserNotAuthenticated() {
	userId := int64(1)

	params := gin.Param{Key: "user_id", Value: strconv.FormatInt(userId, 10)}
	s.requestWithUserAndParams(http.MethodGet, nil, gin.Params{params})

	rq := s.ctx.Request
	s.mockedOAuthService.On("AuthenticateRequest", mock.IsType(rq)).Return(rest_errors.NewAuthorizationError("not authenticated"))

	s.userController.Get(s.ctx)
	s.mockedOAuthService.AssertExpectations(s.T())
	assert.EqualValues(s.T(), http.StatusUnauthorized, s.ctx.Writer.Status())
}

func (s *UCServiceSuite) TestGetUserNotFound() {
	userId := int64(-1)

	params := gin.Param{Key: "user_id", Value: strconv.FormatInt(userId, 10)}
	s.requestWithUserAndParams(http.MethodGet, nil, gin.Params{params})

	rq := s.ctx.Request
	s.mockedUserService.On("GetUser", userId).Return(nil, rest_errors.NewNotFoundError("not found"))
	s.mockedOAuthService.On("AuthenticateRequest", mock.IsType(rq)).Return(nil)

	s.userController.Get(s.ctx)

	s.mockedOAuthService.AssertExpectations(s.T())
	s.mockedUserService.AssertExpectations(s.T())
	assert.EqualValues(s.T(), http.StatusNotFound, s.ctx.Writer.Status())
}

func (s *UCServiceSuite) TestUpdateUserOk() {
	u := models.User{Id: 1, FirstName: "update_me"}
	p := gin.Params{gin.Param{Key: "user_id", Value: strconv.FormatInt(u.Id, 10)}}

	s.requestWithUserAndParams(http.MethodPatch, &u, p)

	s.mockedUserService.On("UpdateUser", true, u).Return(func(bool, models.User) *models.User {
		return &models.User{Id: 1, FirstName: "changed"}
	}, nil)

	s.userController.Update(s.ctx)

	if err := json.Unmarshal(s.response.Body.Bytes(), &u); err != nil {
		s.T().Error("bad body response", s.response.Body)
	}

	s.mockedUserService.AssertExpectations(s.T())

	assert.EqualValues(s.T(), "changed", u.FirstName)
	assert.EqualValues(s.T(), http.StatusOK, s.ctx.Writer.Status())
}

func (s *UCServiceSuite) TestUpdateUserFailed() {
	u := models.User{Id: 1, FirstName: "update_me"}
	p := gin.Params{gin.Param{Key: "user_id", Value: strconv.FormatInt(u.Id, 10)}}

	s.requestWithUserAndParams(http.MethodPatch, &u, p)

	s.mockedUserService.On("UpdateUser", true, u).Return(nil,
		rest_errors.NewInternalServerError("err", errors.New("db error")))

	s.userController.Update(s.ctx)

	s.mockedUserService.AssertExpectations(s.T())
	assert.EqualValues(s.T(), http.StatusInternalServerError, s.ctx.Writer.Status())
}

func (s *UCServiceSuite) TestUpdateUserNotExist() {
	u := models.User{Id: -1}
	p := gin.Params{gin.Param{Key: "user_id", Value: strconv.FormatInt(u.Id, 10)}}

	s.requestWithUserAndParams(http.MethodPatch, &u, p)

	s.mockedUserService.On("UpdateUser", true, u).Return(nil, rest_errors.NewNotFoundError("err"))

	s.userController.Update(s.ctx)

	s.mockedUserService.AssertExpectations(s.T())
	assert.EqualValues(s.T(), http.StatusNotFound, s.ctx.Writer.Status())
}

func (s *UCServiceSuite) TestRemoveUserOk() {
	userId := int64(-1)

	params := gin.Param{Key: "user_id", Value: strconv.FormatInt(userId, 10)}
	s.requestWithUserAndParams(http.MethodGet, nil, gin.Params{params})

	s.mockedUserService.On("DeleteUser", userId).Return(nil)

	s.userController.Delete(s.ctx)
	s.mockedUserService.AssertExpectations(s.T())
	assert.EqualValues(s.T(), http.StatusOK, s.ctx.Writer.Status())
}

func (s *UCServiceSuite) TestRemoveUserNoIdParam() {
	s.userController.Delete(s.ctx)
	s.mockedUserService.AssertExpectations(s.T())
	assert.EqualValues(s.T(), http.StatusBadRequest, s.ctx.Writer.Status())
}

func (s *UCServiceSuite) TestRemoveUserServiceError() {
	userId := int64(-1)

	params := gin.Param{Key: "user_id", Value: strconv.FormatInt(userId, 10)}
	s.requestWithUserAndParams(http.MethodGet, nil, gin.Params{params})

	s.mockedUserService.On("DeleteUser", userId).Return(rest_errors.NewInternalServerError("err", errors.New("db error")))

	s.userController.Delete(s.ctx)
	s.mockedUserService.AssertExpectations(s.T())
	assert.EqualValues(s.T(), http.StatusInternalServerError, s.ctx.Writer.Status())
}

func (s *UCServiceSuite) TestSearchUserOk() {
	s.requestWithQuery(http.MethodGet, "/?status=active")

	result := []models.User{{Id: 1, FirstName: "fname"}}
	s.mockedUserService.On("SearchUsersByStatus", mock.AnythingOfType("string")).Return(result, nil)

	s.userController.Search(s.ctx)
	s.mockedUserService.AssertExpectations(s.T())
	assert.EqualValues(s.T(), http.StatusOK, s.ctx.Writer.Status())
}

func (s *UCServiceSuite) TestLoginOk() {
	lr := models.LoginRequest{Email: "email", Password: "pws"}
	s.mockedUserService.On("LoginUser", lr).Return(&models.User{}, nil)

	s.requestWithLogin(http.MethodPost, lr)
	s.userController.Login(s.ctx)

	s.mockedUserService.AssertExpectations(s.T())
	assert.EqualValues(s.T(), http.StatusOK, s.ctx.Writer.Status())
}

func (s *UCServiceSuite) TestLoginFailed() {
	lr := models.LoginRequest{Email: "email", Password: "pws"}
	s.mockedUserService.On("LoginUser", lr).Return(nil, rest_errors.NewAuthorizationError("not authorized"))

	s.requestWithLogin(http.MethodPost, lr)
	s.userController.Login(s.ctx)

	s.mockedUserService.AssertExpectations(s.T())
	assert.EqualValues(s.T(), http.StatusUnauthorized, s.ctx.Writer.Status())
}

// helpers

func (s *UCServiceSuite) requestWithUserAndParams(httpMethod string, u *models.User, params gin.Params) {
	if params != nil {
		s.ctx.Params = params
	}
	if u != nil {
		jsonUser, _ := json.Marshal(u)
		s.ctx.Request = httptest.NewRequest(httpMethod, "/", strings.NewReader(string(jsonUser)))
	} else {
		s.ctx.Request = httptest.NewRequest(httpMethod, "/", nil)
	}
}

func (s *UCServiceSuite) requestWithLogin(httpMethod string, r models.LoginRequest) {
	js, _ := json.Marshal(r)
	s.ctx.Request = httptest.NewRequest(httpMethod, "/", strings.NewReader(string(js)))
}

func (s *UCServiceSuite) requestWithJson(httpMethod string, rawJson string) {
	s.ctx.Request = httptest.NewRequest(httpMethod, "/", strings.NewReader(rawJson))
}

func (s *UCServiceSuite) requestWithQuery(httpMethod string, query string) {
	s.ctx.Request = httptest.NewRequest(httpMethod, query, nil)
}
