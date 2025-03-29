package auth

import (
	"context"
	"net/http"

	"github.com/raghavyuva/nixopus-api/internal/features/auth/service"
	"github.com/raghavyuva/nixopus-api/internal/features/auth/storage"
	"github.com/raghavyuva/nixopus-api/internal/features/auth/validation"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	"github.com/raghavyuva/nixopus-api/internal/features/notification"
	organization_service "github.com/raghavyuva/nixopus-api/internal/features/organization/service"
	organization_storage "github.com/raghavyuva/nixopus-api/internal/features/organization/storage"
	permissions_service "github.com/raghavyuva/nixopus-api/internal/features/permission/service"
	permissions_storage "github.com/raghavyuva/nixopus-api/internal/features/permission/storage"
	role_service "github.com/raghavyuva/nixopus-api/internal/features/role/service"
	role_storage "github.com/raghavyuva/nixopus-api/internal/features/role/storage"
	shared_storage "github.com/raghavyuva/nixopus-api/internal/storage"
	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
	"github.com/raghavyuva/nixopus-api/internal/utils"
)

type AuthController struct {
	validator    *validation.Validator
	service      service.AuthServiceInterface
	ctx          context.Context
	logger       logger.Logger
	notification *notification.NotificationManager
}

func NewAuthController(
	store *shared_storage.Store,
	ctx context.Context,
	logger logger.Logger,
	notificationManager *notification.NotificationManager,
) *AuthController {
	userStorage := &storage.UserStorage{DB: store.DB, Ctx: ctx}
	permStorage := &permissions_storage.PermissionStorage{DB: store.DB, Ctx: ctx}
	roleStorage := &role_storage.RoleStorage{DB: store.DB, Ctx: ctx}
	orgStorage := &organization_storage.OrganizationStore{DB: store.DB, Ctx: ctx}

	permService := permissions_service.NewPermissionService(store, ctx, logger, permStorage)
	roleService := role_service.NewRoleService(store, ctx, logger, roleStorage)
	orgService := organization_service.NewOrganizationService(store, ctx, logger, orgStorage)

	return &AuthController{
		validator:    validation.NewValidator(),
		service:      service.NewAuthService(userStorage, logger, permService, roleService, orgService, ctx),
		ctx:          ctx,
		logger:       logger,
		notification: notificationManager,
	}
}

func (c *AuthController) parseAndValidate(w http.ResponseWriter, r *http.Request, req any) bool {
	if err := c.validator.ValidateRequest(req); err != nil {
		c.logger.Log(logger.Error, err.Error(), err.Error())
		return false
	}

	if r.URL.Path == "/api/v1/auth/login" {
		return true
	}

	user := utils.GetUser(w, r)
	if err := c.validator.AccessValidator(w, r, user); err != nil {
		c.logger.Log(logger.Error, err.Error(), err.Error())
		return false
	}

	return true
}

func (c *AuthController) Notify(payloadType notification.NotificationPayloadType, user *shared_types.User, r *http.Request) {
	notificationData := notification.NotificationAuthenticationData{
		Email: user.Email,
		NotificationBaseData: notification.NotificationBaseData{
			IP:      r.RemoteAddr,
			Browser: r.UserAgent(),
		},
		UserName: user.Username,
	}

	payload := notification.NewNotificationPayload(
		payloadType,
		user.ID.String(),
		notificationData,
		notification.NotificationCategoryAuthentication,
	)

	c.notification.SendNotification(payload)
}
