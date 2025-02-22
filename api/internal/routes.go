package internal

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/raghavyuva/nixopus-api/internal/controller"
	"github.com/raghavyuva/nixopus-api/internal/controller/auth"
	"github.com/raghavyuva/nixopus-api/internal/controller/organization"
	"github.com/raghavyuva/nixopus-api/internal/middleware"
	"github.com/raghavyuva/nixopus-api/internal/storage"
)

type Router struct {
	app *storage.App
}

func NewRouter(app *storage.App) *Router {
	return &Router{
		app: app,
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

	r.Use(middleware.CorsMiddleware)
	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/health", controller.HealthCheck).Methods("GET", "OPTIONS")

	wsServer, err := NewSocketServer()
	if err != nil {
		log.Fatal(err)
	}

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsServer.HandleHTTP(w, r)
	})

	u := r.PathPrefix("/api/v1").Subrouter()

	// Unauthenticated routes
	authController := auth.NewAuthController(router.app)
	u.HandleFunc("/auth/register", authController.Register).Methods("POST", "OPTIONS")
	u.HandleFunc("/auth/login", authController.Login).Methods("POST", "OPTIONS")

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(func(next http.Handler) http.Handler {
		return middleware.AuthMiddleware(next, router.app)
	})

	// Authenticated routes
	api.HandleFunc("/auth/request-password-reset", authController.GeneratePasswordResetLink).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/reset-password", authController.ResetPassword).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/refresh-token", authController.RefreshToken).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/logout", authController.Logout).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/send-verification-email", authController.SendVerificationEmail).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/verify-email", authController.VerifyEmail).Methods("POST", "OPTIONS")

	// Role based routes
	roleController := organization.NewRolesController(router.app)
	api.HandleFunc("/roles", roleController.CreateRole).Methods("POST", "OPTIONS")
	// api.HandleFunc("/roles/{id}", roleController.GetRole).Methods("GET", "OPTIONS")
	api.HandleFunc("/roles/{id}", roleController.UpdateRole).Methods("PUT", "OPTIONS")
	api.HandleFunc("/roles/{id}", roleController.DeleteRole).Methods("DELETE", "OPTIONS")
	api.HandleFunc("/roles", roleController.GetRoles).Methods("GET", "OPTIONS")

	// Organization Routes
	organizationController := organization.NewOrganizationsController(router.app)
	api.HandleFunc("/organizations", organizationController.CreateOrganization).Methods("POST", "OPTIONS")
	api.HandleFunc("/organizations/{id}", organizationController.GetOrganization).Methods("GET", "OPTIONS")
	api.HandleFunc("/organizations/{id}", organizationController.UpdateOrganization).Methods("PUT", "OPTIONS")
	api.HandleFunc("/organizations/{id}", organizationController.DeleteOrganization).Methods("DELETE", "OPTIONS")
	api.HandleFunc("/organizations", organizationController.GetOrganizations).Methods("GET", "OPTIONS")

	// Permission Routes
	permissionController := organization.NewPermissionsController(router.app)
	api.HandleFunc("/permissions", permissionController.CreatePermission).Methods("POST", "OPTIONS")
	api.HandleFunc("/permission/", permissionController.GetPermission).Methods("GET", "OPTIONS")
	api.HandleFunc("/permissions/update", permissionController.UpdatePermission).Methods("PUT", "OPTIONS")
	api.HandleFunc("/permissions/delete", permissionController.DeletePermission).Methods("DELETE", "OPTIONS")
	api.HandleFunc("/permissions", permissionController.GetPermissions).Methods("GET", "OPTIONS")	
	api.HandleFunc("/add-permission/roles", permissionController.AddPermissionToRole).Methods("POST", "OPTIONS")
	api.HandleFunc("/remove-permission/roles", permissionController.RemovePermissionFromRole).Methods("DELETE", "OPTIONS")

	api.HandleFunc("/permissions/roles", permissionController.GetPermissionsByRole).Methods("GET", "OPTIONS")

	// User Routes
	// userController := controller.NewUserController(router.app)
	// api.HandleFunc("/users", userController.CreateUser).Methods("POST", "OPTIONS")
	// api.HandleFunc("/users/{id}", userController.GetUser).Methods("GET", "OPTIONS")
	// api.HandleFunc("/users/{id}", userController.UpdateUser).Methods("PUT", "OPTIONS")
	// api.HandleFunc("/users/{id}", userController.DeleteUser).Methods("DELETE", "OPTIONS")

	// User based routes

	return r
}
