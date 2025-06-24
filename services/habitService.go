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
