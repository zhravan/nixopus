package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/raghavyuva/nixopus-api/internal/features/deploy/docker"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	"github.com/raghavyuva/nixopus-api/internal/features/ssh"
	"github.com/raghavyuva/nixopus-api/internal/features/update/types"
	"github.com/raghavyuva/nixopus-api/internal/storage"
)

type Environment string

const (
	Production Environment = "production"
	Staging    Environment = "staging"
)

type UpdateService struct {
	storage *storage.App
	logger  *logger.Logger
	ctx     context.Context
	docker  *docker.DockerService
	env     Environment
}

func NewUpdateService(storage *storage.App, logger *logger.Logger, ctx context.Context) *UpdateService {
	var env Environment
	dockerEnv := os.Getenv("DOCKER_CONTEXT")
	if dockerEnv == "nixopus-staging" {
		env = Staging
	} else {
		env = Production
	}
	return &UpdateService{
		storage: storage,
		logger:  logger,
		ctx:     ctx,
		docker:  docker.NewDockerService(),
		env:     env,
	}
}

// CheckForUpdates checks for updates for the current environment and returns the current and latest version
func (s *UpdateService) CheckForUpdates() (*types.UpdateCheckResponse, error) {
	currentVersion, err := s.getCurrentVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get current version: %w", err)
	}

	latestVersion, err := s.fetchLatestVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest version: %w", err)
	}

	updateAvailable := currentVersion != latestVersion && latestVersion != ""

	return &types.UpdateCheckResponse{
		CurrentVersion:  currentVersion,
		LatestVersion:   latestVersion,
		UpdateAvailable: updateAvailable,
		LastChecked:     time.Now(),
		Environment:     string(s.env),
	}, nil
}

// getCurrentVersion gets the current version from the .env file
func (s *UpdateService) getCurrentVersion() (string, error) {
	version := os.Getenv("APP_VERSION")
	if version != "" {
		return version, nil
	}

	return "", fmt.Errorf("APP_VERSION not found in .env file")
}

// fetchLatestVersion fetches the latest version from the appropriate branch from our repo
func (s *UpdateService) fetchLatestVersion() (string, error) {
	branch := "master"
	if s.env == Staging {
		branch = "feat/auto_update" // TODO: Change to develop after testing
	}

	url := fmt.Sprintf("https://raw.githubusercontent.com/raghavyuva/nixopus/refs/heads/%s/version.txt", branch)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch version: status %d", resp.StatusCode)
	}

	versionBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(versionBytes)), nil
}

// PerformUpdate performs an update for the current environment
func (s *UpdateService) PerformUpdate() error {
	ssh := ssh.NewSSH()
	client, err := ssh.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect via SSH: %w", err)
	}
	defer client.Close()

	tempDir := "/tmp/nixopus-update"
	composeFile := "docker-compose.yml"
	if s.env == Staging {
		composeFile = "docker-compose-staging.yml"
	}

	defer func() {
		ssh.RunCommand(fmt.Sprintf("rm -rf %s", tempDir))
	}()

	if _, err := ssh.RunCommand(fmt.Sprintf("rm -rf %s && mkdir -p %s", tempDir, tempDir)); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	cloneCmd := fmt.Sprintf("cd %s && GIT_TERMINAL_PROMPT=0 git clone --verbose https://github.com/raghavyuva/nixopus.git . 2>&1", tempDir)
	cloneOutput, err := ssh.RunCommand(cloneCmd)
	if err != nil {
		s.logger.Log(logger.Error, "Git clone failed", fmt.Sprintf("output: %s, error: %v", cloneOutput, err))
		return fmt.Errorf("failed to clone repository: %w (output: %s)", err, cloneOutput)
	}

	if s.env == Staging {
		checkoutOutput, err := ssh.RunCommand(fmt.Sprintf("cd %s && git checkout feat/auto_update 2>&1", tempDir+"/nixopus"))
		if err != nil {
			s.logger.Log(logger.Error, "Git checkout failed", fmt.Sprintf("output: %s, error: %v", checkoutOutput, err))
			return fmt.Errorf("failed to checkout staging branch: %w (output: %s)", err, checkoutOutput)
		}
	} else {
		checkoutOutput, err := ssh.RunCommand(fmt.Sprintf("cd %s && git checkout master 2>&1", tempDir+"/nixopus"))
		if err != nil {
			s.logger.Log(logger.Error, "Git checkout failed", fmt.Sprintf("output: %s, error: %v", checkoutOutput, err))
			return fmt.Errorf("failed to checkout master branch: %w (output: %s)", err, checkoutOutput)
		}
	}

	tempDir = tempDir + "/nixopus"

	composeCheckOutput, err := ssh.RunCommand(fmt.Sprintf("cd %s && test -f %s && echo 'exists' || echo 'missing'", tempDir, composeFile))
	if err != nil {
		return fmt.Errorf("failed to check compose file: %w", err)
	}
	if strings.TrimSpace(composeCheckOutput) == "missing" {
		return fmt.Errorf("compose file %s not found in repository", composeFile)
	}

	if s.env == Staging {
		restartOutput, err := ssh.RunCommand(fmt.Sprintf("cd %s && docker compose -f %s restart --build 2>&1", tempDir, composeFile))
		if err != nil {
			s.logger.Log(logger.Error, "Container restart failed", fmt.Sprintf("output: %s, error: %v", restartOutput, err))
			return fmt.Errorf("failed to restart staging containers: %w (output: %s)", err, restartOutput)
		}
	} else {
		restartOutput, err := ssh.RunCommand(fmt.Sprintf("cd %s && docker compose -f %s restart 2>&1", tempDir, composeFile))
		if err != nil {
			s.logger.Log(logger.Error, "Container restart failed", fmt.Sprintf("output: %s, error: %v", restartOutput, err))
			return fmt.Errorf("failed to restart production containers: %w (output: %s)", err, restartOutput)
		}
	}

	return nil
}

// GetUserAutoUpdatePreference gets the user's auto update preference from the database
func (s *UpdateService) GetUserAutoUpdatePreference(userID uuid.UUID) (bool, error) {
	var autoUpdate bool

	err := s.storage.Store.DB.NewSelect().
		TableExpr("user_settings").
		Column("auto_update").
		Where("user_id = ?", userID).
		Scan(s.ctx, &autoUpdate)

	if err != nil {
		return false, fmt.Errorf("failed to get user settings: %w", err)
	}

	return autoUpdate, nil
}
