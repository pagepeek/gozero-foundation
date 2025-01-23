package tiktok

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

func (tt *TikTok) GetAccessToken(code, verifier, callbackUrl string) (*AccessToken, error) {
	data := url.Values{
		"client_key":    {tt.ClientID},
		"client_secret": {tt.Secret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {callbackUrl},
		"code_verifier": {verifier},
	}

	req, err := http.NewRequest(
		"POST",
		TikTokOauth2TokenUrl,
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := tt.HttpClient.Do(req)
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
		RefreshExpiresAt: time.Now().Unix() + token.RefreshExpiresIn,
	}}, nil
}

func (t *AccessToken) Refresh() error {
	return nil
}
