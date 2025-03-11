package service

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	docker_types "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	// "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/moby/term"
	"github.com/raghavyuva/nixopus-api/internal/features/deploy/types"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

func (s *DeployService) CreateDeployment(deployment *types.CreateDeploymentRequest, userID uuid.UUID) error {
	application := shared_types.Application{
		ID:                   uuid.New(),
		Name:                 deployment.Name,
		BuildVariables:       GetStringFromMap(deployment.BuildVariables),
		EnvironmentVariables: GetStringFromMap(deployment.EnvironmentVariables),
		Environment:          deployment.Environment,
		BuildPack:            deployment.BuildPack,
		Repository:           deployment.Repository,
		Branch:               deployment.Branch,
		PreRunCommand:        deployment.PreRunCommand,
		PostRunCommand:       deployment.PostRunCommand,
		Port:                 deployment.Port,
		DomainID:             deployment.DomainID,
		UserID:               userID,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	appStatus := shared_types.ApplicationStatus{
		ID:            uuid.New(),
		ApplicationID: application.ID,
		Status:        shared_types.Started,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	appLogs := shared_types.ApplicationLogs{
		ID:            uuid.New(),
		ApplicationID: application.ID,
		Log:           "Deployment process started",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := s.storage.AddApplication(&application)
	if err != nil {
		s.logger.Log(logger.Error, "Failed to create application record: "+err.Error(), "")
		return err
	}
	s.addLog(application.ID, "Application record created successfully")

	err = s.storage.AddApplicationStatus(&appStatus)
	if err != nil {
		s.logger.Log(logger.Error, "Failed to create application status: "+err.Error(), "")
		return err
	}
	s.addLog(application.ID, "Initial application status set to: "+string(shared_types.Started))

	err = s.storage.AddApplicationLogs(&appLogs)
	if err != nil {
		s.logger.Log(logger.Error, "Failed to create application logs: "+err.Error(), "")
		return err
	}

	s.updateStatus(application.ID, shared_types.Cloning, appStatus.ID)
	s.addLog(application.ID, "Starting repository clone process")

	repoID, err := strconv.ParseInt(application.Repository, 10, 64)
	if err != nil {
		s.logger.Log(logger.Error, "Failed to parse repository ID: "+err.Error(), "")
		s.updateStatus(application.ID, shared_types.Failed, appStatus.ID)
		s.addLog(application.ID, "Failed to parse repository ID: "+err.Error())
		return err
	}

	repoPath, err := s.github_service.CloneRepository(uint64(repoID), string(userID.String()), string(application.Environment))
	if err != nil {
		s.logger.Log(logger.Error, "Failed to clone repository: "+err.Error(), "")
		s.updateStatus(application.ID, shared_types.Failed, appStatus.ID)
		s.addLog(application.ID, "Failed to clone repository: "+err.Error())
		return err
	}

	s.logger.Log(logger.Info, "Repository cloned successfully", repoPath)
	s.addLog(application.ID, fmt.Sprintf("Repository cloned successfully to %s", repoPath))

	s.updateStatus(application.ID, shared_types.Building, appStatus.ID)
	s.addLog(application.ID, "Starting container build process")

	err = s.Deployer(application.ID, deployment, userID, repoPath, appStatus.ID)
	if err != nil {
		s.logger.Log(logger.Error, "Failed to create deployment: "+err.Error(), "")
		s.updateStatus(application.ID, shared_types.Failed, appStatus.ID)
		s.addLog(application.ID, "Failed to create deployment: "+err.Error())
		return err
	}

	s.updateStatus(application.ID, shared_types.Deployed, appStatus.ID)
	s.addLog(application.ID, "Deployment completed successfully")

	s.logger.Log(logger.Info, "Deployment created successfully", "")
	return nil
}

// updateStatus updates the application status
func (s *DeployService) updateStatus(applicationID uuid.UUID, status shared_types.Status, id uuid.UUID) {
	appStatus := shared_types.ApplicationStatus{
		ID:            id,
		ApplicationID: applicationID,
		Status:        status,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := s.storage.UpdateApplicationStatus(&appStatus)
	if err != nil {
		s.logger.Log(logger.Error, "Failed to update application status: "+err.Error(), "")
	}
}

// addLog adds a new log entry for the application
func (s *DeployService) addLog(applicationID uuid.UUID, logMessage string) {
	appLog := shared_types.ApplicationLogs{
		ID:            uuid.New(),
		ApplicationID: applicationID,
		Log:           logMessage,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := s.storage.AddApplicationLogs(&appLog)
	if err != nil {
		s.logger.Log(logger.Error, "Failed to add application log: "+err.Error(), "")
	}
}

func GetStringFromMap(m map[string]string) string {
	var result string
	for key, value := range m {
		result += key + "=" + value + " "
	}
	return result
}

// buildImageFromDockerfile builds a Docker image from a Dockerfile in the
// contextPath directory. The Dockerfile is specified by the dockerfile
// parameter, which is relative to the contextPath directory. If force is true,
// the build will be forced, even if the image already exists. The buildArgs
// parameter is a map of build arguments, which are passed to the Dockerfile as
// environment variables. The labels parameter is a map of labels, which are
// applied to the built image. The image_name parameter is the name of the
// built image. The function returns the ID of the built image, or an error if
// the build fails.
func (s *DeployService) buildImageFromDockerfile(applicationID uuid.UUID, contextPath string, dockerfile string, force bool, buildArgs map[string]*string, labels map[string]string, image_name string, statusID uuid.UUID) (string, error) {
	s.addLog(applicationID, "Starting Docker image build from Dockerfile")
	s.updateStatus(applicationID, shared_types.Building, statusID)

	buildContextTar, err := archive.TarWithOptions(contextPath, &archive.TarOptions{})
	if err != nil {
		errMsg := fmt.Sprintf("failed to create tar from build context: %s", err.Error())
		s.addLog(applicationID, errMsg)
		return "", errors.New(errMsg)
	}
	s.addLog(applicationID, "Created build context archive for Docker build")

	dockerfile_path := filepath.Base(dockerfile)
	s.addLog(applicationID, fmt.Sprintf("Using Dockerfile: %s", dockerfile_path))

	buildOptions := docker_types.ImageBuildOptions{
		Dockerfile:  dockerfile_path,
		Remove:      true,
		Tags:        []string{fmt.Sprintf("%s:latest", image_name)},
		NoCache:     force,
		ForceRemove: force,
		BuildArgs:   buildArgs,
		Labels:      labels,
		BuildID:     uuid.New().String(),
	}

	resp, err := s.dockerRepo.BuildImage(buildOptions, buildContextTar)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to start image build: %s", err.Error())
		s.addLog(applicationID, errMsg)
		return "", errors.New(errMsg)
	}
	defer resp.Body.Close()

	logReader := &LogReader{
		Reader:        resp.Body,
		ApplicationID: applicationID,
		DeployService: s,
	}

	termFd, isTerm := term.GetFdInfo(os.Stdout)
	err = jsonmessage.DisplayJSONMessagesStream(logReader, os.Stdout, termFd, isTerm, nil)
	if err != nil {
		errMsg := fmt.Sprintf("Error processing build output: %s", err.Error())
		s.addLog(applicationID, errMsg)
		return "", errors.New(errMsg)
	}

	s.addLog(applicationID, "Docker image build completed successfully")
	return image_name, nil
}

// LogReader is a custom io.Reader that captures logs from Docker operations and adds them to application logs
type LogReader struct {
	Reader        io.Reader
	ApplicationID uuid.UUID
	DeployService *DeployService
	buffer        []byte
}

func (r *LogReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	if n > 0 {
		r.buffer = append(r.buffer, p[:n]...)
		for {
			idx := -1
			for i, b := range r.buffer {
				if b == '\n' {
					idx = i
					break
				}
			}

			if idx == -1 {
				break
			}

			line := r.buffer[:idx]
			r.buffer = r.buffer[idx+1:]

			var jsonMsg jsonmessage.JSONMessage
			if err := json.Unmarshal(line, &jsonMsg); err == nil {
				if jsonMsg.Stream != "" {
					r.DeployService.addLog(r.ApplicationID, "Build: "+jsonMsg.Stream)
				} else if jsonMsg.Status != "" {
					status := jsonMsg.Status
					if jsonMsg.Progress != nil {
						status += " " + jsonMsg.Progress.String()
					}
					r.DeployService.addLog(r.ApplicationID, "Build: "+status)
				} else if jsonMsg.Error != nil {
					r.DeployService.addLog(r.ApplicationID, "Build error: "+jsonMsg.Error.Message)
				}
			} else {
				r.DeployService.addLog(r.ApplicationID, "Build: "+string(line))
			}
		}
	}

	return n, err
}

// RunImage runs a Docker container from the specified image, maps the
// specified port from the container to the host, and sets the specified
// environment variables. The function returns an error if the container
// cannot be started.
func (s *DeployService) RunImage(applicationID uuid.UUID, imageName string, environment_variables map[string]string, port_str string, statusID uuid.UUID) (string, error) {
	if imageName == "" {
		return "", fmt.Errorf("image name is empty")
	}
	s.logger.Log(logger.Info, "Running container from image", imageName)
	s.addLog(applicationID, fmt.Sprintf("Preparing to run container from image %s", imageName))
	s.updateStatus(applicationID, shared_types.Deploying, statusID)

	port, _ := nat.NewPort("tcp", port_str)
	var env_vars []string
	for k, v := range environment_variables {
		env_vars = append(env_vars, fmt.Sprintf("%s=%s", k, v))
	}

	logEnvVars := make([]string, 0)
	for k, v := range environment_variables {
		if containsSensitiveKeyword(k) {
			logEnvVars = append(logEnvVars, fmt.Sprintf("%s=********", k))
		} else {
			logEnvVars = append(logEnvVars, fmt.Sprintf("%s=%s", k, v))
		}
	}
	s.addLog(applicationID, fmt.Sprintf("Environment variables: %v", logEnvVars))
	s.addLog(applicationID, fmt.Sprintf("Container will expose port %s", port_str))

	// images := s.dockerRepo.ListAllImages(image.ListOptions{})
	// var targetImage string
	// for _, image := range images {
	// 	if image.RepoTags[0] == imageName {
	// 		s.logger.Log(logger.Info, "Image already exists",image.ID)
	// 		targetImage = image.ID
	// 		break
	// 	}
	// }

	container_config := container.Config{
		Image:    imageName,
		Hostname: "nixopus",
		ExposedPorts: nat.PortSet{
			port: struct{}{},
		},
		Env: env_vars,
		Labels: map[string]string{
			"com.docker.compose.project": "nixopus",
			"com.docker.compose.version": "0.0.1",
			"com.project.name":           imageName,
			"application.id":             applicationID.String(),
		},
	}

	host_config := container.HostConfig{
		NetworkMode: "bridge",
		PortBindings: map[nat.Port][]nat.PortBinding{
			port: {
				{
					HostIP:   "0.0.0.0",
					HostPort: port_str,
				},
			},
		},
	}

	network_config := network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"bridge": {},
		},
	}

	s.addLog(applicationID, "Creating container...")
	resp, err := s.dockerRepo.CreateContainer(container_config, host_config, network_config, imageName)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create container: %s", err.Error())
		s.addLog(applicationID, errMsg)
		return "", errors.New(errMsg)
	}
	s.addLog(applicationID, fmt.Sprintf("Container created with ID: %s", resp.ID))

	s.addLog(applicationID, "Starting container...")
	err = s.dockerRepo.StartContainer(resp.ID, container.StartOptions{})
	if err != nil {
		errMsg := fmt.Sprintf("Failed to start container: %s", err.Error())
		s.addLog(applicationID, errMsg)
		return "", errors.New(errMsg)
	}
	s.addLog(applicationID, "Container started successfully")

	go s.collectContainerLogs(applicationID, resp.ID)

	return resp.ID, nil
}

// containsSensitiveKeyword checks if a key likely contains sensitive information
func containsSensitiveKeyword(key string) bool {
	sensitiveKeywords := []string{
		"password", "secret", "token", "key", "auth", "credential", "private",
	}

	lowerKey := strings.ToLower(key)
	for _, word := range sensitiveKeywords {
		if strings.Contains(lowerKey, word) {
			return true
		}
	}

	return false
}

// collectContainerLogs collects logs from a running container and adds them to application logs
func (s *DeployService) collectContainerLogs(applicationID uuid.UUID, containerID string) {
	ctx := context.Background()
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: true,
	}

	logs, err := s.dockerRepo.ContainerLogs(ctx, containerID, options)
	if err != nil {
		s.logger.Log(logger.Error, fmt.Sprintf("Failed to attach to container logs: %s", err.Error()), "")
		s.addLog(applicationID, fmt.Sprintf("Failed to attach to container logs: %s", err.Error()))
		return
	}
	defer logs.Close()

	scanner := bufio.NewScanner(logs)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 8 {
			logLine := line[8:]
			if options.Timestamps {
				parts := strings.SplitN(logLine, " ", 2)
				if len(parts) == 2 {
					logLine = parts[0] + " " + parts[1]
				}
			}

			s.addLog(applicationID, fmt.Sprintf("Container: %s", logLine))
		}
	}

	if err := scanner.Err(); err != nil {
		s.logger.Log(logger.Error, fmt.Sprintf("Error reading container logs: %s", err.Error()), "")
		s.addLog(applicationID, fmt.Sprintf("Error reading container logs: %s", err.Error()))
	}
}

