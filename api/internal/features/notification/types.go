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
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
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

func NewSMTPConfig(c *CreateSMTPConfigRequest, userID uuid.UUID) *shared_types.SMTPConfigs {
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
		UserID:    userID,
		ID:        uuid.New(),
	}
}

type Category string

const (
	ActivityCategory Category = "activity"
	SecurityCategory Category = "security"
	UpdateCategory   Category = "update"
)

type PreferenceType struct {
	ID          string `json:"id" validate:"required"`
	Label       string `json:"label" validate:"required"`
	Description string `json:"description" validate:"required"`
	Enabled     bool   `json:"enabled"`
}

type CategoryPreferences struct {
	Category    Category         `json:"category" validate:"required"`
	Preferences []PreferenceType `json:"preferences" validate:"required"`
}

type UpdatePreferenceRequest struct {
	Category string `json:"category" validate:"required,oneof=activity security update"`
	Type     string `json:"type" validate:"required"`
	Enabled  bool   `json:"enabled"`
}

type GetPreferencesResponse struct {
	Activity []PreferenceType `json:"activity"`
	Security []PreferenceType `json:"security"`
	Update   []PreferenceType `json:"update"`
}

type PreferenceItem struct {
	ID           uuid.UUID `json:"id"`
	PreferenceID uuid.UUID `json:"preference_id"`
	Category     string    `json:"category"`
	Type         string    `json:"type"`
	Enabled      bool      `json:"enabled"`
}

func MapToResponse(items []PreferenceItem) GetPreferencesResponse {
	response := GetPreferencesResponse{
		Activity: []PreferenceType{},
		Security: []PreferenceType{},
		Update:   []PreferenceType{},
	}

	typeInfo := map[string]map[string]struct {
		Label       string
		Description string
	}{
		"activity": {
			"team-updates": {
				Label:       "Team Updates",
				Description: "When team members join or leave your team",
			},
		},
		"security": {
			"login-alerts": {
				Label:       "Login Alerts",
				Description: "When a new device logs into your account",
			},
			"password-changes": {
				Label:       "Password Changes",
				Description: "When your password is changed",
			},
			"security-alerts": {
				Label:       "Security Alerts",
				Description: "Important security notifications",
			},
		},
		"update": {
			"product-updates": {
				Label:       "Product Updates",
				Description: "New features and improvements",
			},
			"newsletter": {
				Label:       "Newsletter",
				Description: "Our monthly newsletter with tips and updates",
			},
			"marketing": {
				Label:       "Marketing",
				Description: "Promotions and special offers",
			},
		},
	}

	for _, item := range items {
		info, exists := typeInfo[item.Category][item.Type]
		if !exists {
			continue
		}

		pref := PreferenceType{
			ID:          item.Type,
			Label:       info.Label,
			Description: info.Description,
			Enabled:     item.Enabled,
		}

		switch item.Category {
		case "activity":
			response.Activity = append(response.Activity, pref)
		case "security":
			response.Security = append(response.Security, pref)
		case "update":
			response.Update = append(response.Update, pref)
		}
	}

	return response
}
