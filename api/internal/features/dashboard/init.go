package dashboard

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/raghavyuva/nixopus-api/internal/features/deploy/docker"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	sshpkg "github.com/raghavyuva/nixopus-api/internal/features/ssh"
)

func NewDashboardMonitor(conn *websocket.Conn, log logger.Logger) (*DashboardMonitor, error) {
	ssh_client := sshpkg.NewSSH()
	ctx, cancel := context.WithCancel(context.Background())

	monitor := &DashboardMonitor{
		conn:          conn,
		sshpkg:        ssh_client,
		log:           log,
		ctx:           ctx,
		cancel:        cancel,
		Interval:      time.Second * 10,
		Operations:    AllOperations,
		dockerService: docker.NewDockerService(),
	}

	return monitor, nil
}

func (m *DashboardMonitor) Start() {
	go func() {
		ticker := time.NewTicker(m.Interval)
		defer ticker.Stop()
		client, err := m.sshpkg.Connect()
		if err != nil {
			m.log.Log(logger.Error, "Failed to connect to SSH server", err.Error())
			m.BroadcastError(err.Error(), "ssh_connect")
			return
		}
		m.client = client
		defer client.Close()
		m.HandleAllOperations()

		for {
			select {
			case <-ticker.C:
				m.HandleAllOperations()
			case <-m.ctx.Done():
				m.log.Log(logger.Info, "Dashboard monitor stopped", "")
				return
			}
		}
	}()
}

func (m *DashboardMonitor) Stop() {
	m.log.Log(logger.Info, "Stopping dashboard monitor", "")
	if m.cancel != nil {
		m.cancel()
	}
	m.Close()
}

func (m *DashboardMonitor) HandleAllOperations() {
	for _, operation := range m.Operations {
		select {
		case <-m.ctx.Done():
			return
		default:
			m.HandleOperation(operation)
		}
	}
}

func (m *DashboardMonitor) HandleOperation(operation DashboardOperation) {
	switch operation {
	case GetContainers:
		m.GetContainers()
	case GetSystemStats:
		m.GetSystemStats()
	default:
		m.log.Log(logger.Error, "Unknown operation", string(operation))
	}
}

func (m *DashboardMonitor) UpdateConfig(config MonitoringConfig) {
	m.Interval = config.Interval

	if len(config.Operations) > 0 {
		m.Operations = config.Operations
	}

	m.Stop()
	m.Start()
}

func (m *DashboardMonitor) SetOperations(operations []DashboardOperation) {
	m.Operations = operations

	m.Stop()
	m.Start()
}

func (m *DashboardMonitor) Close() {
	if m.client != nil {
		m.client.Close()
		m.client = nil
	}

	if m.conn != nil {
		m.connMutex.Lock()
		_ = m.conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing connection"),
		)
		_ = m.conn.Close()
		m.conn = nil
		m.connMutex.Unlock()
	}
}
