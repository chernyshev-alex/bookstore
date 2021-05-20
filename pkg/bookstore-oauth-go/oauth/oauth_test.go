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
)

type MockHttpClient struct {
	getFn func(url string) (resp *http.Response, err error)
}

var (
	oauthService = ProvideOAuthClient(HttpConfigClient{
		httpClient: new(MockHttpClient),
		baseURL:    "",
	})
)

func (m MockHttpClient) Get(url string) (resp *http.Response, err error) {
	return m.getFn(url)
}

func TestOauthConstants(t *testing.T) {
	assert.EqualValues(t, "X-Public", headerXPublic)
	assert.EqualValues(t, "X-Client-Id", headerXClientId)
	assert.EqualValues(t, "X-Caller-Id", headerXCallerId)
	assert.EqualValues(t, "access_token", paramAccessToken)
}

func TestIsPublicNilRequest(t *testing.T) {
	assert.True(t, oauthService.IsPublic(nil))
}

func TestIsPublicNoErr(t *testing.T) {
	rq := http.Request{Header: make(http.Header)}
	assert.False(t, oauthService.IsPublic(&rq))
	rq.Header.Add(headerXPublic, "true")
	assert.True(t, oauthService.IsPublic(&rq))
}

func TestGetCallerIdNilRequest(t *testing.T) {
	assert.Equal(t, int64(0), oauthService.GetCallerId(nil))
}

func TestGetCallerId(t *testing.T) {
	rq := http.Request{Header: make(http.Header)}

	rq.Header.Add(headerXCallerId, "not number")
	assert.Equal(t, int64(0), oauthService.GetCallerId(&rq))

	rq.Header.Del(headerXCallerId)
	rq.Header.Add(headerXCallerId, "1")
	assert.Equal(t, int64(1), oauthService.GetCallerId(&rq))
}

func TestGetClientIdNilRequest(t *testing.T) {
	assert.Equal(t, int64(0), oauthService.GetClientId(nil))
}
func TestGetClientId(t *testing.T) {
	rq := http.Request{Header: make(http.Header)}

	rq.Header.Add(headerXClientId, "not number")
	assert.Equal(t, int64(0), oauthService.GetClientId(&rq))

	rq.Header.Del(headerXClientId)
	rq.Header.Add(headerXClientId, "1")
	assert.Equal(t, int64(1), oauthService.GetClientId(&rq))
}
func TestAuthenticateRequestNull(t *testing.T) {
	assert.Nil(t, oauthService.AuthenticateRequest(nil))
}
func TestAuthenticateRequestNoParamToken(t *testing.T) {
	url, _ := url.Parse("http://localhost")
	rq := http.Request{
		Header: make(http.Header),
		URL:    url,
	}
	assert.Nil(t, oauthService.AuthenticateRequest(&rq))
}

func TestAuthenticateRequestOk(t *testing.T) {
	token := accessToken{Id: "AbC123", UserId: 1, ClientId: 100}
	url, _ := url.Parse(fmt.Sprintf("%s?%s=%s", "http://localhost", paramAccessToken, token.Id))
	rq := http.Request{
		Header: make(http.Header),
		URL:    url,
	}
	oauthService := withMock(func(mock *MockHttpClient) {
		mock.getFn = func(url string) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       tokenAsReadCloser(token),
			}, nil
		}
	})
	oauthService.AuthenticateRequest(&rq)

	clientId, _ := strconv.ParseInt(rq.Header.Get(headerXClientId), 10, 64)
	callerId, _ := strconv.ParseInt(rq.Header.Get(headerXCallerId), 10, 64)

	assert.Equal(t, token.ClientId, clientId)
	assert.Equal(t, token.UserId, callerId)
}

func TestAuthenticateRequestFailedTokenNotExist(t *testing.T) {
	token := accessToken{Id: "AbC123", UserId: 0, ClientId: 0}
	url, _ := url.Parse(fmt.Sprintf("%s?%s=%s", "http://localhost", paramAccessToken, token.Id))
	rq := http.Request{
		Header: make(http.Header),
		URL:    url,
	}
	oauthService := withMock(func(mock *MockHttpClient) {
		mock.getFn = func(url string) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       tokenAsReadCloser(token),
			}, nil
		}
	})
	oauthService.AuthenticateRequest(&rq)

	clientId, _ := strconv.ParseInt(rq.Header.Get(headerXClientId), 10, 64)
	callerId, _ := strconv.ParseInt(rq.Header.Get(headerXCallerId), 10, 64)

	assert.Equal(t, int64(0), clientId)
	assert.Equal(t, int64(0), callerId)
}

func TestAuthenticateRequestFailedOAuthServiceNotAvailable(t *testing.T) {
	token := accessToken{Id: "AbC123", UserId: 1, ClientId: 100}
	url, _ := url.Parse(fmt.Sprintf("%s?%s=%s", "http://localhost", paramAccessToken, token.Id))
	rq := http.Request{
		Header: make(http.Header),
		URL:    url,
	}
	oauthService := withMock(func(mock *MockHttpClient) {
		mock.getFn = func(url string) (resp *http.Response, err error) {
			return nil, rest_errors.NewInternalServerError("failed request", errors.New("oauth failed"))
		}
	})
	oauthService.AuthenticateRequest(&rq)

	clientId, _ := strconv.ParseInt(rq.Header.Get(headerXClientId), 10, 64)
	callerId, _ := strconv.ParseInt(rq.Header.Get(headerXCallerId), 10, 64)

	assert.Equal(t, int64(0), clientId)
	assert.Equal(t, int64(0), callerId)
}

func TestGetAccessTokenOk(t *testing.T) {
	token := accessToken{Id: "AbC123", UserId: 1, ClientId: 100}

	oauthService := withMock(func(mock *MockHttpClient) {
		mock.getFn = func(url string) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       tokenAsReadCloser(token),
			}, nil
		}
	})

	accessToken, err := oauthService.GetAccessToken("AbC123")

	assert.Nil(t, err)
	assert.NotNil(t, accessToken)
}

func TestGetAccessTokenServiceNotAvailable(t *testing.T) {
	oauthService := withMock(func(mock *MockHttpClient) {
		mock.getFn = func(url string) (resp *http.Response, err error) {
			return nil, errors.New("bad response")
		}
	})

	accessToken, err := oauthService.GetAccessToken("AbC123")

	assert.Nil(t, accessToken)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestGetAccessTokenBadBodyResponse(t *testing.T) {
	oauthService := withMock(func(mock *MockHttpClient) {
		mock.getFn = func(url string) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("bad json accessToken")),
			}, nil
		}
	})

	accessToken, err := oauthService.GetAccessToken("AbC123")

	assert.Nil(t, accessToken)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

// helpers

func withMock(configFn func(*MockHttpClient)) *OAuthClient {
	httpClient := new(MockHttpClient)
	configFn(httpClient)
	return ProvideOAuthClient(HttpConfigClient{
		httpClient: httpClient,
		baseURL:    ""})
}

func serialize(token accessToken) string {
	bytes, _ := json.Marshal(token)
	return string(bytes)
}

func tokenAsReadCloser(token accessToken) io.ReadCloser {
	return io.NopCloser(strings.NewReader(serialize(token)))
}
