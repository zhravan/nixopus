package internal

import (
	"log"
	"net/http"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"github.com/joho/godotenv"
	audit "github.com/raghavyuva/nixopus-api/internal/features/audit/controller"
	auth "github.com/raghavyuva/nixopus-api/internal/features/auth/controller"
	authService "github.com/raghavyuva/nixopus-api/internal/features/auth/service"
	user_storage "github.com/raghavyuva/nixopus-api/internal/features/auth/storage"
	deploy "github.com/raghavyuva/nixopus-api/internal/features/deploy/controller"
	domain "github.com/raghavyuva/nixopus-api/internal/features/domain/controller"
	file_manager "github.com/raghavyuva/nixopus-api/internal/features/file-manager/controller"
	githubConnector "github.com/raghavyuva/nixopus-api/internal/features/github-connector/controller"
	health "github.com/raghavyuva/nixopus-api/internal/features/health"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	"github.com/raghavyuva/nixopus-api/internal/features/notification"
	notificationController "github.com/raghavyuva/nixopus-api/internal/features/notification/controller"
	organization "github.com/raghavyuva/nixopus-api/internal/features/organization/controller"
	organization_service "github.com/raghavyuva/nixopus-api/internal/features/organization/service"
	organization_storage "github.com/raghavyuva/nixopus-api/internal/features/organization/storage"
	permissions_service "github.com/raghavyuva/nixopus-api/internal/features/permission/service"
	permissions_storage "github.com/raghavyuva/nixopus-api/internal/features/permission/storage"
	role_service "github.com/raghavyuva/nixopus-api/internal/features/role/service"
	role_storage "github.com/raghavyuva/nixopus-api/internal/features/role/storage"
	user "github.com/raghavyuva/nixopus-api/internal/features/user/controller"
	"github.com/raghavyuva/nixopus-api/internal/middleware"
	"github.com/raghavyuva/nixopus-api/internal/realtime"
	"github.com/raghavyuva/nixopus-api/internal/storage"
	api "github.com/raghavyuva/nixopus-api/internal/version-manager"
)

type Router struct {
	app *storage.App
}

func NewRouter(app *storage.App) *Router {
	return &Router{
		app: app,
	}
}

