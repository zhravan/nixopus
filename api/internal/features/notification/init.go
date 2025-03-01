package notification

import (
	"context"
	"fmt"
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
	Email   string
	IP      string
	Browser string
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
	NotificationPayloadTypeRegister NotificationPayloadType = "register"
	NotificationPayloadTypeLogin    NotificationPayloadType = "login"
	NotificationPayloadTypeLogout   NotificationPayloadType = "logout"
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
	Data     interface{}  `json:"data"`
}

func NewNotificationPayload(payloadType NotificationPayloadType, userID string, data interface{}, category NotificationCategory) NotificationPayload {
	return NotificationPayload{
		Type:      payloadType,
		UserID:    userID,
		Timestamp: time.Now(),
		Data:      data,
		Category:  category,
	}
}

func NewNotificationManager(channels *NotificationChannels, db *bun.DB) *NotificationManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &NotificationManager{
		Channels:    channels,
		PayloadChan: make(chan NotificationPayload, 100),
		ctx:         ctx,
		cancel:      cancel,
		db:          db,
	}
}

// Start starts the notification manager to listen for notifications from a go routine
// based on the type of the notification that we receive we can check the notification preferences of the user
// and then send the notification to the respective channel
func (m *NotificationManager) Start() {
	go func() {
		for {
			select {
			case payload := <-m.PayloadChan:
				switch payload.Type {
				case NotificationPayloadTypeRegister:
					// m.Channels.Email.SendRegisterEmail(payload.UserID, payload.OrganizationID, payload.Data)
				case NotificationPayloadTypeLogin:
					fmt.Println("NotificationPayloadTypeLogin", payload)
					// m.Channels.Email.SendLoginEmail(payload.UserID, payload.OrganizationID, payload.Data)
				case NotificationPayloadTypeLogout:
					// m.Channels.Email.SendLogoutEmail(payload.UserID, payload.OrganizationID, payload.Data)
				}
			case <-m.ctx.Done():
				return
			}
		}
	}()
}

func (m *NotificationManager) Stop() {
	m.cancel()
}

// SendNotification sends a notification
func (m *NotificationManager) SendNotification(payload NotificationPayload) {
	m.PayloadChan <- payload
}

// here we can get the notification preferences of the user (like should send to slack/email/discord, how many times to send, what type of contents to send)
func (m *NotificationManager) CheckUserNotificationPreferences(userID string) {

}

// we will categorize the notifications based on the type of the notification
func (m *NotificationManager) GetPreferencesBasedOnCategory() {

}
