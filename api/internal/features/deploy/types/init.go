package types

import (
	"errors"

	"github.com/google/uuid"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

type IsNameAlreadyTakenRequest struct {
	Name string `json:"name"`
}

type IsDomainAlreadyTakenRequest struct {
	Domain string `json:"domain"`
}

type IsDomainValidRequest struct {
	Domain string `json:"domain"`
}

type IsPortAlreadyTakenRequest struct {
	Port int `json:"port"`
}

type CreateDeploymentRequest struct {
	Name                 string                   `json:"name"`
	DomainID             uuid.UUID                `json:"domain_id"`
	Environment          shared_types.Environment `json:"environment"`
	BuildPack            shared_types.BuildPack   `json:"build_pack"`
	Repository           string                   `json:"repository"`
	Branch               string                   `json:"branch"`
	PreRunCommand        string                   `json:"pre_run_command"`
	PostRunCommand       string                   `json:"post_run_command"`
	BuildVariables       map[string]string        `json:"build_variables"`
	EnvironmentVariables map[string]string        `json:"environment_variables"`
	Port                 int                      `json:"port"`
}

var (
	ErrInvalidRequestType = errors.New("invalid request type")
	ErrMissingName        = errors.New("name is required")
	ErrMissingDomainID    = errors.New("domain_id is required")
	ErrMissingRepository  = errors.New("repository is required")
	ErrMissingBranch      = errors.New("branch is required")
	ErrMissingPort        = errors.New("port is required")
	ErrInvalidEnvironment = errors.New("invalid environment")
	ErrInvalidBuildPack   = errors.New("invalid build pack")
	ErrFailedToCreateTarFromContext = errors.New("failed to create tar from context")
	ErrProcessingBuildOutput = errors.New("failed to process build output")
)

var (
	LogDeploymentStarted = "Deployment started"
	LogFailedToCreateApplicationRecord = "Failed to create application record"
	LogFailedToCreateApplicationStatus = "Failed to create application status: "
	LogFailedToCreateApplicationDeployment = "Failed to create application deployment: "
	LogFailedToCreateApplicationDeploymentStatus = "Failed to create application deployment status: "
	LogFailedToCreateApplicationLogs = "Failed to create application logs: "
	LogFailedToParseRepositoryID = "Failed to parse repository ID: "
	LogFailedToCloneRepository = "Failed to clone repository: "
	LogRepositoryClonedSuccessfully = "Repository cloned successfully"
	LogFailedToCreateDeployment = "Failed to create deployment: "
	LogDeploymentCompletedSuccessfully = "Deployment completed successfully"
	LogDockerImageBuiltSuccessfully = "Docker image built successfully"
	LogStartingDockerImageBuild = "Starting Docker image build from Dockerfile"
)