package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/MatTwix/RE-minder/database"
	"github.com/MatTwix/RE-minder/models"
)

func GetHabits(ctx context.Context, optCondition ...Condition) ([]models.Habit, error) {
	whereStatement := ""
	args := []any{}
	if len(optCondition) > 0 {
		whereStatement = fmt.Sprintf("WHERE %s %s $1", optCondition[0].Field, optCondition[0].Operator)
		args = append(args, optCondition[0].Value)
	}

	var habits []models.Habit

	rows, err := database.DB.Query(ctx, fmt.Sprintf(
		"SELECT * FROM habits %s", whereStatement), args...)
	if err != nil {
		return habits, errors.New("Error while getting habits: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var habit models.Habit
		if err := rows.Scan(&habit.ID, &habit.UserId, &habit.Name, &habit.Description, &habit.Frequency, &habit.RemindTime, &habit.Timezone, &habit.CreatedAt, &habit.UpdatedAt); err != nil {
			return habits, errors.New("Error parsing data: " + err.Error())
		}
		habits = append(habits, habit)
	}

	return habits, nil
}

func CreateHabit(ctx context.Context, userId int, name, description, frequency, remindTime, timezone string) (models.Habit, error) {
	if timezone == "" {
		timezone = "UTC"
	}

	habit := models.Habit{
		UserId:      userId,
		Name:        name,
		Description: description,
		Frequency:   frequency,
		RemindTime:  remindTime,
		Timezone:    timezone,
	}

	err := database.DB.QueryRow(context.Background(),
		`INSERT INTO habits 
		(user_id, name, description, frequency, remind_time, timezone) 
		VALUES 
		($1, $2, $3, $4, $5, COALESCE($6, 'UTC')) RETURNING id, created_at, updated_at`,
		userId, name, description, frequency, remindTime, timezone).Scan(&habit.ID, &habit.CreatedAt, &habit.UpdatedAt)
	if err != nil {
		return habit, errors.New("Error creating habit: " + err.Error())
	}

	return habit, nil
}

func UpdateHabit(ctx context.Context, id int, userId int, name, description, frequency, remindTime, timezone string) (models.Habit, error) {
	if timezone == "" {
		timezone = "UTC"
	}

	habit := models.Habit{
		ID:          id,
		UserId:      userId,
		Name:        name,
		Description: description,
		Frequency:   frequency,
		RemindTime:  remindTime,
		Timezone:    timezone,
	}

	result, err := database.DB.Exec(ctx, `
		UPDATE habits
		SET user_id = $1, name = $2, description = $3, frequency = $4, remind_time = $5, timezone = COALESCE($6, 'UTC'), updated_at = NOW()
		WHERE id = $7`, userId, name, description, frequency, remindTime, timezone, id)
	if err != nil {
		return habit, errors.New("Error updating habit: " + err.Error())
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return habit, errors.New("habit not found or no changes made")
	}

	return habit, nil
}

func DeleteHabit(ctx context.Context, id int) error {
	result, err := database.DB.Exec(ctx, "DELETE FROM habits WHERE id = $1", id)
	if err != nil {
		return errors.New("Error deleting habit: " + err.Error())
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("habit not found")
	}

	return nil
}
