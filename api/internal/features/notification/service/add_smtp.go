package service

import (
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	"github.com/raghavyuva/nixopus-api/internal/features/notification"
)

func (s *NotificationService) AddSmtp(SMTPConfigs notification.CreateSMTPConfigRequest) error {
	s.logger.Log(logger.Info, "Adding SMTP configuration", "")
	config := notification.NewSMTPConfig(&SMTPConfigs)
	return s.storage.AddSmtp(config)
}
