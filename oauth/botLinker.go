package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/MatTwix/RE-minder/config"
)

func linkChatToBot(platform, reminderUserIDStr, platformChatID string) error {
	cfg := config.LoadConfig()

	reminderUserID, err := strconv.Atoi(reminderUserIDStr)
	if err != nil {
		return errors.New("invalid reminder user ID: " + reminderUserIDStr)
	}

	payload := map[string]any{
		"reminder_user_id": reminderUserID,
		"platform":         platform,
		"chat_id":          platformChatID,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return errors.New("error marshalling payload: " + err.Error())
	}

	url := cfg.BotsApiUrl + "/api/chats"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return errors.New("error creating request to bot api: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", cfg.InternalApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("error sending request to bot api: " + err.Error())
	}

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Print(string(bodyBytes))
		return errors.New("bot api returned non-2xx status: " + resp.Status)
	}
	return nil
}
