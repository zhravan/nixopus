package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/raghavyuva/nixopus-api/internal/features/domain/types"

	shared_types "github.com/raghavyuva/nixopus-api/internal/types"
)

func (s *DomainsService) CreateDomain(req types.CreateDomainRequest, userID string) (types.CreateDomainResponse, error) {
	existing_domain, err := s.storage.GetDomainByName(req.Name)
	if existing_domain.ID != uuid.Nil {
		return types.CreateDomainResponse{}, types.ErrDomainAlreadyExists
	}

	if err != nil {
		return types.CreateDomainResponse{}, err
	}

	domain := &shared_types.Domain{
		ID:        uuid.New(),
		UserID:    uuid.MustParse(userID),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
		Name:      req.Name,
	}

	if err := s.storage.CreateDomain(domain); err != nil {
		return types.CreateDomainResponse{}, err
	}

	return types.CreateDomainResponse{}, nil
}
