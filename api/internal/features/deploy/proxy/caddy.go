package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/raghavyuva/nixopus-api/internal/features/logger"
)

func NewCaddy(logger *logger.Logger, rootDir string, domain string, port string, fileServerType FileServerType) *Caddy {
	endpoint := os.Getenv("CADDY_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://127.0.0.1:2019"
	}
	return &Caddy{
		Logger:   logger,
		Endpoint: endpoint,
		RootDir:  rootDir,
		Domain:   domain,
		Port:     port,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		FileServerType: fileServerType,
	}
}

func (c *Caddy) Serve() error {
	if err := c.checkCaddyRunning(); err != nil {
		return fmt.Errorf("caddy is not running: %w", err)
	}

	currentConfig, err := c.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get current config: %w", err)
	}

	var handle interface{}
	if c.FileServerType == FileServer {
		handle = FileServerHandle{
			Handler: string(c.FileServerType),
			Root:    c.RootDir,
			Browse:  struct{}{},
		}
	} else {
		subroute := SubrouteHandle{
			Handler: "subroute",
			Routes: []Route{
				{
					Handle: []interface{}{
						ReverseProxyHandle{
							Handler: string(c.FileServerType),
							Upstreams: []Upstream{
								{
									Dial: "http://" + os.Getenv("SSH_HOST") + ":" + c.Port,
								},
							},
						},
					},
				},
			},
		}
		handle = subroute
	}

	routeConfig := Route{
		Match: []Match{
			{
				Host: []string{c.Domain},
			},
		},
		Handle: []interface{}{handle},
	}

	if currentConfig.Apps.HTTP.Servers == nil {
		currentConfig.Apps.HTTP.Servers = make(map[string]Server)
	}
	server := currentConfig.Apps.HTTP.Servers["nixopus"]

	routeExists := false
	for i, route := range server.Routes {
		if len(route.Match) > 0 && len(route.Match[0].Host) > 0 && route.Match[0].Host[0] == c.Domain {
			server.Routes[i] = routeConfig
			routeExists = true
			break
		}
	}

	if !routeExists {
		server.Routes = append(server.Routes, routeConfig)
	}

	currentConfig.Apps.HTTP.Servers["nixopus"] = server

	if err := c.loadConfig(currentConfig); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	c.Logger.Log(logger.Info, "Caddy server started successfully", "")
	return nil
}

func (c *Caddy) checkCaddyRunning() error {
	resp, err := c.client.Get(c.Endpoint + "/config/")
	if err != nil {
		return fmt.Errorf("failed to connect to Caddy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("caddy is not running or not accessible: %s", resp.Status)
	}

	return nil
}

func (c *Caddy) loadConfig(config CaddyConfig) error {
	jsonData, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	req, err := http.NewRequest("POST", c.Endpoint+"/load", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "must-revalidate")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to load config: %s - %s", resp.Status, string(body))
	}

	return nil
}

func (c *Caddy) Stop() error {
	req, err := http.NewRequest("POST", c.Endpoint+"/stop", nil)
	if err != nil {
		return fmt.Errorf("failed to create stop request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send stop request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to stop Caddy: %s", resp.Status)
	}

	return nil
}

func (c *Caddy) GetConfig() (CaddyConfig, error) {
	resp, err := c.client.Get(c.Endpoint + "/config/")
	if err != nil {
		return CaddyConfig{}, fmt.Errorf("failed to get config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CaddyConfig{}, fmt.Errorf("failed to get Caddy config: %s", resp.Status)
	}

	var config CaddyConfig
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return CaddyConfig{}, fmt.Errorf("failed to decode config: %w", err)
	}

	return config, nil
}

func (c *Caddy) UpdateConfig(config CaddyConfig) error {
	jsonData, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	req, err := http.NewRequest("POST", c.Endpoint+"/load", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "must-revalidate")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update Caddy config: %s - %s", resp.Status, string(body))
	}

	return nil
}
