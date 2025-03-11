package service

import (
	"fmt"
	"strconv"
	"time"
	// "github.com/docker/docker/api/types/image"
	"github.com/google/uuid"
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
