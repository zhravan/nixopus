package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"strings"
	"time"

	lxdclient "github.com/canonical/lxd/client"
	lxdapi "github.com/canonical/lxd/shared/api"
	"github.com/raghavyuva/nixopus-api/internal/features/lxd/types"
	configTypes "github.com/raghavyuva/nixopus-api/internal/types"
)

// LXD operations supported by the API

// TODO: refactor/cleanup apis exposure via single interface method for both local and remote connections; temp keeping it separate
type Service interface {
	Create(ctx context.Context, name string, imageAlias string, profiles []string, config map[string]string, devices map[string]map[string]string) (*lxdapi.Instance, error)
	CreateWithServer(ctx context.Context, serverCfg *configTypes.LXDConfig, name string, imageAlias string, profiles []string, config map[string]string, devices map[string]map[string]string) (*lxdapi.Instance, error)
	List(ctx context.Context) ([]lxdapi.Instance, error)
	ListWithServer(ctx context.Context, serverCfg *configTypes.LXDConfig) ([]lxdapi.Instance, error)
	Get(ctx context.Context, name string) (*lxdapi.Instance, error)
	GetWithServer(ctx context.Context, serverCfg *configTypes.LXDConfig, name string) (*lxdapi.Instance, error)
	Start(ctx context.Context, name string) error
	StartWithServer(ctx context.Context, serverCfg *configTypes.LXDConfig, name string) error
	Stop(ctx context.Context, name string, force bool) error
	StopWithServer(ctx context.Context, serverCfg *configTypes.LXDConfig, name string, force bool) error
	Restart(ctx context.Context, name string, timeout time.Duration) error
	RestartWithServer(ctx context.Context, serverCfg *configTypes.LXDConfig, name string, timeout time.Duration) error
	Delete(ctx context.Context, name string) error
	DeleteWithServer(ctx context.Context, serverCfg *configTypes.LXDConfig, name string) error
	DeleteAll(ctx context.Context) error
	DeleteAllWithServer(ctx context.Context, serverCfg *configTypes.LXDConfig) error
}

type ClientService struct {
	client  lxdclient.InstanceServer
	project string
	timeout time.Duration
}

// New creates a new LXD client service from the provided configuration
// Supports both local unix socket and remote HTTPS connections with trust password authentication
func New(cfg configTypes.LXDConfig) (*ClientService, error) {
	// Set default timeout
	opTimeoutSec := cfg.OperationTimeoutSeconds
	if opTimeoutSec <= 0 {
		opTimeoutSec = 60
	}

	var c lxdclient.InstanceServer
	var err error

	// Determine connection type based on protocol and remote address
	if cfg.Protocol == "https" && cfg.RemoteAddress != "" {
		// Remote connection with trust password
		c, err = connectRemote(cfg)
	} else {
		// Local unix socket connection (default)
		c, err = connectLocal(cfg)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to LXD: %w", err)
	}

	// Set project if specified
	if cfg.Project != "" {
		c = c.UseProject(cfg.Project)
	}

	return &ClientService{
		client:  c,
		project: cfg.Project,
		timeout: time.Duration(opTimeoutSec) * time.Second,
	}, nil
}

// connectLocal connects to LXD via local unix socket
func connectLocal(cfg configTypes.LXDConfig) (lxdclient.InstanceServer, error) {
	socketPath := cfg.SocketPath
	if socketPath == "" {
		// common default for snap installations
		socketPath = "/var/snap/lxd/common/lxd/unix.socket"
	}

	c, err := lxdclient.ConnectLXDUnix(socketPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to unix socket %s: %w", socketPath, err)
	}

	return c, nil
}

// connectRemote connects to a remote LXD server using HTTPS with trust password authentication
func connectRemote(cfg configTypes.LXDConfig) (lxdclient.InstanceServer, error) {
	if cfg.RemoteAddress == "" {
		return nil, fmt.Errorf("remote address is required for remote connections")
	}

	// Load system root CA certificates
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		rootCAs = x509.NewCertPool()
	}

	// Configure TLS
	tlsConfig := &tls.Config{
		RootCAs:            rootCAs,
		InsecureSkipVerify: cfg.InsecureSkipVerify,
	}

	// Create custom HTTP transport with TLS config
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	// Prepare connection arguments with custom HTTP client
	args := &lxdclient.ConnectionArgs{
		HTTPClient: &http.Client{
			Transport: transport,
		},
		InsecureSkipVerify: cfg.InsecureSkipVerify,
	}

	// Construct the URL (LXD uses HTTPS on port 8443 by default)
	url := cfg.RemoteAddress
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	// Try to connect to the remote LXD server
	c, err := lxdclient.ConnectLXD(url, args)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to remote LXD at %s: %w", url, err)
	}

	// If trust password is provided, authenticate by adding our certificate
	if cfg.TrustPassword != "" {
		req := lxdapi.CertificatesPost{
			Password: cfg.TrustPassword,
			Type:     "client",
		}

		err = c.CreateCertificate(req)
		if err != nil {
			// If certificate already exists, that's okay
			if !strings.Contains(err.Error(), "already exists") && !strings.Contains(err.Error(), "already trusted") {
				return nil, fmt.Errorf("failed to authenticate with remote LXD server using trust password: %w", err)
			}
		}
	}

	return c, nil
}

