package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/raghavyuva/nixopus-api/internal/features/audit/service"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	"github.com/raghavyuva/nixopus-api/internal/storage"
	"github.com/raghavyuva/nixopus-api/internal/types"
)

// AuditMiddleware is a middleware that captures audit logs for all authenticated requests
func AuditMiddleware(next http.Handler, app *storage.App, l logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auditService := service.NewAuditService(app.Store.DB, app.Ctx, l)

		user, ok := r.Context().Value(types.UserContextKey).(*types.User)
		if !ok {
			l.Log(logger.Debug, "Audit middleware skipped", fmt.Sprintf("No user in context, path: %s", r.URL.Path))
			next.ServeHTTP(w, r)
			return
		}

		orgIDStr, ok := r.Context().Value(types.OrganizationIDKey).(string)
		if !ok {
			l.Log(logger.Debug, "Audit middleware skipped", fmt.Sprintf("No organization ID in context, path: %s, user_id: %s", r.URL.Path, user.ID))
			next.ServeHTTP(w, r)
			return
		}

		orgID, err := uuid.Parse(orgIDStr)
		if err != nil {
			l.Log(logger.Warning, "Audit middleware skipped", fmt.Sprintf("Invalid organization ID, path: %s, user_id: %s, org_id: %s, error: %s", r.URL.Path, user.ID, orgIDStr, err.Error()))
			next.ServeHTTP(w, r)
			return
		}

		auditAction := getAuditActionFromMethod(r.Method)
		auditResourceType := getResourceTypeFromPath(r.URL.Path)

		// Skip if the audit action is access (GET requests typically only read data)
		if auditAction == "access" {
			next.ServeHTTP(w, r)
			return
		}

		auditReq := &service.AuditLogRequest{
			UserID:         user.ID,
			OrganizationID: orgID,
			Action:         types.AuditAction(auditAction),
			ResourceType:   types.AuditResourceType(auditResourceType),
			ResourceID:     uuid.Nil,
			OldValues:      nil,
			NewValues:      nil,
			Metadata: map[string]interface{}{
				"path":   r.URL.Path,
				"method": r.Method,
			},
			IPAddress: r.RemoteAddr,
			UserAgent: r.UserAgent(),
			RequestID: uuid.New(),
		}

		next.ServeHTTP(w, r)

		if err := auditService.LogAction(auditReq); err != nil {
			l.Log(logger.Warning, "Failed to create audit log", fmt.Sprintf("path: %s, method: %s, user_id: %s, org_id: %s, error: %s",
				r.URL.Path, r.Method, user.ID, orgID, err.Error()))
		} else {
			l.Log(logger.Debug, "Audit log created", fmt.Sprintf("path: %s, method: %s, user_id: %s, org_id: %s",
				r.URL.Path, r.Method, user.ID, orgID))
		}
	})
}

func getAuditActionFromMethod(method string) string {
	switch method {
	case http.MethodGet:
		return "access"
	case http.MethodPost:
		return "create"
	case http.MethodPut, http.MethodPatch:
		return "update"
	case http.MethodDelete:
		return "delete"
	default:
		return "access"
	}
}

// getResourceTypeFromPath extracts the resource type from the URL path
func getResourceTypeFromPath(path string) string {
	if len(path) > 8 && path[:8] == "/api/v1/" {
		path = path[8:]
	}
	segments := strings.Split(path, "/")
	if len(segments) == 0 {
		return "application"
	}
	switch segments[0] {
	case "auth":
		return "user"
	case "user":
		if len(segments) > 1 && segments[1] == "organizations" {
			return "organization"
		}
		return "user"
	case "organizations":
		if len(segments) > 1 && segments[1] == "users" {
			return "user"
		}
		return "organization"
	case "roles":
		return "role"
	case "permissions":
		return "permission"
	case "applications":
		return "application"
	case "deploy":
		return "deployment"
	case "deployments":
		return "deployment"
	case "domains":
		return "domain"
	case "github-connector":
		return "github_connector"
	case "smtp":
		return "smtp_config"
	case "file-manager":
		return "application"
	case "audit":
		return "application"
	default:
		return "application"
	}
}
