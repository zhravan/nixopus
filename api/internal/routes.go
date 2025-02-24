package internal

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	auth "github.com/raghavyuva/nixopus-api/internal/features/auth/controller"
	health "github.com/raghavyuva/nixopus-api/internal/features/health"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	organization "github.com/raghavyuva/nixopus-api/internal/features/organization/controller"
	permission "github.com/raghavyuva/nixopus-api/internal/features/permission/controller"
	role "github.com/raghavyuva/nixopus-api/internal/features/role/controller"
	"github.com/raghavyuva/nixopus-api/internal/middleware"
	"github.com/raghavyuva/nixopus-api/internal/storage"
)

type Router struct {
	app            *storage.App
}

func NewRouter(app *storage.App) *Router {
	return &Router{
		app:            app,
	}
}

// Routes returns a new router that handles all API routes, including
// unauthenticated and authenticated routes.
// The following middleware is used:
//
// - middleware.CorsMiddleware: enables CORS
// - middleware.LoggingMiddleware: logs all requests
// - middleware.AuthMiddleware: checks if the request is authenticated
func (router *Router) Routes() *mux.Router {
	r := mux.NewRouter()
	l := logger.NewLogger()
	r.Use(middleware.CorsMiddleware)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RateLimiter)
	
	r.HandleFunc("/health", health.HealthCheck).Methods("GET", "OPTIONS")

	wsServer, err := NewSocketServer()
	if err != nil {
		log.Fatal(err)
	}

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsServer.HandleHTTP(w, r)
	})

	u := r.PathPrefix("/api/v1").Subrouter()

	authCB := u.PathPrefix("/auth").Subrouter()

	authController := auth.NewAuthController(router.app.Store, router.app.Ctx, l)
	authCB.HandleFunc("/register", authController.Register).Methods("POST", "OPTIONS")
	authCB.HandleFunc("/login", authController.Login).Methods("POST", "OPTIONS")

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(func(next http.Handler) http.Handler {
		return middleware.AuthMiddleware(next, router.app)
	})

	authApi := api.PathPrefix("/auth").Subrouter()
	authApi.HandleFunc("/request-password-reset", authController.GeneratePasswordResetLink).Methods("POST", "OPTIONS")
	authApi.HandleFunc("/reset-password", authController.ResetPassword).Methods("POST", "OPTIONS")
	authApi.HandleFunc("/refresh-token", authController.RefreshToken).Methods("POST", "OPTIONS")
	authApi.HandleFunc("/logout", authController.Logout).Methods("POST", "OPTIONS")
	authApi.HandleFunc("/send-verification-email", authController.SendVerificationEmail).Methods("POST", "OPTIONS")
	authApi.HandleFunc("/verify-email", authController.VerifyEmail).Methods("POST", "OPTIONS")

	roleApi := api.PathPrefix("/roles").Subrouter()
	roleController := role.NewRolesController(router.app.Store, router.app.Ctx, l)
	roleApi.HandleFunc("", roleController.CreateRole).Methods("POST", "OPTIONS")
	// roleApi.HandleFunc("/{id}", roleController.GetRole).Methods("GET", "OPTIONS")
	roleApi.HandleFunc("/", roleController.UpdateRole).Methods("PUT", "OPTIONS")
	roleApi.HandleFunc("/", roleController.DeleteRole).Methods("DELETE", "OPTIONS")
	roleApi.HandleFunc("", roleController.GetRoles).Methods("GET", "OPTIONS")

	orgApi := api.PathPrefix("/organizations").Subrouter()
	organizationController := organization.NewOrganizationsController(router.app.Store, router.app.Ctx, l)
	orgApi.HandleFunc("", organizationController.CreateOrganization).Methods("POST", "OPTIONS")
	orgApi.HandleFunc("", organizationController.GetOrganization).Methods("GET", "OPTIONS")
	orgApi.HandleFunc("", organizationController.UpdateOrganization).Methods("PUT", "OPTIONS")
	orgApi.HandleFunc("", organizationController.DeleteOrganization).Methods("DELETE", "OPTIONS")
	orgApi.HandleFunc("", organizationController.GetOrganizations).Methods("GET", "OPTIONS")
	orgApi.HandleFunc("/user", organizationController.AddUserToOrganization).Methods("POST", "OPTIONS")
	orgApi.HandleFunc("/user", organizationController.RemoveUserFromOrganization).Methods("DELETE", "OPTIONS")
	orgApi.HandleFunc("/users", organizationController.GetOrganizationUsers).Methods("GET", "OPTIONS")

	permApi := api.PathPrefix("/permissions").Subrouter()
	permissionController := permission.NewPermissionController(router.app.Store, router.app.Ctx, l)
	permApi.HandleFunc("", permissionController.CreatePermission).Methods("POST", "OPTIONS")
	permApi.HandleFunc("", permissionController.GetPermission).Methods("GET", "OPTIONS")
	permApi.HandleFunc("", permissionController.UpdatePermission).Methods("PUT", "OPTIONS")
	permApi.HandleFunc("", permissionController.DeletePermission).Methods("DELETE", "OPTIONS")
	permApi.HandleFunc("", permissionController.GetPermissions).Methods("GET", "OPTIONS")
	
	rolePermApi := api.PathPrefix("/roles/permission").Subrouter()
	rolePermApi.HandleFunc("", permissionController.AddPermissionToRole).Methods("POST", "OPTIONS")
	rolePermApi.HandleFunc("", permissionController.RemovePermissionFromRole).Methods("DELETE", "OPTIONS")
	rolePermApi.HandleFunc("s", permissionController.GetPermissionsByRole).Methods("GET", "OPTIONS")

	return r
}