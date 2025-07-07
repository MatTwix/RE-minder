package middleware

import (
	"time"

	"github.com/MatTwix/RE-minder/config"
	"github.com/MatTwix/RE-minder/models"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(githubUser models.User) (string, error) {
	cfg := config.LoadConfig()

	claims := jwt.MapClaims{
		"user_id":  githubUser.ID,
		"is_admin": githubUser.IsAdmin,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
