package realtime

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/raghavyuva/nixopus-api/internal/features/dashboard"
	deploy "github.com/raghavyuva/nixopus-api/internal/features/deploy/controller"
	"github.com/raghavyuva/nixopus-api/internal/features/deploy/realtime"
	"github.com/raghavyuva/nixopus-api/internal/features/terminal"
	"github.com/raghavyuva/nixopus-api/internal/types"
	"github.com/uptrace/bun"
)

const (
	maxMessageSize = 10 * 1024 * 1024
	pingInterval   = 10 * time.Second
	pingTimeout    = 60 * time.Second
)

type topics string

const (
	MonitorApplicationDeployment topics = "monitor_application_deployment"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type SocketServer struct {
	conns               *sync.Map
	topicsMu            sync.RWMutex
	topics              map[string]map[*websocket.Conn]bool
	shutdown            chan struct{}
	deployController    *deploy.DeployController
	db                  *bun.DB
	ctx                 context.Context
	postgres_listener   PostgresListener
	terminalMutex       sync.RWMutex
	terminals           map[*websocket.Conn]*terminal.Terminal
	dashboardMonitors   map[*websocket.Conn]*dashboard.DashboardMonitor
	dashboardMutex      sync.Mutex
	applicationMonitors map[*websocket.Conn]*realtime.ApplicationMonitor
	applicationMutex    sync.Mutex
}

// NewSocketServer initializes and returns a new instance of SocketServer.
func NewSocketServer(deployController *deploy.DeployController, db *bun.DB, ctx context.Context) (*SocketServer, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pgListener := NewPostgresListener()

	server := &SocketServer{
		conns:               &sync.Map{},
		shutdown:            make(chan struct{}),
		deployController:    deployController,
		db:                  db,
		ctx:                 ctx,
		topics:              make(map[string]map[*websocket.Conn]bool),
		postgres_listener:   *pgListener,
		terminals:           make(map[*websocket.Conn]*terminal.Terminal),
		dashboardMonitors:   make(map[*websocket.Conn]*dashboard.DashboardMonitor),
		applicationMonitors: make(map[*websocket.Conn]*realtime.ApplicationMonitor),
	}

	notificationChan, err := pgListener.ListenToApplicationChanges(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to listen for PostgreSQL notifications: %w", err)
	}

	go server.handleNotifications(notificationChan)

	return server, nil
}

// HandleHTTP handles incoming HTTP connections and upgrades them to WebSocket connections.
// It verifies the authorization token and authenticates the user.
//
// Parameters:
//   w - the http.ResponseWriter for the response.
//   r - the http.Request from the client.
//
// Returns:
//   - nil if the connection is successfully upgraded.
//   - an error if the connection fails to upgrade or if the token is invalid.
func (s *SocketServer) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	conn.SetReadLimit(maxMessageSize)

	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}

	user, err := s.verifyToken(token)
	if err != nil {
		log.Printf("Auth error: %v", err)
		s.sendError(conn, "Invalid authorization token")
		conn.Close()
		return
	}

	s.conns.Store(conn, user.ID)
	defer s.handleDisconnect(conn)

	log.Printf("User authenticated: %s", user.ID)
	s.readLoop(conn, user)
}

func (s *SocketServer) handleDisconnect(conn *websocket.Conn) {
	userID, _ := s.conns.Load(conn)
	s.conns.Delete(conn)

	s.topicsMu.Lock()
	for topic, connections := range s.topics {
		if _, exists := connections[conn]; exists {
			delete(connections, conn)
			if len(connections) == 0 {
				delete(s.topics, topic)
			}
		}
	}

	s.topicsMu.Unlock()

	conn.Close()
	fmt.Printf("Client disconnected: %s (User ID: %v)\n", conn.RemoteAddr(), userID)
}

func (s *SocketServer) Shutdown() {
	close(s.shutdown)
	s.conns.Range(func(conn, _ interface{}) bool {
		conn.(*websocket.Conn).Close()
		return true
	})
}

func (s *SocketServer) handlePing(conn *websocket.Conn) {
	userID, _ := s.conns.Load(conn)
	fmt.Printf("Received ping from %s (User ID: %v)\n", conn.RemoteAddr(), userID)
}

func (s *SocketServer) sendError(conn *websocket.Conn, message string) {
	conn.WriteJSON(types.Payload{
		Action: "error",
		Data:   message,
	})
}
