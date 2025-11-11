package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	lxdclient "github.com/canonical/lxd/client"
	lxdapi "github.com/canonical/lxd/shared/api"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	"github.com/raghavyuva/nixopus-api/internal/features/lxd/types"
	configTypes "github.com/raghavyuva/nixopus-api/internal/types"
)

type Service interface {
	Create(ctx context.Context, name string, imageAlias string, profiles []string, config map[string]string, devices map[string]map[string]string) (*lxdapi.Instance, error)
	List(ctx context.Context) ([]lxdapi.Instance, error)
	Get(ctx context.Context, name string) (*lxdapi.Instance, error)
	Start(ctx context.Context, name string) error
	Stop(ctx context.Context, name string, force bool) error
	Restart(ctx context.Context, name string, timeout time.Duration) error
	Delete(ctx context.Context, name string) error
	DeleteAll(ctx context.Context) error
}

type ClientService struct {
	client  lxdclient.InstanceServer
	project string
	timeout time.Duration
	logger  logger.Logger
}

// New creates a new LXD client service from the provided configuration
// Supports both local unix socket and remote HTTPS connections with trust password authentication
func New(cfg configTypes.LXDConfig, l logger.Logger) (*ClientService, error) {
	opTimeoutSec := cfg.OperationTimeoutSeconds
	if opTimeoutSec <= 0 {
		opTimeoutSec = 60
	}

	var c lxdclient.InstanceServer
	var err error

	if cfg.Protocol == "https" && cfg.RemoteAddress != "" {
		c, err = connectRemote(cfg, l)
	} else {
		c, err = connectLocal(cfg, l)
	}

	if err != nil {
		l.Log(logger.Error, fmt.Sprintf("failed to connect to LXD: %v", err), "")
		return nil, fmt.Errorf("failed to connect to LXD: %w", err)
	}

	if cfg.Project != "" {
		c = c.UseProject(cfg.Project)
	}

	return &ClientService{
		client:  c,
		project: cfg.Project,
		timeout: time.Duration(opTimeoutSec) * time.Second,
		logger:  l,
	}, nil
}

func connectLocal(cfg configTypes.LXDConfig, l logger.Logger) (lxdclient.InstanceServer, error) {
	socketPath := cfg.SocketPath
	if socketPath == "" {
		var err error
		socketPath, err = detectSocketPath(l)
		if err != nil {
			return nil, fmt.Errorf("failed to detect LXD socket path: %w", err)
		}
	}

	c, err := lxdclient.ConnectLXDUnix(socketPath, nil)
	if err != nil {
		l.Log(logger.Error, fmt.Sprintf("failed to connect to unix socket %s: %v", socketPath, err), "")
		return nil, fmt.Errorf("failed to connect to unix socket %s: %w", socketPath, err)
	}

	return c, nil
}

// checks in all the common paths for unix socket for lxd
func detectSocketPath(l logger.Logger) (string, error) {
	commonPaths := []string{
		"/var/snap/lxd/common/lxd/unix.socket",
		"/var/lib/lxd/unix.socket",
		"/run/lxd.socket",
		"/var/run/lxd.socket",
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			l.Log(logger.Info, fmt.Sprintf("detected LXD socket at %s", path), "")
			return path, nil
		}
	}

	return "", fmt.Errorf("LXD socket not found in common locations: %v. Please configure socket_path explicitly", commonPaths)
}

func connectRemote(cfg configTypes.LXDConfig, l logger.Logger) (lxdclient.InstanceServer, error) {
	if cfg.RemoteAddress == "" {
		l.Log(logger.Error, "remote address is required for remote connections", "")
		return nil, fmt.Errorf("remote address is required for remote connections")
	}

	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		rootCAs = x509.NewCertPool()
	}

	tlsConfig := &tls.Config{
		RootCAs:            rootCAs,
		InsecureSkipVerify: cfg.InsecureSkipVerify,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	args := &lxdclient.ConnectionArgs{
		HTTPClient: &http.Client{
			Transport: transport,
		},
		InsecureSkipVerify: cfg.InsecureSkipVerify,
	}

	url := cfg.RemoteAddress
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	c, err := lxdclient.ConnectLXD(url, args)
	if err != nil {
		l.Log(logger.Error, fmt.Sprintf("failed to connect to remote LXD at %s: %v", url, err), "")
		return nil, fmt.Errorf("failed to connect to remote LXD at %s: %w", url, err)
	}

	if cfg.TrustPassword != "" {
		req := lxdapi.CertificatesPost{
			Password: cfg.TrustPassword,
			Type:     "client",
		}

		err = c.CreateCertificate(req)
		if err != nil {
			if isCertificateAlreadyExistsError(err) {
				l.Log(logger.Info, "certificate already exists or is already trusted, continuing", "")
			} else {
				l.Log(logger.Error, fmt.Sprintf("failed to authenticate with remote LXD server using trust password: %v", err), "")
				return nil, fmt.Errorf("failed to authenticate with remote LXD server using trust password: %w", err)
			}
		}
	}

	return c, nil
}

func isCertificateAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	patterns := []string{
		"already exists",
		"already trusted",
		"certificate already",
		"already present",
		"duplicate certificate",
	}

	for _, pattern := range patterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	return false
}

func (s *ClientService) isImageCached(alias string) bool {
	images, err := s.client.GetImages()
	if err != nil {
		return false
	}

	for _, img := range images {
		for _, a := range img.Aliases {
			if a.Name == alias {
				return true
			}
		}
	}
	return false
}

