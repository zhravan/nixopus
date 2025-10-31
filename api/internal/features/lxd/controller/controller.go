package controller

import (
	"context"
	"net/http"
	"time"

	lxdapi "github.com/canonical/lxd/shared/api"
	"github.com/go-fuego/fuego"
	"github.com/raghavyuva/nixopus-api/internal/features/lxd/types"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

// TODO: cleanup direct controllers

func (c *Controller) Create(ctx fuego.ContextWithBody[types.CreateRequest]) (*shared_types.Response, error) {
	body, err := ctx.Body()
	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusBadRequest}
	}
	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), 120*time.Second)
	defer cancel()

	var inst *lxdapi.Instance
	// Use custom server config if provided, otherwise use default
	if body.ServerConfig != nil {
		cfg := body.ServerConfig.ToLXDConfig()
		inst, err = c.svc.CreateWithServer(reqCtx, &cfg, body.Name, body.Image, body.Profiles, body.Config, body.Devices)
	} else {
		inst, err = c.svc.Create(reqCtx, body.Name, body.Image, body.Profiles, body.Config, body.Devices)
	}

	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return &shared_types.Response{Status: "success", Message: "created", Data: inst}, nil
}

func (c *Controller) List(ctx fuego.ContextWithBody[types.ListRequest]) (*shared_types.Response, error) {
	body, err := ctx.Body()
	if err != nil {
		// If no body provided, use default server
		body = types.ListRequest{}
	}
	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), 30*time.Second)
	defer cancel()

	var list []lxdapi.Instance
	// Use custom server config if esle use default
	if body.ServerConfig != nil {
		cfg := body.ServerConfig.ToLXDConfig()
		list, err = c.svc.ListWithServer(reqCtx, &cfg)
	} else {
		list, err = c.svc.List(reqCtx)
	}

	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return &shared_types.Response{Status: "success", Data: list}, nil
}

func (c *Controller) Get(ctx fuego.ContextWithBody[types.GetRequest]) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	body, err := ctx.Body()
	if err != nil {
		// If no body provided, use default server
		body = types.GetRequest{}
	}
	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), 30*time.Second)
	defer cancel()

	var inst *lxdapi.Instance
	// Use custom server config if provided, otherwise use default
	if body.ServerConfig != nil {
		cfg := body.ServerConfig.ToLXDConfig()
		inst, err = c.svc.GetWithServer(reqCtx, &cfg, name)
	} else {
		inst, err = c.svc.Get(reqCtx, name)
	}

	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusNotFound}
	}
	return &shared_types.Response{Status: "success", Data: inst}, nil
}

func (c *Controller) Start(ctx fuego.ContextWithBody[types.StartRequest]) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	body, err := ctx.Body()
	if err != nil {
		// If no body provided, use default server
		body = types.StartRequest{}
	}
	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), 60*time.Second)
	defer cancel()

	// Use custom server config if provided, otherwise use default
	if body.ServerConfig != nil {
		cfg := body.ServerConfig.ToLXDConfig()
		err = c.svc.StartWithServer(reqCtx, &cfg, name)
	} else {
		err = c.svc.Start(reqCtx, name)
	}

	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return &shared_types.Response{Status: "success", Message: "started"}, nil
}

func (c *Controller) Stop(ctx fuego.ContextWithBody[types.StopRequest]) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	body, err := ctx.Body()
	if err != nil {
		// If no body provided, use default server
		body = types.StopRequest{}
	}
	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), 60*time.Second)
	defer cancel()

	// Use custom server config if provided, otherwise use default
	if body.ServerConfig != nil {
		cfg := body.ServerConfig.ToLXDConfig()
		err = c.svc.StopWithServer(reqCtx, &cfg, name, true)
	} else {
		err = c.svc.Stop(reqCtx, name, true)
	}

	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return &shared_types.Response{Status: "success", Message: "stopped"}, nil
}

func (c *Controller) Restart(ctx fuego.ContextWithBody[types.RestartRequest]) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	body, err := ctx.Body()
	if err != nil {
		// If no body provided, use default server
		body = types.RestartRequest{}
	}
	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), 90*time.Second)
	defer cancel()

	// Use custom server config if provided, otherwise use default
	if body.ServerConfig != nil {
		cfg := body.ServerConfig.ToLXDConfig()
		err = c.svc.RestartWithServer(reqCtx, &cfg, name, 0)
	} else {
		err = c.svc.Restart(reqCtx, name, 0)
	}

	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return &shared_types.Response{Status: "success", Message: "restarted"}, nil
}

func (c *Controller) Delete(ctx fuego.ContextWithBody[types.DeleteRequest]) (*shared_types.Response, error) {
	name := ctx.PathParam("name")
	body, err := ctx.Body()
	if err != nil {
		// If no body provided, use default server
		body = types.DeleteRequest{}
	}
	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), 90*time.Second)
	defer cancel()

	// Use custom server config if provided, otherwise use default
	if body.ServerConfig != nil {
		cfg := body.ServerConfig.ToLXDConfig()
		err = c.svc.DeleteWithServer(reqCtx, &cfg, name)
	} else {
		err = c.svc.Delete(reqCtx, name)
	}

	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return &shared_types.Response{Status: "success", Message: "deleted"}, nil
}

func (c *Controller) DeleteAll(ctx fuego.ContextWithBody[types.DeleteAllRequest]) (*shared_types.Response, error) {
	body, err := ctx.Body()
	if err != nil {
		// If no body provided, use default server
		body = types.DeleteAllRequest{}
	}
	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), 5*time.Minute)
	defer cancel()

	// Use custom server config if provided, otherwise use default
	if body.ServerConfig != nil {
		cfg := body.ServerConfig.ToLXDConfig()
		err = c.svc.DeleteAllWithServer(reqCtx, &cfg)
	} else {
		err = c.svc.DeleteAll(reqCtx)
	}

	if err != nil {
		return nil, fuego.HTTPError{Err: err, Status: http.StatusInternalServerError}
	}
	return &shared_types.Response{Status: "success", Message: "deleted-all"}, nil
}
