package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/MatTwix/RE-minder/database"
	"github.com/MatTwix/RE-minder/models"
	"github.com/jackc/pgx/v5"
)

func GetUsers(ctx context.Context, optCondition ...Condition) ([]models.User, error) {
	whereStatement := ""
	args := []any{}
	if len(optCondition) > 0 {
		whereStatement = fmt.Sprintf("WHERE %s %s $1", optCondition[0].Field, optCondition[0].Operator)
		args = append(args, optCondition[0].Value)
	}

	var users []models.User

	rows, err := database.DB.Query(ctx, "SELECT * FROM users "+whereStatement, args...)
	if err != nil {
		return users, errors.New("Error while getting users: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.TelegramId, &user.GithubId, &user.CreatedAt, &user.UpdatedAt, &user.IsAdmin); err != nil {
			return users, errors.New("Error parsing data: " + err.Error())
		}
		users = append(users, user)
	}

	return users, nil
}

func CreateUser(ctx context.Context, username string, telegramId *int, githubId *int) (models.User, error) {
	user := models.User{
		Username:   username,
		TelegramId: telegramId,
		GithubId:   githubId,
	}

	err := database.DB.QueryRow(ctx,
		"INSERT INTO users (username, telegram_id, github_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at",
		username, telegramId, githubId).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, errors.New("Error creating user: " + err.Error())
	}

	return user, nil
}

func UpdateUser(ctx context.Context, id int, username string, telegramId *int, githubId *int) (models.User, error) {
	user := models.User{
		ID:         id,
		Username:   username,
		TelegramId: telegramId,
		GithubId:   githubId,
	}

	err := database.DB.QueryRow(ctx, `
		UPDATE users 
		SET username = $1, telegram_id = $2, github_id = $3, updated_at = NOW() 
		WHERE id = $4
		RETURNING created_at, updated_at`,
		username, telegramId, githubId, id).Scan(&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, errors.New("user not found")
		}
		return user, errors.New("Error updating user: " + err.Error())
	}

	return user, nil
}

func CreateOrUpdateUser(user *models.User) (models.User, error) {
	var existingUserSingle models.User

	existingUser, err := GetUsers(context.Background(), Condition{
		Field:    "github_id",
		Operator: Equal,
		Value:    user.GithubId,
	})

	if err != nil {
		return existingUserSingle, errors.New("Error checking existing user:" + err.Error())
	}

	if len(existingUser) == 0 {
		_, err = CreateUser(context.Background(), user.Username, nil, user.GithubId)
		if err != nil {
			return existingUserSingle, errors.New("Error creating new user: " + err.Error())
		}
	} else {
		_, err = UpdateUser(context.Background(), existingUser[0].ID, user.Username, nil, user.GithubId)
		if err != nil {
			return existingUserSingle, errors.New("Error updating user: " + err.Error())
		}
	}

	return existingUserSingle, nil
}

func SetUserStatus(ctx context.Context, id int, isAdmin bool) (models.User, error) {
	user := models.User{
		ID:      id,
		IsAdmin: isAdmin,
	}

	err := database.DB.QueryRow(ctx, `
		UPDATE users 
		SET is_admin = $1, updated_at = NOW() 
		WHERE id = $2
		RETURNING username, telegram_id, github_id, created_at, updated_at`,
		isAdmin, id).Scan(&user.Username, &user.TelegramId, &user.GithubId, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, errors.New("user not found")
		}
		return user, errors.New("Error updating user status: " + err.Error())
	}

	return user, nil
}

func DeleteUser(ctx context.Context, id int) error {
	result, err := database.DB.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return errors.New("Error deleting user: " + err.Error())
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
