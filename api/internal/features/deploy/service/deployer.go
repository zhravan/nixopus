package service

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	// "github.com/docker/docker/api/types/image"

	"github.com/google/uuid"
	"github.com/raghavyuva/nixopus-api/internal/features/deploy/proxy"
	"github.com/raghavyuva/nixopus-api/internal/features/deploy/types"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	"github.com/raghavyuva/nixopus-api/internal/features/ssh"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

type DeployerConfig struct {
	application       shared_types.Application
	deployment        *shared_types.DeploymentRequestConfig
	userID            uuid.UUID
	contextPath       string
	appStatus         shared_types.ApplicationDeploymentStatus
	deployment_config *shared_types.ApplicationDeployment
}

func (s *DeployService) runCommands(applicationID uuid.UUID, deploymentConfigID uuid.UUID,
	command string, commandType string) error {
	s.addLog(applicationID, fmt.Sprintf("Running %s commands %v", commandType, command), deploymentConfigID)

	if command == "" {
		return nil
	}

	client := ssh.NewSSH()
	output, err := client.RunCommand(command)
	if err != nil {
		s.addLog(applicationID, fmt.Sprintf("Error while running %s command %v", commandType, output), deploymentConfigID)
		return err
	}

	if output != "" {
		s.addLog(applicationID, fmt.Sprintf("%s command resulted in %v", commandType, output), deploymentConfigID)
	}

	return nil
}

func (s *DeployService) PrerunCommands(d DeployerConfig) error {
	return s.runCommands(d.application.ID, d.deployment_config.ID,
		d.application.PreRunCommand, "pre run")
}

func (s *DeployService) PostRunCommands(d DeployerConfig) error {
	return s.runCommands(d.application.ID, d.deployment_config.ID,
		d.application.PostRunCommand, "post run")
}

// Deployer deploys an application using the specified build pack.
//
// The build pack is determined by the Application field of the DeployerConfig
// argument. If the build pack is Dockerfile, the function calls
// handleDockerfileDeployment. If the build pack is Docker Compose, the function
// calls handleDockerComposeDeployment. If the build pack is neither of these,
// the function returns ErrInvalidBuildPack.
func (s *DeployService) Deployer(d DeployerConfig) error {
	s.addLog(d.application.ID, fmt.Sprintf(types.LogDeploymentBuildPack, d.application.BuildPack), d.deployment_config.ID)
	s.PrerunCommands(d)

	var err error
	switch d.application.BuildPack {
	case shared_types.DockerFile:
		err = s.handleDockerfileDeployment(d)
	case shared_types.DockerCompose:
		err = s.handleDockerComposeDeployment(d)
	case shared_types.Static:
		err = s.handleStaticDeployment(d)
	default:
		return types.ErrInvalidBuildPack
	}

	if err != nil {
		return err
	}

	return s.PostRunCommands(d)
}

// handleDockerfileDeployment processes Dockerfile-based deployments
func (s *DeployService) handleDockerfileDeployment(d DeployerConfig) error {
	s.addLog(d.application.ID, types.LogUsingDockerfileStrategy, d.deployment_config.ID)
	s.addLog(d.application.ID, fmt.Sprintf(types.LogBuildContextPath, d.contextPath), d.deployment_config.ID)
	if err := s.buildAndRunDockerImage(d); err != nil {
		return err
	}

	return nil
}

