package realtime

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/raghavyuva/nixopus-api/internal/features/dashboard"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	"github.com/raghavyuva/nixopus-api/internal/types"
)

func (s *SocketServer) handleStopDashboardMonitor(conn *websocket.Conn) {
	s.dashboardMutex.Lock()
	defer s.dashboardMutex.Unlock()
	if monitor, exists := s.dashboardMonitors[conn]; exists {
		monitor.Stop()
		delete(s.dashboardMonitors, conn)

		conn.WriteJSON(types.Payload{
			Action: "dashboard_monitor_stopped",
			Data:   nil,
		})
	}
}

func (s *SocketServer) handleDashboardMonitor(conn *websocket.Conn, msg types.Payload) {
	s.dashboardMutex.Lock()
	monitor, exists := s.dashboardMonitors[conn]
	if !exists {
		newMonitor, err := dashboard.NewDashboardMonitor(conn, logger.NewLogger())
		if err != nil {
			s.dashboardMutex.Unlock()
			s.sendError(conn, "Failed to create dashboard monitor")
			return
		}

		s.dashboardMonitors[conn] = newMonitor
		monitor = newMonitor
	}
	s.dashboardMutex.Unlock()

	if msg.Data != nil {
		dataMap, ok := msg.Data.(map[string]interface{})
		if !ok {
			s.sendError(conn, "Invalid dashboard monitor configuration")
			return
		}

		var interval time.Duration
		if intervalSec, ok := dataMap["interval"].(float64); ok {
			interval = time.Duration(intervalSec) * time.Second
		} else {
			interval = 10 * time.Second
		}

		var operations []dashboard.DashboardOperation
		if ops, ok := dataMap["operations"].([]interface{}); ok {
			for _, op := range ops {
				if opStr, ok := op.(string); ok {
					operations = append(operations, dashboard.DashboardOperation(opStr))
				}
			}
		}

		if len(operations) == 0 {
			operations = dashboard.AllOperations
		}

		config := dashboard.MonitoringConfig{
			Interval:   interval,
			Operations: operations,
		}

		monitor.Interval = config.Interval
		monitor.Operations = config.Operations

		monitor.Start()

		response := types.Payload{
			Action: "dashboard_monitor_started",
			Data: map[string]interface{}{
				"interval":   interval.Seconds(),
				"operations": operations,
			},
		}

		jsonData, err := json.Marshal(response)
		if err != nil {
			s.sendError(conn, "Failed to marshal response")
			return
		}

		conn.WriteMessage(websocket.TextMessage, jsonData)
	} else {
		monitor.Stop()
		response := types.Payload{
			Action: "dashboard_monitor_stopped",
			Data:   nil,
		}

		jsonData, err := json.Marshal(response)
		if err != nil {
			s.sendError(conn, "Failed to marshal response")
			return
		}

		conn.WriteMessage(websocket.TextMessage, jsonData)
	}
}
