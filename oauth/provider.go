package oauth

import (
	"net/http"

	"github.com/MatTwix/RE-minder/config"
	"golang.org/x/oauth2"
)

var cfg = config.LoadConfig()

type UserInfo struct {
	ID       string
	Username string
	Email    string
}

type OAuthProvider interface {
	GetConfig() *oauth2.Config
	GetUserInfo(client *http.Client) (UserInfo, error)
	Platform() string
}