func (s *ClientService) Create(ctx context.Context, name string, imageAlias string, profiles []string, config map[string]string, devices map[string]map[string]string) (*lxdapi.Instance, error) {
	if name == "" {
		return nil, types.ErrMissingName
	}
	if imageAlias == "" {
		return nil, types.ErrMissingImageAlias
	}

	source := lxdapi.InstanceSource{
		Type:  "image",
		Alias: imageAlias,
	}

	if !s.isImageCached(imageAlias) {
		source.Server = "https://images.lxd.canonical.com"
		source.Protocol = "simplestreams"
		source.Mode = "pull"
	}

	req := lxdapi.InstancesPost{
		Name: name,
		InstancePut: lxdapi.InstancePut{
			Config:   config,
			Devices:  mapToDevices(devices),
			Profiles: profiles,
		},
		Source: source,
	}

	op, err := s.client.CreateInstance(req)
	if err != nil {
		s.logger.Log(logger.Error, fmt.Sprintf("failed to create instance %s: %v", name, err), "")
		return nil, err
	}
	if err := waitOp(ctx, op, s.timeout); err != nil {
		s.logger.Log(logger.Error, fmt.Sprintf("failed to wait for instance creation %s: %v", name, err), "")
		return nil, err
	}
	inst, _, err := s.client.GetInstance(name)
	if err != nil {
		s.logger.Log(logger.Error, fmt.Sprintf("failed to get instance %s after creation: %v", name, err), "")
		return nil, err
	}
	return inst, nil
}

func (s *ClientService) List(ctx context.Context) ([]lxdapi.Instance, error) {
	instances, err := s.client.GetInstances(lxdapi.InstanceTypeAny)
	if err != nil {
		s.logger.Log(logger.Error, fmt.Sprintf("failed to list instances: %v", err), "")
		return nil, err
	}
	return instances, nil
}

func (s *ClientService) Get(ctx context.Context, name string) (*lxdapi.Instance, error) {
	inst, _, err := s.client.GetInstance(name)
	if err != nil {
		s.logger.Log(logger.Error, fmt.Sprintf("failed to get instance %s: %v", name, err), "")
		return nil, err
	}
	return inst, nil
}

func (s *ClientService) Start(ctx context.Context, name string) error {
	req := lxdapi.InstanceStatePut{Action: "start", Timeout: int(s.timeout.Seconds()), Force: false, Stateful: false}
	op, err := s.client.UpdateInstanceState(name, req, "")
	if err != nil {
		s.logger.Log(logger.Error, fmt.Sprintf("failed to start instance %s: %v", name, err), "")
		return err
	}
	if err := waitOp(ctx, op, s.timeout); err != nil {
		s.logger.Log(logger.Error, fmt.Sprintf("failed to wait for instance start %s: %v", name, err), "")
		return err
	}
	return nil
}

func (s *ClientService) Stop(ctx context.Context, name string, force bool) error {
	req := lxdapi.InstanceStatePut{Action: "stop", Timeout: int(s.timeout.Seconds()), Force: force}
	op, err := s.client.UpdateInstanceState(name, req, "")
	if err != nil {
		s.logger.Log(logger.Error, fmt.Sprintf("failed to stop instance %s: %v", name, err), "")
		return err
	}
	if err := waitOp(ctx, op, s.timeout); err != nil {
		s.logger.Log(logger.Error, fmt.Sprintf("failed to wait for instance stop %s: %v", name, err), "")
		return err
	}
	return nil
}

func (s *ClientService) Restart(ctx context.Context, name string, timeout time.Duration) error {
	to := int(s.timeout.Seconds())
	if timeout > 0 {
		to = int(timeout.Seconds())
	}
	req := lxdapi.InstanceStatePut{Action: "restart", Timeout: to, Force: true}
	op, err := s.client.UpdateInstanceState(name, req, "")
	if err != nil {
		s.logger.Log(logger.Error, fmt.Sprintf("failed to restart instance %s: %v", name, err), "")
		return err
	}
	if err := waitOp(ctx, op, time.Duration(to)*time.Second); err != nil {
		s.logger.Log(logger.Error, fmt.Sprintf("failed to wait for instance restart %s: %v", name, err), "")
		return err
	}
	return nil
}

func (s *ClientService) Delete(ctx context.Context, name string) error {
	stopErr := s.Stop(ctx, name, true)
	if stopErr != nil {
		s.logger.Log(logger.Warning, fmt.Sprintf("failed to stop instance %s before deletion (attempting deletion anyway): %v", name, stopErr), "")
	}

	op, err := s.client.DeleteInstance(name)
	if err != nil {
		s.logger.Log(logger.Error, fmt.Sprintf("failed to delete instance %s: %v", name, err), "")
		return err
	}
	if err := waitOp(ctx, op, s.timeout); err != nil {
		s.logger.Log(logger.Error, fmt.Sprintf("failed to wait for instance deletion %s: %v", name, err), "")
		return err
	}
	return nil
}

func (s *ClientService) DeleteAll(ctx context.Context) error {
	instances, err := s.client.GetInstances(lxdapi.InstanceTypeAny)
	if err != nil {
		s.logger.Log(logger.Error, fmt.Sprintf("failed to list instances for deletion: %v", err), "")
		return err
	}
	var errs []string
	for _, inst := range instances {
		if err := s.Delete(ctx, inst.Name); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", inst.Name, err))
		}
	}
	if len(errs) > 0 {
		errMsg := fmt.Sprintf("failed to delete some instances: %s", strings.Join(errs, ", "))
		s.logger.Log(logger.Error, errMsg, "")
		return fmt.Errorf(errMsg)
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
