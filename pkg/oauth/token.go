package oauth

import (
	"net/http"
	"time"
)

type BaseAccessToken struct {
	UserID           string
	Token            string
	RefreshToken     string
	ExpiresAt        int64
	RefreshExpiresAt int64
	HttpClient       *http.Client
}

func (t BaseAccessToken) String() string {
	return t.Token
}

func (t BaseAccessToken) Valid() bool {
	now := time.Now().Unix()
	// token过期,且无法刷新token
	return t.ExpiresAt >= now || t.RefreshExpiresAt >= now
}

func (t BaseAccessToken) GetUserID() string {
	return t.UserID
}

func (t BaseAccessToken) GetToken() string {
	return t.Token
}

func (t BaseAccessToken) GetExpiresAt() int64 {
	return t.ExpiresAt
}

func (t BaseAccessToken) GetRefreshToken() string {
	return t.RefreshToken
}

func (t BaseAccessToken) GetRefreshExpiresAt() int64 {
	return t.RefreshExpiresAt
}
