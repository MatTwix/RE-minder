package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MatTwix/RE-minder/database"
	"github.com/MatTwix/RE-minder/models"
	"github.com/gofiber/fiber/v3"

	"github.com/jackc/pgx/v5"
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
		"SELECT id, user_id, name, description, frequency, remind_time, start_date, timezone, created_at, updated_at FROM habits %s", whereStatement), args...)
	if err != nil {
		return habits, errors.New("Error while getting habits: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var habit models.Habit
		if err := rows.Scan(&habit.ID, &habit.UserId, &habit.Name, &habit.Description, &habit.Frequency, &habit.RemindTime, &habit.StartDate, &habit.Timezone, &habit.CreatedAt, &habit.UpdatedAt); err != nil {
			return habits, errors.New("Error parsing data: " + err.Error())
		}
		habits = append(habits, habit)
	}

	return habits, nil
}

func GetUserHabits(c fiber.Ctx, userId int) ([]models.Habit, error) {
	existingUser, err := GetUsers(c.Context(), Condition{
		Field:    "id",
		Operator: Equal,
		Value:    userId,
	})
	if err != nil {
		return nil, errors.New("Error while getting user: " + err.Error())
	}
	if len(existingUser) == 0 {
		return nil, errors.New("user not found")
	}

	habits, err := GetHabits(c.Context(), Condition{
		Field:    "user_id",
		Operator: Equal,
		Value:    userId,
	})
	if err != nil {
		return nil, errors.New("Error while getting user habits: " + err.Error())
	}
	return habits, nil
}

func GetHabitsForNotification() ([]models.Habit, error) {
	var habits []models.Habit

	rows, err := database.DB.Query(context.Background(), `
		SELECT h.id, h.user_id, h.name, h.description, h.start_date, h.frequency
		FROM habits h
		WHERE (NOW() AT TIME ZONE h.timezone)::time BETWEEN h.remind_time AND h.remind_time + INTERVAL '59 seconds'`)
	if err != nil {
		return habits, errors.New("Error while getting habits for notification: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var habit models.Habit
		if err := rows.Scan(&habit.ID, &habit.UserId, &habit.Name, &habit.Description, &habit.StartDate, &habit.Frequency); err != nil {
			log.Printf("Error scanning habit row: %v", err)
			continue
		}
		switch habit.Frequency {
		case "daily":
			if !time.Now().Before(habit.StartDate) {
				habits = append(habits, habit)
			}
		case "weekly":
			days := int(time.Since(habit.StartDate).Hours() / 24)
			if days >= 0 && days%7 == 0 {
				habits = append(habits, habit)
			}
		case "monthly":
			if time.Now().Day() == habit.StartDate.Day() && !time.Now().Before(habit.StartDate) {
				habits = append(habits, habit)
			}
		}
	}

	return habits, nil
}

func GetCreatorId(c fiber.Ctx) (int, error) {
	var userId int
	habitID := c.Params("id")
	err := database.DB.QueryRow(context.Background(), "SELECT user_id FROM habits WHERE id = $1", habitID).Scan(&userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errors.New("habit not found")
		}
		return 0, errors.New("Error getting creator ID: " + err.Error())
	}
	return userId, nil
}

func CreateHabit(ctx context.Context, userId int, name, description, frequency, remindTime, timezone string, startDate time.Time) (models.Habit, error) {
	if timezone == "" {
		timezone = "UTC"
	}

	if startDate.IsZero() {
		startDate = time.Now().In(time.FixedZone(timezone, 0))
	}

	habit := models.Habit{
		UserId:      userId,
		Name:        name,
		Description: description,
		Frequency:   frequency,
		RemindTime:  remindTime,
		Timezone:    timezone,
		StartDate:   startDate,
	}

	err := database.DB.QueryRow(context.Background(),
		`INSERT INTO habits 
		(user_id, name, description, frequency, remind_time, timezone, start_date) 
		VALUES 
		($1, $2, $3, $4, $5, COALESCE($6, 'UTC'), $7) RETURNING id, created_at, updated_at`,
		userId, name, description, frequency, remindTime, timezone, startDate).Scan(&habit.ID, &habit.CreatedAt, &habit.UpdatedAt)
	if err != nil {
		return habit, errors.New("Error creating habit: " + err.Error())
	}

	return habit, nil
}

func UpdateHabit(ctx context.Context, id int, name, description, frequency, remindTime, timezone string) (models.Habit, error) {
	if timezone == "" {
		timezone = "UTC"
	}

	habit := models.Habit{
		ID:          id,
		Name:        name,
		Description: description,
		Frequency:   frequency,
		RemindTime:  remindTime,
		Timezone:    timezone,
	}

	err := database.DB.QueryRow(ctx, `
		UPDATE habits
		SET name = $1, description = $2, frequency = $3, remind_time = $4, timezone = COALESCE($5, 'UTC'), updated_at = NOW()
		WHERE id = $6
		RETURNING user_id, created_at, updated_at`,
		name, description, frequency, remindTime, timezone, id).Scan(&habit.UserId, &habit.CreatedAt, &habit.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return habit, errors.New("habit not found")
		}
		return habit, errors.New("Error updating habit: " + err.Error())
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
