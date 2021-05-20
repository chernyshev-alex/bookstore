package access_token

import (
	"fmt"
	"strings"
	"time"

	"github.com/chernyshev-alex/bookstore_oauth-api/utils/crypto_utils"
	"github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
)

const (
	expirationTime             = 24
	grantTypePassword          = "password"
	grantTypeClientCredentials = "client_credentials"
)

type AccessTokenRequest struct {
	GrantType string `json:"access_token"`
	Scope     string `json:"scope"`

	// user for password grant type
	Username string `json:"username"`
	Password string `json:"password"`

	// used for client_credentials grant type
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (at *AccessTokenRequest) Validate() *rest_errors.RestErr {
	switch at.GrantType {
	case grantTypePassword:
		break
	case grantTypeClientCredentials:
		break
	default:
		return rest_errors.NewBadRequestError("bad GrantType parameter")
	}
	return nil
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

func GetNewAccessToken(userId int64) AccessToken {
	return AccessToken{
		UserId:  userId,
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at *AccessToken) Validate() *rest_errors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.UserId <= 0 {
		return rest_errors.NewBadRequestError("bad user id")
	}

	if at.ClientId <= 0 {
		return rest_errors.NewBadRequestError("bad client id")
	}

	if at.Expires <= 0 {
		return rest_errors.NewBadRequestError("bad expires")
	}
	return nil
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}

func (at *AccessToken) Generate() {
	at.AccessToken = crypto_utils.GetMD5(fmt.Sprintf("at-%d-%d-ran", at.UserId, at.Expires))
}
