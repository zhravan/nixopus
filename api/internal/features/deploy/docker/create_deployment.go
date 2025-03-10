package docker

import (
	"fmt"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"github.com/google/uuid"
	deploy_types "github.com/raghavyuva/nixopus-api/internal/features/deploy/types"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

func (s *DockerService) buildImageFromDockerfile(contextPath string, dockerfile string, force bool, buildArgs map[string]*string, labels map[string]string) (string, error) {
	buildContextTar, err := archive.TarWithOptions(contextPath, &archive.TarOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create tar from build context: %w", err)
	}
	
	relativeDockerfilePath := filepath.Base(dockerfile)
	
	s.logger.Log(logger.Info, "Using relative Dockerfile path", relativeDockerfilePath)

	buildOptions := types.ImageBuildOptions{
		Dockerfile:  relativeDockerfilePath,
		Remove:      true,
		Tags:        []string{"latest"},
		NoCache:     force,
		ForceRemove: force,
		BuildArgs:   buildArgs,
		Labels:      labels,
		BuildID:     uuid.New().String(),
	}

	return s.BuildImage(buildOptions, buildContextTar)
}

// CreateDeployment creates a new deployment based on the given request.
func (s *DockerService) CreateDeployment(deployment *deploy_types.CreateDeploymentRequest, userID uuid.UUID, contextPath string) error {
	s.logger.Log(logger.Info, "Creating deployment", contextPath)

	switch deployment.BuildPack {
	case shared_types.DockerFile:
		s.logger.Log(logger.Info, "Dockerfile building", "")
		
		buildArgs := make(map[string]*string)
		for k, v := range deployment.BuildVariables {
			value := v
			buildArgs[k] = &value
		}

		labels := make(map[string]string)
		for k, v := range deployment.EnvironmentVariables {
			labels[k] = v
		}

		dockerfilePath := "Dockerfile"

		s.logger.Log(logger.Info, "Build context path", contextPath)
		s.logger.Log(logger.Info, "Using Dockerfile", dockerfilePath)

		id, err := s.buildImageFromDockerfile(
			contextPath,    
			dockerfilePath, 
			false,        
			buildArgs,
			labels,
		)
		if err != nil {
			return fmt.Errorf("failed to build Docker image: %w", err)
		}
		
		s.logger.Log(logger.Info, "Dockerfile built successfully", id)

	case shared_types.DockerCompose:
		s.logger.Log(logger.Info, "Docker compose building", "")
		return nil
	}

	return nil
}