package models

import (
	"encoding/json"
	"strings"

	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/rest_errors"
)

const (
	STATUS_ACTIVE = "active"
)

type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:"password"`
}

type Users []User

func (user *User) Validate() rest_errors.RestErr {

	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Email = strings.Trim(strings.ToLower(user.Email), "")

	if user.Email == "" {
		return rest_errors.NewBadRequestError("invalid email")
	}

	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" {
		return rest_errors.NewBadRequestError("empty password")
	}

	return nil
}

type PublicUser struct {
	Id          int64  `json:"id"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
}

type PrivateUser struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
}

func (u *User) Marshall(isPublic bool) interface{} {
	if isPublic {
		return PublicUser{
			Id:          u.Id,
			DateCreated: u.DateCreated,
			Status:      u.Status,
		}
	}
	userJson, _ := json.Marshal(u)
	var privUser PrivateUser
	json.Unmarshal(userJson, &privUser)
	return privUser
}

func (users Users) Marshall(isPublic bool) []interface{} {
	result := make([]interface{}, len(users))
	for i, u := range users {
		result[i] = u.Marshall(isPublic)
	}
	return result
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
