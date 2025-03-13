package service

import (
	"github.com/google/uuid"
	"github.com/raghavyuva/nixopus-api/internal/features/deploy/types"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
	"strconv"
	"time"
)

// CreateDeployment creates a new application deployment in the database
// and starts the deployment process in a separate goroutine.
// It takes a pointer to a CreateDeploymentRequest and a user ID as parameters,
// and returns the created Application struct and an error.
// If the deployment process fails, it returns an error.
func (s *DeployService) CreateDeployment(deployment *types.CreateDeploymentRequest, userID uuid.UUID) (shared_types.Application, error) {
	application := createApplicationFromDeploymentRequest(deployment, userID, nil)
	deployment_config := createDeploymentConfig(application)
	appLogs := createAppLogs(application, deployment_config)
	deployment_status := createDeploymentStatus(deployment_config)

	operations := []struct {
		operation  func() error
		errMessage string
	}{
		{
			operation: func() error {
				return s.storage.AddApplication(&application)
			},
			errMessage: types.LogFailedToCreateApplicationRecord,
		},
		{
			operation: func() error {
				return s.storage.AddApplicationDeployment(&deployment_config)
			},
			errMessage: types.LogFailedToCreateApplicationDeployment,
		},
		{
			operation: func() error {
				return s.storage.AddApplicationLogs(&appLogs)
			},
			errMessage: types.LogFailedToCreateApplicationLogs,
		},
		{
			operation: func() error {
				return s.storage.AddApplicationDeploymentStatus(&deployment_status)
			},
			errMessage: types.LogFailedToCreateApplicationDeploymentStatus,
		},
	}

	for _, op := range operations {
		if err := s.executeDBOperations(op.operation, op.errMessage); err != nil {
			return shared_types.Application{}, err
		}
	}

	deployment_request := DeploymentRequestConfig{
		BuildPack:            deployment.BuildPack,
		BuildVariables:       deployment.BuildVariables,
		EnvironmentVariables: deployment.EnvironmentVariables,
		Name:                 deployment.Name,
		Port:                 deployment.Port,
		Type:                 DeploymentTypeCreate,
	}

	go s.StartDeploymentInBackground(application, &deployment_request, userID, deployment_status, &deployment_config)

	s.logger.Log(logger.Info, types.LogDeploymentStarted, "")
	return application, nil
}

// UpdateDeployment updates an existing application deployment
// in the database and starts the deployment process in a separate goroutine.
// It takes a pointer to a CreateDeploymentRequest and a user ID as parameters,
// and returns the updated Application struct and an error.
// If the deployment process fails, it returns an error.
func (s *DeployService) UpdateDeployment(deployment *types.UpdateDeploymentRequest, userID uuid.UUID) (shared_types.Application, error) {
	application, err := s.storage.GetApplicationById(deployment.ID.String())
	if err != nil {
		return shared_types.Application{}, err
	}
	application = createApplicationFromExistingApplicationAndUpdateRequest(application, deployment)
	deployment_config := createDeploymentConfig(application)
	appLogs := createAppLogs(application, deployment_config)
	deployment_status := createDeploymentStatus(deployment_config)

	operations := []struct {
		operation  func() error
		errMessage string
	}{
		{
			operation: func() error {
				return s.storage.UpdateApplication(&application)
			},
			errMessage: types.LogFailedToCreateApplicationRecord,
		},
		{
			operation: func() error {
				return s.storage.AddApplicationDeployment(&deployment_config)
			},
			errMessage: types.LogFailedToCreateApplicationDeployment,
		},
		{
			operation: func() error {
				return s.storage.AddApplicationLogs(&appLogs)
			},
			errMessage: types.LogFailedToCreateApplicationLogs,
		},
		{
			operation: func() error {
				return s.storage.AddApplicationDeploymentStatus(&deployment_status)
			},
			errMessage: types.LogFailedToCreateApplicationDeploymentStatus,
		},
	}

	for _, op := range operations {
		if err := s.executeDBOperations(op.operation, op.errMessage); err != nil {
			return shared_types.Application{}, err
		}
	}

	deployment_request := DeploymentRequestConfig{
		BuildPack:            application.BuildPack,
		BuildVariables:       GetMapFromString(application.BuildVariables),
		EnvironmentVariables: GetMapFromString(application.EnvironmentVariables),
		Name:                 application.Name,
		Port:                 application.Port,
		Type:                 DeploymentTypeUpdate,
	}

	go s.StartDeploymentInBackground(application, &deployment_request, userID, deployment_status, &deployment_config)

	s.logger.Log(logger.Info, types.LogDeploymentStarted, "")
	return application, nil
}

// executeDBOperations executes a database operation and logs an error if it fails.
// The first parameter is a function that performs the database operation.
// The second parameter is an error message prefix that is used when logging the error.
// If the operation fails, it logs the error message and returns the error.
// Otherwise, it returns nil.
func (s *DeployService) executeDBOperations(fn func() error, errMessage string) error {
	err := fn()
	if err != nil {
		s.logger.Log(logger.Error, errMessage+err.Error(), "")
		return err
	}
	return nil
}

