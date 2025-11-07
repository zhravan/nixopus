package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/go-fuego/fuego"
	"github.com/raghavyuva/nixopus-api/internal/features/lxd/types"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

func (c *Controller) getServerConfig(reqConfig *types.ServerConfig) shared_types.LXDConfig {
	if reqConfig != nil {
		return reqConfig.ToLXDConfig()
	}
	return c.defaultCfg
}

func (c *Controller) withTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}

func (c *Controller) successResponse(data interface{}, message string) *shared_types.Response {
	resp := &shared_types.Response{Status: "success"}
	if data != nil {
		resp.Data = data
	}
	if message != "" {
		resp.Message = message
	}
	return resp
}

func (c *Controller) Create(ctx fuego.ContextWithBody[types.CreateRequest]) (*shared_types.Response, error) {
	body, err := ctx.Body()
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusBadRequest}
	}
	reqCtx, cancel := c.withTimeout(ctx.Request().Context(), 120*time.Second)
	defer cancel()

	cfg := c.getServerConfig(body.ServerConfig)
	inst, err := c.svc.CreateWithServer(reqCtx, &cfg, body.Name, body.Image, body.Profiles, body.Config, body.Devices)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return c.successResponse(inst, "created"), nil
}

func (c *Controller) List(ctx fuego.ContextWithBody[types.ListRequest]) (*shared_types.Response, error) {
	body, _ := ctx.Body()
	reqCtx, cancel := c.withTimeout(ctx.Request().Context(), 30*time.Second)
	defer cancel()

	cfg := c.getServerConfig(body.ServerConfig)
	list, err := c.svc.ListWithServer(reqCtx, &cfg)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return c.successResponse(list, ""), nil
}

func (c *Controller) Get(ctx fuego.ContextWithBody[types.GetRequest]) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	body, _ := ctx.Body()
	reqCtx, cancel := c.withTimeout(ctx.Request().Context(), 30*time.Second)
	defer cancel()

	cfg := c.getServerConfig(body.ServerConfig)
	inst, err := c.svc.GetWithServer(reqCtx, &cfg, name)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusNotFound}
	}
	return c.successResponse(inst, ""), nil
}

func (c *Controller) Start(ctx fuego.ContextWithBody[types.StartRequest]) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	body, _ := ctx.Body()
	reqCtx, cancel := c.withTimeout(ctx.Request().Context(), 60*time.Second)
	defer cancel()

	cfg := c.getServerConfig(body.ServerConfig)
	err := c.svc.StartWithServer(reqCtx, &cfg, name)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return c.successResponse(nil, "started"), nil
}

func (c *Controller) Stop(ctx fuego.ContextWithBody[types.StopRequest]) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	body, _ := ctx.Body()
	reqCtx, cancel := c.withTimeout(ctx.Request().Context(), 60*time.Second)
	defer cancel()

	cfg := c.getServerConfig(body.ServerConfig)
	err := c.svc.StopWithServer(reqCtx, &cfg, name, true)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return c.successResponse(nil, "stopped"), nil
}

func (c *Controller) Restart(ctx fuego.ContextWithBody[types.RestartRequest]) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	body, _ := ctx.Body()
	reqCtx, cancel := c.withTimeout(ctx.Request().Context(), 90*time.Second)
	defer cancel()

	cfg := c.getServerConfig(body.ServerConfig)
	err := c.svc.RestartWithServer(reqCtx, &cfg, name, 0)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return c.successResponse(nil, "restarted"), nil
}

func (c *Controller) Delete(ctx fuego.ContextWithBody[types.DeleteRequest]) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	body, _ := ctx.Body()
	reqCtx, cancel := c.withTimeout(ctx.Request().Context(), 90*time.Second)
	defer cancel()

	cfg := c.getServerConfig(body.ServerConfig)
	err := c.svc.DeleteWithServer(reqCtx, &cfg, name)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return c.successResponse(nil, "deleted"), nil
}

func (c *Controller) DeleteAll(ctx fuego.ContextWithBody[types.DeleteAllRequest]) (*shared_types.Response, error) {
	body, _ := ctx.Body()
	reqCtx, cancel := c.withTimeout(ctx.Request().Context(), 5*time.Minute)
	defer cancel()

	cfg := c.getServerConfig(body.ServerConfig)
	err := c.svc.DeleteAllWithServer(reqCtx, &cfg)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return c.successResponse(nil, "deleted-all"), nil
}
