package tools

// GetContainerLogsInput is the input structure for the MCP tool
type GetContainerLogsInput struct {
	ID             string  `json:"id" jsonschema:"required"`
	OrganizationID string  `json:"organization_id" jsonschema:"required"`
	Follow         bool    `json:"follow,omitempty"`
	Tail           *int    `json:"tail,omitempty"`
	Since          *string `json:"since,omitempty"`
	Until          *string `json:"until,omitempty"`
	Stdout         bool    `json:"stdout,omitempty"`
	Stderr         bool    `json:"stderr,omitempty"`
}

// GetOrganizationID implements OrganizationIDExtractor interface
func (i GetContainerLogsInput) GetOrganizationID() string {
	return i.OrganizationID
}

// GetContainerLogsOutput is the output structure for the MCP tool
type GetContainerLogsOutput struct {
	Logs string `json:"logs"`
}
