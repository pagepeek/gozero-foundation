package instagram

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/pagepeek/gozero-foundation/pkg/oauth"
	"github.com/pagepeek/gozero-foundation/pkg/utils"
)

type AccessToken struct {
	*oauth.BaseAccessToken
}

func (ins *Instagram) GetAccessToken(code, verifier, callbackUrl string) (*AccessToken, error) {
	data := url.Values{
		"client_id":     {ins.ClientID},
		"client_secret": {ins.Secret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {callbackUrl},
	}

	req, err := http.NewRequest(
		"POST",
		InstagramOauth2TokenUrl,
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := ins.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	token, err := utils.JsonDecode[AccessTokenResp](resp.Body)
	if err != nil {
		return nil, err
	}

	if token.ErrMsg != "" {
		return nil, errors.New(token.ErrMsg)
	}

	auth, _ := ins.Auth(&AccessToken{&oauth.BaseAccessToken{Token: token.AccessToken}})

	return auth.(*AuthClient).LongLiveToken(token.AccessToken)
}

func (t *AccessToken) Refresh() error {
	return nil
}

func (t *AccessToken) Valid() bool {
	// 未设置过期时间,视为短效token,不做校验
	if t.ExpiresAt == 0 || t.RefreshExpiresAt == 0 {
		return true
	}

	return t.BaseAccessToken.Valid()
}
