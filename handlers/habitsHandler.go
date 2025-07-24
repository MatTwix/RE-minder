package handlers

import (
	"strconv"
	"time"

	"github.com/MatTwix/RE-minder/config"
	"github.com/MatTwix/RE-minder/services"
	"github.com/gofiber/fiber/v3"
)

type habitsInput struct {
	Name        string    `json:"name" validate:"required,min=2,max=64"`
	Description string    `json:"description,omitempty" validate:"max=256"`
	Frequency   string    `json:"frequency" validate:"required,oneof=daily weekly monthly"`
	RemindTime  string    `json:"remind_time" validate:"required"`
	Timezone    string    `json:"timezone"`
	StartDate   time.Time `json:"start_date,omitempty"`
}

func GetHabits(c fiber.Ctx) error {
	habits, err := services.GetHabits(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error while getting habits: " + err.Error()})
	}

	return c.JSON(habits)
}

func GetUserHabits(c fiber.Ctx) error {
	userIDRaw := c.Params("id")
	userID, err := strconv.Atoi(userIDRaw)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID: " + err.Error()})
	}

	habits, err := services.GetUserHabits(c, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error while getting user habits: " + err.Error()})
	}

	return c.JSON(habits)
}

func GetHabit(c fiber.Ctx) error {
	habit, err := services.GetHabits(c.Context(), services.Condition{
		Field:    "id",
		Operator: services.Equal,
		Value:    c.Params("id"),
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error while getting habit: " + err.Error()})
	}
	if len(habit) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Habit not found"})
	}
	singleHabit := habit[0]

	return c.JSON(singleHabit)
}

func CreateHabit(c fiber.Ctx) error {
	var input habitsInput

	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Incorrect data format: " + err.Error()})
	}

	if err := config.Validator.Struct(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Validation error: " + err.Error()})
	}

	userIDRaw := c.Locals("user_id")
	userID, ok := userIDRaw.(int)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	createdHabit, err := services.CreateHabit(c.Context(), userID, input.Name, input.Description, input.Frequency, input.RemindTime, input.Timezone, input.StartDate)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error creating habit: " + err.Error()})
	}

	return c.JSON(createdHabit)
}

func UpdateHabit(c fiber.Ctx) error {
	idRaw := c.Params("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid habit ID: " + err.Error()})
	}

	var input habitsInput

	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Incorrect data format: " + err.Error()})
	}

	if err := config.Validator.Struct(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Validation error: " + err.Error()})
	}

	updatedHabit, err := services.UpdateHabit(c.Context(), id, input.Name, input.Description, input.Frequency, input.RemindTime, input.Timezone)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error updating habit: " + err.Error()})
	}

	return c.JSON(updatedHabit)
}

func DeleteHabit(c fiber.Ctx) error {
	idRaw := c.Params("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid habit ID: " + err.Error()})
	}

	if err := services.DeleteHabit(c.Context(), id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error deleting habit: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Habit succesfully deleted"})
}
