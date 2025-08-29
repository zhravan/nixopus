package service

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/raghavyuva/nixopus-api/internal/config"
	auth_storage "github.com/raghavyuva/nixopus-api/internal/features/auth/storage"
	auth_utils "github.com/raghavyuva/nixopus-api/internal/features/auth/utils"
	inv_store "github.com/raghavyuva/nixopus-api/internal/features/invitations/storage"
	inv_types "github.com/raghavyuva/nixopus-api/internal/features/invitations/types"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	emailhelper "github.com/raghavyuva/nixopus-api/internal/features/notification/helpers/email"
	org_service "github.com/raghavyuva/nixopus-api/internal/features/organization/service"
	org_types "github.com/raghavyuva/nixopus-api/internal/features/organization/types"
	role_service "github.com/raghavyuva/nixopus-api/internal/features/role/service"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
	"github.com/uptrace/bun"
)

type Service struct {
	Invitations *inv_store.InvitationStore
	Users       auth_storage.AuthRepository
	Roles       *role_service.RoleService
	Orgs        *org_service.OrganizationService
	Email       *emailhelper.EmailManager
	Logger      logger.Logger
	DB          *bun.DB
}

type CreateInviteRequest struct {
	Email          string `json:"email"`
	Name           string `json:"name"`
	Role           string `json:"role"`
	OrganizationID string `json:"organization_id"`
}

type AcceptInviteResponse struct {
	Status string `json:"status"`
}

