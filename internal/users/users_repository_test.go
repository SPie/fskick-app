package users

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/players"
	"github.com/stretchr/testify/assert"
)

type mockConnection struct {
	executedQueries []struct {
		query string
		args  []any
	}
	expectedResult sql.Result
	expectedErr    error
}

func (conn mockConnection) Query(query string, args ...any) (*sql.Rows, error) {
	// TODO
	return nil, nil
}

func (conn mockConnection) QueryRow(query string, args ...any) *sql.Row {
	// TODO
	return nil
}

func (conn *mockConnection) Exec(query string, args ...any) (sql.Result, error) {
	conn.executedQueries = append(
		conn.executedQueries,
		struct {
			query string
			args  []any
		}{
			query: query,
			args:  args,
		},
	)

	return conn.expectedResult, conn.expectedErr
}

func (conn mockConnection) Begin() (*sql.Tx, error) {
	// TODO
	return nil, nil
}

func (conn mockConnection) Close() error {
	// TODO
	return nil
}

func TestUsersRepository_CreateUser(t *testing.T) {
	tests := map[string]struct {
		user       User
		setUpMocks func() (UsersRepository, *mockConnection)
		assertions []func(t *testing.T, user User, conn *mockConnection, err error)
	}{
		"with user created": {
			user: User{
				Player: players.Player{
					Model: db.Model{
						ID:        23,
						UUID:      "uuid123",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Name: "test_player",
				},
				Email:    "email@example.com",
				Password: "hashedpassword",
			},
			setUpMocks: func() (UsersRepository, *mockConnection) {
				conn := &mockConnection{}
				return UsersRepository{conn: conn}, conn
			},
			assertions: []func(t *testing.T, user User, conn *mockConnection, err error){
				func(t *testing.T, user User, conn *mockConnection, err error) {
					assert.Equal(
						t,
						"UPDATE players SET email = ?, password = ?, updated_at = ? WHERE id = ?",
						conn.executedQueries[0].query,
					)
					assert.Equal(t, "email@example.com", conn.executedQueries[0].args[0])
					assert.Equal(t, "hashedpassword", conn.executedQueries[0].args[1])
					assert.Equal(t, uint(23), conn.executedQueries[0].args[3])
				},
			},
		},
		"with error on update": {
			user: User{
				Player: players.Player{
					Model: db.Model{
						ID:        23,
						UUID:      "uuid123",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Name: "test_player",
				},
				Email:    "email@example.com",
				Password: "hashedpassword",
			},
			setUpMocks: func() (UsersRepository, *mockConnection) {
				conn := &mockConnection{expectedErr: errors.New("some db error")}
				return UsersRepository{conn: conn}, conn
			},
			assertions: []func(t *testing.T, user User, conn *mockConnection, err error){
				func(t *testing.T, user User, conn *mockConnection, err error) {
					assert.ErrorContains(t, err, "Error update player to create user")
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			usersRepository, conn := tt.setUpMocks()
			err := usersRepository.CreateUser(&tt.user)

			for _, assertion := range tt.assertions {
				assertion(t, tt.user, conn, err)
			}
		})
	}
}
