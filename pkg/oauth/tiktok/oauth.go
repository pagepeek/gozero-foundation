package tiktok

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"strings"

	"github.com/pagepeek/gozero-foundation/pkg/oauth"
	"github.com/pagepeek/gozero-foundation/pkg/utils"
)

const (
	TikTokHost           = "open.tiktokapis.com"
	TikTokOauthUrl       = "https://www.tiktok.com/v2/auth/authorize/"
	TikTokOauth2TokenUrl = "https://open.tiktokapis.com/v2/oauth/token/"
)

func (tt *TikTok) GetAuthUrl(redirectUrl string, scopes ...string) *oauth.OauthRedirectUrl {
	hash := sha256.New()
	hash.Write([]byte(utils.RandStr(43)))
	verifier := hex.EncodeToString(hash.Sum(nil))

	return &oauth.OauthRedirectUrl{
		Target:   TikTokOauthUrl,
		Provider: "tiktok",
		Verifier: verifier,
		Params: url.Values{
			"client_key":            {tt.ClientID},
			"scope":                 {strings.Join(scopes, ",")},
			"redirect_uri":          {redirectUrl},
			"state":                 {"state"},
			"response_type":         {"code"},
			"code_challenge":        {hex.EncodeToString(hash.Sum(nil))},
			"code_challenge_method": {"S256"},
		},
	}
}

func (tt *TikTok) CodeToUser(code, verifier, callbackUrl string) (*oauth.OauthUser, error) {
	token, err := tt.GetAccessToken(code, verifier, callbackUrl)
	if err != nil {
		return nil, err
	}

	auth, err := tt.Auth(token)
	if err != nil {
		return nil, err
	}

	resp, err := auth.Get("/v2/user/info/", map[string]string{
		"fields": strings.Join([]string{"open_id", "union_id", "avatar_url", "display_name"}, ","),
	})
	if err != nil {
		return nil, err
	}

	user, err := utils.JsonDecode[UserData](resp.Body)
	if err != nil {
		return nil, err
	}

	if user.Error.Code != "ok" {
		return nil, errors.New(user.Error.Message)
	}

	oauthUser := &oauth.OauthUser{
		ID:          user.Data.User.OpenId,
		Name:        user.Data.User.Name,
		Provider:    "tiktok",
		AccessToken: token,
	}

	return oauthUser, nil
}
