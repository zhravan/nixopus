package controller

import (
	"net/http"

	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	"github.com/raghavyuva/nixopus-api/internal/utils"
)


func (c *OrganizationsController) GetOrganizations(w http.ResponseWriter, r *http.Request) {
	organization, err := c.service.GetOrganizations();
	if err != nil {
		c.logger.Log(logger.Error, err.Error(), "")
		utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendJSONResponse(w, "success", "Organizations fetched successfully", organization)
}