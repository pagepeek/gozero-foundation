package instagram

import (
	"errors"
	"net/url"
	"strings"

	"github.com/pagepeek/gozero-foundation/pkg/oauth"
	"github.com/pagepeek/gozero-foundation/pkg/utils"
)

func (ins *Instagram) GetAuthUrl(redirectUrl string, scopes ...string) *oauth.OauthRedirectUrl {
	verifier := utils.RandStr(43)

	return &oauth.OauthRedirectUrl{
		Target:   InstagramOauthUrl,
		Provider: "instagram",
		Verifier: verifier,
		Params: url.Values{
			"client_id":     {ins.ClientID},
			"scope":         {strings.Join(scopes, ",")},
			"redirect_uri":  {redirectUrl},
			"state":         {verifier},
			"response_type": {"code"},
		},
	}
}

func (ins *Instagram) CodeToUser(code, verifier, callbackUrl string) (*oauth.OauthUser, error) {
	token, err := ins.GetAccessToken(code, verifier, callbackUrl)
	if err != nil {
		return nil, err
	}

	auth, err := ins.Auth(token)
	if err != nil {
		return nil, err
	}

	resp, err := auth.Get("/v21.0/me", map[string]string{
		"fields": strings.Join([]string{"id", "user_id", "name", "username"}, ","),
	})
	if err != nil {
		return nil, err
	}

	user, err := utils.JsonDecode[UserData](resp.Body)
	if err != nil {
		return nil, err
	}

	if user.ErrMsg != "" {
		return nil, errors.New(user.ErrMsg)
	}

	oauthUser := &oauth.OauthUser{
		ID:          user.UserId,
		Name:        user.Name,
		Provider:    "instagram",
		AccessToken: token,
	}

	return oauthUser, nil
}
