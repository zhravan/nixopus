package tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	authService "github.com/raghavyuva/nixopus-api/internal/features/auth/service"
	user_storage "github.com/raghavyuva/nixopus-api/internal/features/auth/storage"
	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	organization_service "github.com/raghavyuva/nixopus-api/internal/features/organization/service"
	organization_storage "github.com/raghavyuva/nixopus-api/internal/features/organization/storage"
	permissions_service "github.com/raghavyuva/nixopus-api/internal/features/permission/service"
	permissions_storage "github.com/raghavyuva/nixopus-api/internal/features/permission/storage"
	role_service "github.com/raghavyuva/nixopus-api/internal/features/role/service"
	role_storage "github.com/raghavyuva/nixopus-api/internal/features/role/storage"
	dbstorage "github.com/raghavyuva/nixopus-api/internal/storage"
	"github.com/raghavyuva/nixopus-api/internal/types"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var (
	testDB *bun.DB
	ctx    context.Context
)

func init() {
	ctx = context.Background()

	envFiles := []string{
		"../../../../env.test",
		"../../../../../env.test",
		"env.test",
	}

	envLoaded := false
	for _, file := range envFiles {
		if err := godotenv.Load(file); err == nil {
			envLoaded = true
			break
		}
	}

	if !envLoaded {
		fmt.Println("Warning: Could not load env.test file from any location")
	}

	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5433")
	dbUser := getEnvOrDefault("DB_USER", "nixopus")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "nixopus")
	dbName := getEnvOrDefault("DB_NAME", "nixopus_test")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	fmt.Printf("Connecting to test database: %s\n", connStr)

	config, err := pgx.ParseConfig(connStr)
	if err != nil {
		fmt.Printf("Failed to parse config: %v\n", err)
		os.Exit(1)
	}

	sqldb := stdlib.OpenDB(*config)
	testDB = bun.NewDB(sqldb, pgdialect.New())

	if err := testDB.Ping(); err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully connected to test database")

	testDB.RegisterModel((*types.RolePermissions)(nil))
	testDB.RegisterModel((*types.OrganizationUsers)(nil))

	store := dbstorage.NewStore(testDB)
	if err := store.Init(ctx); err != nil {
		fmt.Printf("Failed to initialize store: %v\n", err)
		os.Exit(1)
	}

	if err := dbstorage.ResetMigrations(testDB); err != nil {
		fmt.Printf("Failed to reset migrations: %v\n", err)
		os.Exit(1)
	}

	if err := dbstorage.MigrateDownAll(testDB, "../../../../migrations"); err != nil {
		fmt.Printf("Failed to migrate down all: %v\n", err)
		os.Exit(1)
	}

	if err := dbstorage.RunMigrations(testDB, "../../../../migrations"); err != nil {
		fmt.Printf("Failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully ran migrations")
}

func TestMain(m *testing.M) {
	code := m.Run()

	testDB.Close()

	os.Exit(code)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func cleanDatabase() error {
	store := dbstorage.NewStore(testDB)
	if err := store.DropAllTables(ctx); err != nil {
		return fmt.Errorf("failed to drop all tables: %w", err)
	}

	if err := dbstorage.ResetMigrations(testDB); err != nil {
		return fmt.Errorf("failed to reset migrations: %w", err)
	}

	if err := dbstorage.RunMigrations(testDB, "../../../../migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func GetTestStorage() (*user_storage.UserStorage, *authService.AuthService) {
	if testDB == nil {
		panic("testDB is nil - database not initialized")
	}
	if ctx == nil {
		panic("ctx is nil - context not initialized")
	}

	if err := cleanDatabase(); err != nil {
		panic(fmt.Sprintf("failed to clean database: %v", err))
	}

	l := logger.NewLogger()
	userStorage := &user_storage.UserStorage{DB: testDB, Ctx: ctx}
	permStorage := &permissions_storage.PermissionStorage{DB: testDB, Ctx: ctx}
	roleStorage := &role_storage.RoleStorage{DB: testDB, Ctx: ctx}
	orgStorage := &organization_storage.OrganizationStore{DB: testDB, Ctx: ctx}
	permService := permissions_service.NewPermissionService(&dbstorage.Store{DB: testDB}, ctx, l, permStorage)
	roleService := role_service.NewRoleService(&dbstorage.Store{DB: testDB}, ctx, l, roleStorage)
	orgService := organization_service.NewOrganizationService(&dbstorage.Store{DB: testDB}, ctx, l, orgStorage)
	authService := authService.NewAuthService(userStorage, l, permService, roleService, orgService, ctx)
	return userStorage, authService
}
