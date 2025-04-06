package notification

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
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
				shouldSend, err := m.CheckUserNotificationPreferences(
					payload.UserID,
					NotificationCategory(payload.Category),
					NotificationPayloadType(payload.Type),
				)
				if err != nil {
					log.Printf("Failed to check notification preferences: %v", err)
					continue
				}

				if !shouldSend {
					log.Printf("Notification skipped due to user preferences: %+v", payload)
					continue
				}

				switch payload.Category {
				case NotificationCategoryAuthentication:
					fmt.Printf("Authentication Notification - %+v", payload)
					if payload.Type == NotificationPayloadTypePasswordReset {
						fmt.Printf("Password Reset Notification - %+v", payload)
						if data, ok := payload.Data.(NotificationPasswordResetData); ok {
							m.SendPasswordResetEmail(payload.UserID, data.Token)
						}
					}
					if payload.Type == NotificationPayloadTypeVerificationEmail {
						fmt.Printf("Verification Email Notification - %+v", payload)
						if data, ok := payload.Data.(NotificationVerificationEmailData); ok {
							m.SendVerificationEmail(payload.UserID, data.Token)
						}
					}
				case NotificationCategoryOrganization:
					fmt.Printf("Organization Notification - %+v", payload)
					if payload.Type == NotificationPayloadTypeUpdateUserRole {
						fmt.Printf("Update User Role Notification - %+v", payload)
						if data, ok := payload.Data.(NotificationOrganizationData); ok {
							m.SendUpdateUserRoleEmail(payload.UserID, data.OrganizationID, data.UserID)
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

// here we can get the notification preferences of the user (like should send to slack/email/discord, how many times to send, what type of contents to send)
func (m *NotificationManager) CheckUserNotificationPreferences(userID string, category NotificationCategory, notificationType NotificationPayloadType) (bool, error) {
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		return false, fmt.Errorf("invalid user ID: %w", err)
	}

	var preferenceID uuid.UUID
	err = m.db.NewSelect().
		Model((*shared_types.NotificationPreferences)(nil)).
		Column("id").
		Where("user_id = ?", uuidUserID).
		Where("deleted_at IS NULL").
		Scan(m.ctx, &preferenceID)

	if err != nil {
		return false, fmt.Errorf("failed to fetch user preferences: %w", err)
	}

	var storageCategory string
	switch category {
	case NotificationCategoryAuthentication:
		storageCategory = "security"
	case NotificationCategoryOrganization:
		storageCategory = "activity"
	default:
		return false, fmt.Errorf("unsupported notification category: %s", category)
	}

	var storageType string
	switch notificationType {
	case NotificationPayloadTypePasswordReset:
		storageType = "password-changes"
	case NotificationPayloadTypeVerificationEmail:
		storageType = "security-alerts"
	case NotificationPayloadTypeUpdateUserRole:
		storageType = "team-updates"
	default:
		return false, fmt.Errorf("unsupported notification type: %s", notificationType)
	}

	var preferenceItem shared_types.PreferenceItem
	err = m.db.NewSelect().
		Model(&preferenceItem).
		Where("preference_id = ?", preferenceID).
		Where("category = ?", storageCategory).
		Where("type = ?", storageType).
		Scan(m.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to fetch preference item: %w", err)
	}

	return preferenceItem.Enabled, nil
}

func (m *NotificationManager) shouldSendEmail(userID string, category string, notificationType string) (bool, error) {
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		return false, fmt.Errorf("invalid user ID: %w", err)
	}

	var preferenceID uuid.UUID
	err = m.db.NewSelect().
		Model((*shared_types.NotificationPreferences)(nil)).
		Column("id").
		Where("user_id = ?", uuidUserID).
		Where("deleted_at IS NULL").
		Scan(m.ctx, &preferenceID)

	if err != nil {
		return false, fmt.Errorf("failed to fetch user preferences: %w", err)
	}

	var preferenceItem shared_types.PreferenceItem
	err = m.db.NewSelect().
		Model(&preferenceItem).
		Where("preference_id = ?", preferenceID).
		Where("category = ?", category).
		Where("type = ?", notificationType).
		Scan(m.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to fetch preference item: %w", err)
	}

	return preferenceItem.Enabled, nil
}

func (m *NotificationManager) SendEmailWithTemplate(userID string, emailData EmailData) error {
	shouldSend, err := m.shouldSendEmail(userID, emailData.Category, emailData.Type)
	if err != nil {
		return fmt.Errorf("failed to check email preferences: %w", err)
	}

	if !shouldSend {
		return nil
	}

	smtpConfig, err := m.GetSmtp(userID)
	if err != nil {
		return fmt.Errorf("smtp error: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %w", err)
	}

	tmpl, err := template.ParseFiles(filepath.Join(wd, "internal/features/notification/templates/"+emailData.Template))
	if err != nil {
		return fmt.Errorf("template parsing error: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, emailData.Data); err != nil {
		return fmt.Errorf("template execution error: %w", err)
	}

	from := smtpConfig.Username
	to := []string{smtpConfig.FromEmail}

	msg := []byte(fmt.Sprintf("Subject: %s\r\n"+
		"From: %s\r\n"+
		"To: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: %s\r\n"+
		"\r\n"+
		"%s", emailData.Subject, from, smtpConfig.FromEmail, emailData.ContentType, body.String()))

	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)
	addr := fmt.Sprintf("%s:%d", smtpConfig.Host, smtpConfig.Port)

	if err := smtp.SendMail(addr, auth, from, to, msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (m *NotificationManager) SendPasswordResetEmail(userID string, token string) error {
	resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", token)
	data := ResetEmailData{
		ResetURL: resetURL,
	}

	emailData := EmailData{
		Subject:     "Password Reset Request",
		Template:    "password_reset.html",
		Data:        data,
		ContentType: "text/html; charset=UTF-8",
		Category:    string(shared_types.SecurityCategory),
		Type:        "password-changes",
	}

	if err := m.SendEmailWithTemplate(userID, emailData); err != nil {
		log.Printf("Failed to send password reset email: %s", err)
		return err
	}

	log.Printf("Password reset email sent successfully")
	return nil
}

func (m *NotificationManager) SendVerificationEmail(userID string, token string) error {
	resetURL := fmt.Sprintf("http://localhost:3000/verify-email?token=%s", token)
	data := ResetEmailData{
		ResetURL: resetURL,
	}

	emailData := EmailData{
		Subject:     "Verification Email",
		Template:    "verification_email.html",
		Data:        data,
		ContentType: "text/html; charset=UTF-8",
		Category:    string(shared_types.SecurityCategory),
		Type:        "security-alerts",
	}

	if err := m.SendEmailWithTemplate(userID, emailData); err != nil {
		log.Printf("Failed to send verification email: %s", err)
		return err
	}

	log.Printf("Verification email sent successfully")
	return nil
}

func (m *NotificationManager) SendUpdateUserRoleEmail(userID string, organizationID string, updatedUserID string) error {
	data := UpdateUserRoleData{
		OrganizationID: organizationID,
		UserID:         updatedUserID,
	}

	emailData := EmailData{
		Subject:     "User Role Updated",
		Template:    "update_user_role.html",
		Data:        data,
		ContentType: "text/html; charset=UTF-8",
		Category:    string(shared_types.ActivityCategory),
		Type:        "team-updates",
	}

	if err := m.SendEmailWithTemplate(userID, emailData); err != nil {
		log.Printf("Failed to send update user role email: %s", err)
		return err
	}

	log.Printf("Update user role email sent successfully")
	return nil
}

func (m *NotificationManager) GetSmtp(ID string) (*shared_types.SMTPConfigs, error) {
	config := &shared_types.SMTPConfigs{}
	err := m.db.NewSelect().Model(config).Where("user_id = ?", ID).Scan(m.ctx)
	if err != nil {
		return nil, err
	}
	return config, nil
}
