package service

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/raghavyuva/nixopus-api/internal/features/deploy/types"
	"github.com/raghavyuva/nixopus-api/internal/features/github-connector/service"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

// CreateDeployment creates a new application deployment in the database
// and starts the deployment process in a separate goroutine.
func (s *DeployService) CreateDeployment(deployment *types.CreateDeploymentRequest, userID uuid.UUID) (shared_types.Application, error) {
	application := createApplicationFromDeploymentRequest(deployment, userID, nil)

	deploy_config, err := s.prepareDeploymentConfig(application, userID, shared_types.DeploymentTypeCreate, false, false)
	if err != nil {
		return shared_types.Application{}, err
	}

	go s.StartDeploymentInBackground(deploy_config)
	s.logger.Log(logger.Info, types.LogDeploymentStarted, "")
	return application, nil
}

// UpdateDeployment updates an existing application deployment
// in the database and starts the deployment process in a separate goroutine.
func (s *DeployService) UpdateDeployment(deployment *types.UpdateDeploymentRequest, userID uuid.UUID) (shared_types.Application, error) {
	application, err := s.storage.GetApplicationById(deployment.ID.String())
	if err != nil {
		return shared_types.Application{}, err
	}

	application = createApplicationFromExistingApplicationAndUpdateRequest(application, deployment)

	deploy_config, err := s.prepareDeploymentConfig(application, userID, shared_types.DeploymentTypeUpdate, false, false)
	if err != nil {
		return shared_types.Application{}, err
	}

	go s.StartDeploymentInBackground(deploy_config)
	s.logger.Log(logger.Info, types.LogDeploymentStarted, "")
	return application, nil
}

// ReDeployApplication redeploys an existing application
func (s *DeployService) ReDeployApplication(redeployRequest *types.ReDeployApplicationRequest, userID uuid.UUID) (shared_types.Application, error) {
	application, err := s.storage.GetApplicationById(redeployRequest.ID.String())
	if err != nil {
		return shared_types.Application{}, err
	}

	deploy_config, err := s.prepareDeploymentConfig(application, userID, shared_types.DeploymentTypeReDeploy, redeployRequest.Force, redeployRequest.ForceWithoutCache)
	if err != nil {
		return shared_types.Application{}, err
	}

	go s.StartDeploymentInBackground(deploy_config)
	s.logger.Log(logger.Info, types.LogDeploymentStarted, "")
	return application, nil
}

// StartDeploymentInBackground starts the deployment process in a separate goroutine.
// It takes the application, deployment request, user ID, application status, and
// deployment configuration as parameters.
// It logs any errors that occur during the deployment process and updates the
// application deployment status accordingly.
// It also logs the start of the deployment process and adds a new log entry for
// the application.
func (s *DeployService) StartDeploymentInBackground(
	d DeployerConfig,
) {
	handleError := func(errorMessage string, err error) {
		errMsg := errorMessage + err.Error()
		s.logger.Log(logger.Error, errMsg, "")
		s.updateStatus(d.deployment_config.ID, shared_types.Failed, d.appStatus.ID)
		s.addLog(d.application.ID, errMsg, d.deployment_config.ID)
	}

	// if the deployment type is restart then just restart the container
	if d.deployment.Type == shared_types.DeploymentTypeRestart {
		s.RestartContainer(d)
		return
	}

	s.updateStatus(d.deployment_config.ID, shared_types.Cloning, d.appStatus.ID)
	s.addLog(d.application.ID, types.LogDeploymentStarted, d.deployment_config.ID)

	repoID, err := strconv.ParseInt(d.application.Repository, 10, 64)
	if err != nil {
		handleError(types.LogFailedToParseRepositoryID, err)
		return
	}

	cloneRepositoryConfig := service.CloneRepositoryConfig{
		RepoID:         uint64(repoID),
		UserID:         string(d.userID.String()),
		Environment:    string(d.application.Environment),
		DeploymentID:   d.deployment_config.ID.String(),
		DeploymentType: string(d.deployment.Type),
		Branch:         d.application.Branch,
		ApplicationID:  d.application.ID.String(),
	}

	// we will pass the commit hash to the clone repository function for rollback feature
	repoPath, err := s.github_service.CloneRepository(cloneRepositoryConfig, &d.deployment_config.CommitHash)
	if err != nil {
		handleError(types.LogFailedToCloneRepository, err)
		return
	}

	s.logger.Log(logger.Info, types.LogRepositoryClonedSuccessfully, repoPath)
	s.updateStatus(d.deployment_config.ID, shared_types.Building, d.appStatus.ID)

	// based on the deployment type we will get the path of the repository where it is present for the deployment, (the context for the build basically) (till this point the context path will be empty)
	d.contextPath = repoPath
	err = s.Deployer(d)
	if err != nil {
		handleError(types.LogFailedToCreateDeployment, err)
		return
	}

	s.updateStatus(d.deployment_config.ID, shared_types.Deployed, d.appStatus.ID)
	s.addLog(d.application.ID, types.LogDeploymentCompletedSuccessfully, d.deployment_config.ID)
}

