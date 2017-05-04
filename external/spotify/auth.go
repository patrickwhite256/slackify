package spotify

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

var spotifyEndpoint = oauth2.Endpoint{
	AuthURL:  "https://accounts.spotify.com/authorize",
	TokenURL: "https://accounts.spotify.com/api/token",
}

func oauthClient(clientId string, clientSecret string, refreshToken string) *http.Client {
	token := &oauth2.Token{
		RefreshToken: refreshToken,
		Expiry:       time.Now().Add(-time.Minute),
	}
	conf := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8000",
		Scopes: []string{
			"user-read-playback-state",
		},
		Endpoint: spotifyEndpoint,
	}
	return conf.Client(context.Background(), token)
}
