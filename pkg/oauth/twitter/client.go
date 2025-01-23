package twitter

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pagepeek/gozero-foundation/pkg/oauth"
	"github.com/pagepeek/gozero-foundation/pkg/utils"
)

type Twitter struct {
	ClientID   string
	Secret     string
	HttpClient *http.Client
}

type AuthClient struct {
	*Twitter
	*oauth.BaseClient

	AccessToken oauth.AccessToken
}

func (x *Twitter) Auth(token oauth.AccessToken) (oauth.AuthClient, error) {
	if !token.Valid() {
		return nil, errors.New("Twitter authorization invalid")
	}

	base := oauth.JsonClient(TwitterV2Host, x.HttpClient, func(req *http.Request) {
		req.Header.Add("Idempotency-Key", utils.RandStr(16))
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	})

	return &AuthClient{x, base, token}, nil
}
