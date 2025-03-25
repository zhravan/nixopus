package validation

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/raghavyuva/nixopus-api/internal/features/notification"
	"github.com/raghavyuva/nixopus-api/internal/features/notification/storage"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

type Validator struct {
	storage storage.NotificationRepository
}

func NewValidator(storage storage.NotificationRepository) *Validator {
	return &Validator{
		storage: storage,
	}
}

func (v *Validator) ValidateRequest(req interface{}) error {
	switch r := req.(type) {
	case *notification.CreateSMTPConfigRequest:
		return v.validateCreateSMTPConfigRequest(*r)
	case *notification.GetSMTPConfigRequest:
		return v.validateGetSMTPConfigRequest(*r)
	case *notification.UpdateSMTPConfigRequest:
		return v.validateUpdateSMTPConfigRequest(*r)
	case *notification.DeleteSMTPConfigRequest:
		return v.validateDeleteSMTPConfigRequest(*r)
	case *notification.UpdatePreferenceRequest:
		return v.validateUpdatePreferenceRequest(*r)
	default:
		return notification.ErrInvalidRequestType
	}
}

func (v *Validator) ParseRequestBody(req interface{}, body io.ReadCloser, decoded interface{}) error {
	return json.NewDecoder(body).Decode(decoded)
}

func (v *Validator) validateCreateSMTPConfigRequest(req notification.CreateSMTPConfigRequest) error {
	if req.Host == "" {
		return notification.ErrMissingHost
	}
	if req.Port == 0 {
		return notification.ErrMissingPort
	}
	if req.Username == "" {
		return notification.ErrMissingUsername
	}
	if req.Password == "" {
		return notification.ErrMissingPassword
	}
	return nil
}

func (v *Validator) validateUpdateSMTPConfigRequest(req notification.UpdateSMTPConfigRequest) error {
	if req.ID.String() == "" {
		return notification.ErrMissingID
	}

	return nil
}

func (v *Validator) validateDeleteSMTPConfigRequest(req notification.DeleteSMTPConfigRequest) error {
	if req.ID.String() == "" {
		return notification.ErrMissingID
	}

	return nil
}

func (v *Validator) validateGetSMTPConfigRequest(req notification.GetSMTPConfigRequest) error {
	if req.ID.String() == "" {
		return notification.ErrMissingID
	}
	return nil
}

func (v *Validator) validateUpdatePreferenceRequest(req notification.UpdatePreferenceRequest) error {
	if req.Type == "" {
		return notification.ErrMissingType
	}
	if req.Category == "" {
		return notification.ErrMissingCategory
	}

	if req.Category != "activity" && req.Category != "security" && req.Category != "update" {
		return notification.ErrInvalidRequestType
	}

	return nil
}

func (v *Validator) AccessValidator(w http.ResponseWriter, r *http.Request, user *shared_types.User) error {
	path := r.URL.Path

	// only admin has access to create, update, and delete smtp
	if path == "/api/v1/notification/smtp" && (r.Method == "POST" || r.Method == "DELETE" || r.Method == "PUT") && user.Type != shared_types.RoleAdmin {
		return notification.ErrAccessDenied
	}

	smtp, err := v.storage.GetSmtp(r.URL.Query().Get("id"))

	// check if admin is the one who created the smtp
	if path == "/api/v1/notification/smtp" && (r.Method == "DELETE" || r.Method == "PUT") {
		if err != nil {
			return err
		}
		if smtp.UserID != user.ID {
			return notification.ErrPermissionDenied
		}
	}

	// for getting smtp configuration user must be either admin, or the one who created the smtp, or have read access to the organization
	if path == "/api/v1/notification/smtp" && r.Method == "GET" && user.Type != shared_types.RoleAdmin {
		if err != nil {
			return err
		}
		err = v.checkIfUserBelongsToOrganization(user.Organizations, smtp.OrganizationID)
		if err != nil {
			return err
		}

		err = v.checkIfUserHasReadAccess(user.OrganizationUsers, smtp.OrganizationID)
		if err != nil {
			return err
		}
	}

	// for preferences we do not check access validation since they are specific to the user not for org
	return nil
}

func (v *Validator) checkIfUserBelongsToOrganization(user_orgs []shared_types.Organization, org_id uuid.UUID) error {
	for _, org := range user_orgs {
		if org.ID == org_id {
			return nil
		}
	}
	return notification.ErrUserDoesNotBelongToOrganization
}

func (v *Validator) checkIfUserHasReadAccess(user_orgs []shared_types.OrganizationUsers, org_id uuid.UUID) error {
	for _, orgUser := range user_orgs {
		if orgUser.OrganizationID == org_id {
			if orgUser.Role != nil {
				// Admin and Member roles automatically have read access
				if orgUser.Role.Name == shared_types.RoleAdmin || orgUser.Role.Name == shared_types.RoleMember {
					return nil
				}

				// Check if the user has specific read permissions for resources
				if orgUser.Role.Permissions != nil {
					for _, permission := range orgUser.Role.Permissions {
						// Verify if the permission allows reading smtp or organization
						if permission.Name == "read" && permission.Resource == "organization" {
							return nil
						}
					}
				}
			}
		}
	}

	return notification.ErrUserDoesNotHavePermissionForTheResource
}
