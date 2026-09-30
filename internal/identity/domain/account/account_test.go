package account

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAccount(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		uuid string
		email string
		passwordHash string
		role Role
		wantErr bool
	}{
		{
			name: "missing email",
			uuid: "acc-1",
			email: "",
			passwordHash: "password",
			role: RoleOfficer,
			wantErr: true,
		},
		{
			name: "missing password",
			uuid: "acc-1",
			email: "udin@gmail.com",
			passwordHash: "",
			role: RoleOfficer,
			wantErr: true,
		},
		{
			name: "happy case",
			uuid: "acc-1",
			email: "udin@gmail.com",
			passwordHash: "password123",
			role: RoleOfficer,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewAccount(tc.uuid, tc.email, tc.passwordHash, tc.role)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}
