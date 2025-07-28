package oauth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"
)

type vkProvider struct {
	config *oauth2.Config
}

func NewVKProider() OAuthProvider {
	return &vkProvider{
		config: &oauth2.Config{
			ClientID:     cfg.VKClientID,
			ClientSecret: cfg.VKClientSecret,
			RedirectURL:  cfg.VKRedirectUrl,
			Scopes:       nil,
			Endpoint:     vk.Endpoint,
		},
	}
}

func (p *vkProvider) GetConfig() *oauth2.Config {
	return p.config
}

func (p *vkProvider) Platform() string {
	return "vk"
}

func (p *vkProvider) GetUserInfo(client *http.Client) (UserInfo, error) {
	var userInfo UserInfo

	url := "https://api.vk.com/method/users.get?fields=screen_name&v=5.199"

	res, err := client.Get(url)
	if err != nil {
		return userInfo, err
	}
	defer res.Body.Close()

	var vkResponce struct {
		Response []struct {
			ID         int    `json:"id"`
			FirstName  string `json:"first_name"`
			LastName   string `json:"last_name"`
			ScreenName string `json:"screen_name"`
		} `json:"response"`
	}

	if err := json.NewDecoder(res.Body).Decode(&vkResponce); err != nil {
		return userInfo, err
	}

	if len(vkResponce.Response) == 0 {
		return userInfo, errors.New("vk api returned no user")
	}

	vkUser := vkResponce.Response[0]

	userInfo.ID = strconv.Itoa(vkUser.ID)
	userInfo.Username = vkUser.ScreenName

	return userInfo, nil
}
