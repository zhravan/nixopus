package types

import (
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

type UserOrganizationsResponse struct {
	Organization shared_types.Organization
	Role         RolesResponse
}

type RolesResponse struct {
	shared_types.Role
	Permissions []shared_types.Permission
}