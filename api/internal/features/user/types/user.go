package types

import (
	"errors"

	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

type UserOrganizationsResponse struct {
	Organization shared_types.Organization `json:"organization"`
	Role         shared_types.Role         `json:"role"`
}

type UpdateUserNameRequest struct {
	Name string `json:"name"`
}

var (
	ErrUserDoesNotExist   = errors.New("user does not exist")
	ErrFailedToUpdateUser = errors.New("failed to update user")
	ErrInvalidRequestType = errors.New("invalid request type")
	ErrInvalidAccess      = errors.New("invalid access")
	ErrUserNameIsEmpty    = errors.New("user name is empty")
	ErrSameUserName       = errors.New("user name is same")
	ErrUserNameTooLong    = errors.New("user name is too long")
	ErrUserNameContainsSpaces = errors.New("user name contains spaces")
	ErrUsernameTooShort = errors.New("user name is too short")
)
