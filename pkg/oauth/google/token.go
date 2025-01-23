package google

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pagepeek/gozero-foundation/pkg/oauth"
	"github.com/pagepeek/gozero-foundation/pkg/utils"
)

type AccessToken struct {
	*oauth.BaseAccessToken
}

func (g *Google) GetAccessToken(code, verifier, callbackUrl string) (*AccessToken, error) {
	data := url.Values{
		"client_id":     {g.ClientID},
		"client_secret": {g.Secret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {callbackUrl},
	}

	req, err := http.NewRequest(
		http.MethodPost,
		GoogleTokenUrl,
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := g.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	token, err := utils.JsonDecode[AccessTokenResp](resp.Body)
	if err != nil {
		return nil, err
	}

	if token.Error != "" {
		return nil, errors.New(token.ErrDesc)
	}

	return &AccessToken{&oauth.BaseAccessToken{
		Token:            token.AccessToken,
		ExpiresAt:        time.Now().Unix() + token.ExpiresIn,
		RefreshToken:     token.RefreshToken,
		RefreshExpiresAt: time.Now().Unix() + token.ExpiresIn,
	}}, nil
}

func (t *AccessToken) Refresh() error {
	return nil
}
