package passwords

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordService_HashPassword(t *testing.T) {
	tests := map[string]struct {
		plaintextPassword string
		expectsErr        bool
	}{
		"successfully hashes a simple password": {
			plaintextPassword: "supersecretpassword",
			expectsErr:        false,
		},
		"successfully hashes a password with special characters": {
			plaintextPassword: "Pa$$w0rd!@#",
			expectsErr:        false,
		},
		"successfully hashes an empty password": {
			plaintextPassword: "",
			expectsErr:        false,
		},
		"error on hashing 73-byte long password": {
			plaintextPassword: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+={}|[]\\:;\"'<>,.?/~`",
			expectsErr:        true,
		},
	}

	passwordService := PasswordService{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			hashedPassword, err := passwordService.HashPassword([]byte(tt.plaintextPassword))

			if tt.expectsErr {
				assert.Error(t, err)
			} else {
				compareErr := bcrypt.CompareHashAndPassword(hashedPassword, []byte(tt.plaintextPassword))
				assert.Nil(t, compareErr, "bcrypt.CompareHashAndPassword should return nil, indicating a match")
			}
		})
	}
}
