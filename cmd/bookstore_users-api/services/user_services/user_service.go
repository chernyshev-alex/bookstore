package user_services

import (
	"database/sql"
	"time"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/dao/mysql/gen"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/dao/user_dao"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/models"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/date_utils"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/rest_errors"
)

type usersService struct {
	UserService
	userDao user_dao.UserDao
}

func NewService(userDao user_dao.UserDao) UserService {
	return &usersService{
		userDao: userDao,
	}
}

func (s *usersService) GetUser(userId int64) (*models.User, rest_errors.RestErr) {
	result, err := s.userDao.Get(userId)
	if err != nil {
		return nil, err
	}
	return &models.User{
		Id:          int64(result.ID),
		FirstName:   result.FirstName.String,
		LastName:    result.LastName.String,
		Email:       result.Email,
		DateCreated: date_utils.Time2String(result.DateCreated),
		Status:      result.Status.String,
	}, nil
}

func (s *usersService) CreateUser(u models.User) (*models.User, rest_errors.RestErr) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	insertUser := gen.InsertUserParams{
		FirstName:   nillableStr(u.FirstName),
		LastName:    nillableStr(u.LastName),
		Email:       u.Email,
		DateCreated: time.Now(),
		Status:      nillableStr(u.Status),
		Password:    nillableStr(u.Password),
	}

	userId, err := s.userDao.Save(insertUser)
	if err != nil {
		return nil, err
	}

	u.Id = userId
	return &u, nil
}

func (s *usersService) UpdateUser(isPartial bool, u models.User) (*models.User, rest_errors.RestErr) {
	cu, err := s.GetUser(u.Id)
	if err != nil {
		return nil, err
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	if isPartial {
		if u.FirstName != "" {
			cu.FirstName = u.FirstName
		}
		if u.LastName != "" {
			cu.LastName = u.LastName
		}
		if u.Email != "" {
			cu.Email = u.Email
		}

	} else {
		cu.FirstName = u.FirstName
		cu.LastName = u.LastName
		cu.Email = u.Email
	}

	updateUser := gen.UpdateUserParams{
		FirstName: nillableStr(cu.FirstName),
		LastName:  nillableStr(cu.FirstName),
		Email:     u.Email,
		ID:        int32(u.Id)}

	if err := s.userDao.Update(updateUser); err != nil {
		return nil, err
	}
	return cu, nil
}

func (s *usersService) DeleteUser(userId int64) rest_errors.RestErr {
	return s.userDao.Delete(userId)
}

func (s *usersService) SearchUsersByStatus(status string) ([]models.User, rest_errors.RestErr) {
	result, err := s.userDao.FindByStatus(status)
	if err != nil {
		return nil, err
	}
	ls := make([]models.User, 0, len(result))
	for _, rec := range result {
		u := models.User{
			Id:        int64(rec.ID),
			FirstName: rec.FirstName.String,
			LastName:  rec.LastName.String,
			Email:     rec.Email,
			//	DateCreated: date_utils.rec.DateCreated,  // TODO
			Status: rec.Status.String,
		}
		ls = append(ls, u)
	}
	return ls, nil
}

func (s *usersService) LoginUser(rq models.LoginRequest) (*models.User, rest_errors.RestErr) {
	input := gen.FindByEMailAndPswParams{
		Email:    rq.Email,
		Password: nillableStr(rq.Password),
		Status:   nillableStr(models.STATUS_ACTIVE),
	}

	result, err := s.userDao.FindByEmailAndPsw(input)
	if err != nil {
		return nil, err
	}

	u := models.User{
		Id:          int64(result.ID),
		FirstName:   result.FirstName.String,
		LastName:    result.LastName.String,
		Email:       result.Email,
		DateCreated: date_utils.Time2String(result.DateCreated),
		Status:      result.Status.String,
	}
	return &u, nil
}

func nillableStr(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}
