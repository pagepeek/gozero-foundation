package google

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pagepeek/gozero-foundation/pkg/oauth"
)

type Google struct {
	ClientID   string
	Secret     string
	HttpClient *http.Client
}

type AuthClient struct {
	*Google
	*oauth.BaseClient

	AccessToken oauth.AccessToken
}

func New(clientID, secret string, httpClient *http.Client) *Google {
	return &Google{
		ClientID:   clientID,
		Secret:     secret,
		HttpClient: httpClient,
	}
}

func (g *Google) Auth(token oauth.AccessToken) (oauth.AuthClient, error) {
	if !token.Valid() {
		return nil, errors.New("Google authorization invalid")
	}

	base := oauth.JsonClient(GoogleHost, g.HttpClient, func(req *http.Request) {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	})

	return &AuthClient{g, base, token}, nil
}