func (s *DeployService) GetDeploymentById(deploymentID string) (shared_types.ApplicationDeployment, error) {
	return s.storage.GetApplicationDeploymentById(deploymentID)
}

func (s *DeployService) DeleteDeployment(deployment *types.DeleteDeploymentRequest, userID uuid.UUID) error {
	return s.storage.DeleteDeployment(deployment, userID)
}

func (s *DeployService) RollbackDeployment(deployment *types.RollbackDeploymentRequest, userID uuid.UUID) error {
	deployment_details, err := s.storage.GetApplicationDeploymentById(deployment.ID.String())
	if err != nil {
		return err
	}

	application_details, err := s.storage.GetApplicationById(string(deployment_details.ApplicationID.String()))
	if err != nil {
		return err
	}

	deployStatus := createDeploymentStatus(deployment.ID)
	s.updateStatus(deployment_details.ID, shared_types.Deploying, deployStatus.ID)

	deploy_config, err := s.prepareDeploymentConfig(application_details, userID, shared_types.DeploymentTypeRollback, false, false)
	if err != nil {
		return err
	}

	go s.StartDeploymentInBackground(deploy_config)
	return nil
}

func (s *DeployService) RestartDeployment(deployment *types.RestartDeploymentRequest, userID uuid.UUID) error {
	deployment_details, err := s.storage.GetApplicationDeploymentById(deployment.ID.String())
	if err != nil {
		return err
	}

	application_details, err := s.storage.GetApplicationById(string(deployment_details.ApplicationID.String()))
	if err != nil {
		return err
	}

	deployStatus := createDeploymentStatus(deployment.ID)
	s.updateStatus(deployment_details.ID, shared_types.Deploying, deployStatus.ID)

	deploy_config, err := s.prepareDeploymentConfig(application_details, userID, shared_types.DeploymentTypeRestart, false, false)
	if err != nil {
		return err
	}

	s.StartDeploymentInBackground(deploy_config)
	return nil
}

func (s *DeployService) prepareDeploymentConfig(application shared_types.Application, userID uuid.UUID, deploymentType shared_types.DeploymentType, force, forceWithoutCache bool) (DeployerConfig, error) {
	deployRequest, deployStatus, deployment_config, err := s.createAndPrepareDeployment(application, shared_types.DeploymentRequestConfig{
		Type:              deploymentType,
		Force:             force,
		ForceWithoutCache: forceWithoutCache,
	})
	if err != nil {
		return DeployerConfig{}, err
	}

	return DeployerConfig{
		application:       application,
		deployment:        deployRequest,
		userID:            userID,
		contextPath:       "",
		appStatus:         deployStatus,
		deployment_config: &deployment_config,
	}, nil
}
