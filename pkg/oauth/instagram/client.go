package instagram

import (
	"errors"
	"net/http"
	"time"

	"github.com/pagepeek/gozero-foundation/pkg/oauth"
	"github.com/pagepeek/gozero-foundation/pkg/utils"
)

const (
	InstagramVersion        = "v21.0"
	InstagramHost           = "graph.instagram.com"
	InstagramOauthUrl       = "https://api.instagram.com/oauth/authorize"
	InstagramOauth2TokenUrl = "https://api.instagram.com/oauth/access_token"
	InstagramUserTokenUrl   = "https://graph.instagram.com/access_token"
)

type Instagram struct {
	ClientID   string
	Secret     string
	HttpClient *http.Client
}

type AuthClient struct {
	*Instagram
	*oauth.BaseClient

	AccessToken oauth.AccessToken
}

func (ins *Instagram) Auth(token oauth.AccessToken) (oauth.AuthClient, error) {
	if !token.Valid() {
		return nil, errors.New("Instagram authorization invalid")
	}

	base := oauth.QueryClient(InstagramHost, ins.HttpClient, func(req *http.Request) {
		query := req.URL.Query()
		query.Add("access_token", token.String())

		req.URL.RawQuery = query.Encode()
	})

	return &AuthClient{ins, base, token}, nil
}

func (auth *AuthClient) LongLiveToken(accessToken string) (*AccessToken, error) {
	resp, err := auth.Get(InstagramUserTokenUrl, map[string]string{
		"grant_type":    "ig_exchange_token",
		"client_secret": auth.Secret,
	})

	if err != nil {
		return nil, err
	}

	token, err := utils.JsonDecode[UserTokenResp](resp.Body)
	if err != nil {
		return nil, err
	}

	return &AccessToken{&oauth.BaseAccessToken{
		Token:            token.AccessToken,
		ExpiresAt:        time.Now().Unix() + token.ExpiresIn,
		RefreshToken:     token.AccessToken,
		RefreshExpiresAt: time.Now().Unix() + token.ExpiresIn,
	}}, nil
}
