package loader

import (
	"context"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"golang.org/x/crypto/bcrypt"

	"github.com/raghavyuva/nixopus-api/internal/types"
)

type FixtureLoader struct {
	db *bun.DB
}

func NewFixtureLoader(db *bun.DB) *FixtureLoader {
	return &FixtureLoader{db: db}
}

// LoadFixtures loads the fixtures from the given fixture path and options
// helps us to create the fixtures with the correct functions we have used in the actual implementation
// for example:
// - now: returns the current time in RFC3339Nano format
// - uuid: returns a new UUID
// - hashPassword: hashes the given password using bcrypt
// we can use these functions in the fixtures to create the fixtures with the correct values
func (fl *FixtureLoader) LoadFixtures(ctx context.Context, fixturePath string, options ...dbfixture.FixtureOption) error {
	fl.registerModels()

	funcMap := template.FuncMap{
		"now": func() string {
			return time.Now().Format(time.RFC3339Nano)
		},
		"uuid": func() string {
			return uuid.New().String()
		},
		"hashPassword": func(password string) string {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				panic(fmt.Sprintf("failed to hash password: %v", err))
			}
			return string(hashedPassword)
		},
	}

	fixture := dbfixture.New(fl.db, append(options, dbfixture.WithTemplateFuncs(funcMap))...)

	err := fixture.Load(ctx, os.DirFS("."), fixturePath)
	if err != nil {
		return fmt.Errorf("failed to load fixture %s: %w", fixturePath, err)
	}

	return nil
}

// LoadFixturesWithRecreate loads the fixtures with the recreate tables option
// it will drop the tables and create them again
func (fl *FixtureLoader) LoadFixturesWithRecreate(ctx context.Context, fixturePath string) error {
	return fl.LoadFixtures(ctx, fixturePath, dbfixture.WithRecreateTables())
}

// LoadFixturesWithTruncate loads the fixtures with the truncate tables option
// it will truncate the tables and load the fixtures
func (fl *FixtureLoader) LoadFixturesWithTruncate(ctx context.Context, fixturePath string) error {
	return fl.LoadFixtures(ctx, fixturePath, dbfixture.WithTruncateTables())
}

// registerModels registers the models with the database
// it will register the models with the database so that the fixtures can be loaded
func (fl *FixtureLoader) registerModels() {
	fl.db.RegisterModel((*types.OrganizationUsers)(nil))
	fl.db.RegisterModel((*types.RolePermissions)(nil))
	fl.db.RegisterModel((*types.User)(nil))
	fl.db.RegisterModel((*types.Organization)(nil))
	fl.db.RegisterModel((*types.Role)(nil))
	fl.db.RegisterModel((*types.Permission)(nil))
	fl.db.RegisterModel((*types.Domain)(nil))
	fl.db.RegisterModel((*types.Application)(nil))
	fl.db.RegisterModel((*types.ApplicationStatus)(nil))
	fl.db.RegisterModel((*types.ApplicationLogs)(nil))
	fl.db.RegisterModel((*types.ApplicationDeployment)(nil))
	fl.db.RegisterModel((*types.ApplicationDeploymentStatus)(nil))
	fl.db.RegisterModel((*types.FeatureFlag)(nil))
}

// GetFixtureData gets the fixture data for the given model and row ID
func (fl *FixtureLoader) GetFixtureData(fixture *dbfixture.Fixture, modelName, rowID string) interface{} {
	return fixture.MustRow(fmt.Sprintf("%s.%s", modelName, rowID))
}
