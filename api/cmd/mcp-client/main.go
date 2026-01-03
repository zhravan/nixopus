package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()

	// Get the path to the MCP server binary
	serverPath := os.Getenv("MCP_SERVER_PATH")
	if serverPath == "" {
		log.Fatalf("MCP_SERVER_PATH environment variable is required.\nExample: export MCP_SERVER_PATH=./nixopus-mcp-server")
	}

	// Create a new MCP client
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "nixopus-mcp-client",
		Version: "1.0.0",
	}, nil)

	// Connect to the server over stdin/stdout using CommandTransport
	transport := &mcp.CommandTransport{
		Command: exec.Command(serverPath),
	}

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatalf("Failed to connect to MCP server: %v", err)
	}
	defer session.Close()

	fmt.Println("Connected to MCP server")

	// List available tools
	fmt.Println("\nListing available tools...")
	toolsResult, err := session.ListTools(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	fmt.Printf("Found %d tool(s):\n", len(toolsResult.Tools))
	for _, tool := range toolsResult.Tools {
		fmt.Printf("  - %s: %s\n", tool.Name, tool.Description)
	}

	// Example: Call get_container_logs tool
	if len(os.Args) > 1 && os.Args[1] == "test" {
		fmt.Println("\nTesting get_container_logs tool...")

		// Get container ID and organization ID from environment or use defaults
		containerID := os.Getenv("CONTAINER_ID")
		organizationID := os.Getenv("ORGANIZATION_ID")
		authToken := os.Getenv("AUTH_TOKEN")

		if containerID == "" {
			containerID = "test-container-id"
		}
		if organizationID == "" {
			organizationID = "test-org-id"
		}
		if authToken == "" {
			fmt.Println("Warning: AUTH_TOKEN not set. Authentication will fail.")
			fmt.Println("   Set AUTH_TOKEN environment variable with a valid SuperTokens session token.")
		}

		params := &mcp.CallToolParams{
			Name: "get_container_logs",
			Arguments: map[string]any{
				"id":              containerID,
				"organization_id": organizationID,
				"follow":          false,
				"tail":            100,
				"stdout":          true,
				"stderr":          true,
			},
		}

		// Add auth token to metadata if provided
		if authToken != "" {
			params.Meta = mcp.Meta{
				"auth_token": authToken,
			}
		}

		res, err := session.CallTool(ctx, params)
		if err != nil {
			log.Fatalf("CallTool failed: %v", err)
		}

		if res.IsError {
			fmt.Printf("Tool returned an error:\n")
			if len(res.Content) > 0 {
				for _, c := range res.Content {
					if textContent, ok := c.(*mcp.TextContent); ok {
						fmt.Printf("  Error: %s\n", textContent.Text)
					}
				}
			}
			os.Exit(1)
		}

		fmt.Println("Tool call successful:")
		for _, c := range res.Content {
			switch content := c.(type) {
			case *mcp.TextContent:
				fmt.Printf("  Text: %s\n", content.Text)
			case *mcp.ImageContent:
				fmt.Printf("  Image: %d bytes (mime: %s)\n", len(content.Data), content.MIMEType)
			case *mcp.AudioContent:
				fmt.Printf("  Audio: %d bytes (mime: %s)\n", len(content.Data), content.MIMEType)
			default:
				// Try to marshal unknown content types
				jsonBytes, err := json.MarshalIndent(content, "  ", "  ")
				if err == nil {
					fmt.Printf("  Content: %s\n", string(jsonBytes))
				} else {
					fmt.Printf("  Content: %+v\n", content)
				}
			}
		}

		// Also check structured content if available
		if res.StructuredContent != nil {
			fmt.Println("\nStructured Content:")
			jsonBytes, err := json.MarshalIndent(res.StructuredContent, "  ", "  ")
			if err == nil {
				fmt.Printf("%s\n", string(jsonBytes))
			} else {
				fmt.Printf("%+v\n", res.StructuredContent)
			}
		}
	} else {
		fmt.Println("\nUsage:")
		fmt.Println("  Set MCP_SERVER_PATH environment variable:")
		fmt.Println("    export MCP_SERVER_PATH=./nixopus-mcp-server")
		fmt.Println("")
		fmt.Println("  To test a tool call, run:")
		fmt.Println("    CONTAINER_ID=<id> ORGANIZATION_ID=<org-id> AUTH_TOKEN=<token> go run ./cmd/mcp-client test")
		fmt.Println("")
		fmt.Println("  Note: AUTH_TOKEN is required for authentication.")
	}
}
