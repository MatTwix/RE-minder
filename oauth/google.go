package oauth

import (
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleProvider struct {
	config *oauth2.Config
}

func NewGoogleProvider() OAuthProvider {
	return &googleProvider{
		config: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientID,
			RedirectURL:  cfg.GoogleRedirectUrl,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (p *googleProvider) GetConfig() *oauth2.Config {
	return p.config
}

func (p *googleProvider) Platform() string {
	return "google"
}

func (p *googleProvider) GetUserInfo(client *http.Client) (UserInfo, error) {
	var userInfo UserInfo

	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return userInfo, err
	}
	defer res.Body.Close()

	var googleUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(res.Body).Decode(&googleUser); err != nil {
		return userInfo, err
	}

	userInfo.ID = googleUser.ID
	userInfo.Email = googleUser.Email
	userInfo.Username = googleUser.Name

	return userInfo, nil
}
