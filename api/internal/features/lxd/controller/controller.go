package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/go-fuego/fuego"
	"github.com/raghavyuva/nixopus-api/internal/features/lxd/types"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

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

	inst, err := c.svc.Create(reqCtx, body.Name, body.Image, body.Profiles, body.Config, body.Devices)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return c.successResponse(inst, "created"), nil
}

func (c *Controller) List(ctx fuego.ContextNoBody) (*shared_types.Response, error) {
	reqCtx, cancel := c.withTimeout(ctx.Request().Context(), 30*time.Second)
	defer cancel()

	list, err := c.svc.List(reqCtx)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return c.successResponse(list, ""), nil
}

func (c *Controller) Get(ctx fuego.ContextNoBody) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	reqCtx, cancel := c.withTimeout(ctx.Request().Context(), 30*time.Second)
	defer cancel()

	inst, err := c.svc.Get(reqCtx, name)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusNotFound}
	}
	return c.successResponse(inst, ""), nil
}

func (c *Controller) Start(ctx fuego.ContextNoBody) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	reqCtx, cancel := c.withTimeout(ctx.Request().Context(), 60*time.Second)
	defer cancel()

	err := c.svc.Start(reqCtx, name)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return c.successResponse(nil, "started"), nil
}

func (c *Controller) Stop(ctx fuego.ContextNoBody) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	reqCtx, cancel := c.withTimeout(ctx.Request().Context(), 60*time.Second)
	defer cancel()

	err := c.svc.Stop(reqCtx, name, true)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return c.successResponse(nil, "stopped"), nil
}

func (c *Controller) Restart(ctx fuego.ContextNoBody) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	reqCtx, cancel := c.withTimeout(ctx.Request().Context(), 90*time.Second)
	defer cancel()

	err := c.svc.Restart(reqCtx, name, 0)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return c.successResponse(nil, "restarted"), nil
}

func (c *Controller) Delete(ctx fuego.ContextNoBody) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	reqCtx, cancel := c.withTimeout(ctx.Request().Context(), 90*time.Second)
	defer cancel()

	err := c.svc.Delete(reqCtx, name)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return c.successResponse(nil, "deleted"), nil
}

func (c *Controller) DeleteAll(ctx fuego.ContextNoBody) (*shared_types.Response, error) {
	reqCtx, cancel := c.withTimeout(ctx.Request().Context(), 5*time.Minute)
	defer cancel()

	err := c.svc.DeleteAll(reqCtx)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return c.successResponse(nil, "deleted-all"), nil
}
