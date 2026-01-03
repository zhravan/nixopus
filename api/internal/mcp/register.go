package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	container_tools "github.com/raghavyuva/nixopus-api/internal/features/container/tools"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	mcp_middleware "github.com/raghavyuva/nixopus-api/internal/mcp/middleware"
	shared_storage "github.com/raghavyuva/nixopus-api/internal/storage"
)

// RegisterTool registers an MCP tool with automatic authentication and authorization middleware
func RegisterTool[Input mcp_middleware.OrganizationIDExtractor, Output any](
	server *mcp.Server,
	store *shared_storage.Store,
	ctx context.Context,
	l logger.Logger,
	tool *mcp.Tool,
	handler func(context.Context, *mcp.CallToolRequest, Input) (*mcp.CallToolResult, Output, error),
) {
	// Automatically wrap handler with auth middleware
	wrappedHandler := mcp_middleware.WithAuth(store, l, handler)
	mcp.AddTool(server, tool, wrappedHandler)
}

// RegisterTools registers all MCP tools with the server
func RegisterTools(
	server *mcp.Server,
	store *shared_storage.Store,
	ctx context.Context,
	l logger.Logger,
) {
	containerLogsHandler := container_tools.GetContainerLogsHandler(store, ctx, l)
	RegisterTool(server, store, ctx, l, &mcp.Tool{
		Name:        "get_container_logs",
		Description: "Get logs from a Docker container. Requires container ID and organization ID.",
	}, containerLogsHandler)
}
