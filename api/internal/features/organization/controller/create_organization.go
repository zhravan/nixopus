package controller

import (
	"net/http"

	"github.com/go-fuego/fuego"
	"github.com/raghavyuva/nixopus-api/internal/features/notification"
	"github.com/raghavyuva/nixopus-api/internal/features/organization/types"
	"github.com/raghavyuva/nixopus-api/internal/utils"

	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

func (c *OrganizationsController) CreateOrganization(f fuego.ContextWithBody[types.CreateOrganizationRequest]) (*shared_types.Response, error) {
	organization, err := f.Body()
	if err != nil {
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	w, r := f.Response(), f.Request()
	if !c.parseAndValidate(w, r, &organization) {
		return nil, fuego.HTTPError{
			Err:    nil,
			Status: http.StatusBadRequest,
		}
	}

	loggedInUser := utils.GetUser(w, r)
	if loggedInUser == nil {
		return nil, fuego.HTTPError{
			Err:    nil,
			Status: http.StatusUnauthorized,
		}
	}

	createdOrganization, err := c.service.CreateOrganization(&organization)
	if err != nil {
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}

	c.Notify(notification.NortificationPayloadTypeCreateOrganization, loggedInUser, r)

	return &shared_types.Response{
		Status:  "success",
		Message: "Organization created successfully",
		Data:    createdOrganization,
	}, nil
}
