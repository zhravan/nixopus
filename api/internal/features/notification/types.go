package notification

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
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

type NotificationBaseData struct {
	IP      string
	Browser string
}

type NotificationAuthenticationData struct {
	NotificationBaseData
	Email    string
	UserName string
}

type NotificationOrganizationData struct {
	NotificationBaseData
	OrganizationID string
	UserID         string
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
	NotificationPayloadTypeRegister               NotificationPayloadType = "register"
	NotificationPayloadTypeLogin                  NotificationPayloadType = "login"
	NotificationPayloadTypeLogout                 NotificationPayloadType = "logout"
	NotificationPayloadTypePasswordReset          NotificationPayloadType = "password_reset"
	NortificationPayloadTypeAddUserToOrganization NotificationPayloadType = "add_user_to_organization"
)

const (
	NortificationPayloadTypeCreateOrganization         NotificationPayloadType = "create_organization"
	NortificationPayloadTypeRemoveUserFromOrganization NotificationPayloadType = "remove_user_from_organization"
	NotificationPayloadTypeDeleteOrganization          NotificationPayloadType = "delete_organization"
	NotificationPayloadTypeUpdateOrganization          NotificationPayloadType = "update_organization"
)

type NotificationCategory string

const (
	NotificationCategoryAuthentication NotificationCategory = "authentication"
	NotificationCategoryOrganization   NotificationCategory = "organization"
)

type NotificationPayload struct {
	Category  NotificationCategory    `json:"category"`
	Type      NotificationPayloadType `json:"type"`
	UserID    string                  `json:"user_id"`
	Timestamp time.Time               `json:"timestamp"`
	Data      interface{}             `json:"data"`
}

type CreateSMTPConfigRequest struct {
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	FromName  string    `json:"from_name"`
	FromEmail string    `json:"from_email"`
	UserID    uuid.UUID `json:"user_id"`
}

type UpdateSMTPConfigRequest struct {
	ID        uuid.UUID `json:"id"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	FromName  string    `json:"from_name"`
	FromEmail string    `json:"from_email"`
}

type DeleteSMTPConfigRequest struct {
	ID uuid.UUID `json:"id"`
}

type GetSMTPConfigRequest struct {
	ID uuid.UUID `json:"id"`
}

var (
	ErrInvalidRequestType = errors.New("invalid request type")
	ErrMissingHost        = errors.New("host is required")
	ErrMissingPort        = errors.New("port is required")
	ErrMissingUsername    = errors.New("username is required")
	ErrMissingPassword    = errors.New("password is required")
	ErrMissingID          = errors.New("id is required")
)

func NewSMTPConfig(c *CreateSMTPConfigRequest) *shared_types.SMTPConfigs {
	if c.FromEmail == "" {
		c.FromEmail = c.Username
	}

	if c.FromName == "" {
		c.FromName = strings.Split(c.Username, "@")[0]
	}
	return &shared_types.SMTPConfigs{
		Host:      c.Host,
		Port:      c.Port,
		Username:  c.Username,
		Password:  c.Password,
		FromName:  c.FromName,
		FromEmail: c.FromEmail,
		UserID:    c.UserID,
	}
}
