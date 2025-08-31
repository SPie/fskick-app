package players

import (
	"testing"

	"github.com/spie/fskick/internal/db"
	"github.com/stretchr/testify/assert"
)

type mockPlayerRepository struct {
	player Player
	err    error
}

func (repo mockPlayerRepository) FindPlayerByName(name string) (Player, error) {
	return repo.player, repo.err
}

func (repo mockPlayerRepository) CreatePlayer(player *Player) error {
	// Mock implementation: do nothing and return nil error
	return nil
}

func (repo mockPlayerRepository) FindPlayerByUUID(uuid string) (Player, error) {
	// Mock implementation: return zero value and nil error
	return Player{}, nil
}

func (repo mockPlayerRepository) FindPlayersByNames(names []string) ([]Player, error) {
	// Mock implementation: return empty slice and nil error
	return []Player{}, nil
}

func TestPlayersManager_GetPlayerByName(t *testing.T) {
	tests := map[string]struct {
		playerName string
		setupMocks func() Manager
		assertions []func(t *testing.T, player Player, err error)
	}{
		"with player found": {
			playerName: "test_player",
			setupMocks: func() Manager {
				playerRepository := mockPlayerRepository{
					player: Player{
						Model: db.Model{
							ID:   23,
							UUID: "someuuid123",
						},
						Name: "test_player",
					},
				}

				return NewManager(playerRepository)
			},
			assertions: []func(t *testing.T, player Player, err error){
				func(t *testing.T, player Player, err error) {
					assert.Equal(t, "test_player", player.Name)
					assert.Equal(t, uint(23), player.ID)
					assert.Equal(t, "someuuid123", player.UUID)
					assert.NoError(t, err)
				},
			},
		},
		"with player not found": {
			playerName: "test_player",
			setupMocks: func() Manager {
				playerRepository := mockPlayerRepository{
					err: ErrPlayerNotFound,
				}

				return NewManager(playerRepository)
			},
			assertions: []func(t *testing.T, player Player, err error){
				func(t *testing.T, player Player, err error) {
					assert.Zero(t, player)
					assert.ErrorIs(t, err, ErrPlayerNotFound)
					assert.ErrorContains(t, err, "get player by name: ")
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			manager := tt.setupMocks()

			player, err := manager.GetPlayerByName(tt.playerName)

			for _, assertion := range tt.assertions {
				assertion(t, player, err)
			}
		})
	}
}
