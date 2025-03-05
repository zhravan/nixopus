package storage

import (
	"context"
	// shared_types "github.com/raghavyuva/nixopus-api/internal/types"
	"github.com/uptrace/bun"
)

type DeployStorage struct {
	DB  *bun.DB
	Ctx context.Context
}

type DeployRepository interface {
	// CreateDeploy(deploy *shared_types.Application) error
	// GetDeploy(id string) (*shared_types.Application, error)
	// UpdateDeploy(ID string, Name string) error
	// DeleteDeploy(deploy *shared_types.Application) error
	// GetDeploys() ([]shared_types.Application, error)
}