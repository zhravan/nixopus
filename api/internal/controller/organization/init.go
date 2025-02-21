package organization

import "github.com/raghavyuva/nixopus-api/internal/storage"

type OrganizationsController struct {
	app *storage.App
}

func NewOrganizationsController(app *storage.App) *OrganizationsController {
	return &OrganizationsController{
		app: app,
	}
}
