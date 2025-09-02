package controller

import (
    "net/http"
    "strings"

    "github.com/go-fuego/fuego"
    "github.com/google/uuid"
    "github.com/raghavyuva/nixopus-api/internal/config"
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
    // On successful acceptance, redirect to hosted frontend URL instead of returning JSON.
    // Uses CORS.AllowedOrigin as the public frontend origin.
    frontend := strings.TrimSpace(config.AppConfig.CORS.AllowedOrigin)
    if frontend != "" {
        frontend = strings.TrimRight(frontend, "/")
        http.Redirect(ctx.Response(), ctx.Request(), frontend+"/?invite=accepted", http.StatusFound)
        return nil, nil
    }
    // Fallback to JSON if frontend origin is not configured
    return &shared_types.Response{Status: "success", Message: "Invitation accepted", Data: res}, nil
}

// GetOrganizationUsersWithInviteStatus returns users in org and pending invites
func (c *Controller) GetOrganizationUsersWithInviteStatus(ctx fuego.ContextNoBody) (*shared_types.Response, error) {
	id := ctx.QueryParam("id")
	if id == "" {
		return nil, fuego.HTTPError{Status: http.StatusBadRequest}
	}
	enriched, err := c.svc.GetOrganizationUsersWithInviteStatus(id)
	if err != nil {
		return nil, fuego.HTTPError{Status: http.StatusInternalServerError, Err: err}
	}
	return &shared_types.Response{Status: "success", Message: "Fetched Org Users with invite statuses", Data: enriched}, nil
}
