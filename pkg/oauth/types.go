package oauth

import (
	"net/url"
)

type OauthUser struct {
	ID          string      `json:"openid"`
	Name        string      `json:"name"`
	Email       string      `json:"email"`
	Avatar      string      `json:"avatar"`
	Provider    string      `json:"provider"`
	AccessToken AccessToken `json:"-"`
}

type OauthRedirectUrl struct {
	Params   url.Values
	Target   string
	Provider string
	Verifier string
}

func (u OauthRedirectUrl) String() string {
	target, _ := url.Parse(u.Target)

	target.RawQuery = u.Params.Encode()

	return target.String()
}
