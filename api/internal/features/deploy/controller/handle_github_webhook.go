package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-fuego/fuego"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

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
