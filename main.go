package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	// Read environment variables
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	if botToken == "" {
		log.Fatal("SLACK_BOT_TOKEN environment variable is required")
	}
	if appToken == "" {
		log.Fatal("SLACK_APP_TOKEN environment variable is required")
	}

	// Create Slack client
	api := slack.New(
		botToken,
		slack.OptionAppLevelToken(appToken),
		slack.OptionDebug(false),
	)

	// Create Socket Mode client
	socketClient := socketmode.New(api)

	// Register slash command handler
	go func() {
		for evt := range socketClient.Events {
			switch evt.Type {
			case socketmode.EventTypeSlashCommand:
				cmd, ok := evt.Data.(slack.SlashCommand)
				if !ok {
					log.Printf("Could not type cast slash command")
					continue
				}

				// Handle the /pomodoro command
				if cmd.Command == "/pomodoro" {
					handlePomodoroCommand(socketClient, api, &evt, cmd)
				} else {
					socketClient.Ack(*evt.Request, map[string]interface{}{"text": "Unknown command"})
				}
			}
		}
	}()

	log.Println("Tomatobot is running...")
	if err := socketClient.Run(); err != nil {
		log.Fatal(err)
	}
}

func handlePomodoroCommand(socketClient *socketmode.Client, api *slack.Client, evt *socketmode.Event, cmd slack.SlashCommand) {
	// Parse command text
	text := strings.TrimSpace(cmd.Text)
	minutes := 25 // default
	label := "your pomodoro"

	if text != "" {
		parts := strings.Fields(text)

		// Try to parse first word as a number
		if len(parts) > 0 {
			if num, err := strconv.Atoi(parts[0]); err == nil {
				minutes = num
				// Remaining parts form the label
				if len(parts) > 1 {
					label = strings.Join(parts[1:], " ")
				}
			} else {
				// No leading number, entire text is the label
				label = text
			}
		}
	}

	// Validate duration
	if minutes < 1 || minutes > 120 {
		responseText := "Duration must be between 1 and 120 minutes."
		socketClient.Ack(*evt.Request, map[string]interface{}{"text": responseText})
		return
	}

	// Calculate schedule time
	scheduleAt := time.Now().Add(time.Duration(minutes) * time.Minute).Unix()
	scheduleAtStr := strconv.FormatInt(scheduleAt, 10)

	// Schedule the reminder message
	reminderText := fmt.Sprintf("Hey <@%s>, it's been %d minutes since you started '%s'", cmd.UserID, minutes, label)

	_, _, err := api.ScheduleMessage(
		cmd.ChannelID,
		scheduleAtStr,
		slack.MsgOptionText(reminderText, false),
		slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
			UnfurlLinks: false,
		}),
	)

	if err != nil {
		responseText := fmt.Sprintf("Failed to schedule reminder: %v", err)
		socketClient.Ack(*evt.Request, map[string]interface{}{"text": responseText})
		log.Printf("Error scheduling message: %v", err)
		return
	}

	// Send acknowledgment
	ackText := fmt.Sprintf("Got it! I'll remind you in %d minutes about '%s'", minutes, label)
	socketClient.Ack(*evt.Request, map[string]interface{}{"text": ackText})

	log.Printf("Scheduled pomodoro reminder for user %s: %d minutes, label: %s", cmd.UserID, minutes, label)
}
