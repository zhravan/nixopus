package tests

import (
	"testing"

	"github.com/raghavyuva/nixopus-api/internal/features/auth/types"
	"github.com/stretchr/testify/assert"
)

func TestLogout(t *testing.T) {
	_, authService := GetTestStorage()

	registerRequest := types.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Username: "testuser",
		Type:     "viewer",
	}

	registerResponse, err := authService.Register(registerRequest)
	assert.NoError(t, err)

	tests := []struct {
		name          string
		refreshToken  string
		expectError   bool
		errorContains string
	}{
		{
			name:         "successful logout",
			refreshToken: registerResponse.RefreshToken,
			expectError:  false,
		},
		{
			name:          "invalid refresh token",
			refreshToken:  "invalid_token",
			expectError:   true,
			errorContains: "invalid refresh token",
		},
		{
			name:          "empty refresh token",
			refreshToken:  "",
			expectError:   true,
			errorContains: "refresh token is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.Logout(tt.refreshToken)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestLogoutWithAlreadyRevokedToken(t *testing.T) {
	_, authService := GetTestStorage()

	registerRequest := types.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Username: "testuser",
		Type:     "viewer",
	}

	registerResponse, err := authService.Register(registerRequest)
	assert.NoError(t, err)

	err = authService.Logout(registerResponse.RefreshToken)
	assert.NoError(t, err)

	err = authService.Logout(registerResponse.RefreshToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), types.ErrInvalidRefreshToken.Error())
}
