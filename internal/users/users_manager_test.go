package users

import (
	"errors"
	"testing"
	"time"

	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/players"
	"github.com/stretchr/testify/assert"
)

type mockUsersRepository struct {
	updatedAt time.Time
	err       error
}

func (mockUserRepository mockUsersRepository) CreateUser(user *User) error {
	if mockUserRepository.err != nil {
		return mockUserRepository.err
	}

	user.UpdatedAt = mockUserRepository.updatedAt

	return nil
}

type mockPlayersManager struct {
	player players.Player
	err    error
}

func (mockPlayersManager mockPlayersManager) GetPlayerByName(name string) (players.Player, error) {
	return mockPlayersManager.player, mockPlayersManager.err
}

type mockPasswordService struct {
	hashedPassword []byte
	err            error
}

func (mockPasswordService mockPasswordService) HashPassword(password []byte) ([]byte, error) {
	return mockPasswordService.hashedPassword, mockPasswordService.err
}

func TestCreateUserFromPlayer(t *testing.T) {
	type args struct {
		playerName        string
		email             string
		plaintextPassword string
	}
	tests := map[string]struct {
		args struct {
			playerName        string
			email             string
			plaintextPassword string
		}
		setupMocks func() Manager
		assertions []func(t *testing.T, user User, err error)
	}{
		"successfully created user": {
			args: args{
				playerName:        "test_player",
				email:             "test@example.com",
				plaintextPassword: "password123",
			},
			setupMocks: func() Manager {
				return Manager{
					usersRepository: mockUsersRepository{
						updatedAt: time.Date(2021, time.September, 8, 0, 0, 0, 0, time.UTC),
					},
					playersManager: mockPlayersManager{player: players.Player{
						Name: "test_player",
						Model: db.Model{
							ID:        23,
							UUID:      "uuid1234",
							CreatedAt: time.Date(2009, time.July, 28, 0, 0, 0, 0, time.UTC),
							UpdatedAt: time.Date(2009, time.July, 28, 0, 0, 0, 0, time.UTC),
						},
					}},
					passwordService: mockPasswordService{hashedPassword: []byte("hashedpassword")},
				}
			},
			assertions: []func(t *testing.T, user User, err error){
				func(t *testing.T, user User, err error) {
					assert.Equal(t, uint(23), user.ID)
					assert.Equal(t, "test@example.com", user.Email)
					assert.Equal(t, "test_player", user.Name)
					assert.Equal(t, "hashedpassword", user.Password)
					assert.NotNil(t, user.CreatedAt)
					assert.Equal(t, time.Date(2021, time.September, 8, 0, 0, 0, 0, time.UTC), user.UpdatedAt)
					assert.NotNil(t, user.UUID)
					assert.NoError(t, err)
				},
			},
		},
		"player not found": {
			args: args{
				playerName:        "test_player",
				email:             "test@example.com",
				plaintextPassword: "password123",
			},
			setupMocks: func() Manager {
				return Manager{playersManager: mockPlayersManager{err: players.ErrPlayerNotFound}}
			},
			assertions: []func(t *testing.T, user User, err error){
				func(t *testing.T, user User, err error) {
					assert.Zero(t, user)
					assert.ErrorContains(t, err, "Player test_player not found")
				},
			},
		},
		"error on player retrieve": {
			args: args{
				playerName:        "test_player",
				email:             "test@example.com",
				plaintextPassword: "password123",
			},
			setupMocks: func() Manager {
				return Manager{playersManager: mockPlayersManager{err: errors.New("some error")}}
			},
			assertions: []func(t *testing.T, user User, err error){
				func(t *testing.T, user User, err error) {
					assert.Zero(t, user)
					assert.ErrorContains(t, err, "Get player for CreateUserFromPlayer")
					assert.ErrorContains(t, err, "some error")
				},
			},
		},
		"with error on storing user": {
			args: args{
				playerName:        "test_player",
				email:             "test@example.com",
				plaintextPassword: "password123",
			},
			setupMocks: func() Manager {
				return Manager{
					usersRepository: mockUsersRepository{
						err: errors.New("some error"),
					},
					playersManager:  mockPlayersManager{player: players.Player{}},
					passwordService: mockPasswordService{hashedPassword: []byte("hashedpassword")},
				}
			},
			assertions: []func(t *testing.T, user User, err error){
				func(t *testing.T, user User, err error) {
					assert.Zero(t, user)
					assert.Error(t, err)
				},
			},
		},
		"with error on password hashing": {
			args: args{
				playerName:        "test_player",
				email:             "test@example.com",
				plaintextPassword: "password123",
			},
			setupMocks: func() Manager {
				return Manager{
					playersManager:  mockPlayersManager{player: players.Player{}},
					passwordService: mockPasswordService{err: errors.New("some password error")},
				}
			},
			assertions: []func(t *testing.T, user User, err error){
				func(t *testing.T, user User, err error) {
					assert.Zero(t, user)
					assert.ErrorContains(t, err, "some password error")
					assert.ErrorContains(t, err, "Hash password for create user from player: ")
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			manager := tt.setupMocks()

			user, err := manager.CreateUserFromPlayer(
				tt.args.playerName,
				tt.args.email,
				tt.args.plaintextPassword,
			)

			for _, assertion := range tt.assertions {
				assertion(t, user, err)
			}
		})
	}
}
