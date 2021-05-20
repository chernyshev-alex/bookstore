package services

import (
	"strings"

	"github.com/chernyshev-alex/bookstore_oauth-api/domain/access_token"
	"github.com/chernyshev-alex/bookstore_oauth-api/repository/db"
	"github.com/chernyshev-alex/bookstore_oauth-api/repository/rest"
	"github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
)

type Service interface {
	GetById(string) (*access_token.AccessToken, rest_errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, rest_errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken) rest_errors.RestErr
}

type service struct {
	dbRepo        db.DbRepository
	restUsersRepo rest.RestUsersRepository
}

func NewService(usersRepo rest.RestUsersRepository, dbRepo db.DbRepository) Service {
	return &service{
		restUsersRepo: usersRepo,
		dbRepo:        dbRepo,
	}
}

func (s *service) GetById(tokenId string) (*access_token.AccessToken, rest_errors.RestErr) {
	tokenId = strings.TrimSpace(tokenId)
	if len(tokenId) == 0 {
		return nil, rest_errors.NewBadRequestError("empty token")
	}
	accessToken, err := s.dbRepo.GetById(tokenId)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *service) Create(rq access_token.AccessTokenRequest) (*access_token.AccessToken, rest_errors.RestErr) {
	if err := rq.Validate(); err != nil {
		return nil, err
	}
	user, err := s.restUsersRepo.LoginUser(rq.Username, rq.Password)
	if err != nil {
		return nil, err
	}

	at := access_token.GetNewAccessToken(user.Id)
	at.Generate()

	if err := s.dbRepo.Create(at); err != nil {
		return nil, err
	}

	return &at, nil
}

func (s *service) UpdateExpirationTime(at access_token.AccessToken) rest_errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.dbRepo.UpdateExpirationTime(at)
}
