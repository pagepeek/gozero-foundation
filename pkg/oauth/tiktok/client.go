package tiktok

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pagepeek/gozero-foundation/pkg/oauth"
	"github.com/pagepeek/gozero-foundation/pkg/utils"
)

type TikTok struct {
	ClientID   string
	Secret     string
	HttpClient *http.Client
}

type AuthClient struct {
	*TikTok
	*oauth.BaseClient

	AccessToken oauth.AccessToken
}

func (tt *TikTok) Auth(token oauth.AccessToken) (oauth.AuthClient, error) {
	if !token.Valid() {
		return nil, errors.New("TikTok authorization invalid")
	}

	base := oauth.JsonClient(TikTokHost, tt.HttpClient, func(req *http.Request) {
		req.Header.Add("Idempotency-Key", utils.RandStr(16))
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	})

	return &AuthClient{tt, base, token}, nil
}
