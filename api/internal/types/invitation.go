package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Invitation represents an invitation to join an organization/team.
type Invitation struct {
	bun.BaseModel  `bun:"table:invitations,alias:inv"`
	ID             uuid.UUID  `json:"id" bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Email          string     `json:"email" bun:"email,type:text,notnull"`
	Name           string     `json:"name" bun:"name,type:text"`
	Role           string     `json:"role" bun:"role,type:text,notnull"`
	Token          string     `json:"token" bun:"token,type:text,notnull,unique"`
	ExpiresAt      time.Time  `json:"expires_at" bun:"expires_at,type:timestamp,notnull"`
	AcceptedAt     *time.Time `json:"accepted_at,omitempty" bun:"accepted_at"`
	CreatedAt      time.Time  `json:"created_at" bun:"created_at,type:timestamp,notnull,default:now()"`
	UpdatedAt      time.Time  `json:"updated_at" bun:"updated_at,type:timestamp,notnull,default:now()"`
	InviterUserID  uuid.UUID  `json:"inviter_user_id" bun:"inviter_user_id,type:uuid,notnull"`
	OrganizationID uuid.UUID  `json:"organization_id" bun:"organization_id,type:uuid,notnull"`
	UserID         uuid.UUID  `json:"user_id" bun:"user_id,type:uuid,notnull"`
}
