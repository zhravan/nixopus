package notification

import (
	"context"
	"sync"
	"time"

	"github.com/slack-go/slack"
	"github.com/uptrace/bun"
)

type NotificationChannels struct {
	Email   *Email   `json:"email"`
	Slack   *Slack   `json:"slack"`
	Discord *Discord `json:"discord"`
}

type Email struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
}

type Slack struct {
	SlackClient *slack.Client
	ChannelID   string
}

type Discord struct {
	WebhookUrl string `json:"webhook_url"`
}

func NewNotificationChannels() *NotificationChannels {
	return &NotificationChannels{
		Email:   &Email{},
		Slack:   &Slack{},
		Discord: &Discord{},
	}
}

type NotificationAuthenticationData struct {
	Email    string
	IP       string
	Browser  string
	UserName string
}

type NotificationManager struct {
	sync.RWMutex
	Channels    *NotificationChannels
	PayloadChan chan NotificationPayload
	ctx         context.Context
	cancel      context.CancelFunc
	db          *bun.DB
}

type NotificationPayloadType string

const (
	NotificationPayloadTypeRegister      NotificationPayloadType = "register"
	NotificationPayloadTypeLogin         NotificationPayloadType = "login"
	NotificationPayloadTypeLogout        NotificationPayloadType = "logout"
	NotificationPayloadTypePasswordReset NotificationPayloadType = "password_reset"
)

type NotificationCategory string

const (
	NotificationCategoryAuthentication NotificationCategory = "authentication"
)

type NotificationPayload struct {
	Category  NotificationCategory    `json:"category"`
	Type      NotificationPayloadType `json:"type"`
	UserID    string                  `json:"user_id"`
	Timestamp time.Time               `json:"timestamp"`
	Data      interface{}             `json:"data"`
}
