package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/rest_errors"
)

const (
	headerXPublic   = "X-Public"
	headerXClientId = "X-Client-Id"
	headerXCallerId = "X-Caller-Id"

	paramAccessToken = "access_token"
)

type accessToken struct {
	Id       string `json:"id"`
	UserId   int64  `json:"user_id"`
	ClientId int64  `json:"client_id"`
}
type OAuthInterface interface {
	IsPublic(*http.Request) bool
	GetCallerId(*http.Request) int64
	GetClientId(*http.Request) int64
	AuthenticateRequest(*http.Request) rest_errors.RestErr
}
type HttpClientInterface interface {
	Get(url string) (resp *http.Response, err error)
}
type OAuthClient struct {
	baseURL    string
	httpClient HttpClientInterface
}

func ProvideOAuthClient(httpClient HttpClientInterface, baseURL string) *OAuthClient {
	return &OAuthClient{httpClient: httpClient, baseURL: baseURL}
}

func (oa OAuthClient) IsPublic(req *http.Request) bool {
	if req == nil {
		return true
	}
	return req.Header.Get(headerXPublic) == "true"
}

func (oa OAuthClient) GetCallerId(req *http.Request) int64 {
	if req == nil {
		return 0
	}
	callerId, err := strconv.ParseInt(req.Header.Get(headerXCallerId), 10, 64)
	if err != nil {
		return 0
	}
	return callerId
}

func (oa OAuthClient) GetClientId(req *http.Request) int64 {
	if req == nil {
		return 0
	}
	clientId, err := strconv.ParseInt(req.Header.Get(headerXClientId), 10, 64)
	if err != nil {
		return 0
	}
	return clientId
}

func (oa OAuthClient) AuthenticateRequest(req *http.Request) rest_errors.RestErr {
	if req == nil {
		return nil
	}

	accessTokenId := strings.TrimSpace(req.URL.Query().Get(paramAccessToken))
	if accessTokenId == "" {
		return nil
	}

	oa.cleanRequest(req)

	at, err := oa.GetAccessToken(accessTokenId)
	if err != nil {
		if err.Status() == http.StatusNotFound {
			return nil
		}
		return nil
	}

	req.Header.Add(headerXClientId, fmt.Sprintf("%v", at.ClientId))
	req.Header.Add(headerXCallerId, fmt.Sprintf("%v", at.UserId))
	return nil
}

func (oa OAuthClient) cleanRequest(req *http.Request) {
	if req == nil {
		return
	}
	req.Header.Del(headerXClientId)
	req.Header.Del(headerXCallerId)
}

func (oa OAuthClient) GetAccessToken(accessTokenId string) (*accessToken, rest_errors.RestErr) {
	resp, err := oa.httpClient.Get(makeURL(oa.baseURL, "oauth/access_token", accessTokenId))
	if err != nil {
		return nil, rest_errors.NewInternalServerError("timeout error", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			var restError rest_errors.RestErr
			if err = json.Unmarshal(body, &restError); err != nil {
				return nil, rest_errors.NewInternalServerError("error unmarshal failed", err)
			}
			return nil, restError
		}
	}

	var token accessToken
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, rest_errors.NewInternalServerError("response failed", err)
	}
	if err = json.Unmarshal(body, &token); err != nil {
		return nil, rest_errors.NewInternalServerError("token unmarshal failed", err)
	}
	return &token, nil
}

func makeURL(base_url string, parts ...string) string {
	var sb strings.Builder
	sb.Grow(255)
	sb.WriteString(base_url)
	for _, part := range parts {
		sb.WriteString("/")
		sb.WriteString(part)
	}
	return sb.String()
}
