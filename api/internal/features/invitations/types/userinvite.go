package types

import (
	"time"

	"github.com/google/uuid"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

type UserWithInvite struct {
	shared_types.OrganizationUsers
	ExpiresAt   *time.Time `json:"expires_at"`
	AcceptedAt  *time.Time `json:"accepted_at"`
	InvitedBy   *uuid.UUID `json:"invited_by"`
	InviteEmail *string    `json:"invite_email"`
	InviteName  *string    `json:"invite_name"`
	InviteRole  *string    `json:"invite_role"`
}
