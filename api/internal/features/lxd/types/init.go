package types

import (
	"errors"

	configTypes "github.com/raghavyuva/nixopus-api/internal/types"
)

// ServerConfig allows overriding LXD server connection per request (optional)
type ServerConfig struct {
	Protocol           string `json:"protocol,omitempty"`             // unix or https
	SocketPath         string `json:"socket_path,omitempty"`          // for local unix socket
	RemoteAddress      string `json:"remote_address,omitempty"`       // for remote (e.g., "10.0.0.1:8443")
	TrustPassword      string `json:"trust_password,omitempty"`       // trust password for remote auth
	Project            string `json:"project,omitempty"`              // LXD project name
	InsecureSkipVerify bool   `json:"insecure_skip_verify,omitempty"` // skip TLS verification
	Timeout            int    `json:"timeout,omitempty"`              // operation timeout in seconds
}

// ToLXDConfig converts ServerConfig to types.LXDConfig
func (sc *ServerConfig) ToLXDConfig() configTypes.LXDConfig {
	timeout := sc.Timeout
	if timeout <= 0 {
		timeout = 60 // default
	}
	return configTypes.LXDConfig{
		Enabled:                 true,
		SocketPath:              sc.SocketPath,
		Project:                 sc.Project,
		OperationTimeoutSeconds: timeout,
		RemoteAddress:           sc.RemoteAddress,
		Protocol:                sc.Protocol,
		TrustPassword:           sc.TrustPassword,
		InsecureSkipVerify:      sc.InsecureSkipVerify,
	}
}

// Request types
type CreateRequest struct {
	Name         string                       `json:"name"`
	Image        string                       `json:"image"`
	Profiles     []string                     `json:"profiles"`
	Config       map[string]string            `json:"config"`
	Devices      map[string]map[string]string `json:"devices"`
	ServerConfig *ServerConfig                `json:"server_config,omitempty"` // optional: override default server
}

type ListRequest struct {
	ServerConfig *ServerConfig `json:"server_config,omitempty"` // optional: override default server
}

type GetRequest struct {
	ServerConfig *ServerConfig `json:"server_config,omitempty"` // optional: override default server
}

type StartRequest struct {
	ServerConfig *ServerConfig `json:"server_config,omitempty"` // optional: override default server
}

type StopRequest struct {
	ServerConfig *ServerConfig `json:"server_config,omitempty"` // optional: override default server
}

type RestartRequest struct {
	ServerConfig *ServerConfig `json:"server_config,omitempty"` // optional: override default server
}

type DeleteRequest struct {
	ServerConfig *ServerConfig `json:"server_config,omitempty"` // optional: override default server
}

type DeleteAllRequest struct {
	ServerConfig *ServerConfig `json:"server_config,omitempty"` // optional: override default server
}

// Error variables
var (
	ErrMissingName       = errors.New("name is required")
	ErrMissingImageAlias = errors.New("image alias is required")
)