func randPassword() (string, error) {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (s *Service) CreateInvite(inviterID string, req CreateInviteRequest) (*shared_types.Invitation, string, error) {
	if req.Email == "" || req.Role == "" || req.OrganizationID == "" {
		return nil, "", errors.New("missing required fields")
	}

	roleName := strings.ToLower(req.Role)
	role, err := s.Roles.GetRoleByName(roleName)
	if err != nil || role == nil {
		return nil, "", fmt.Errorf("invalid role")
	}

	// create user with generated password if not exists
	dbUser, findErr := s.Users.FindUserByEmail(req.Email)
	if findErr != nil && !errors.Is(findErr, sql.ErrNoRows) {
		s.Logger.Log(logger.Error, "user lookup failed", fmt.Sprintf("email=%s err=%v", req.Email, findErr))
		return nil, "", fmt.Errorf("failed to lookup user: %w", findErr)
	}
	var user shared_types.User
	var generatedPassword string
	if dbUser == nil || dbUser.ID == uuid.Nil {
		s.Logger.Log(logger.Info, "invitation user branch", fmt.Sprintf("email=%s action=create-new-user", req.Email))
		gen, perr := randPassword()
		if perr != nil {
			s.Logger.Log(logger.Error, "password generation failed", perr.Error())
			return nil, "", fmt.Errorf("failed to generate password: %w", perr)
		}
		generatedPassword = gen
		hashed, err := auth_utils.HashPassword(generatedPassword)
		if err != nil {
			return nil, "", err
		}
		user = shared_types.NewUser(req.Email, hashed, req.Email, "", roleName, false)
		if err := s.Users.CreateUser(&user); err != nil {
			return nil, "", err
		}
	} else {
		s.Logger.Log(logger.Info, "invitation user branch", fmt.Sprintf("email=%s action=existing-user", req.Email))
		user = *dbUser
		// reinvite flow: rotate a new password and update the user so email carries a fresh password
		gen, perr := randPassword()
		if perr != nil {
			s.Logger.Log(logger.Error, "password generation failed", perr.Error())
			return nil, "", fmt.Errorf("failed to generate password: %w", perr)
		}
		generatedPassword = gen
		hashed, err := auth_utils.HashPassword(generatedPassword)
		if err != nil {
			return nil, "", err
		}
		user.Password = hashed
		user.UpdatedAt = time.Now()
		if err := s.Users.UpdateUser(&user); err != nil {
			return nil, "", err
		}
	}

	// ensure only one invitation per (user_id, organization_id)
	token := uuid.New().String()
	orgID := uuid.MustParse(req.OrganizationID)
	existingInv, err := s.Invitations.GetInvitationByUserAndOrg(user.ID, orgID)
	if err != nil {
		return nil, "", err
	}
	var inv *shared_types.Invitation
	if existingInv != nil && existingInv.ID != uuid.Nil {
		// Update token/expiry/name/role for reinvite
		if err := s.Invitations.UpdateInvitationForReinvite(existingInv.ID, token, time.Now().Add(72*time.Hour), req.Name, roleName); err != nil {
			return nil, "", err
		}
		// reload updated invitation (optional) or reuse with updated fields
		existingInv.Token = token
		existingInv.ExpiresAt = time.Now().Add(72 * time.Hour)
		existingInv.Name = req.Name
		existingInv.Role = roleName
		existingInv.UpdatedAt = time.Now()
		inv = existingInv
	} else {
		// create new invitation
		inv = &shared_types.Invitation{
			ID:             uuid.New(),
			Email:          req.Email,
			Name:           req.Name,
			Role:           roleName,
			Token:          token,
			ExpiresAt:      time.Now().Add(72 * time.Hour),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			InviterUserID:  uuid.MustParse(inviterID),
			OrganizationID: orgID,
			UserID:         user.ID,
		}
		if err := s.Invitations.CreateInvitation(inv); err != nil {
			return nil, "", err
		}
	}

	// send invitation email to the invitee's email with link and generated password
	base := strings.TrimSpace(config.AppConfig.App.APIURL)
	base = strings.TrimRight(base, "/")
	path := "/api/v1/invitations/accept"
	acceptURL := fmt.Sprintf("%s%s?token=%s", base, path, token)
	data := struct {
		Name      string
		Email     string
		Password  string
		AcceptURL string
	}{Name: req.Name, Email: req.Email, Password: generatedPassword, AcceptURL: acceptURL}

	s.Logger.Log(logger.Info, "invite email debug", fmt.Sprintf("to=%s generatedPassword=%q empty=%t", req.Email, generatedPassword, generatedPassword == ""))

	s.Logger.Log(logger.Info, "sending invite email", fmt.Sprintf("to=%s org=%s", req.Email, req.OrganizationID))
	if err := s.Email.SendEmailToAddress(inviterID, req.Email, emailhelper.EmailData{
		Subject:     "You're invited to Nixopus",
		Template:    "invitation_email.html",
		Data:        data,
		ContentType: "text/html; charset=UTF-8",
		Category:    string(shared_types.SecurityCategory),
		Type:        "security-alerts",
	}); err != nil {
		s.Logger.Log(logger.Error, "invite email failed", err.Error())
	} else {
		s.Logger.Log(logger.Info, "invite email sent", req.Email)
	}

	return inv, generatedPassword, nil
}

func (s *Service) AcceptInvite(token string) (*AcceptInviteResponse, error) {
	inv, err := s.Invitations.GetInvitationByToken(token)
	if err != nil || inv == nil {
		return nil, errors.New("invalid invitation")
	}
	if time.Now().After(inv.ExpiresAt) {
		return nil, errors.New("invitation expired")
	}
	// mark user verified
	if err := s.Users.UpdateUserEmailVerification(inv.UserID.String(), true); err != nil {
		return nil, err
	}
	// add user to org with role
	role, err := s.Roles.GetRoleByName(inv.Role)
	if err != nil || role == nil {
		return nil, fmt.Errorf("invalid role on invite")
	}

	addReq := org_types.AddUserToOrganizationRequest{
		UserID:         inv.UserID.String(),
		OrganizationID: inv.OrganizationID.String(),
		RoleId:         role.ID.String(),
	}
	if s.Orgs != nil {
		_ = s.Orgs.AddUserToOrganization(addReq)
	} else {
		_ = s.Invitations.AddUserToOrganization(inv.UserID, inv.OrganizationID, role.ID)
	}
	if err := s.Invitations.MarkAccepted(inv.ID); err != nil {
		return nil, err
	}

	inviter, _ := s.Users.FindUserByID(inv.InviterUserID.String())
	if inviter != nil {
		_ = s.Email.SendEmailToAddress(inviter.ID.String(), inviter.Email, emailhelper.EmailData{
			Subject:     "Invitation accepted",
			Template:    "invitation_accepted.html",
			Data:        struct{ Name, Email string }{Name: inv.Name, Email: inv.Email},
			ContentType: "text/html; charset=UTF-8",
			Category:    string(shared_types.ActivityCategory),
			Type:        "team-updates",
		})
	}
	return &AcceptInviteResponse{Status: "accepted"}, nil
}

// GetOrganizationUsersWithInviteStatus returns org members wth invite status
func (s *Service) GetOrganizationUsersWithInviteStatus(orgID string) ([]inv_types.UserWithInvite, error) {
	if s.Orgs == nil || s.Invitations == nil {
		return nil, fmt.Errorf("service not initialized: missing Orgs or Invitations store")
	}

	users, err := s.Orgs.GetOrganizationUsers(orgID)
	if err != nil {
		return nil, err
	}

	latestByUser, err := s.Invitations.GetLatestInvitationsMapByOrganization(orgID)
	if err != nil {
		return nil, err
	}

	enriched := make([]inv_types.UserWithInvite, 0, len(users))
	presentUserIDs := make(map[uuid.UUID]struct{}, len(users))

	for _, u := range users {
		presentUserIDs[u.UserID] = struct{}{}
		row := inv_types.UserWithInvite{OrganizationUsers: u}
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

	// Append pending invites for users who are not yet members of the organization
	for _, inv := range latestByUser {
		if _, exists := presentUserIDs[inv.UserID]; exists {
			continue
		}
		pending := inv_types.UserWithInvite{
			OrganizationUsers: shared_types.OrganizationUsers{
				UserID:         inv.UserID,
				OrganizationID: inv.OrganizationID,
				CreatedAt:      inv.CreatedAt,
				UpdatedAt:      inv.UpdatedAt,
			},
		}
		if !inv.ExpiresAt.IsZero() {
			t := inv.ExpiresAt
			pending.ExpiresAt = &t
		}
		if inv.AcceptedAt != nil {
			pending.AcceptedAt = inv.AcceptedAt
		}
		if inv.InviterUserID != uuid.Nil {
			inviter := inv.InviterUserID
			pending.InvitedBy = &inviter
		}
		email := inv.Email
		name := inv.Name
		role := inv.Role
		pending.InviteEmail = &email
		pending.InviteName = &name
		pending.InviteRole = &role
		enriched = append(enriched, pending)
	}

	return enriched, nil
}
