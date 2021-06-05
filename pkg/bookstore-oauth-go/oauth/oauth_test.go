package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

//go:generate mockery  --name=OAuthInterface --output ../mocks
//go:generate mockery  --name=HTTPClientInterface --output ../mocks

type MockHttpClient struct {
	getFn func(url string) (resp *http.Response, err error)
}

func (m MockHttpClient) Get(url string) (resp *http.Response, err error) {
	return m.getFn(url)
}

type OAuthTestSuite struct {
	suite.Suite
	httpClient  *MockHttpClient
	oauthClient *OAuthClient
}

func TestOAuthTestSuite(t *testing.T) {
	suite.Run(t, new(OAuthTestSuite))
}

func (s *OAuthTestSuite) SetupTest() {
	s.httpClient = new(MockHttpClient)
	s.oauthClient = ProvideOAuthClient(s.httpClient, "")
}

func (s *OAuthTestSuite) TestOauthConstants() {
	assert.EqualValues(s.T(), "X-Public", headerXPublic)
	assert.EqualValues(s.T(), "X-Client-Id", headerXClientId)
	assert.EqualValues(s.T(), "X-Caller-Id", headerXCallerId)
	assert.EqualValues(s.T(), "access_token", paramAccessToken)
}

func (s *OAuthTestSuite) TestIsPublicNilRequest() {
	assert.True(s.T(), s.oauthClient.IsPublic(nil))
}

func (s *OAuthTestSuite) TestIsPublicNoErr() {
	rq := http.Request{Header: make(http.Header)}
	assert.False(s.T(), s.oauthClient.IsPublic(&rq))
	rq.Header.Add(headerXPublic, "true")
	assert.True(s.T(), s.oauthClient.IsPublic(&rq))
}

func (s *OAuthTestSuite) TestGetCallerIdNilRequest() {
	assert.Equal(s.T(), int64(0), s.oauthClient.GetCallerId(nil))
}

func (s *OAuthTestSuite) TestGetCallerId() {
	rq := http.Request{Header: make(http.Header)}

	rq.Header.Add(headerXCallerId, "not number")
	assert.Equal(s.T(), int64(0), s.oauthClient.GetCallerId(&rq))

	rq.Header.Del(headerXCallerId)
	rq.Header.Add(headerXCallerId, "1")
	assert.Equal(s.T(), int64(1), s.oauthClient.GetCallerId(&rq))
}

func (s *OAuthTestSuite) TestGetClientIdNilRequest() {
	assert.Equal(s.T(), int64(0), s.oauthClient.GetClientId(nil))
}
func (s *OAuthTestSuite) TestGetClientId() {
	rq := http.Request{Header: make(http.Header)}

	rq.Header.Add(headerXClientId, "not number")
	assert.Equal(s.T(), int64(0), s.oauthClient.GetClientId(&rq))

	rq.Header.Del(headerXClientId)
	rq.Header.Add(headerXClientId, "1")
	assert.Equal(s.T(), int64(1), s.oauthClient.GetClientId(&rq))
}
func (s *OAuthTestSuite) TestAuthenticateRequestNull() {
	assert.Nil(s.T(), s.oauthClient.AuthenticateRequest(nil))
}
func (s *OAuthTestSuite) TestAuthenticateRequestNoParamToken() {
	url, _ := url.Parse("http://localhost")
	rq := http.Request{
		Header: make(http.Header),
		URL:    url,
	}
	assert.Nil(s.T(), s.oauthClient.AuthenticateRequest(&rq))
}

