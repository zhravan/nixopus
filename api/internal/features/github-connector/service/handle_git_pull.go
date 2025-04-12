package service

import (
	"fmt"

	"github.com/raghavyuva/nixopus-api/internal/features/logger"
)

func (s *GithubConnectorService) handleGitPull(authenticatedURL, clonePath string, userID string) error {
	hasChanges, err := s.gitClient.HasUncommittedChanges(clonePath)
	if err != nil {
		s.logger.Log(logger.Error, fmt.Sprintf("Failed to check for uncommitted changes: %s", err.Error()), userID)
		return err
	}

	if hasChanges {
		s.logger.Log(logger.Info, "Stashing uncommitted changes", userID)
		stashID, err := s.gitClient.Stash(clonePath)
		if err != nil {
			s.logger.Log(logger.Error, fmt.Sprintf("Failed to stash changes: %s", err.Error()), userID)
			return err
		}

		defer func() {
			if stashID != "" {
				s.logger.Log(logger.Info, "Applying stashed changes", userID)
				if err := s.gitClient.ApplyStash(clonePath, stashID); err != nil {
					s.logger.Log(logger.Error, fmt.Sprintf("Failed to apply stash: %s", err.Error()), userID)
				}
			}
		}()
	}

	s.logger.Log(logger.Info, "Pulling latest changes", userID)
	if err := s.gitClient.Pull(authenticatedURL, clonePath); err != nil {
		s.logger.Log(logger.Error, fmt.Sprintf("Failed to pull repository: %s", err.Error()), userID)
		return err
	}

	return nil
}