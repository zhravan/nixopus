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
				switch payload.Category {
				case NotificationCategoryAuthentication:
					fmt.Printf("Authentication Notification - %+v", payload)
					// m.SendEmail(payload.UserID, "login successfully")
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

func (m *NotificationManager) SendEmail(userId string, body string) {
	smtpConfig, err := m.GetSmtp(userId)
	fmt.Println(smtpConfig)
	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
	from := smtpConfig.Username
	pass := smtpConfig.Password
	to := smtpConfig.FromEmail

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		body

	err = smtp.SendMail(smtpConfig.Host+":"+fmt.Sprint(smtpConfig.Port),
		smtp.PlainAuth("", from, pass, smtpConfig.Host),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
	log.Println("Successfully sended to " + to)
}

type ResetEmailData struct {
	ResetURL string
}

func (m *NotificationManager) SendPasswordResetEmail(userID string, token string) {
	smtpConfig, err := m.GetSmtp(userID)
	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Printf("error getting working directory: %s", err)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join(wd, "internal/features/notification/templates/password_reset.html"))
	if err != nil {
		log.Printf("template parsing error: %s", err)
		return
	}

	resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", token)
	data := ResetEmailData{
		ResetURL: resetURL,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Printf("template execution error: %s", err)
		return
	}

	from := smtpConfig.Username
	to := []string{smtpConfig.FromEmail}
	subject := "Password Reset Request"

	msg := []byte(fmt.Sprintf("Subject: %s\r\n"+
		"From: %s\r\n"+
		"To: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", subject, from, smtpConfig.FromEmail, body.String()))

	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)
	addr := fmt.Sprintf("%s:%d", smtpConfig.Host, smtpConfig.Port)

	if err := smtp.SendMail(addr, auth, from, to, msg); err != nil {
		log.Printf("Failed to send password reset email: %s", err)
		return
	}

	log.Printf("Password reset email sent successfully to %s", smtpConfig.FromEmail)
}

func (s *NotificationManager) GetSmtp(ID string) (*shared_types.SMTPConfigs, error) {
	config := &shared_types.SMTPConfigs{}
	err := s.db.NewSelect().Model(config).Where("user_id = ?", ID).Scan(s.ctx)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (m *NotificationManager) SendVerificationEmail(userID string, token string) {
	smtpConfig, err := m.GetSmtp(userID)
	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Printf("error getting working directory: %s", err)
		return
	}
	tmpl, err := template.ParseFiles(filepath.Join(wd, "internal/features/notification/templates/verification_email.html"))
	if err != nil {
		log.Printf("template parsing error: %s", err)
		return
	}

	resetURL := fmt.Sprintf("http://localhost:3000/verify-email?token=%s", token)
	data := ResetEmailData{
		ResetURL: resetURL,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Printf("template execution error: %s", err)
		return
	}

	from := smtpConfig.Username
	to := []string{smtpConfig.FromEmail}
	subject := "Verification Email"

	msg := []byte(fmt.Sprintf("Subject: %s\r\n"+
		"From: %s\r\n"+
		"To: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", subject, from, smtpConfig.FromEmail, body.String()))

	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)
	addr := fmt.Sprintf("%s:%d", smtpConfig.Host, smtpConfig.Port)

	if err := smtp.SendMail(addr, auth, from, to, msg); err != nil {
		log.Printf("Failed to send verification email: %s", err)
		return
	}

	log.Printf("Verification email sent successfully to %s", smtpConfig.FromEmail)
}