func (s *ClientService) Create(ctx context.Context, name string, imageAlias string, profiles []string, config map[string]string, devices map[string]map[string]string) (*lxdapi.Instance, error) {
	if name == "" {
		return nil, types.ErrMissingName
	}
	if imageAlias == "" {
		return nil, types.ErrMissingImageAlias
	}

	req := lxdapi.InstancesPost{
		Name: name,
		InstancePut: lxdapi.InstancePut{
			Config:   config,
			Devices:  mapToDevices(devices),
			Profiles: profiles,
		},
		Source: lxdapi.InstanceSource{
			Type:  "image",
			Alias: imageAlias,
		},
	}

	op, err := s.client.CreateInstance(req)
	if err != nil {
		return nil, err
	}
	if err := waitOp(ctx, op, s.timeout); err != nil {
		return nil, err
	}
	inst, _, err := s.client.GetInstance(name)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

func (s *ClientService) List(ctx context.Context) ([]lxdapi.Instance, error) {
	instances, err := s.client.GetInstances(lxdapi.InstanceTypeAny)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func (s *ClientService) Get(ctx context.Context, name string) (*lxdapi.Instance, error) {
	inst, _, err := s.client.GetInstance(name)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

func (s *ClientService) Start(ctx context.Context, name string) error {
	req := lxdapi.InstanceStatePut{Action: "start", Timeout: int(s.timeout.Seconds()), Force: false, Stateful: false}
	op, err := s.client.UpdateInstanceState(name, req, "")
	if err != nil {
		return err
	}
	return waitOp(ctx, op, s.timeout)
}

func (s *ClientService) Stop(ctx context.Context, name string, force bool) error {
	req := lxdapi.InstanceStatePut{Action: "stop", Timeout: int(s.timeout.Seconds()), Force: force}
	op, err := s.client.UpdateInstanceState(name, req, "")
	if err != nil {
		return err
	}
	return waitOp(ctx, op, s.timeout)
}

func (s *ClientService) Restart(ctx context.Context, name string, timeout time.Duration) error {
	to := int(s.timeout.Seconds())
	if timeout > 0 {
		to = int(timeout.Seconds())
	}
	req := lxdapi.InstanceStatePut{Action: "restart", Timeout: to, Force: true}
	op, err := s.client.UpdateInstanceState(name, req, "")
	if err != nil {
		return err
	}
	return waitOp(ctx, op, time.Duration(to)*time.Second)
}

func (s *ClientService) Delete(ctx context.Context, name string) error {
	// Ensure stopped
	_ = s.Stop(ctx, name, true)
	op, err := s.client.DeleteInstance(name)
	if err != nil {
		return err
	}
	return waitOp(ctx, op, s.timeout)
}

func (s *ClientService) DeleteAll(ctx context.Context) error {
	instances, err := s.client.GetInstances(lxdapi.InstanceTypeAny)
	if err != nil {
		return err
	}
	var errs []string
	for _, inst := range instances {
		if err := s.Delete(ctx, inst.Name); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", inst.Name, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("failed to delete some instances: %s", strings.Join(errs, ", "))
	}
	return nil
}

func mapToDevices(in map[string]map[string]string) map[string]map[string]string {
	if in == nil {
		return map[string]map[string]string{}
	}
	return in
}

// waitOp waits on the LXD operation with a timeout respecting ctx
func waitOp(ctx context.Context, op lxdclient.Operation, timeout time.Duration) error {
	done := make(chan error, 1)
	go func() {
		done <- op.Wait()
	}()
	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return context.DeadlineExceeded
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Methods that accept custom server config in the request

func (s *ClientService) CreateWithServer(ctx context.Context, serverCfg *configTypes.LXDConfig, name string, imageAlias string, profiles []string, config map[string]string, devices map[string]map[string]string) (*lxdapi.Instance, error) {
	tempSvc, err := New(*serverCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary LXD client: %w", err)
	}
	return tempSvc.Create(ctx, name, imageAlias, profiles, config, devices)
}

func (s *ClientService) ListWithServer(ctx context.Context, serverCfg *configTypes.LXDConfig) ([]lxdapi.Instance, error) {
	tempSvc, err := New(*serverCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary LXD client: %w", err)
	}
	return tempSvc.List(ctx)
}

func (s *ClientService) GetWithServer(ctx context.Context, serverCfg *configTypes.LXDConfig, name string) (*lxdapi.Instance, error) {
	tempSvc, err := New(*serverCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary LXD client: %w", err)
	}
	return tempSvc.Get(ctx, name)
}

func (s *ClientService) StartWithServer(ctx context.Context, serverCfg *configTypes.LXDConfig, name string) error {
	tempSvc, err := New(*serverCfg)
	if err != nil {
		return fmt.Errorf("failed to create temporary LXD client: %w", err)
	}
	return tempSvc.Start(ctx, name)
}

func (s *ClientService) StopWithServer(ctx context.Context, serverCfg *configTypes.LXDConfig, name string, force bool) error {
	tempSvc, err := New(*serverCfg)
	if err != nil {
		return fmt.Errorf("failed to create temporary LXD client: %w", err)
	}
	return tempSvc.Stop(ctx, name, force)
}

func (s *ClientService) RestartWithServer(ctx context.Context, serverCfg *configTypes.LXDConfig, name string, timeout time.Duration) error {
	tempSvc, err := New(*serverCfg)
	if err != nil {
		return fmt.Errorf("failed to create temporary LXD client: %w", err)
	}
	return tempSvc.Restart(ctx, name, timeout)
}

func (s *ClientService) DeleteWithServer(ctx context.Context, serverCfg *configTypes.LXDConfig, name string) error {
	tempSvc, err := New(*serverCfg)
	if err != nil {
		return fmt.Errorf("failed to create temporary LXD client: %w", err)
	}
	return tempSvc.Delete(ctx, name)
}

func (s *ClientService) DeleteAllWithServer(ctx context.Context, serverCfg *configTypes.LXDConfig) error {
	tempSvc, err := New(*serverCfg)
	if err != nil {
		return fmt.Errorf("failed to create temporary LXD client: %w", err)
	}
	return tempSvc.DeleteAll(ctx)
}
