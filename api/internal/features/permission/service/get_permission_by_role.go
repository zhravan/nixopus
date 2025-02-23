package service

import (
	"github.com/google/uuid"
	"github.com/raghavyuva/nixopus-api/internal/features/permission/types"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

func (p *PermissionService) GetPermissionByRole(roleID string) ([]shared_types.RolePermissions, error) {
	existingRole, err := p.storage.GetRole(roleID)
	if err == nil && existingRole.ID == uuid.Nil {
		return nil, types.ErrRoleDoesNotExist
	}

	permissions, err := p.storage.GetPermissionsByRole(roleID)
	if err != nil {
		return nil, types.ErrFailedToGetPermissionsByRole
	}

	return permissions, nil
}
