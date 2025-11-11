package controller

import (
	"github.com/raghavyuva/nixopus-api/internal/features/lxd/service"
)

type Controller struct {
	svc *service.ClientService
}

func NewController(svc *service.ClientService) *Controller {
	return &Controller{
		svc: svc,
	}
}
