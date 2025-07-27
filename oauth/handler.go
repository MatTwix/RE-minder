package oauth

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func RedirectToProvider(c fiber.Ctx) error {
	platform := c.Params("platform")
	provider, ok := GetProvider(platform)
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Provider not found"})
	}

	userID := c.Locals("user_id").(int)
	url := provider.GetConfig().AuthCodeURL(fmt.Sprintf("%d", userID))

	return c.Redirect().Status(http.StatusTemporaryRedirect).To(url)
}

func HandleCallback(c fiber.Ctx) error {
	platform := c.Params("platform")
	provider, ok := GetProvider(platform)
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Provider not found"})
	}

	userID := c.Locals("user_id").(int)
	expectedState := fmt.Sprintf("%d", userID)
	if c.Query("state") != expectedState {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid state parameter"})
	}

	config := provider.GetConfig()
	token, err := config.Exchange(c.Context(), c.Query("code"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to exchange token", "details": err.Error()})
	}

	client := config.Client(c.Context(), token)
	userInfo, err := provider.GetUserInfo(client)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user info", "details": err.Error()})
	}

	err = linkChatToBot(provider.Platform(), expectedState, userInfo.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to link chat to bot", "details": err.Error()})
	}

	return c.Redirect().Status(http.StatusFound).To(fmt.Sprintf("/profile?linked=%s", provider.Platform()))
}
