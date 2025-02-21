package controller

import (
	"encoding/json"
	"net/http"

	"github.com/raghavyuva/nixopus-api/internal/storage"
	"github.com/raghavyuva/nixopus-api/internal/types"
	"github.com/raghavyuva/nixopus-api/internal/utils"
)

type RolesController struct {
	app *storage.App
}

func NewRolesController(app *storage.App) *RolesController {
	return &RolesController{
		app: app,
	}
}

func validateCreateRoleRequest(role types.CreateRoleRequest) error {
	if role.Name == "" {
		return types.ErrRoleNameRequired
	}
	return nil
}

func validateGetRoleRequest(id string) error {
	if id == "" {
		return types.ErrRoleIDRequired
	}
	return nil
}

func validateUpdateRoleRequest(role types.UpdateRoleRequest) error {
	if role.ID == "" {
		return types.ErrRoleIDRequired
	}
	if role.Name == "" && role.Description == "" {
		return types.ErrRoleEmptyFields
	}
	return nil
}

// CreateRole creates a new role in the database
func (c *RolesController) CreateRole(w http.ResponseWriter, r *http.Request) {
	var role types.CreateRoleRequest

	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		utils.SendErrorResponse(w, types.ErrFailedToDecodeRequest.Error(), http.StatusBadRequest)
		return
	}

	if err := validateCreateRoleRequest(role); err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := storage.CreateRole(c.app.Store.DB, role, c.app.Ctx); err != nil {
		utils.SendErrorResponse(w, types.ErrFailedToCreateRole.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, "success", "Role created successfully", nil)
}

// GetRoles returns all roles that are active in the database
// passing is_disabled as true will return all roles
func (c *RolesController) GetRoles(w http.ResponseWriter, r *http.Request) {
	isDisabled := r.URL.Query().Get("is_disabled") == "true"

	roles, err := storage.GetRoles(c.app.Store.DB, isDisabled, c.app.Ctx)
	if err != nil {
		utils.SendErrorResponse(w, types.ErrFailedToGetRoles.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, "success", "Roles fetched successfully", roles)
}

// GetRole returns a single role for the given id in the database
func (c *RolesController) GetRole(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if err := validateGetRoleRequest(id); err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	role, err := storage.GetRole(c.app.Store.DB, id, c.app.Ctx)
	if err != nil {
		utils.SendErrorResponse(w, types.ErrFailedToGetRole.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, "success", "Role fetched successfully", role)
}

// UpdateRole updates a role in the database
// Takes in four parameters: id and name, description, isDeleted is optional
func (c *RolesController) UpdateRole(w http.ResponseWriter, r *http.Request) {
	var role types.UpdateRoleRequest

	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		utils.SendErrorResponse(w, types.ErrFailedToDecodeRequest.Error(), http.StatusBadRequest)
		return
	}

	if err := validateUpdateRoleRequest(role); err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := storage.UpdateRole(c.app.Store.DB, &role, c.app.Ctx); err != nil {
		utils.SendErrorResponse(w, types.ErrFailedToUpdateRole.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, "success", "Role updated successfully", nil)
}

// DeleteRole deletes a role from the database
func (c *RolesController) DeleteRole(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if err := validateGetRoleRequest(id); err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := storage.DeleteRole(c.app.Store.DB, id, c.app.Ctx); err != nil {
		utils.SendErrorResponse(w, types.ErrFailedToDeleteRole.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, "success", "Role deleted successfully", nil)
}