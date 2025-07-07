package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/MatTwix/RE-minder/config"
	"github.com/MatTwix/RE-minder/models"
	"github.com/MatTwix/RE-minder/services"
	"github.com/gofiber/fiber/v3"
)

const githubTokenURL = "https://github.com/login/oauth/access_token"
const githubUserURL = "https://api.github.com/user"

var cfg = config.LoadConfig()

func RedirectToGithub(c fiber.Ctx) error {
	url := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s", cfg.GithubClient)
	return c.Redirect().To(url)
}

func GithubCallback(c fiber.Ctx) error {
	code := c.Query("code")

	resp, err := http.PostForm(githubTokenURL, map[string][]string{
		"client_id":     {cfg.GithubClient},
		"client_secret": {cfg.GithubClientSecret},
		"code":          {code},
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error getting token: " + err.Error()})
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error reading token response: " + err.Error()})
	}

	values, err := url.ParseQuery(string(bodyBytes))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Invalid token response: " + err.Error()})
	}

	accessToken := values.Get("access_token")
	if accessToken == "" {
		return c.Status(500).JSON(fiber.Map{"error": "Access token not found"})
	}

	req, _ := http.NewRequest("GET", githubUserURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	userResp, err := client.Do(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error getting user data: " + err.Error()})
	}
	defer userResp.Body.Close()

	var githubUser struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
	}

	if err := json.NewDecoder(userResp.Body).Decode(&githubUser); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error decoding user data: " + err.Error()})
	}

	user := &models.User{
		GithubId: &githubUser.ID,
		Username: githubUser.Login,
	}

	currUser, err := services.CreateOrUpdateUser(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error saving user to DB: " + err.Error()})
	}

	signedToken, err := GenerateJWT(currUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error signing token: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"token": signedToken,
		"user":  githubUser.Login,
	})
}
