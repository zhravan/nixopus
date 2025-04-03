package tests

import (
	"testing"

	"github.com/raghavyuva/nixopus-api/internal/features/auth/types"
	"github.com/stretchr/testify/assert"
)

func TestVerifyToken(t *testing.T) {
	_, authService := GetTestStorage()

	registerRequest := types.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Username: "testuser",
		Type:     "viewer",
	}

	registerResponse, err := authService.Register(registerRequest)
	assert.NoError(t, err)

	verificationToken, err := authService.GenerateVerificationToken(registerResponse.User.ID.String())
	assert.NoError(t, err)
	assert.NotEmpty(t, verificationToken)

	tests := []struct {
		name          string
		token         string
		expectError   bool
		errorContains string
	}{
		{
			name:        "successful verification",
			token:       verificationToken,
			expectError: false,
		},
		{
			name:          "invalid token",
			token:         "invalid_token",
			expectError:   true,
			errorContains: "verification token is already used",
		},
		{
			name:          "empty token",
			token:         "",
			expectError:   true,
			errorContains: "verification token is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := authService.VerifyToken(tt.token)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, registerResponse.User.ID.String(), userID)

			user, err := authService.GetUserByID(userID)
			assert.NoError(t, err)
			assert.True(t, user.IsVerified)
		})
	}
}

func TestVerifyTokenWithExpiredToken(t *testing.T) {
	_, authService := GetTestStorage()

	registerRequest := types.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Username: "testuser",
		Type:     "viewer",
	}

	registerResponse, err := authService.Register(registerRequest)
	assert.NoError(t, err)

	verificationToken, err := authService.GenerateVerificationToken(registerResponse.User.ID.String())
	assert.NoError(t, err)
	assert.NotEmpty(t, verificationToken)

	userID, err := authService.VerifyToken(verificationToken)
	assert.NoError(t, err)
	assert.Equal(t, registerResponse.User.ID.String(), userID)

	_, err = authService.VerifyToken(verificationToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "verification token is already used")
}
