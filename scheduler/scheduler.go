package scheduler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/MatTwix/RE-minder/services"
	"github.com/robfig/cron/v3"
)

type NotificationTask struct {
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
}

func StartScheduler() {
	log.Println("Starting scheduler...")

	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("0 * * * * *", checkHabitsForReminder)
	if err != nil {
		log.Fatalf("Could not add cron job: %v", err)
	}

	go c.Start()

	log.Println("Scheduler started successfully.")
}

func checkHabitsForReminder() {
	log.Println("Running cron job: Checking habits for notification..")

	habits, err := services.GetHabitsForNotification()
	if err != nil {
		log.Printf("Error getting habits from service: %v", err)
		return
	}

	if len(habits) == 0 {
		log.Println("There are no habits to remind at this time.")
	}

	for _, habit := range habits {
		message := fmt.Sprintf("Remind: %s", habit.Name)
		if habit.Description != "" {
			message = fmt.Sprintf("%s (%s)", message, habit.Description)
		}
		task := NotificationTask{
			UserID:  habit.UserId,
			Message: message,
		}

		taskBody, err := json.Marshal(task)
		if err != nil {
			log.Printf("Error marshaling task for user %d: %v", habit.UserId, err)
			continue
		}

		log.Printf("SIMULATING PUBLISH TO RabbitMQ: %s", string(taskBody))
	}
}
