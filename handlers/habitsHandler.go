package handlers

import (
	"strconv"

	"github.com/MatTwix/RE-minder/models"
	"github.com/MatTwix/RE-minder/services"
	"github.com/gofiber/fiber/v3"
)

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
	habit := new(models.Habit)
	if err := c.Bind().Body(habit); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Incorrect data format: " + err.Error()})
	}

	userIDRaw := c.Locals("user_id")
	userID, ok := userIDRaw.(int)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	habit.UserId = userID

	createdHabit, err := services.CreateHabit(c.Context(), habit.UserId, habit.Name, habit.Description, habit.Frequency, habit.RemindTime, habit.Timezone, habit.StartDate)
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

	habit := new(models.Habit)

	if err := c.Bind().Body(habit); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Incorrect data format: " + err.Error()})
	}

	updatedHabit, err := services.UpdateHabit(c.Context(), id, habit.Name, habit.Description, habit.Frequency, habit.RemindTime, habit.Timezone)
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
