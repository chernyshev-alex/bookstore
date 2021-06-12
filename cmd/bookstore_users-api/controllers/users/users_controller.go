package users

import (
	"net/http"
	"strconv"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/domain/users"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/services"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore-oauth-go/oauth"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/rest_errors"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	usersService services.UsersServiceInterface
	oauthService oauth.OAuthInterface
}

func ProvideUserController(userServiceInterface services.UsersServiceInterface,
	oauthService oauth.OAuthInterface) *UserController {
	return &UserController{
		usersService: userServiceInterface,
		oauthService: oauthService,
	}
}

func getUserId(userIdParam string) (int64, rest_errors.RestErr) {
	userId, parseErr := strconv.ParseInt(userIdParam, 10, 64)
	if parseErr != nil {
		return 0, rest_errors.NewBadRequestError("parse error user id")
	}
	return userId, nil
}

func (uc UserController) Create(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}
	result, err := uc.usersService.CreateUser(user)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func (uc UserController) Get(c *gin.Context) {
	if err := uc.oauthService.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.Status(), err)
		return
	}

	userId, err := getUserId(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	result, getErr := uc.usersService.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	if uc.oauthService.GetCallerId(c.Request) == userId {
		c.JSON(http.StatusOK, result.Marshall(false))
	}
	c.JSON(http.StatusOK, result.Marshall(uc.oauthService.IsPublic(c.Request)))
}

func (uc UserController) Update(c *gin.Context) {
	userId, err := getUserId(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	var u users.User
	if err := c.ShouldBindJSON(&u); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	u.Id = userId

	isPartial := c.Request.Method == http.MethodPatch

	result, err := uc.usersService.UpdateUser(isPartial, u)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func (uc UserController) Delete(c *gin.Context) {
	userId, err := getUserId(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	if err := uc.usersService.DeleteUser(userId); err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (uc UserController) Search(c *gin.Context) {

	status := c.Query("status")

	users, err := uc.usersService.SearchUser(status)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, users.Marshall(c.GetHeader("X-Public") == "true"))
}

func (uc UserController) Login(c *gin.Context) {

	var req users.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	u, err := uc.usersService.LoginUser(req)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, u.Marshall(c.GetHeader("X-Public") == "true"))
}