// handleDockerComposeDeployment processes Docker Compose-based deployments
func (s *DeployService) handleDockerComposeDeployment(d DeployerConfig) error {
	s.addLog(d.application.ID, types.LogUsingDockerComposeStrategy, d.deployment_config.ID)

	// For monorepo setups, we need to consider the base path
	composeContextPath := d.contextPath
	if d.application.BasePath != "" && d.application.BasePath != "/" {
		composeContextPath = filepath.Join(d.contextPath, d.application.BasePath)
	}

	// Handle docker-compose.yml path relative to build context
	composeFilePath := "docker-compose.yml"
	if d.application.DockerfilePath != "" {
		// Docker Compose file path should be relative to the build context
		composeFilePath = d.application.DockerfilePath
		if strings.HasPrefix(composeFilePath, "/") {
			composeFilePath = composeFilePath[1:]
		}
	}
	absComposePath := filepath.Join(composeContextPath, composeFilePath)

	s.addLog(d.application.ID, fmt.Sprintf(types.LogBuildContextPath, composeContextPath), d.deployment_config.ID)
	s.addLog(d.application.ID, fmt.Sprintf("Using docker-compose file: %s", absComposePath), d.deployment_config.ID)

	envVars := make(map[string]string)
	for k, v := range GetMapFromString(d.application.EnvironmentVariables) {
		envVars[k] = v
	}
	for k, v := range GetMapFromString(d.application.BuildVariables) {
		envVars[k] = v
	}
	s.addLog(d.application.ID, "Building Docker Compose services...", d.deployment_config.ID)
	if err := s.dockerRepo.ComposeBuild(absComposePath, envVars); err != nil {
		s.addLog(d.application.ID, fmt.Sprintf("Error building Docker Compose services: %v", err), d.deployment_config.ID)
		return fmt.Errorf("%w: %v", types.ErrDockerComposeCommandFailed, err)
	}

	s.addLog(d.application.ID, "Starting Docker Compose services...", d.deployment_config.ID)
	if err := s.dockerRepo.ComposeUp(absComposePath, envVars); err != nil {
		s.addLog(d.application.ID, fmt.Sprintf("Error starting Docker Compose services: %v", err), d.deployment_config.ID)
		return fmt.Errorf("%w: %v", types.ErrDockerComposeCommandFailed, err)
	}

	s.addLog(d.application.ID, "Docker Compose deployment completed successfully", d.deployment_config.ID)
	return nil
}

func (s *DeployService) handleStaticDeployment(d DeployerConfig) error {
	s.addLog(d.application.ID, "Using static file deployment strategy", d.deployment_config.ID)

	caddyProxy := proxy.NewCaddy(&s.logger, d.contextPath, d.application.Domain, strconv.Itoa(d.application.Port), proxy.FileServer)
	if err := caddyProxy.Serve(); err != nil {
		s.addLog(d.application.ID, fmt.Sprintf("Failed to start Caddy proxy: %v", err), d.deployment_config.ID)
		return err
	}
	s.addLog(d.application.ID, "Caddy proxy started successfully", d.deployment_config.ID)

	s.updateStatus(d.deployment_config.ID, shared_types.Deployed, d.appStatus.ID)
	s.addLog(d.application.ID, types.LogDeploymentCompletedSuccessfully, d.deployment_config.ID)
	s.addLog(d.application.ID, fmt.Sprintf("Application %s is available at %s", d.application.Name, d.application.Domain), d.deployment_config.ID)
	return nil
}

// buildAndRunDockerImage handles the Docker image building and running process
func (s *DeployService) buildAndRunDockerImage(d DeployerConfig) error {
	_, err := s.buildImageFromDockerfile(d)
	if err != nil {
		s.addLog(d.application.ID, fmt.Sprintf(types.LogFailedToBuildDockerImage, err.Error()), d.deployment_config.ID)
		return fmt.Errorf("%w: %v", types.ErrBuildDockerImage, err)
	}

	s.logger.Log(logger.Info, types.LogDockerImageBuiltSuccessfully, d.application.Name)
	s.addLog(d.application.ID, types.LogDockerImageBuiltSuccessfully, d.deployment_config.ID)
	containerID, err := s.AtomicUpdateContainer(d)
	if err != nil {
		s.addLog(d.application.ID, fmt.Sprintf(types.LogFailedToRunDockerImage, err.Error()), d.deployment_config.ID)
		return fmt.Errorf("%w: %v", types.ErrRunDockerImage, err)
	}

	s.addLog(d.application.ID, fmt.Sprintf(types.LogContainerRunning, containerID), d.deployment_config.ID)
	s.addLog(d.application.ID, fmt.Sprintf(types.LogApplicationExposed, d.application.Port), d.deployment_config.ID)

	return nil
}
