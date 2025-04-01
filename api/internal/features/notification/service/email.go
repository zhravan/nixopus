package service

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

type EmailService struct {
	logger logger.Logger
}

func NewEmailService(logger logger.Logger) *EmailService {
	return &EmailService{
		logger: logger,
	}
}

func (s *EmailService) SendPasswordResetEmail(user *shared_types.User, resetToken string) error {
	s.logger.Log(logger.Info, "sending password reset email", user.Email)

	// Get SMTP configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")
	fromName := os.Getenv("SMTP_FROM_NAME")
	frontendURL := os.Getenv("FRONTEND_URL")

	if smtpHost == "" || smtpPort == "" || smtpUsername == "" || smtpPassword == "" || fromEmail == "" || fromName == "" || frontendURL == "" {
		s.logger.Log(logger.Error, "missing SMTP configuration", "")
		return fmt.Errorf("missing SMTP configuration")
	}

	// Construct the reset link
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, resetToken)

	// Email template
	subject := "Password Reset Request"
	body := fmt.Sprintf(`Hello %s,

We received a request to reset your password. Click the link below to set a new password:

%s

This link will expire in 5 minutes.

If you didn't request this password reset, you can safely ignore this email.

Best regards,
%s Team`, user.Username, resetLink, fromName)

	// Prepare email headers
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", fromName, fromEmail)
	headers["To"] = user.Email
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=\"utf-8\""

	// Build email message
	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + body

	// Connect to SMTP server
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		fromEmail,
		[]string{user.Email},
		[]byte(message),
	)

	if err != nil {
		s.logger.Log(logger.Error, "failed to send email", err.Error())
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Log(logger.Info, "password reset email sent successfully", user.Email)
	return nil
}
