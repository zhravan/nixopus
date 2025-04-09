package email

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/raghavyuva/nixopus-api/internal/features/notification/helpers/preferences"
	"github.com/raghavyuva/nixopus-api/internal/types"
	"github.com/uptrace/bun"

	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

type EmailManager struct {
	db          *bun.DB
	ctx         context.Context
	prefManager *preferences.PreferenceManager
}

func NewEmailManager(db *bun.DB, ctx context.Context) *EmailManager {
	return &EmailManager{
		db:          db,
		ctx:         ctx,
		prefManager: preferences.NewPreferenceManager(db, ctx),
	}
}

type EmailData struct {
	Subject     string
	Template    string
	Data        interface{}
	ContentType string
	Category    string
	Type        string
}

type ResetEmailData struct {
	ResetURL string `json:"reset_url"`
}

type VerificationEmailData struct {
	VerifyURL string `json:"verify_url"`
}

type UpdateUserRoleData struct {
	OrganizationID string `json:"organization_id"`
	UserID         string `json:"user_id"`
}

func (m *EmailManager) SendEmailWithTemplate(userID string, data EmailData) error {
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	var user shared_types.User
	err = m.db.NewSelect().
		Model(&user).
		Where("id = ?", uuidUserID).
		Where("deleted_at IS NULL").
		Scan(m.ctx)

	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	// TODO: Implement actual email sending logic here
	log.Printf("Sending email to %s with subject %s", user.Email, data.Subject)
	return nil
}

func (m *EmailManager) SendPasswordResetEmail(userID string, token string) error {
	shouldSend, err := m.prefManager.CheckUserNotificationPreferences(userID, string(types.SecurityCategory), "password-changes")
	if err != nil {
		return fmt.Errorf("failed to check notification preferences: %w", err)
	}

	if !shouldSend {
		return nil
	}

	resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", token)
	data := ResetEmailData{
		ResetURL: resetURL,
	}

	emailData := EmailData{
		Subject:     "Password Reset Request",
		Template:    "password_reset.html",
		Data:        data,
		ContentType: "text/html; charset=UTF-8",
		Category:    string(types.SecurityCategory),
		Type:        "password-changes",
	}

	if err := m.SendEmailWithTemplate(userID, emailData); err != nil {
		log.Printf("Failed to send password reset email: %s", err)
		return err
	}

	log.Printf("Password reset email sent successfully")
	return nil
}

func (m *EmailManager) SendVerificationEmail(userID string, token string) error {
	shouldSend, err := m.prefManager.CheckUserNotificationPreferences(userID, string(types.SecurityCategory), "security-alerts")
	if err != nil {
		return fmt.Errorf("failed to check notification preferences: %w", err)
	}

	if !shouldSend {
		return nil
	}

	verifyURL := fmt.Sprintf("http://localhost:3000/verify-email?token=%s", token)
	data := VerificationEmailData{
		VerifyURL: verifyURL,
	}

	emailData := EmailData{
		Subject:     "Verification Email",
		Template:    "verification_email.html",
		Data:        data,
		ContentType: "text/html; charset=UTF-8",
		Category:    string(types.SecurityCategory),
		Type:        "security-alerts",
	}

	if err := m.SendEmailWithTemplate(userID, emailData); err != nil {
		log.Printf("Failed to send verification email: %s", err)
		return err
	}

	log.Printf("Verification email sent successfully")
	return nil
}

func (m *EmailManager) SendUpdateUserRoleEmail(userID string, organizationID string, updatedUserID string) error {
	shouldSend, err := m.prefManager.CheckUserNotificationPreferences(userID, string(types.ActivityCategory), "team-updates")
	if err != nil {
		return fmt.Errorf("failed to check notification preferences: %w", err)
	}

	if !shouldSend {
		return nil
	}

	data := UpdateUserRoleData{
		OrganizationID: organizationID,
		UserID:         updatedUserID,
	}

	emailData := EmailData{
		Subject:     "User Role Updated",
		Template:    "update_user_role.html",
		Data:        data,
		ContentType: "text/html; charset=UTF-8",
		Category:    string(types.ActivityCategory),
		Type:        "team-updates",
	}

	if err := m.SendEmailWithTemplate(userID, emailData); err != nil {
		log.Printf("Failed to send update user role email: %s", err)
		return err
	}

	log.Printf("Update user role email sent successfully")
	return nil
}
