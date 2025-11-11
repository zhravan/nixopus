package types

import "errors"

type CreateRequest struct {
	Name     string                       `json:"name"`
	Image    string                       `json:"image"`
	Profiles []string                     `json:"profiles"`
	Config   map[string]string            `json:"config"`
	Devices  map[string]map[string]string `json:"devices"`
}

type UpdateRequest struct {
	Profiles []string                     `json:"profiles,omitempty"`
	Config   map[string]string            `json:"config,omitempty"`
	Devices  map[string]map[string]string `json:"devices,omitempty"`
}

var (
	ErrMissingName       = errors.New("name is required")
	ErrMissingImageAlias = errors.New("image alias is required")
)
