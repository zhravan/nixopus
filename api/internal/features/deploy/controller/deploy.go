package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
	"github.com/raghavyuva/nixopus-api/internal/features/deploy/types"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	"github.com/raghavyuva/nixopus-api/internal/utils"

	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

func (c *DeployController) HandleDeploy(f fuego.ContextWithBody[types.CreateDeploymentRequest]) (*shared_types.Response, error) {
	c.logger.Log(logger.Info, "starting deployment process", "")

	data, err := f.Body()
	if err != nil {
		c.logger.Log(logger.Error, "failed to read request body", err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	c.logger.Log(logger.Info, "request body parsed successfully", "name: "+data.Name)

	if err := c.validator.ValidateRequest(&data); err != nil {
		c.logger.Log(logger.Error, "request validation failed", "name: "+data.Name+", error: "+err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	user := utils.GetUser(f.Response(), f.Request())
	if user == nil {
		c.logger.Log(logger.Error, "user authentication failed", "name: "+data.Name)
		return nil, fuego.HTTPError{
			Err:    nil,
			Status: http.StatusUnauthorized,
		}
	}

	organizationID := utils.GetOrganizationID(f.Request())
	if organizationID == uuid.Nil {
		c.logger.Log(logger.Error, "organization not found", "name: "+data.Name)
		return nil, fuego.HTTPError{
			Err:    nil,
			Status: http.StatusUnauthorized,
		}
	}

	c.logger.Log(logger.Info, "attempting to create deployment", "name: "+data.Name+", user_id: "+user.ID.String())

	application, err := c.service.CreateDeployment(&data, user.ID, organizationID)
	if err != nil {
		c.logger.Log(logger.Error, "failed to create deployment", "name: "+data.Name+", error: "+err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}

	c.logger.Log(logger.Info, "deployment created successfully", "name: "+data.Name)
	return &shared_types.Response{
		Status:  "success",
		Message: "Deployment created successfully",
		Data:    application,
	}, nil
}

func (c *DeployController) UpdateApplication(f fuego.ContextWithBody[types.UpdateDeploymentRequest]) (*shared_types.Response, error) {
	c.logger.Log(logger.Info, "starting application update process", "")

	data, err := f.Body()
	if err != nil {
		if err == io.EOF {
			c.logger.Log(logger.Error, "empty request body received", "id is required for update")
			return nil, fuego.HTTPError{
				Err:    types.ErrMissingID,
				Status: http.StatusBadRequest,
			}
		}
		c.logger.Log(logger.Error, "failed to read request body", err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	c.logger.Log(logger.Info, "request body parsed successfully", "id: "+data.ID.String())

	if err := c.validator.ValidateRequest(&data); err != nil {
		c.logger.Log(logger.Error, "request validation failed", "id: "+data.ID.String()+", error: "+err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	user := utils.GetUser(f.Response(), f.Request())
	if user == nil {
		c.logger.Log(logger.Error, "user authentication failed", "id: "+data.ID.String())
		return nil, fuego.HTTPError{
			Err:    nil,
			Status: http.StatusUnauthorized,
		}
	}

	organizationID := utils.GetOrganizationID(f.Request())
	if organizationID == uuid.Nil {
		c.logger.Log(logger.Error, "organization not found", "id: "+data.ID.String())
		return nil, fuego.HTTPError{
			Err:    nil,
			Status: http.StatusUnauthorized,
		}
	}

	c.logger.Log(logger.Info, "attempting to update application", "id: "+data.ID.String()+", user_id: "+user.ID.String())

	application, err := c.service.UpdateDeployment(&data, user.ID, organizationID)
	if err != nil {
		c.logger.Log(logger.Error, "failed to update application", "id: "+data.ID.String()+", error: "+err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}

	c.logger.Log(logger.Info, "application updated successfully", "id: "+data.ID.String())
	return &shared_types.Response{
		Status:  "success",
		Message: "Application updated successfully",
		Data:    application,
	}, nil
}

func (c *DeployController) DeleteApplication(f fuego.ContextWithBody[types.DeleteDeploymentRequest]) (*shared_types.Response, error) {
	c.logger.Log(logger.Info, "starting application deletion process", "")

	data, err := f.Body()
	if err != nil {
		if err == io.EOF {
			c.logger.Log(logger.Error, "empty request body received", "id is required for deletion")
			return nil, fuego.HTTPError{
				Err:    types.ErrMissingID,
				Status: http.StatusBadRequest,
			}
		}
		c.logger.Log(logger.Error, "failed to read request body", err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	c.logger.Log(logger.Info, "request body parsed successfully", "id: "+data.ID.String())

	if err := c.validator.ValidateRequest(&data); err != nil {
		c.logger.Log(logger.Error, "request validation failed", "id: "+data.ID.String()+", error: "+err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	user := utils.GetUser(f.Response(), f.Request())
	if user == nil {
		c.logger.Log(logger.Error, "user authentication failed", "id: "+data.ID.String())
		return nil, fuego.HTTPError{
			Err:    nil,
			Status: http.StatusUnauthorized,
		}
	}

	c.logger.Log(logger.Info, "attempting to delete application", "id: "+data.ID.String()+", user_id: "+user.ID.String())

	err = c.service.DeleteDeployment(&data, user.ID)
	if err != nil {
		c.logger.Log(logger.Error, "failed to delete application", "id: "+data.ID.String()+", error: "+err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}

	c.logger.Log(logger.Info, "application deleted successfully", "id: "+data.ID.String())
	return &shared_types.Response{
		Status:  "success",
		Message: "Application deleted successfully",
		Data:    nil,
	}, nil
}

func (c *DeployController) ReDeployApplication(f fuego.ContextWithBody[types.ReDeployApplicationRequest]) (*shared_types.Response, error) {
	c.logger.Log(logger.Info, "starting application redeployment process", "")

	data, err := f.Body()
	if err != nil {
		if err == io.EOF {
			c.logger.Log(logger.Error, "empty request body received", "id is required for redeployment")
			return nil, fuego.HTTPError{
				Err:    types.ErrMissingID,
				Status: http.StatusBadRequest,
			}
		}
		c.logger.Log(logger.Error, "failed to read request body", err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	c.logger.Log(logger.Info, "request body parsed successfully", "id: "+data.ID.String())

	if err := c.validator.ValidateRequest(&data); err != nil {
		c.logger.Log(logger.Error, "request validation failed", "id: "+data.ID.String()+", error: "+err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	user := utils.GetUser(f.Response(), f.Request())
	if user == nil {
		c.logger.Log(logger.Error, "user authentication failed", "id: "+data.ID.String())
		return nil, fuego.HTTPError{
			Err:    nil,
			Status: http.StatusUnauthorized,
		}
	}

	organizationID := utils.GetOrganizationID(f.Request())
	if organizationID == uuid.Nil {
		c.logger.Log(logger.Error, "organization not found", "id: "+data.ID.String())
		return nil, fuego.HTTPError{
			Err:    nil,
			Status: http.StatusUnauthorized,
		}
	}

	c.logger.Log(logger.Info, "attempting to redeploy application", "id: "+data.ID.String()+", user_id: "+user.ID.String())

	application, err := c.service.ReDeployApplication(&data, user.ID, organizationID)
	if err != nil {
		c.logger.Log(logger.Error, "failed to redeploy application", "id: "+data.ID.String()+", error: "+err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}

	c.logger.Log(logger.Info, "application redeployed successfully", "id: "+data.ID.String())
	return &shared_types.Response{
		Status:  "success",
		Message: "Application redeployed successfully",
		Data:    application,
	}, nil
}

func (c *DeployController) GetDeploymentById(f fuego.ContextNoBody) (*shared_types.Response, error) {
	deploymentID := f.PathParam("deployment_id")

	deployment, err := c.service.GetDeploymentById(deploymentID)
	if err != nil {
		c.logger.Log(logger.Error, err.Error(), "")
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}

	return &shared_types.Response{
		Status:  "success",
		Message: "Deployment Retrieved successfully",
		Data:    deployment,
	}, nil
}

func (c *DeployController) HandleRollback(f fuego.ContextWithBody[types.RollbackDeploymentRequest]) (*shared_types.Response, error) {
	c.logger.Log(logger.Info, "starting application rollback process", "")

	data, err := f.Body()
	if err != nil {
		if err == io.EOF {
			c.logger.Log(logger.Error, "empty request body received", "id is required for rollback")
			return nil, fuego.HTTPError{
				Err:    types.ErrMissingID,
				Status: http.StatusBadRequest,
			}
		}
		c.logger.Log(logger.Error, "failed to read request body", err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	c.logger.Log(logger.Info, "request body parsed successfully", "id: "+data.ID.String())

	if err := c.validator.ValidateRequest(&data); err != nil {
		c.logger.Log(logger.Error, "request validation failed", "id: "+data.ID.String()+", error: "+err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	user := utils.GetUser(f.Response(), f.Request())
	if user == nil {
		c.logger.Log(logger.Error, "user authentication failed", "id: "+data.ID.String())
		return nil, fuego.HTTPError{
			Err:    nil,
			Status: http.StatusUnauthorized,
		}
	}

	organizationID := utils.GetOrganizationID(f.Request())
	if organizationID == uuid.Nil {
		c.logger.Log(logger.Error, "organization not found", "id: "+data.ID.String())
		return nil, fuego.HTTPError{
			Err:    nil,
			Status: http.StatusUnauthorized,
		}
	}

	c.logger.Log(logger.Info, "attempting to rollback application", "id: "+data.ID.String()+", user_id: "+user.ID.String())

	err = c.service.RollbackDeployment(&data, user.ID, organizationID)
	if err != nil {
		c.logger.Log(logger.Error, "failed to rollback application", "id: "+data.ID.String()+", error: "+err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}

	c.logger.Log(logger.Info, "application rolled back successfully", "id: "+data.ID.String())
	return &shared_types.Response{
		Status:  "success",
		Message: "Application rolled back successfully",
		Data:    nil,
	}, nil
}

func (c *DeployController) HandleRestart(f fuego.ContextWithBody[types.RestartDeploymentRequest]) (*shared_types.Response, error) {
	c.logger.Log(logger.Info, "starting application restart process", "")

	data, err := f.Body()
	if err != nil {
		if err == io.EOF {
			c.logger.Log(logger.Error, "empty request body received", "id is required for restart")
			return nil, fuego.HTTPError{
				Err:    types.ErrMissingID,
				Status: http.StatusBadRequest,
			}
		}
		c.logger.Log(logger.Error, "failed to read request body", err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	c.logger.Log(logger.Info, "request body parsed successfully", "id: "+data.ID.String())

	if err := c.validator.ValidateRequest(&data); err != nil {
		c.logger.Log(logger.Error, "request validation failed", "id: "+data.ID.String()+", error: "+err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	user := utils.GetUser(f.Response(), f.Request())
	if user == nil {
		c.logger.Log(logger.Error, "user authentication failed", "id: "+data.ID.String())
		return nil, fuego.HTTPError{
			Err:    nil,
			Status: http.StatusUnauthorized,
		}
	}

	c.logger.Log(logger.Info, "attempting to restart application", "id: "+data.ID.String()+", user_id: "+user.ID.String())

	organizationID := utils.GetOrganizationID(f.Request())
	if organizationID == uuid.Nil {
		c.logger.Log(logger.Error, "organization not found", "id: "+data.ID.String())
		return nil, fuego.HTTPError{
			Err:    nil,
			Status: http.StatusUnauthorized,
		}
	}
	err = c.service.RestartDeployment(&data, user.ID, organizationID)
	if err != nil {
		c.logger.Log(logger.Error, "failed to restart application", "id: "+data.ID.String()+", error: "+err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}

	c.logger.Log(logger.Info, "application restarted successfully", "id: "+data.ID.String())
	return &shared_types.Response{
		Status:  "success",
		Message: "Application restarted successfully",
		Data:    nil,
	}, nil
}

func (c *DeployController) HandleGithubWebhook(f fuego.ContextNoBody) (*shared_types.Response, error) {
	c.logger.Log(logger.Info, "handling github webhook", "")

	payload, err := io.ReadAll(f.Request().Body)
	if err != nil {
		c.logger.Log(logger.Error, "failed to read webhook payload", err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	signature := f.Request().Header.Get("X-Hub-Signature-256")
	if signature == "" {
		c.logger.Log(logger.Error, "missing webhook signature", "")
		return nil, fuego.HTTPError{
			Err:    fmt.Errorf("missing webhook signature"),
			Status: http.StatusUnauthorized,
		}
	}

	eventType := f.Request().Header.Get("X-GitHub-Event")
	if eventType != "push" {
		c.logger.Log(logger.Info, "ignoring non-push event", eventType)
		return &shared_types.Response{
			Status:  "success",
			Message: "Ignored non-push event",
			Data:    nil,
		}, nil
	}

	fmt.Printf("payload: %+v\n", string(payload))

	var webhookPayload struct {
		Repository struct {
			ID       uint64 `json:"id"`
			FullName string `json:"full_name"`
		} `json:"repository"`
		Ref    string `json:"ref"`
		Before string `json:"before"`
		After  string `json:"after"`
		Pusher struct {
			Name string `json:"name"`
		} `json:"pusher"`
	}

	if err := json.Unmarshal(payload, &webhookPayload); err != nil {
		c.logger.Log(logger.Error, "failed to parse webhook payload", err.Error())
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	// branch := strings.TrimPrefix(webhookPayload.Ref, "refs/heads/")

	parts := strings.Split(webhookPayload.Repository.FullName, "/")
	if len(parts) != 2 {
		c.logger.Log(logger.Error, "invalid repository name format", webhookPayload.Repository.FullName)
		return nil, fuego.HTTPError{
			Err:    fmt.Errorf("invalid repository name format"),
			Status: http.StatusBadRequest,
		}
	}

	c.logger.Log(logger.Info, "deployment created successfully", webhookPayload.Repository.FullName)
	return &shared_types.Response{
		Status:  "success",
		Message: "Deployment created successfully",
		Data:    nil,
	}, nil
}
