package handlers

import (
	"context"

	"github.com/MatTwix/RE-minder/database"
	"github.com/MatTwix/RE-minder/models"
	"github.com/gofiber/fiber/v3"
)

func GetHabits(c fiber.Ctx) error {
	rows, err := database.DB.Query(context.Background(), "SELECT * FROM habits")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error while getting habits: " + err.Error()})
	}

	defer rows.Close()

	var habits []models.Habit

	for rows.Next() {
		var habit models.Habit
		if err := rows.Scan(&habit.ID, &habit.UserId, &habit.Name, &habit.Description, &habit.Frequency, &habit.RemindTime, &habit.Timezone, &habit.CreatedAt, &habit.UpdatedAt); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error parsing data: " + err.Error()})
		}
		habits = append(habits, habit)
	}

	return c.JSON(habits)
}

func GetHabit(c fiber.Ctx) error {
	id := c.Params("id")

	var habit models.Habit
	err := database.DB.QueryRow(context.Background(),
		"SELECT * FROM habits WHERE id = $1", id).
		Scan(&habit.ID, &habit.UserId, &habit.Name, &habit.Description, &habit.Frequency, &habit.RemindTime, &habit.Timezone, &habit.CreatedAt, &habit.UpdatedAt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Habit not found: " + err.Error()})
	}

	return c.JSON(habit)
}

func CreateHabit(c fiber.Ctx) error {
	habit := new(models.Habit)
	if err := c.Bind().Body(habit); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Incorrect data format: " + err.Error()})
	}

	_, err := database.DB.Exec(context.Background(),
		`INSERT INTO habits 
		(user_id, name, description, frequency, remind_time, timezone) 
		VALUES 
		($1, $2, $3, $4, $5, $6)`,
		habit.UserId, habit.Name, habit.Description, habit.Frequency, habit.RemindTime, habit.Timezone)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error creating habit: " + err.Error()})
	}

	return c.JSON(habit)
}

func UpdateHabit(c fiber.Ctx) error {
	id := c.Params("id")
	habit := new(models.Habit)

	if err := c.Bind().Body(habit); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Incorrect data format: " + err.Error()})
	}

	_, err := database.DB.Exec(context.Background(), `
		UPDATE habits
		SET user_id = $1, name = $2, description = $3, frequency = $4, remind_time = $5, timezone = $6, updated_at = NOW()
		WHERE id = $7`,
		habit.UserId, habit.Name, habit.Description, habit.Frequency, habit.RemindTime, habit.Timezone, id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error updating habit: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Habit seccesfully updated"})
}

func DeleteHabit(c fiber.Ctx) error {
	id := c.Params("id")

	_, err := database.DB.Exec(context.Background(), "DELETE FROM habits WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error deleting habit: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Habit succesfully deleted"})
}
