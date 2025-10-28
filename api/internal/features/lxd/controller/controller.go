package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/go-fuego/fuego"
	"github.com/raghavyuva/nixopus-api/internal/features/lxd/types"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

func (c *Controller) Create(ctx fuego.ContextWithBody[types.CreateRequest]) (*shared_types.Response, error) {
	body, err := ctx.Body()
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusBadRequest}
	}
	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), 120*time.Second)
	defer cancel()
	inst, err := c.svc.Create(reqCtx, body.Name, body.Image, body.Profiles, body.Config, body.Devices)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return &shared_types.Response{Status: "success", Message: "created", Data: inst}, nil
}

func (c *Controller) List(ctx fuego.ContextNoBody) (*shared_types.Response, error) {
	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), 30*time.Second)
	defer cancel()
	list, err := c.svc.List(reqCtx)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return &shared_types.Response{Status: "success", Data: list}, nil
}

func (c *Controller) Get(ctx fuego.ContextNoBody) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), 30*time.Second)
	defer cancel()
	inst, err := c.svc.Get(reqCtx, name)
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusNotFound}
	}
	return &shared_types.Response{Status: "success", Data: inst}, nil
}

func (c *Controller) Start(ctx fuego.ContextNoBody) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), 60*time.Second)
	defer cancel()
	if err := c.svc.Start(reqCtx, name); err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return &shared_types.Response{Status: "success", Message: "started"}, nil
}

func (c *Controller) Stop(ctx fuego.ContextNoBody) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), 60*time.Second)
	defer cancel()
	if err := c.svc.Stop(reqCtx, name, true); err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return &shared_types.Response{Status: "success", Message: "stopped"}, nil
}

func (c *Controller) Restart(ctx fuego.ContextNoBody) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), 90*time.Second)
	defer cancel()
	if err := c.svc.Restart(reqCtx, name, 0); err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return &shared_types.Response{Status: "success", Message: "restarted"}, nil
}

func (c *Controller) Delete(ctx fuego.ContextNoBody) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), 90*time.Second)
	defer cancel()
	if err := c.svc.Delete(reqCtx, name); err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return &shared_types.Response{Status: "success", Message: "deleted"}, nil
}

func (c *Controller) DeleteAll(ctx fuego.ContextNoBody) (*shared_types.Response, error) {
	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), 5*time.Minute)
	defer cancel()
	if err := c.svc.DeleteAll(reqCtx); err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return &shared_types.Response{Status: "success", Message: "deleted-all"}, nil
}
