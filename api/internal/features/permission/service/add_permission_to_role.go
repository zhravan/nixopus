package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/raghavyuva/nixopus-api/internal/features/permission/types"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

func (c *PermissionService) AddPermissionToRole(permissionID string, roleID string) error {
	existingRole, err := c.storage.GetRole(roleID)
	if err == nil && existingRole.ID == uuid.Nil {
		return types.ErrRoleDoesNotExist
	}

	existingPermission, err := c.storage.GetPermission(permissionID)
	if err == nil && existingPermission.ID == uuid.Nil {
		return types.ErrPermissionDoesNotExist
	}

	rolePermissionToCreate := shared_types.RolePermissions{
		ID:           uuid.New(),
		RoleID:       existingRole.ID,
		PermissionID: existingPermission.ID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := c.storage.AddPermissionToRole(rolePermissionToCreate); err != nil {
		return types.ErrFailedToAddPermissionToRole
	}

	return nil
}