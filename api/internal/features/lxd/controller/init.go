package controller

import (
	"github.com/raghavyuva/nixopus-api/internal/features/lxd/service"
	configTypes "github.com/raghavyuva/nixopus-api/internal/types"
)

type Controller struct {
	svc         *service.ClientService
	defaultCfg  configTypes.LXDConfig
}

func NewController(svc *service.ClientService, defaultCfg configTypes.LXDConfig) *Controller {
	return &Controller{
		svc:        svc,
		defaultCfg: defaultCfg,
	}
}
