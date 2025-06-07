package auth

import (
	"net/http"
	"testing"

	. "github.com/Eun/go-hit"
	"github.com/raghavyuva/nixopus-api/internal/tests"
	"github.com/raghavyuva/nixopus-api/internal/testutils"
)

func TestIsAdminRegistered(t *testing.T) {
	testCases := []struct {
		name           string
		setup          func() *testutils.TestSetup
		expectedStatus int
	}{
		{
			name: "Successfully check if admin is registered",
			setup: func() *testutils.TestSetup {
				setup := testutils.NewTestSetup()
				setup.CreateTestUserAndOrg()
				return setup
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Admin is not registered",
			setup: func() *testutils.TestSetup {
				return testutils.NewTestSetup()
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			Test(t,
				Description(tc.name),
				Get(tests.GetIsAdminRegisteredURL()),
				Expect().Status().Equal(int64(tc.expectedStatus)),
			)
		})
	}
}
