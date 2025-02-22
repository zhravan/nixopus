package organization

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/raghavyuva/nixopus-api/internal/storage"
	"github.com/raghavyuva/nixopus-api/internal/types"
	"github.com/raghavyuva/nixopus-api/internal/utils"
)

func validateGetOrganizationRequest(id string) error {
	if id == "" {
		return types.ErrMissingOrganizationID
	}
	return nil
}

func validateCreateOrganizationRequest(organization types.CreateOrganizationRequest) error {
	if organization.Name == "" {
		return types.ErrMissingOrganizationName
	}
	return nil
}

func validateUpdateOrganizationRequest(organization types.UpdateOrganizationRequest) error {
	if organization.Name == "" {
		return types.ErrMissingOrganizationName
	}
	return nil
}


func (c *OrganizationsController) GetOrganizations(w http.ResponseWriter, r *http.Request) {
	organization, err := storage.GetOrganizations(c.app.Store.DB, c.app.Ctx)
	if err != nil {
		utils.SendErrorResponse(w, types.ErrFailedToGetOrganizations.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendJSONResponse(w, "success", "Organizations fetched successfully", organization)
}

func (c *OrganizationsController) GetOrganization(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if err := validateGetOrganizationRequest(id); err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	organization, err := storage.GetOrganization(c.app.Store.DB, id, c.app.Ctx)
	if err != nil {
		utils.SendErrorResponse(w, types.ErrFailedToGetOrganization.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendJSONResponse(w, "success", "Organization fetched successfully", organization)
}

func (c *OrganizationsController) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	var organization types.CreateOrganizationRequest

	if err := json.NewDecoder(r.Body).Decode(&organization); err != nil {
		utils.SendErrorResponse(w, types.ErrFailedToDecodeRequest.Error(), http.StatusBadRequest)
		return
	}

	if err := validateCreateOrganizationRequest(organization); err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingOrganization, err := storage.GetOrganizationByName(c.app.Store.DB, organization.Name, c.app.Ctx)
	if err == nil && existingOrganization.ID != uuid.Nil {
		utils.SendErrorResponse(w, types.ErrOrganizationAlreadyExists.Error(), http.StatusBadRequest)
		return
	}

	if err := storage.CreateOrganization(c.app.Store.DB, organization, c.app.Ctx); err != nil {
		utils.SendErrorResponse(w, types.ErrFailedToCreateOrganization.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, "success", "Organization created successfully", nil)
}

func (c *OrganizationsController) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	var organization types.UpdateOrganizationRequest

	if err := json.NewDecoder(r.Body).Decode(&organization); err != nil {
		utils.SendErrorResponse(w, types.ErrFailedToDecodeRequest.Error(), http.StatusBadRequest)
		return
	}

	if err := validateUpdateOrganizationRequest(organization); err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingOrganization, err := storage.GetOrganization(c.app.Store.DB, organization.ID, c.app.Ctx)
	if err == nil && existingOrganization.ID == uuid.Nil {
		utils.SendErrorResponse(w, types.ErrOrganizationDoesNotExist.Error(), http.StatusBadRequest)
		return
	}

	if err := storage.UpdateOrganization(c.app.Store.DB, &organization, c.app.Ctx); err != nil {
		utils.SendErrorResponse(w, types.ErrFailedToUpdateOrganization.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, "success", "Organization updated successfully", nil)
}

func (c *OrganizationsController) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	var organization types.DeleteOrganizationRequest

	if err := json.NewDecoder(r.Body).Decode(&organization); err != nil {
		utils.SendErrorResponse(w, types.ErrFailedToDecodeRequest.Error(), http.StatusBadRequest)
		return
	}

	if err := validateGetOrganizationRequest(organization.ID); err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingOrganization, err := storage.GetOrganization(c.app.Store.DB, organization.ID, c.app.Ctx)
	if err == nil && existingOrganization.ID == uuid.Nil {
		utils.SendErrorResponse(w, types.ErrOrganizationDoesNotExist.Error(), http.StatusBadRequest)
		return
	}

	if err := storage.DeleteOrganization(c.app.Store.DB, organization.ID, c.app.Ctx); err != nil {
		utils.SendErrorResponse(w, types.ErrFailedToDeleteOrganization.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, "success", "Organization deleted successfully", nil)
}
