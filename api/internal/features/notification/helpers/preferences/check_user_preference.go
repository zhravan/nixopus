package preferences

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
	"github.com/uptrace/bun"
)

type PreferenceManager struct {
	db  *bun.DB
	ctx context.Context
}

func NewPreferenceManager(db *bun.DB, ctx context.Context) *PreferenceManager {
	return &PreferenceManager{
		db:  db,
		ctx: ctx,
	}
}

func (m *PreferenceManager) CheckUserNotificationPreferences(userID string, category string, notificationType string) (bool, error) {
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
	case "security":
		storageCategory = "security"
	case "activity":
		storageCategory = "activity"
	default:
		return false, fmt.Errorf("unsupported notification category: %s", category)
	}

	var preferenceItem shared_types.PreferenceItem
	err = m.db.NewSelect().
		Model(&preferenceItem).
		Where("preference_id = ?", preferenceID).
		Where("category = ?", storageCategory).
		Where("type = ?", notificationType).
		Scan(m.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to fetch preference item: %w", err)
	}

	return preferenceItem.Enabled, nil
}
