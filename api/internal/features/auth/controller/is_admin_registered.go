package auth

import (
	"net/http"

	"github.com/go-fuego/fuego"

	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

func (c *AuthController) IsAdminRegistered(s fuego.ContextNoBody) (shared_types.Response, error) {
	isAdminRegistered, err := c.service.IsAdminRegistered()
	if err != nil {
		return shared_types.Response{}, fuego.HTTPError{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}

	if !isAdminRegistered {
		return shared_types.Response{}, fuego.HTTPError{
			Err:    nil,
			Status: http.StatusUnauthorized,
		}
	}

	return shared_types.Response{
		Status:  "success",
		Message: "Admin registered",
		Data:    isAdminRegistered,
	}, nil
}
