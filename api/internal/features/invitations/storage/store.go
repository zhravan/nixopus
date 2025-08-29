package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
	"github.com/uptrace/bun"
)

type InvitationStore struct {
	DB  *bun.DB
	Ctx context.Context
	tx  *bun.Tx
}

func (s *InvitationStore) WithTx(tx bun.Tx) *InvitationStore {
	return &InvitationStore{DB: s.DB, Ctx: s.Ctx, tx: &tx}
}

func (s *InvitationStore) getDB() bun.IDB {
	if s.tx != nil {
		return *s.tx
	}
	return s.DB
}

func (s *InvitationStore) BeginTx() (bun.Tx, error) {
	return s.DB.BeginTx(s.Ctx, nil)
}

func (s *InvitationStore) CreateInvitation(inv *shared_types.Invitation) error {
	_, err := s.getDB().NewInsert().Model(inv).Exec(s.Ctx)
	return err
}

// GetInvitationByUserAndOrg returns an invitation matching the given user and organization.
// Returns (nil, nil) when not found.
func (s *InvitationStore) GetInvitationByUserAndOrg(userID uuid.UUID, orgID uuid.UUID) (*shared_types.Invitation, error) {
	inv := &shared_types.Invitation{}
	err := s.getDB().NewSelect().Model(inv).
		Where("user_id = ? AND organization_id = ?", userID, orgID).
		Scan(s.Ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return inv, err
}

// UpdateInvitationForReinvite updates an invitation's token/expiry and metadata for reinvite flow.
func (s *InvitationStore) UpdateInvitationForReinvite(invID uuid.UUID, token string, expiresAt time.Time, name string, role string) error {
	now := time.Now()
	_, err := s.getDB().NewUpdate().Model(&shared_types.Invitation{}).
		Set("token = ?", token).
		Set("expires_at = ?", expiresAt).
		Set("updated_at = ?", now).
		Set("name = ?", name).
		Set("role = ?", role).
		Where("id = ?", invID).
		Exec(s.Ctx)
	return err
}

func (s *InvitationStore) GetInvitationByToken(token string) (*shared_types.Invitation, error) {
	inv := &shared_types.Invitation{}
	err := s.getDB().NewSelect().Model(inv).Where("token = ?", token).Scan(s.Ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return inv, err
}

func (s *InvitationStore) MarkAccepted(invID uuid.UUID) error {
	now := time.Now()
	_, err := s.getDB().NewUpdate().Model(&shared_types.Invitation{}).Set("accepted_at = ?", now).Where("id = ?", invID).Exec(s.Ctx)
	return err
}

func (s *InvitationStore) AddUserToOrganization(userID uuid.UUID, orgID uuid.UUID, roleID uuid.UUID) error {
	ou := shared_types.OrganizationUsers{
		ID:             uuid.New(),
		UserID:         userID,
		OrganizationID: orgID,
		RoleID:         roleID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	_, err := s.getDB().NewInsert().Model(&ou).Exec(s.Ctx)
	return err
}

func (s *InvitationStore) ListInvitationsByOrganization(orgID string) ([]shared_types.Invitation, error) {
	var invs []shared_types.Invitation
	err := s.getDB().NewSelect().Model(&invs).Where("organization_id = ?", orgID).Scan(s.Ctx)
	if err == sql.ErrNoRows {
		return []shared_types.Invitation{}, nil
	}
	return invs, err
}

func (s *InvitationStore) GetLatestInvitationsMapByOrganization(orgID string) (map[uuid.UUID]*shared_types.Invitation, error) {
	var invs []shared_types.Invitation
	err := s.getDB().NewSelect().
		Model(&invs).
		Where("organization_id = ?", orgID).
		OrderExpr("updated_at DESC").
		Scan(s.Ctx)
	if err == sql.ErrNoRows {
		return map[uuid.UUID]*shared_types.Invitation{}, nil
	}
	if err != nil {
		return nil, err
	}
	out := make(map[uuid.UUID]*shared_types.Invitation, len(invs))
	for i := range invs {
		inv := &invs[i]
		if _, exists := out[inv.UserID]; !exists {
			out[inv.UserID] = inv
		}
	}
	return out, nil
}
