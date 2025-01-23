package twitter

import (
	"errors"
	"net/url"
	"strings"

	"github.com/pagepeek/gozero-foundation/pkg/oauth"
	"github.com/pagepeek/gozero-foundation/pkg/utils"
)

const (
	TwitterV2Host         = "api.x.com"
	TwitterOauthUrl       = "https://twitter.com/i/oauth2/authorize"
	TwitterOauth2TokenUrl = "https://api.x.com/2/oauth2/token"
)

func (x *Twitter) GetAuthUrl(redirectUrl string, scopes ...string) *oauth.OauthRedirectUrl {
	verifier := utils.RandStr(43)

	return &oauth.OauthRedirectUrl{
		Target:   TwitterOauthUrl,
		Provider: "twitter",
		Verifier: verifier,
		Params: url.Values{
			"response_type":         {"code"},
			"client_id":             {x.ClientID},
			"redirect_uri":          {redirectUrl},
			"scope":                 {strings.Join(scopes, " ")},
			"state":                 {"state"},
			"code_challenge":        {verifier},
			"code_challenge_method": {"plain"},
		},
	}
}

func (x *Twitter) CodeToUser(code, verifier, callbackUrl string) (*oauth.OauthUser, error) {
	token, err := x.GetAccessToken(code, verifier, callbackUrl)

	if err != nil {
		return nil, err
	}

	auth, err := x.Auth(token)
	if err != nil {
		return nil, err
	}

	resp, err := auth.Get("/2/users/me", nil)
	if err != nil {
		return nil, err
	}

	user, err := utils.JsonDecode[UserData](resp.Body)
	if err != nil {
		return nil, err
	}

	if user.Error != "" {
		return nil, errors.New(user.Error)
	}

	oauthUser := &oauth.OauthUser{
		ID:          user.Data.ID,
		Name:        user.Data.Name,
		Provider:    "twitter",
		AccessToken: token,
	}

	return oauthUser, nil
}
