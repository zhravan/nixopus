package notification

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/raghavyuva/nixopus-api/internal/features/notification/helpers/discord"
	"github.com/raghavyuva/nixopus-api/internal/features/notification/helpers/email"
	"github.com/raghavyuva/nixopus-api/internal/features/notification/helpers/preferences"
	slackhelper "github.com/raghavyuva/nixopus-api/internal/features/notification/helpers/slack"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
	"github.com/uptrace/bun"
)

func NewNotificationPayload(payloadType NotificationPayloadType, userID string, data interface{}, category NotificationCategory) NotificationPayload {
	return NotificationPayload{
		Type:      payloadType,
		UserID:    userID,
		Timestamp: time.Now(),
		Data:      data,
		Category:  category,
	}
}

func NewNotificationManager(db *bun.DB) *NotificationManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &NotificationManager{
		db:             db,
		ctx:            ctx,
		cancel:         cancel,
		PayloadChan:    make(chan NotificationPayload, 100),
		prefManager:    preferences.NewPreferenceManager(db, ctx),
		emailManager:   email.NewEmailManager(db, ctx),
		slackManager:   slackhelper.NewSlackManager(),
		discordManager: discord.NewDiscordManager(),
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
				switch payload.Category {
				case NotificationCategoryAuthentication:
					fmt.Printf("Authentication Notification - %+v", payload)
					if payload.Type == NotificationPayloadTypePasswordReset {
						fmt.Printf("Password Reset Notification - %+v", payload)
						if data, ok := payload.Data.(NotificationPasswordResetData); ok {
							m.emailManager.SendPasswordResetEmail(payload.UserID, data.Token)
							m.slackManager.SendNotification("Password reset requested")
							m.discordManager.SendNotification("Password reset requested")
						}
					}
					if payload.Type == NotificationPayloadTypeVerificationEmail {
						fmt.Printf("Verification Email Notification - %+v", payload)
						if data, ok := payload.Data.(NotificationVerificationEmailData); ok {
							m.emailManager.SendVerificationEmail(payload.UserID, data.Token)
							m.slackManager.SendNotification("Email verification requested")
							m.discordManager.SendNotification("Email verification requested")
						}
					}
					if payload.Type == NotificationPayloadTypeLogin {
						fmt.Printf("Login Notification - %+v", payload)
						if data, ok := payload.Data.(NotificationAuthenticationData); ok {
							shouldSend, err := m.prefManager.CheckUserNotificationPreferences(payload.UserID, string(NotificationCategoryAuthentication), "security-alerts")
							if err != nil {
								log.Printf("Failed to check notification preferences: %s", err)
							}
							if shouldSend {
								err := m.emailManager.SendEmailWithTemplate(payload.UserID, email.EmailData{
									Subject:  "Login Notification",
									Template: "login_notification.html",
									Data: map[string]interface{}{
										"IP":       data.IP,
										"Browser":  data.Browser,
										"Email":    data.Email,
										"UserName": data.UserName,
									},
									Type:        "security-alerts",
									ContentType: "text/html; charset=UTF-8",
									Category:    string(shared_types.SecurityCategory),
								})
								if err != nil {
									log.Printf("Failed to send login notification email: %s", err)
								}
								m.discordManager.SendNotification(fmt.Sprintf("User %s logged in from %s", data.Email, data.IP))
							}
						}
					}
				case NotificationCategoryOrganization:
					fmt.Printf("Organization Notification - %+v", payload)
					if payload.Type == NotificationPayloadTypeUpdateUserRole {
						fmt.Printf("Update User Role Notification - %+v", payload)
						if data, ok := payload.Data.(NotificationOrganizationData); ok {
							m.emailManager.SendUpdateUserRoleEmail(payload.UserID, data.OrganizationID, data.UserID)
							m.slackManager.SendNotification(fmt.Sprintf("User role updated in organization %s", data.OrganizationID))
							m.discordManager.SendNotification(fmt.Sprintf("User role updated in organization %s", data.OrganizationID))
						}
					}
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
