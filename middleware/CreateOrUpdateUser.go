package middleware

import (
	"context"
	"log"

	"github.com/MatTwix/RE-minder/database"
	"github.com/MatTwix/RE-minder/models"
)

func CreateOrUpdateUser(user *models.User) error {
	var existingUser models.User
	err := database.DB.QueryRow(context.Background(), "SELECT id, github_id, username FROM users WHERE github_id=$1", user.GithubId).
		Scan(&existingUser.ID, &existingUser.GithubId, &existingUser.Username)

	if err != nil {
		if err.Error() == "no rows in result set" {
			githubEmail := ""
			githubPassword := ""

			_, err = database.DB.Exec(context.Background(), `
				INSERT INTO users (github_id, username, email, password_hash)
				VALUES ($1, $2, $3, $4)`, user.GithubId, user.Username, githubEmail, githubPassword)
			if err != nil {
				log.Println("Error creating new user:", err)
				return err
			}
			log.Println("User created successfully!")
			return nil
		}
		log.Println("Error checking user:", err)
		return err
	}

	_, err = database.DB.Exec(context.Background(), `
		UPDATE users SET username=$1, updated_at=NOW() WHERE github_id=$2`,
		user.Username, user.GithubId)
	if err != nil {
		log.Println("Error updating user:", err)
		return err
	}
	log.Println("User updated successfully!")
	return nil
}
