package storage

import (
	"context"
	"time"

	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
	"github.com/uptrace/bun"
)

type GithubConnectorStorage struct {
	DB  *bun.DB
	Ctx context.Context
}

func (s *GithubConnectorStorage) CreateConnector(connector *shared_types.GithubConnector) error {
	tx, err := s.DB.BeginTx(s.Ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NewInsert().Model(connector).Exec(s.Ctx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *GithubConnectorStorage) UpdateConnector(InstallationID string) error {
	tx, err := s.DB.BeginTx(s.Ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var connector shared_types.GithubConnector
	connector.UpdatedAt = time.Now()
	_, err = tx.NewUpdate().Model(&connector).Where("installation_id = ?", InstallationID).Exec(s.Ctx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *GithubConnectorStorage) GetConnector(ConnectorID string) (*shared_types.GithubConnector, error) {
	var connector shared_types.GithubConnector
	err := s.DB.NewSelect().Model(&connector).Where("id = ?", ConnectorID).Scan(s.Ctx)
	return &connector, err
}

func (s *GithubConnectorStorage) GetAllConnectors(UserID string) ([]shared_types.GithubConnector, error) {
	var connectors []shared_types.GithubConnector
	err := s.DB.NewSelect().Model(&connectors).Where("user_id = ?", UserID).Scan(s.Ctx)
	return connectors, err
}
