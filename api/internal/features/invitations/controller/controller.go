package controller

import (
	"net/http"
	"time"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
	inv_service "github.com/raghavyuva/nixopus-api/internal/features/invitations/service"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	"github.com/raghavyuva/nixopus-api/internal/features/notification"
	org_service "github.com/raghavyuva/nixopus-api/internal/features/organization/service"
	"github.com/raghavyuva/nixopus-api/internal/storage"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
	"github.com/raghavyuva/nixopus-api/internal/utils"
)

type Controller struct {
	svc           *inv_service.Service
	orgs          *org_service.OrganizationService
	logger        logger.Logger
	notifications *notification.NotificationManager
}

func NewController(s *storage.Store, svc *inv_service.Service, orgs *org_service.OrganizationService, l logger.Logger, n *notification.NotificationManager) *Controller {
	return &Controller{svc: svc, orgs: orgs, logger: l, notifications: n}
}

// CreateInviteRequest wires to service
type CreateInviteRequest struct {
	Email          string `json:"email"`
	Name           string `json:"name"`
	Role           string `json:"role"`
	OrganizationID string `json:"organization_id"`
}

func (c *Controller) CreateInvite(ctx fuego.ContextWithBody[CreateInviteRequest]) (*shared_types.Response, error) {
	w, r := ctx.Response(), ctx.Request()
	inviter := utils.GetUser(w, r)
	if inviter == nil {
		return nil, fuego.HTTPError{Status: http.StatusUnauthorized}
	}
	// Require organization from context (set by auth middleware)
	orgID := utils.GetOrganizationID(r)
	if orgID == uuid.Nil {
		return nil, fuego.HTTPError{Status: http.StatusBadRequest, Detail: "missing X-Organization-Id"}
	}
	req, err := ctx.Body()
	if err != nil {
		return nil, fuego.HTTPError{Status: http.StatusBadRequest, Err: err}
	}
	// If client provided organization_id in body, it must match header/context
	if req.OrganizationID != "" && req.OrganizationID != orgID.String() {
		return nil, fuego.HTTPError{Status: http.StatusBadRequest, Detail: "organization_id must match X-Organization-Id"}
	}
	// Force organization from header/context
	req.OrganizationID = orgID.String()
	inv, _, err := c.svc.CreateInvite(inviter.ID.String(), inv_service.CreateInviteRequest(req))
	if err != nil {
		return nil, fuego.HTTPError{Status: http.StatusBadRequest, Err: err}
	}
	// notify inviter that invite has been sent (optional)
	c.notifications.SendNotification(notification.NewNotificationPayload(notification.NotificationPayloadTypeUpdateOrganization, inviter.ID.String(), map[string]string{"message": "Invite sent"}, notification.NotificationCategoryOrganization))
	return &shared_types.Response{Status: "success", Message: "Invitation created", Data: inv}, nil
}

func (c *Controller) AcceptInvite(ctx fuego.ContextNoBody) (*shared_types.Response, error) {
	token := ctx.QueryParam("token")
	if token == "" {
		return nil, fuego.HTTPError{Status: http.StatusBadRequest, Detail: "token is required"}
	}
	res, err := c.svc.AcceptInvite(token)
	if err != nil {
		return nil, fuego.HTTPError{Status: http.StatusBadRequest, Err: err}
	}
	return &shared_types.Response{Status: "success", Message: "Invitation accepted", Data: res}, nil
}

// GetOrganizationUsersWithInviteStatus returns users in org and pending invites
func (c *Controller) GetOrganizationUsersWithInviteStatus(ctx fuego.ContextNoBody) (*shared_types.Response, error) {
	id := ctx.QueryParam("id")
	if id == "" {
		return nil, fuego.HTTPError{Status: http.StatusBadRequest}
	}
	users, err := c.orgs.GetOrganizationUsers(id)
	if err != nil {
		return nil, fuego.HTTPError{Status: http.StatusInternalServerError, Err: err}
	}
	// map latest invitation per user to enrich the users payload
	latestByUser, err := c.svc.Invitations.GetLatestInvitationsMapByOrganization(id)
	if err != nil {
		return nil, fuego.HTTPError{Status: http.StatusInternalServerError, Err: err}
	}

	type UserWithInvite struct {
		shared_types.OrganizationUsers
		ExpiresAt   *time.Time `json:"expires_at"`
		AcceptedAt  *time.Time `json:"accepted_at"`
		InvitedBy   *uuid.UUID `json:"invited_by"`
		InviteEmail *string    `json:"invite_email"`
		InviteName  *string    `json:"invite_name"`
		InviteRole  *string    `json:"invite_role"`
	}

	enriched := make([]UserWithInvite, 0, len(users))
	for _, u := range users {
		row := UserWithInvite{OrganizationUsers: u}
		if inv, ok := latestByUser[u.UserID]; ok {
			if !inv.ExpiresAt.IsZero() {
				t := inv.ExpiresAt
				row.ExpiresAt = &t
			}
			if inv.AcceptedAt != nil {
				row.AcceptedAt = inv.AcceptedAt
			}
			if inv.InviterUserID != uuid.Nil {
				id := inv.InviterUserID
				row.InvitedBy = &id
			}
			email := inv.Email
			name := inv.Name
			role := inv.Role
			row.InviteEmail = &email
			row.InviteName = &name
			row.InviteRole = &role
		}
		enriched = append(enriched, row)
	}

	return &shared_types.Response{Status: "success", Message: "Fetched Org Users with invite statuses", Data: enriched}, nil
}
