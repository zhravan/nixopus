package organization

import (
	"net/http"

	"github.com/raghavyuva/nixopus-api/internal/storage"
	"github.com/raghavyuva/nixopus-api/internal/types"
	"github.com/raghavyuva/nixopus-api/internal/utils"
)

// func validateAddUserToOrganizationRequest(user types.AddUserToOrganizationRequest) error {
// 	if user.OrganizationID == "" {
// 		return types.ErrMissingOrganizationID
// 	}
// 	return nil
// }

func (c *OrganizationsController) GetOrganizationUsers(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if err := validateGetOrganizationRequest(id); err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	users, err := storage.GetOrganizationUsers(c.app.Store.DB, id, c.app.Ctx)
	if err != nil {
		utils.SendErrorResponse(w, types.ErrFailedToGetOrganizationUsers.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendJSONResponse(w, "success", "Organization users fetched successfully", users)
}

// func (c *OrganizationsController) AddUserToOrganization(w http.ResponseWriter, r *http.Request) {
// 	var user types.AddUserToOrganizationRequest

// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		utils.SendErrorResponse(w, types.ErrFailedToDecodeRequest.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if err := validateAddUserToOrganizationRequest(user); err != nil {
// 		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if err := storage.AddUserToOrganization(c.app.Store.DB, user, c.app.Ctx); err != nil {
// 		utils.SendErrorResponse(w, types.ErrFailedToAddUserToOrganization.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	utils.SendJSONResponse(w, "success", "User added to organization successfully", nil)
// }