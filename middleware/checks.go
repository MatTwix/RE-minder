package middleware

import (
	"strconv"

	"github.com/MatTwix/RE-minder/services"
	"github.com/gofiber/fiber/v3"
)

func IsAdmin() fiber.Handler {
	return func(c fiber.Ctx) error {
		isAdminRaw := c.Locals("is_admin")

		isAdmin := false
		if isAdminRaw != nil {
			isAdmin = isAdminRaw.(bool)
		}

		if isAdmin {
			return c.Next()
		}

		return fiber.ErrForbidden
	}
}

func IsSelf() fiber.Handler {
	return func(c fiber.Ctx) error {
		currUserIDRaw := c.Locals("user_id")
		currUserID := 0
		if currUserIDRaw != nil {
			currUserID = currUserIDRaw.(int)
		}

		userID := c.Params("id")

		if strconv.Itoa(currUserID) == userID {
			return c.Next()
		}

		return fiber.ErrForbidden
	}
}

func IsSelfOrAdmin() fiber.Handler {
	return func(c fiber.Ctx) error {
		isAdminRaw := c.Locals("is_admin")
		isAdmin := false
		if isAdminRaw != nil {
			isAdmin = isAdminRaw.(bool)
		}

		if isAdmin {
			return c.Next()
		}

		return IsSelf()(c)
	}
}

func IsCreatorOrAdmin() fiber.Handler {
	return func(c fiber.Ctx) error {
		userIDRaw := c.Locals("user_id")
		userID := 0
		if userIDRaw != nil {
			userID = userIDRaw.(int)
		}

		isAdminRaw := c.Locals("is_admin")
		isAdmin := false
		if isAdminRaw != nil {
			isAdmin = isAdminRaw.(bool)
		}

		authorID, err := services.GetCreatorId(c)
		if err != nil {
			return fiber.NewError(fiber.StatusForbidden, "Error getting creator ID: "+err.Error())
		}
		if userID == authorID || isAdmin {
			return c.Next()
		}

		return fiber.ErrForbidden
	}
}
