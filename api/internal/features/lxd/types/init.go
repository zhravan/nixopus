package types

import "errors"

// Request types
type CreateRequest struct {
	Name     string                       `json:"name"`
	Image    string                       `json:"image"`
	Profiles []string                     `json:"profiles"`
	Config   map[string]string            `json:"config"`
	Devices  map[string]map[string]string `json:"devices"`
}

// Error variables
var (
	ErrMissingName       = errors.New("name is required")
	ErrMissingImageAlias = errors.New("image alias is required")
)
