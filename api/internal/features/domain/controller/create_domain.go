package controller

import (
	"net/http"

	"github.com/raghavyuva/nixopus-api/internal/features/domain/types"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	"github.com/raghavyuva/nixopus-api/internal/utils"
)

// @Summary Create a new domain
// @Description Creates a new domain in the application.
// @Tags domain
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param domain body types.CreateDomainRequest true "Domain creation request"
// @Success 200 {object} types.CreateDomainResponse "Success response with domain"
// @Failure 400 {object} types.Response "Bad request"
// @Failure 500 {object} types.Response "Internal server error"
// @Router /domain [post]
func (c *DomainsController) CreateDomain(w http.ResponseWriter, r *http.Request) {
	var domainRequest types.CreateDomainRequest

	if !c.parseAndValidate(w, r, &domainRequest) {
		return
	}

	user := utils.GetUser(w, r)

	if user == nil {
		return
	}

	created, err := c.service.CreateDomain(domainRequest, user.ID.String())

	if err != nil {
		c.logger.Log(logger.Error, err.Error(), "")
		utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, "success", "Domain created successfully", created)
}
