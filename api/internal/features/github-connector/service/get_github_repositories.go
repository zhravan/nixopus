package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

// GetGithubRepositoriesPaginated fetches repositories for the user's GitHub installation with pagination.
func (c *GithubConnectorService) GetGithubRepositoriesPaginated(userID string, page int, pageSize int) ([]shared_types.GithubRepository, int, error) {
	connectors, err := c.storage.GetAllConnectors(userID)
	if err != nil {
		c.logger.Log(logger.Error, err.Error(), "")
		return nil, 0, err
	}

	if len(connectors) == 0 {
		c.logger.Log(logger.Error, "No connectors found for user", userID)
		return []shared_types.GithubRepository{}, 0, nil
	}

	installation_id := connectors[0].InstallationID
	jwt := GenerateJwt(&connectors[0])

	accessToken, err := c.getInstallationToken(jwt, installation_id)
	if err != nil {
		c.logger.Log(logger.Error, fmt.Sprintf("Failed to get installation token: %s", err.Error()), "")
		return nil, 0, err
	}

	client := &http.Client{}
	url := fmt.Sprintf("https://api.github.com/installation/repositories?per_page=%d&page=%d", pageSize, page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.logger.Log(logger.Error, err.Error(), "")
		return nil, 0, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "nixopus")

	resp, err := client.Do(req)
	if err != nil {
		c.logger.Log(logger.Error, err.Error(), "")
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.logger.Log(logger.Error, fmt.Sprintf("GitHub API error: %s - %s", resp.Status, string(bodyBytes)), "")
		return nil, 0, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var response struct {
		TotalCount   int                             `json:"total_count"`
		Repositories []shared_types.GithubRepository `json:"repositories"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logger.Log(logger.Error, err.Error(), "")
		return nil, 0, err
	}

	return response.Repositories, response.TotalCount, nil
}
