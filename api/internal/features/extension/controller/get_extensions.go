package controller

import (
	"net/http"

	"github.com/go-fuego/fuego"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	"github.com/raghavyuva/nixopus-api/internal/types"
)

func (c *ExtensionsController) GetExtensions(ctx fuego.ContextNoBody) ([]types.Extension, error) {
	categoryParam := ctx.QueryParam("category")

	var category *types.ExtensionCategory
	if categoryParam != "" {
		cat := types.ExtensionCategory(categoryParam)
		category = &cat
	}

	extensions, err := c.service.ListExtensions(category)
	if err != nil {
		c.logger.Log(logger.Error, err.Error(), "")
		return nil, fuego.HTTPError{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}

	return extensions, nil
}

func (c *ExtensionsController) GetExtension(ctx fuego.ContextNoBody) (types.Extension, error) {
	id := ctx.PathParam("id")
	if id == "" {
		return types.Extension{}, fuego.HTTPError{
			Err:    nil,
			Status: http.StatusBadRequest,
		}
	}

	extension, err := c.service.GetExtension(id)
	if err != nil {
		if err.Error() == "extension not found" {
			return types.Extension{}, fuego.HTTPError{
				Err:    err,
				Status: http.StatusNotFound,
			}
		}
		c.logger.Log(logger.Error, err.Error(), "")
		return types.Extension{}, fuego.HTTPError{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}

	return *extension, nil
}

func (c *ExtensionsController) GetExtensionByExtensionID(ctx fuego.ContextNoBody) (types.Extension, error) {
	extensionID := ctx.PathParam("extension_id")
	if extensionID == "" {
		return types.Extension{}, fuego.HTTPError{
			Err:    nil,
			Status: http.StatusBadRequest,
		}
	}

	extension, err := c.service.GetExtensionByID(extensionID)
	if err != nil {
		if err.Error() == "extension not found" {
			return types.Extension{}, fuego.HTTPError{
				Err:    err,
				Status: http.StatusNotFound,
			}
		}
		c.logger.Log(logger.Error, err.Error(), "")
		return types.Extension{}, fuego.HTTPError{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}

	return *extension, nil
}
