package organization

import (
	"encoding/json"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/raghavyuva/nixopus-api/internal/types"
)

const (
	MaxNameLength        = 50
	MaxDescriptionLength = 100
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateID(id string, errorType error) error {
	if id == "" {
		return errorType
	}
	if _, err := uuid.Parse(id); err != nil {
		return errorType
	}
	return nil
}

func (v *Validator) ValidateName(name string) error {
	switch {
	case name == "":
		return types.ErrMissingOrganizationName
	case utf8.RuneCountInString(name) > MaxNameLength:
		return types.ErrOrganizationNameTooLong
	case strings.Contains(name, " "):
		return types.ErrOrganizationNameContainsSpaces
	default:
		return nil
	}
}

func (v *Validator) ValidateDescription(description string) error {
	if utf8.RuneCountInString(description) > MaxDescriptionLength {
		return types.ErrOrganizationDescriptionTooLong
	}
	return nil
}

// ValidateRequest validates a request object against a set of predefined rules.
// It returns an error if the request object is invalid.
//
// The supported request types are:
// - types.CreateOrganizationRequest
// - types.UpdateOrganizationRequest
// - types.DeleteOrganizationRequest
// - types.AddUserToOrganizationRequest
// - types.RemoveUserFromOrganizationRequest
//
// If the request object is not of one of the above types, it returns
// types.ErrInvalidRequestType.
func (v *Validator) ValidateRequest(req interface{}) error {
	switch r := req.(type) {
	case types.CreateOrganizationRequest:
		return v.validateCreate(r)
	case types.UpdateOrganizationRequest:
		return v.validateUpdate(r)
	case types.DeleteOrganizationRequest:
		return v.validateDelete(r)
	case types.AddUserToOrganizationRequest:
		return v.validateAddUser(r)
	case types.RemoveUserFromOrganizationRequest:
		return v.validateRemoveUser(r)
	default:
		return types.ErrInvalidRequestType
	}
}

func (v *Validator) validateCreate(req types.CreateOrganizationRequest) error {
	if err := v.ValidateName(req.Name); err != nil {
		return err
	}
	return v.ValidateDescription(req.Description)
}

func (v *Validator) validateUpdate(req types.UpdateOrganizationRequest) error {
	if err := v.ValidateName(req.Name); err != nil {
		return err
	}
	return v.ValidateDescription(req.Description)
}

func (v *Validator) validateDelete(req types.DeleteOrganizationRequest) error {
	return v.ValidateID(req.ID, types.ErrMissingOrganizationID)
}

func (v *Validator) validateAddUser(req types.AddUserToOrganizationRequest) error {
	if err := v.ValidateID(req.OrganizationID, types.ErrMissingOrganizationID); err != nil {
		return err
	}
	if err := v.ValidateID(req.UserID, types.ErrMissingUserID); err != nil {
		return err
	}
	return v.ValidateID(req.RoleId, types.ErrMissingRoleID)
}

func (v *Validator) validateRemoveUser(req types.RemoveUserFromOrganizationRequest) error {
	if err := v.ValidateID(req.OrganizationID, types.ErrMissingOrganizationID); err != nil {
		return err
	}
	return v.ValidateID(req.UserID, types.ErrMissingUserID)
}

func (v *Validator) ParseRequestBody(req interface{}, body io.ReadCloser, decoded interface{}) error {
	return json.NewDecoder(body).Decode(decoded)
}
