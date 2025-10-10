package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	lxdclient "github.com/canonical/lxd/client"
	lxdapi "github.com/canonical/lxd/shared/api"
)

// LXD operations supported by the API
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
}

// New creates a new LXD client service using local unix socket and optional project
func New(socketPath string, project string, opTimeoutSec int) (*ClientService, error) {
	if socketPath == "" {
		// common default for snap
		socketPath = "/var/snap/lxd/common/lxd/unix.socket"
	}
	if opTimeoutSec <= 0 {
		opTimeoutSec = 60
	}

	c, err := lxdclient.ConnectLXDUnix(socketPath, nil)
	if err != nil {
		return nil, err
	}

	if project != "" {
		c = c.UseProject(project)
	}

	return &ClientService{client: c, project: project, timeout: time.Duration(opTimeoutSec) * time.Second}, nil
}

func (s *ClientService) Create(ctx context.Context, name string, imageAlias string, profiles []string, config map[string]string, devices map[string]map[string]string) (*lxdapi.Instance, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if imageAlias == "" {
		return nil, errors.New("image alias is required")
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
