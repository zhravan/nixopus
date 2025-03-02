package service

import (
	"github.com/google/uuid"
	"github.com/raghavyuva/nixopus-api/internal/features/notification"
)

func(s *NotificationService) UpdatePreference(preferences notification.UpdatePreferenceRequest, userID uuid.UUID) error {
	s.storage.UpdatePreference(s.Ctx, preferences, userID)
	return nil
}