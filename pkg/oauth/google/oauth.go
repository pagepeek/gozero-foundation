package google

import (
	"errors"
	"io"
	"net/url"
	"strings"

	"github.com/pagepeek/gozero-foundation/pkg/oauth"
	"github.com/pagepeek/gozero-foundation/pkg/utils"
)

const (
	GoogleHost     = "www.googleapis.com"
	GoogleOauthUrl = "https://accounts.google.com/o/oauth2/v2/auth"
	GoogleTokenUrl = "https://oauth2.googleapis.com/token"
)

func (g *Google) GetAuthUrl(redirectUrl string, scopes ...string) *oauth.OauthRedirectUrl {
	verifier := utils.RandStr(43)

	return &oauth.OauthRedirectUrl{
		Provider: "google",
		Verifier: verifier,
		Target:   GoogleOauthUrl,
		Params: url.Values{
			"response_type": {"code"},
			"client_id":     {g.ClientID},
			"redirect_uri":  {redirectUrl},
			"scope":         {strings.Join(scopes, " ")},
			"state":         {"state"},
			// "code_challenge":        {verifier},
			// "code_challenge_method": {"plain"},
		},
	}
}

func (g *Google) CodeToUser(code, verifier, callbackUrl string) (*oauth.OauthUser, error) {
	token, err := g.GetAccessToken(code, verifier, callbackUrl)

	if err != nil {
		return nil, err
	}

	auth, err := g.Auth(token)
	if err != nil {
		return nil, err
	}

	resp, err := auth.Get("/oauth2/v2/userinfo", map[string]string{"alt": "json"})
	if err != nil && err != io.EOF {
		return nil, err
	}

	user, err := utils.JsonDecode[UserData](resp.Body)
	if err != nil {
		return nil, err
	}

	if user.Error.Code != 0 {
		return nil, errors.New(user.Error.Message)
	}

	oauthUser := &oauth.OauthUser{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Avatar:      user.Picture,
		Provider:    "google",
		AccessToken: token,
	}

	return oauthUser, nil
}