// Deployer handles deployment processes based on the specified build pack type.
//
// This method logs the start of the deployment process and executes different
// deployment strategies depending on the build pack type specified in the
// CreateDeploymentRequest. For Dockerfile build packs, it builds a Docker image
// using the provided build variables and environment variables, then runs the
// Docker image. For DockerCompose build packs, it logs the intended operation
// and returns without performing any actions.
//
// Parameters:
//
//	applicationID - the UUID of the application.
//	deployment - a pointer to the CreateDeploymentRequest containing the deployment details.
//	userID - the UUID of the user initiating the deployment.
//	contextPath - the path to the build context directory.
//
// Returns:
//
//	error - an error if the deployment process fails at any step, otherwise nil.
func (s *DeployService) Deployer(applicationID uuid.UUID, deployment *types.CreateDeploymentRequest, userID uuid.UUID, contextPath string, statusID uuid.UUID) error {
	s.logger.Log(logger.Info, "Creating deployment", contextPath)
	s.addLog(applicationID, fmt.Sprintf("Starting deployment process for build pack: %s", deployment.BuildPack))

	switch deployment.BuildPack {
	case shared_types.DockerFile:
		s.addLog(applicationID, "Using Dockerfile build strategy")
		s.logger.Log(logger.Info, "Dockerfile building", "")

		buildArgs := make(map[string]*string)
		for k, v := range deployment.BuildVariables {
			value := v
			buildArgs[k] = &value
		}
		s.addLog(applicationID, fmt.Sprintf("Using %d build arguments", len(buildArgs)))

		labels := make(map[string]string)
		for k, v := range deployment.EnvironmentVariables {
			labels[k] = v
		}

		dockerfilePath := "Dockerfile"

		s.logger.Log(logger.Info, "Build context path", contextPath)
		s.addLog(applicationID, fmt.Sprintf("Build context path: %s", contextPath))
		s.logger.Log(logger.Info, "Using Dockerfile", dockerfilePath)

		_, err := s.buildImageFromDockerfile(
			applicationID,
			contextPath,
			dockerfilePath,
			false,
			buildArgs,
			labels,
			deployment.Name,
			statusID,
		)
		if err != nil {
			s.addLog(applicationID, fmt.Sprintf("Failed to build Docker image: %s", err.Error()))
			return fmt.Errorf("failed to build Docker image: %w", err)
		}

		s.logger.Log(logger.Info, "Dockerfile built successfully", deployment.Name)
		s.addLog(applicationID, "Docker image built successfully")

		containerID, err := s.RunImage(applicationID, deployment.Name, deployment.EnvironmentVariables, fmt.Sprintf("%d", deployment.Port), statusID)
		if err != nil {
			s.addLog(applicationID, fmt.Sprintf("Failed to run Docker image: %s", err.Error()))
			return fmt.Errorf("failed to run Docker image: %w", err)
		}

		s.addLog(applicationID, fmt.Sprintf("Container is running with ID: %s", containerID))
		s.addLog(applicationID, fmt.Sprintf("Application exposed on port: %d", deployment.Port))

	case shared_types.DockerCompose:
		s.logger.Log(logger.Info, "Docker compose building", "")
		s.addLog(applicationID, "Docker Compose deployment strategy selected")
		s.addLog(applicationID, "Docker Compose deployment not implemented yet")
		return nil
	}

	return nil
}
