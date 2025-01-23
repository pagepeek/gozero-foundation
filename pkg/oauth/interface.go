package oauth

import (
	"net/http"
)

type ClientCreator func() Client

type RequestOption func(*http.Request)

type Client interface {
	Auth(token AccessToken) (AuthClient, error)

	GetAuthUrl(redirectUrl string, scopes ...string) *OauthRedirectUrl

	CodeToUser(code, verifier, callbackUrl string) (*OauthUser, error)
}

type AuthClient interface {
	Get(path string, data map[string]string, options ...RequestOption) (*http.Response, error)

	Post(path string, data any, options ...RequestOption) (*http.Response, error)

	Put(path string, data any, options ...RequestOption) (*http.Response, error)

	Delete(path string, data any, options ...RequestOption) (*http.Response, error)

	Send(method, path string, data any, options ...RequestOption) (*http.Response, error)
}

type AccessToken interface {
	String() string

	Valid() bool

	Refresh() error

	GetUserID() string

	GetToken() string

	GetExpiresAt() int64

	GetRefreshToken() string

	GetRefreshExpiresAt() int64
}
