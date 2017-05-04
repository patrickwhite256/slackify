package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	*http.Client
	token string
}

type profile struct {
	StatusText  string `json:"status_text"`
	StatusEmoji string `json:"status_emoji"`
}

type profileResponse struct {
	OK      bool    `json:"ok"`
	Profile profile `json:"profile"`
	Error   string  `json:"error"`
}

func NewClient(token string) *Client {
	return &Client{&http.Client{}, token}
}

func (c *Client) GetStatus() (string, string, error) {
	resp, err := c.Get("https://slack.com/api/users.profile.get?token=" + c.token)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	var responseData profileResponse
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return "", "", err
	}
	if !responseData.OK {
		return "", "", fmt.Errorf(responseData.Error)
	}
	return responseData.Profile.StatusText, responseData.Profile.StatusEmoji, nil
}

func (c *Client) SetStatus(status, emoji string) error {
	profileValue := profile{
		StatusText:  status,
		StatusEmoji: emoji,
	}

	profileString, err := json.Marshal(profileValue)
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set("token", c.token)
	data.Set("profile", string(profileString))
	resp, err := c.PostForm("https://slack.com/api/users.profile.set", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var responseData profileResponse
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return err
	}
	if !responseData.OK {
		return fmt.Errorf(responseData.Error)
	}
	return nil
}
