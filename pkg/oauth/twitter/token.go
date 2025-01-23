package twitter

import (
	"encoding/base64"
	"errors"
	"fmt"
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

func (x *Twitter) GetAccessToken(code, verifier, callbackUrl string) (*AccessToken, error) {
	data := url.Values{
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {callbackUrl},
		"code_verifier": {verifier},
	}

	req, err := http.NewRequest(
		http.MethodPost,
		TwitterOauth2TokenUrl,
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return nil, err
	}

	basicToken := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%s", x.ClientID, x.Secret)),
	)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+basicToken)

	resp, err := x.HttpClient.Do(req)
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