func (s *OAuthTestSuite) TestAuthenticateRequestOk() {
	token := accessToken{Id: "AbC123", UserId: 1, ClientId: 100}
	url, _ := url.Parse(fmt.Sprintf("%s?%s=%s", "http://localhost", paramAccessToken, token.Id))
	rq := http.Request{
		Header: make(http.Header),
		URL:    url,
	}

	s.httpClient.getFn = func(url string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       tokenAsReadCloser(token),
		}, nil
	}

	s.oauthClient.AuthenticateRequest(&rq)

	clientId, _ := strconv.ParseInt(rq.Header.Get(headerXClientId), 10, 64)
	callerId, _ := strconv.ParseInt(rq.Header.Get(headerXCallerId), 10, 64)

	assert.Equal(s.T(), token.ClientId, clientId)
	assert.Equal(s.T(), token.UserId, callerId)
}

func (s *OAuthTestSuite) TestAuthenticateRequestFailedTokenNotExist() {
	token := accessToken{Id: "AbC123", UserId: 0, ClientId: 0}
	url, _ := url.Parse(fmt.Sprintf("%s?%s=%s", "http://localhost", paramAccessToken, token.Id))
	rq := http.Request{
		Header: make(http.Header),
		URL:    url,
	}
	s.httpClient.getFn = func(url string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       tokenAsReadCloser(token),
		}, nil
	}

	s.oauthClient.AuthenticateRequest(&rq)

	clientId, _ := strconv.ParseInt(rq.Header.Get(headerXClientId), 10, 64)
	callerId, _ := strconv.ParseInt(rq.Header.Get(headerXCallerId), 10, 64)

	assert.Equal(s.T(), int64(0), clientId)
	assert.Equal(s.T(), int64(0), callerId)
}

func (s *OAuthTestSuite) TestAuthenticateRequestFailedOAuthServiceNotAvailable() {
	token := accessToken{Id: "AbC123", UserId: 1, ClientId: 100}
	url, _ := url.Parse(fmt.Sprintf("%s?%s=%s", "http://localhost", paramAccessToken, token.Id))
	rq := http.Request{
		Header: make(http.Header),
		URL:    url,
	}

	s.httpClient.getFn = func(url string) (resp *http.Response, err error) {
		return nil, rest_errors.NewInternalServerError("failed request", errors.New("oauth failed"))
	}
	s.oauthClient.AuthenticateRequest(&rq)

	clientId, _ := strconv.ParseInt(rq.Header.Get(headerXClientId), 10, 64)
	callerId, _ := strconv.ParseInt(rq.Header.Get(headerXCallerId), 10, 64)

	assert.Equal(s.T(), int64(0), clientId)
	assert.Equal(s.T(), int64(0), callerId)
}

func (s *OAuthTestSuite) TestGetAccessTokenOk() {
	token := accessToken{Id: "AbC123", UserId: 1, ClientId: 100}

	s.httpClient.getFn = func(url string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       tokenAsReadCloser(token),
		}, nil
	}

	accessToken, err := s.oauthClient.GetAccessToken("AbC123")

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), accessToken)
}

func (s *OAuthTestSuite) TestGetAccessTokenServiceNotAvailable() {
	s.httpClient.getFn = func(url string) (resp *http.Response, err error) {
		return nil, errors.New("bad response")
	}

	accessToken, err := s.oauthClient.GetAccessToken("AbC123")

	assert.Nil(s.T(), accessToken)
	assert.NotNil(s.T(), err)
	assert.EqualValues(s.T(), http.StatusInternalServerError, err.Status())
}

func (s *OAuthTestSuite) TestGetAccessTokenBadBodyResponse() {
	s.httpClient.getFn = func(url string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("bad json accessToken")),
		}, nil
	}

	accessToken, err := s.oauthClient.GetAccessToken("AbC123")

	assert.Nil(s.T(), accessToken)
	assert.NotNil(s.T(), err)
	assert.EqualValues(s.T(), http.StatusInternalServerError, err.Status())
}

// helpers

func serialize(token accessToken) string {
	bytes, _ := json.Marshal(token)
	return string(bytes)
}

func tokenAsReadCloser(token accessToken) io.ReadCloser {
	return io.NopCloser(strings.NewReader(serialize(token)))
}
