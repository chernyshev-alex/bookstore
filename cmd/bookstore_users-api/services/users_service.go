package services

import (
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/domain/users"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/crypto_utils"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/date_utils"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/rest_errors"
)

type UsersServiceInterface interface {
	GetUser(int64) (*users.User, rest_errors.RestErr)
	CreateUser(users.User) (*users.User, rest_errors.RestErr)
	UpdateUser(bool, users.User) (*users.User, rest_errors.RestErr)
	DeleteUser(int64) rest_errors.RestErr
	SearchUser(string) (users.Users, rest_errors.RestErr)
	LoginUser(users.LoginRequest) (*users.User, rest_errors.RestErr)
}

type UsersService interface {
	UsersServiceInterface
}
type usersService struct {
	usersDAO users.UsersDAOInterface
}

func NewService(userDaoInterface users.UsersDAOInterface) UsersService {
	return &usersService{
		usersDAO: userDaoInterface,
	}
}

// deprecated
func ProvideUserService(userDaoInterface users.UsersDAOInterface) UsersServiceInterface {
	return &usersService{
		usersDAO: userDaoInterface,
	}
}

func (s *usersService) GetUser(userId int64) (*users.User, rest_errors.RestErr) {
	u, err := s.usersDAO.Get(userId)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *usersService) CreateUser(user users.User) (*users.User, rest_errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	user.Status = users.StatusActive
	user.DateCreated = date_utils.GetNowDbFormat()
	user.Password = crypto_utils.GetMD5(user.Password)
	if err := s.usersDAO.Save(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *usersService) UpdateUser(isPartial bool, u users.User) (*users.User, rest_errors.RestErr) {
	currentUser, err := s.GetUser(u.Id)
	if err != nil {
		return nil, err
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	if isPartial {
		if u.FirstName != "" {
			currentUser.FirstName = u.FirstName
		}
		if u.LastName != "" {
			currentUser.LastName = u.LastName
		}
		if u.Email != "" {
			currentUser.Email = u.Email
		}

	} else {
		currentUser.FirstName = u.FirstName
		currentUser.LastName = u.LastName
		currentUser.Email = u.Email
	}

	if err := s.usersDAO.Update(currentUser); err != nil {
		return nil, err
	}
	return currentUser, nil
}

func (s *usersService) DeleteUser(userId int64) rest_errors.RestErr {
	return s.usersDAO.Delete(&users.User{Id: userId})
}

func (s *usersService) SearchUser(status string) (users.Users, rest_errors.RestErr) {
	return s.usersDAO.FindByStatus(status)
}

func (s *usersService) LoginUser(req users.LoginRequest) (*users.User, rest_errors.RestErr) {
	u := &users.User{
		Email:    req.Email,
		Password: crypto_utils.GetMD5(req.Password),
	}
	if err := s.usersDAO.FindByEmailAndPsw(u); err != nil {
		return nil, err
	}
	return u, nil
}
