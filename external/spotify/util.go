package spotify

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/patrickwhite256/slackify/config"
)

type Client struct {
	*http.Client
}

type Song struct {
	Artists string
	Name    string
}

type playingResponse struct {
	IsPlaying bool        `json:"is_playing"`
	Item      itemDetails `json:item"`
}

type itemDetails struct {
	Artists   []artistDetails `json:"artists"`
	TrackName string          `json:"name"`
}

type artistDetails struct {
	Name string `json:"name"`
}

func NewClient(conf *config.Config, refreshToken string) *Client {
	return &Client{oauthClient(conf.SpotifyClientId, conf.SpotifyClientSecret, refreshToken)}
}

// Returns nil if no song is being played
func (c *Client) GetCurrentlyPlaying() (*Song, error) {
	resp, err := c.Get("https://api.spotify.com/v1/me/player/currently-playing")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var responseData playingResponse
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}
	if !responseData.IsPlaying { // covers both not playing and paused
		return nil, nil
	} else {
		artists := make([]string, len(responseData.Item.Artists))
		for i, artist := range responseData.Item.Artists {
			artists[i] = artist.Name
		}
		song := &Song{
			Name:    responseData.Item.TrackName,
			Artists: strings.Join(artists, ", "),
		}
		return song, nil
	}
}
