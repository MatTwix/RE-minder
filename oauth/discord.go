package oauth

import (
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
)

type discordProvider struct {
	config *oauth2.Config
}

func NewDiscordProvider() OAuthProvider {
	return &discordProvider{
		config: &oauth2.Config{
			ClientID:     cfg.DiscordClientID,
			ClientSecret: cfg.DiscordClientSecret,
			RedirectURL:  cfg.DiscordRedirectUrl,

			Scopes: []string{"identify"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://discord.com/api/oauth2/authorize",
				TokenURL: "https://discord.com/api/oauth2/token",
			},
		},
	}
}

func (p *discordProvider) GetConfig() *oauth2.Config {
	return p.config
}

func (p *discordProvider) Platform() string {
	return "discord"
}

func (p *discordProvider) GetUserInfo(client *http.Client) (UserInfo, error) {
	var userInfo UserInfo

	res, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		return userInfo, err
	}
	defer res.Body.Close()

	var discordUser struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(res.Body).Decode(&discordUser); err != nil {
		return userInfo, err
	}

	userInfo.ID = discordUser.ID
	userInfo.Username = discordUser.Username
	userInfo.Email = discordUser.Email

	return userInfo, nil
}