func (router *Router) Routes() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	PORT := os.Getenv("PORT")

	docs := api.NewVersionDocumentation()
	if err := docs.Save("api/versions.json"); err != nil {
		log.Printf("Warning: Failed to save version documentation: %v", err)
	}

	l := logger.NewLogger()
	server := fuego.NewServer(
		fuego.WithGlobalMiddlewares(
			middleware.RecoveryMiddleware,
			middleware.CorsMiddleware,
			middleware.LoggingMiddleware,
			api.VersionMiddleware,
			api.MigrationMiddleware,
			// middleware.RateLimiter
		),
		fuego.WithSecurity(openapi3.SecuritySchemes{
			"bearerAuth": &openapi3.SecuritySchemeRef{
				Value: openapi3.NewSecurityScheme().
					WithType("http").
					WithScheme("bearer").
					WithBearerFormat("JWT").
					WithDescription("Enter your JWT token in the format: Bearer <token>"),
			},
		}),
		fuego.WithAddr(":"+PORT),
	)

	apiV1 := api.NewVersion(api.CurrentVersion)
	healthGroup := fuego.Group(server, apiV1.Path+"/health")
	router.BasicRoutes(healthGroup)

	notificationManager := notification.NewNotificationManager(router.app.Store.DB)
	notificationManager.Start()
	deployController := deploy.NewDeployController(router.app.Store, router.app.Ctx, l, notificationManager)
	router.WebSocketServer(server, deployController)

	userStorage := &user_storage.UserStorage{DB: router.app.Store.DB, Ctx: router.app.Ctx}
	permStorage := &permissions_storage.PermissionStorage{DB: router.app.Store.DB, Ctx: router.app.Ctx}
	roleStorage := &role_storage.RoleStorage{DB: router.app.Store.DB, Ctx: router.app.Ctx}
	orgStorage := &organization_storage.OrganizationStore{DB: router.app.Store.DB, Ctx: router.app.Ctx}
	permService := permissions_service.NewPermissionService(router.app.Store, router.app.Ctx, l, permStorage)
	roleService := role_service.NewRoleService(router.app.Store, router.app.Ctx, l, roleStorage)
	orgService := organization_service.NewOrganizationService(router.app.Store, router.app.Ctx, l, orgStorage)
	authService := authService.NewAuthService(userStorage, l, permService, roleService, orgService, router.app.Ctx)
	authController := auth.NewAuthController(router.app.Ctx, l, notificationManager, *authService)
	authGroup := fuego.Group(server, apiV1.Path+"/auth")
	router.AuthRoutes(authController, authGroup)

	fuego.Use(server, func(next http.Handler) http.Handler {
		return middleware.AuthMiddleware(next, router.app)
	})

	fuego.Use(server, func(next http.Handler) http.Handler {
		return middleware.AuditMiddleware(next, router.app, l)
	})

	apiV1Group := fuego.Group(server, apiV1.Path)

	authProtectedGroup := fuego.Group(server, apiV1.Path+"/auth")
	router.AuthenticatedAuthRoutes(authProtectedGroup, authController)

	userController := user.NewUserController(router.app.Store, router.app.Ctx, l)
	userGroup := fuego.Group(apiV1Group, "/user")
	fuego.Use(userGroup, func(next http.Handler) http.Handler {
		return middleware.RBACMiddleware(next, router.app, "user")
	})
	router.UserRoutes(userGroup, userController)

	domainController := domain.NewDomainsController(router.app.Store, router.app.Ctx, l, notificationManager)
	domainGroup := fuego.Group(apiV1Group, "/domain")
	domainsAllGroup := fuego.Group(apiV1Group, "/domains")
	fuego.Use(domainGroup, func(next http.Handler) http.Handler {
		return middleware.RBACMiddleware(next, router.app, "domain")
	})
	fuego.Use(domainsAllGroup, func(next http.Handler) http.Handler {
		return middleware.RBACMiddleware(next, router.app, "domain")
	})
	router.DomainRoutes(domainGroup, domainsAllGroup, domainController)

	githubConnectorController := githubConnector.NewGithubConnectorController(router.app.Store, router.app.Ctx, l, notificationManager)
	githubConnectorGroup := fuego.Group(apiV1Group, "/github-connector")
	fuego.Use(githubConnectorGroup, func(next http.Handler) http.Handler {
		return middleware.RBACMiddleware(next, router.app, "github-connector")
	})
	router.GithubConnectorRoutes(githubConnectorGroup, githubConnectorController)

	notifController := notificationController.NewNotificationController(router.app.Store, router.app.Ctx, l, notificationManager)
	notificationGroup := fuego.Group(apiV1Group, "/notification")
	fuego.Use(notificationGroup, func(next http.Handler) http.Handler {
		return middleware.RBACMiddleware(next, router.app, "notification")
	})
	router.NotificationRoutes(notificationGroup, notifController)

	organizationController := organization.NewOrganizationsController(router.app.Store, router.app.Ctx, l, notificationManager)
	organizationGroup := fuego.Group(apiV1Group, "/organizations")
	fuego.Use(organizationGroup, func(next http.Handler) http.Handler {
		return middleware.RBACMiddleware(next, router.app, "organization")
	})
	router.OrganizationRoutes(organizationGroup, organizationController)

	fileManagerController := file_manager.NewFileManagerController(router.app.Ctx, l, notificationManager)
	fileManagerGroup := fuego.Group(apiV1Group, "/file-manager")
	fuego.Use(fileManagerGroup, func(next http.Handler) http.Handler {
		return middleware.RBACMiddleware(next, router.app, "file-manager")
	})
	router.FileManagerRoutes(fileManagerGroup, fileManagerController)

	deployGroup := fuego.Group(apiV1Group, "/deploy")
	fuego.Use(deployGroup, func(next http.Handler) http.Handler {
		return middleware.RBACMiddleware(next, router.app, "deploy")
	})
	router.DeployRoutes(deployGroup, deployController)

	auditController := audit.NewAuditController(router.app.Store.DB, router.app.Ctx, l)
	auditGroup := fuego.Group(apiV1Group, "/audit")
	fuego.Use(auditGroup, func(next http.Handler) http.Handler {
		return middleware.RBACMiddleware(next, router.app, "audit")
	})
	router.AuditRoutes(auditGroup, auditController)

	server.Run()
}

func (s *Router) BasicRoutes(fs *fuego.Server) {
	fuego.Get(fs, "", health.HealthCheck)
	versionGroup := fuego.Group(fs, "/versions")
	fuego.Get(versionGroup, "", func(c fuego.ContextNoBody) (interface{}, error) {
		docs := api.NewVersionDocumentation()
		if err := docs.Load("api/versions.json"); err != nil {
			return nil, err
		}
		return docs, nil
	})
}

func (router *Router) WebSocketServer(f *fuego.Server, deployController *deploy.DeployController) {
	wsServer, err := realtime.NewSocketServer(deployController, router.app.Store.DB, router.app.Ctx)
	if err != nil {
		log.Fatal(err)
	}
	wsHandler := func(c fuego.ContextNoBody) (interface{}, error) {
		log.Printf("WebSocket connection attempt from: %s", c.Request().RemoteAddr)

		wsServer.HandleHTTP(c.Response(), c.Request())
		return nil, nil
	}

	fuego.Get(f, "/ws", wsHandler)
}

func (router *Router) AuthRoutes(authController *auth.AuthController, s *fuego.Server) {
	fuego.Post(s, "/register", authController.Register)
	fuego.Post(s, "/login", authController.Login)
	fuego.Post(s, "/refresh-token", authController.RefreshToken)
}

func (router *Router) AuthenticatedAuthRoutes(s *fuego.Server, authController *auth.AuthController) {
	fuego.Post(s, "/request-password-reset", authController.GeneratePasswordResetLink)
	fuego.Post(s, "/reset-password", authController.ResetPassword)
	fuego.Post(s, "/logout", authController.Logout)
	fuego.Post(s, "/send-verification-email", authController.SendVerificationEmail)
	fuego.Get(s, "/verify-email", authController.VerifyEmail)
	fuego.Post(s, "/create-user", authController.CreateUser)
}

