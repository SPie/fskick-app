package passwords

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockBcrypt struct {
	generateFromPassword func(password []byte, cost int) ([]byte, error)
}

func (m *mockBcrypt) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	if m.generateFromPassword != nil {
		return m.generateFromPassword(password, cost)
	}
	return bcrypt.GenerateFromPassword(password, cost)
}

func TestPasswordService_HashPassword(t *testing.T) {
	tests := map[string]struct {
		plaintextPassword string
		setupMocks        func() PasswordService
		assertions        []func(t *testing.T, hashedPassword []byte, err error)
	}{
		"with correct password hash": {
			plaintextPassword: "test",
			setupMocks: func() PasswordService {
				return PasswordService{
					bcrypt: &mockBcrypt{},
				}
			},
			assertions: []func(t *testing.T, hashedPassword []byte, err error) {
				func(t *testing.T, hashedPassword []byte, err error) {
					assert.Nil(t, err)
					assert.NotEmpty(t, hashedPassword)
				},
			},
		},
		"with error hashing password": {
			plaintextPassword: "test",
			setupMocks: func() PasswordService {
				return PasswordService{
					bcrypt: &mockBcrypt{
						generateFromPassword: func(password []byte, cost int) ([]byte, error) {
							return nil, fmt.Errorf("error hashing password")
						},
					},
				}
			},
			assertions: []func(t *testing.T, hashedPassword []byte, err error) {
				func(t *testing.T, hashedPassword []byte, err error) {
					assert.NotNil(t, err)
					assert.Empty(t, hashedPassword)
					assert.EqualError(t, err, "error hashing password")
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			passwordService := tt.setupMocks()

			hashedPassword, err := passwordService.HashPassword([]byte(tt.plaintextPassword))

			for _, assertion := range tt.assertions {
				assertion(t, hashedPassword, err)
			}
		})
	}
}