// StartDeploymentInBackground starts the deployment process in a separate goroutine.
// It takes the application, deployment request, user ID, application status, and
// deployment configuration as parameters.
// It logs any errors that occur during the deployment process and updates the
// application deployment status accordingly.
// It also logs the start of the deployment process and adds a new log entry for
// the application.
func (s *DeployService) StartDeploymentInBackground(
	application shared_types.Application,
	deployment *DeploymentRequestConfig,
	userID uuid.UUID,
	appStatus shared_types.ApplicationDeploymentStatus,
	deployment_config *shared_types.ApplicationDeployment,
) {
	handleError := func(errorMessage string, err error) {
		errMsg := errorMessage + err.Error()
		s.logger.Log(logger.Error, errMsg, "")
		s.updateStatus(deployment_config.ID, shared_types.Failed, appStatus.ID)
		s.addLog(application.ID, errMsg, deployment_config.ID)
	}

	s.updateStatus(deployment_config.ID, shared_types.Cloning, appStatus.ID)
	s.addLog(application.ID, types.LogDeploymentStarted, deployment_config.ID)

	repoID, err := strconv.ParseInt(application.Repository, 10, 64)
	if err != nil {
		handleError(types.LogFailedToParseRepositoryID, err)
		return
	}

	repoPath, err := s.github_service.CloneRepository(uint64(repoID), string(userID.String()), string(application.Environment), deployment_config.ID.String())
	if err != nil {
		handleError(types.LogFailedToCloneRepository, err)
		return
	}

	s.logger.Log(logger.Info, types.LogRepositoryClonedSuccessfully, repoPath)
	s.updateStatus(deployment_config.ID, shared_types.Building, appStatus.ID)

	deployer_config := DeployerConfig{
		application.ID,
		&DeploymentRequestConfig{
			BuildVariables:       deployment.BuildVariables,
			EnvironmentVariables: deployment.EnvironmentVariables,
			BuildPack:            deployment.BuildPack,
			Name:                 deployment.Name,
			Port:                 deployment.Port,
			Type:                 deployment.Type,
		},
		userID,
		repoPath,
		appStatus.ID,
		deployment_config,
	}

	err = s.Deployer(deployer_config)
	if err != nil {
		handleError(types.LogFailedToCreateDeployment, err)
		return
	}

	s.updateStatus(deployment_config.ID, shared_types.Deployed, appStatus.ID)
	s.addLog(application.ID, types.LogDeploymentCompletedSuccessfully, deployment_config.ID)
}

// createApplicationFromDeploymentRequest creates an application from a CreateDeploymentRequest
// and a user ID. It populates the application's fields with the corresponding
// values from the request, and sets the CreatedAt and UpdatedAt fields to the
// current time.
func createApplicationFromDeploymentRequest(deployment *types.CreateDeploymentRequest, userID uuid.UUID, createdAt *time.Time) shared_types.Application {
	timeValue := time.Now()
	if createdAt != nil {
		timeValue = *createdAt
	}

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
		CreatedAt:            timeValue,
		UpdatedAt:            time.Now(),
	}

	return application
}

func createApplicationFromExistingApplicationAndUpdateRequest(application shared_types.Application, deployment *types.UpdateDeploymentRequest) shared_types.Application {
	if deployment.Name != "" {
		application.Name = deployment.Name
	}

	if deployment.BuildVariables != nil {
		application.BuildVariables = GetStringFromMap(deployment.BuildVariables)
	}

	if deployment.EnvironmentVariables != nil {
		application.EnvironmentVariables = GetStringFromMap(deployment.EnvironmentVariables)
	}

	if deployment.PreRunCommand != "" {
		application.PreRunCommand = deployment.PreRunCommand
	}

	if deployment.PostRunCommand != "" {
		application.PostRunCommand = deployment.PostRunCommand
	}

	if deployment.Port != 0 {
		application.Port = deployment.Port
	}

	application.UpdatedAt = time.Now()

	return application
}

// createDeploymentConfig creates an ApplicationDeployment from an Application.
// It sets the CreatedAt and UpdatedAt fields with the current time and returns
// the created ApplicationDeployment.
func createDeploymentConfig(application shared_types.Application) shared_types.ApplicationDeployment {
	deployment_config := shared_types.ApplicationDeployment{
		ID:            uuid.New(),
		ApplicationID: application.ID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return deployment_config
}

// createAppLogs creates an ApplicationLogs with the given application ID and the log
// message LogDeploymentStarted. The CreatedAt and UpdatedAt fields are set
// to the current time.
func createAppLogs(application shared_types.Application, deployment_config shared_types.ApplicationDeployment) shared_types.ApplicationLogs {
	app_logs := shared_types.ApplicationLogs{
		ID:                      uuid.New(),
		ApplicationID:           application.ID,
		Log:                     types.LogDeploymentStarted,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
		ApplicationDeploymentID: deployment_config.ID,
	}

	return app_logs
}

// createDeploymentStatus creates an ApplicationDeploymentStatus from an ApplicationDeployment.
// It sets the Status to Started and populates the CreatedAt and UpdatedAt fields
// with the current time.
func createDeploymentStatus(deployment_config shared_types.ApplicationDeployment) shared_types.ApplicationDeploymentStatus {
	deployment_status := shared_types.ApplicationDeploymentStatus{
		ID:                      uuid.New(),
		ApplicationDeploymentID: deployment_config.ID,
		Status:                  shared_types.Started,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}

	return deployment_status
}
