package tests

import (
	"testing"

	"github.com/google/uuid"
	domainService "github.com/raghavyuva/nixopus-api/internal/features/domain/service"
	domainStorage "github.com/raghavyuva/nixopus-api/internal/features/domain/storage"
	domainTypes "github.com/raghavyuva/nixopus-api/internal/features/domain/types"
	"github.com/raghavyuva/nixopus-api/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestUpdateDomain(t *testing.T) {
	t.Run("should update domain successfully", func(t *testing.T) {
		setup := testutils.NewTestSetup()
		storage := &domainStorage.DomainStorage{DB: setup.DB, Ctx: setup.Ctx}
		service := domainService.NewDomainsService(setup.Store, setup.Ctx, setup.Logger, storage)

		user, org, err := setup.CreateTestUserAndOrg()
		assert.NoError(t, err)

		req := domainTypes.CreateDomainRequest{
			Name:           "test.domain.com",
			OrganizationID: org.ID,
		}

		resp, err := service.CreateDomain(req, user.ID.String())
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.ID)

		newName := "updated.domain.com"
		updated, err := service.UpdateDomain(newName, user.ID.String(), resp.ID)
		assert.NoError(t, err)
		assert.Equal(t, newName, updated.Name)

		domains, err := service.GetDomains(org.ID.String(), user.ID)
		assert.NoError(t, err)
		assert.Len(t, domains, 1)
		assert.Equal(t, newName, domains[0].Name)
	})

	t.Run("should not update domain with invalid domain ID", func(t *testing.T) {
		setup := testutils.NewTestSetup()
		storage := &domainStorage.DomainStorage{DB: setup.DB, Ctx: setup.Ctx}
		service := domainService.NewDomainsService(setup.Store, setup.Ctx, setup.Logger, storage)

		user, _, err := setup.CreateTestUserAndOrg()
		assert.NoError(t, err)

		_, err = service.UpdateDomain("new.domain.com", user.ID.String(), "invalid-uuid")
		assert.Error(t, err)
		assert.ErrorIs(t, err, domainTypes.ErrInvalidDomainID)
	})

	t.Run("should not update domain with non-existent domain ID", func(t *testing.T) {
		setup := testutils.NewTestSetup()
		storage := &domainStorage.DomainStorage{DB: setup.DB, Ctx: setup.Ctx}
		service := domainService.NewDomainsService(setup.Store, setup.Ctx, setup.Logger, storage)

		user, _, err := setup.CreateTestUserAndOrg()
		assert.NoError(t, err)

		nonExistentID := uuid.New().String()
		_, err = service.UpdateDomain("new.domain.com", user.ID.String(), nonExistentID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domainTypes.ErrDomainNotFound)
	})
}
