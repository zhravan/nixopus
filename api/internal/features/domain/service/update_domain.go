package service

import (
	"github.com/google/uuid"
	"github.com/raghavyuva/nixopus-api/internal/features/domain/types"
)

func (s *DomainsService) UpdateDomain(newName string,userID string,domainID string) (types.CreateDomainResponse, error) {
	existing_domain, err := s.storage.GetDomain(domainID)
	if existing_domain.ID == uuid.Nil {
		return types.CreateDomainResponse{}, types.ErrDomainAlreadyExists
	}

	if err != nil {
		return types.CreateDomainResponse{}, err
	}

	if err := s.storage.UpdateDomain(newName,domainID); err != nil {
		return types.CreateDomainResponse{}, err
	}

	return types.CreateDomainResponse{}, nil
}
