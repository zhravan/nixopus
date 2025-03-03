package service

import "github.com/raghavyuva/nixopus-api/internal/features/github-connector/types"

func (c *GithubConnectorService) UpdateGithubConnectorRequest(ConnectorID string, InstallationID string) error {
	connector, err := c.storage.GetConnector(ConnectorID)
	if err != nil {
		return err
	}
	if connector == nil {
		return types.ErrConnectorDoesNotExist
	}
	err = c.storage.UpdateConnector(InstallationID)
	if err != nil {
		return err
	}
	return nil
}