func (router *Router) UserRoutes(s *fuego.Server, userController *user.UserController) {
	fuego.Get(s, "", userController.GetUserDetails)
	fuego.Patch(s, "/name", userController.UpdateUserName)
	fuego.Get(s, "/organizations", userController.GetUserOrganizations)
}

func (router *Router) NotificationRoutes(s *fuego.Server, notificationController *notificationController.NotificationController) {
	smtpGroup := fuego.Group(s, "/smtp")
	fuego.Post(smtpGroup, "", notificationController.AddSmtp)
	fuego.Get(smtpGroup, "", notificationController.GetSmtp)
	fuego.Put(smtpGroup, "", notificationController.UpdateSmtp)
	fuego.Delete(smtpGroup, "", notificationController.DeleteSmtp)

	preferenceGroup := fuego.Group(s, "/preferences")
	fuego.Post(preferenceGroup, "", notificationController.UpdatePreference)
	fuego.Get(preferenceGroup, "", notificationController.GetPreferences)

	webhookGroup := fuego.Group(s, "/webhook")
	fuego.Post(webhookGroup, "", notificationController.CreateWebhookConfig)
	fuego.Get(webhookGroup, "/{type}", notificationController.GetWebhookConfig)
	fuego.Put(webhookGroup, "", notificationController.UpdateWebhookConfig)
	fuego.Delete(webhookGroup, "", notificationController.DeleteWebhookConfig)
}

func (router *Router) DomainRoutes(s *fuego.Server, domainsGroup *fuego.Server, domainController *domain.DomainsController) {
	fuego.Post(s, "", domainController.CreateDomain)
	fuego.Put(s, "", domainController.UpdateDomain)
	fuego.Delete(s, "", domainController.DeleteDomain)
	fuego.Get(s, "/generate", domainController.GenerateRandomSubDomain)
	fuego.Get(domainsGroup, "", domainController.GetDomains)
}

func (router *Router) GithubConnectorRoutes(s *fuego.Server, githubConnectorController *githubConnector.GithubConnectorController) {
	fuego.Post(s, "", githubConnectorController.CreateGithubConnector)
	fuego.Put(s, "", githubConnectorController.UpdateGithubConnectorRequest)
	fuego.Get(s, "/all", githubConnectorController.GetGithubConnectors)
	fuego.Get(s, "/repositories", githubConnectorController.GetGithubRepositories)
}

func (router *Router) DeployRoutes(f *fuego.Server, deployController *deploy.DeployController) {
	fuego.Get(f, "/applications", deployController.GetApplications)
	deploy_application_group := fuego.Group(f, "/application")
	router.DeployApplicationRoutes(deploy_application_group, deployController)
}

func (router *Router) DeployApplicationRoutes(f *fuego.Server, deployController *deploy.DeployController) {
	fuego.Post(f, "", deployController.HandleDeploy)
	fuego.Get(f, "", deployController.GetApplicationById)
	fuego.Delete(f, "", deployController.DeleteApplication)
	fuego.Put(f, "", deployController.UpdateApplication)
	fuego.Post(f, "/redeploy", deployController.ReDeployApplication)
	fuego.Get(f, "/deployments/{deployment_id}", deployController.GetDeploymentById)
	fuego.Post(f, "/rollback", deployController.HandleRollback)
	fuego.Post(f, "/restart", deployController.HandleRestart)
}

func (router *Router) FileManagerRoutes(f *fuego.Server, fileManagerController *file_manager.FileManagerController) {
	fuego.Get(f, "", fileManagerController.ListFiles)
	fuego.Post(f, "/create-directory", fileManagerController.CreateDirectory)
	fuego.Post(f, "/move-directory", fileManagerController.MoveDirectory)
	fuego.Post(f, "/copy-directory", fileManagerController.CopyDirectory)
	fuego.Post(f, "/upload", fileManagerController.UploadFile)
	fuego.Delete(f, "/delete-directory", fileManagerController.DeleteDirectory)
}

func (router *Router) OrganizationRoutes(f *fuego.Server, organizationController *organization.OrganizationsController) {
	fuego.Get(f, "/users", organizationController.GetOrganizationUsers)
	fuego.Post(f, "/add-user", organizationController.AddUserToOrganization)
	fuego.Post(f, "/remove-user", organizationController.RemoveUserFromOrganization)
	fuego.Post(f, "/update-user-role", organizationController.UpdateUserRole)
	fuego.Get(f, "/roles", organizationController.GetRoles)
	fuego.Get(f, "/resources", organizationController.GetResources)
	fuego.Put(f, "", organizationController.UpdateOrganization)
	fuego.Post(f, "", organizationController.CreateOrganization)
	fuego.Delete(f, "", organizationController.DeleteOrganization)
}

func (router *Router) AuditRoutes(s *fuego.Server, auditController *audit.AuditController) {
	fuego.Get(s, "/logs", auditController.GetRecentAuditLogs)
}
