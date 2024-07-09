package games

import (
	"time"

	"github.com/spie/fskick/internal/seasons"
)

type Manager struct {
	gameRepository GamesRepository
	seasonsManager seasons.Manager
}

func NewManager(gameRepository GamesRepository, seasonsManager seasons.Manager) Manager {
	return Manager{
		gameRepository: gameRepository,
		seasonsManager: seasonsManager,
	}
}

func (manager Manager) CreateGame(playedAt time.Time, ) (*Game, error) {
	activeSeason, err := manager.seasonsManager.ActiveSeason()
	if err != nil {
		return &Game{}, err
	}

	if playedAt.IsZero() {
		playedAt = time.Now()
	}

	game := &Game{Season: &activeSeason, PlayedAt: playedAt}

	err = manager.gameRepository.CreateGame(game)
	if err != nil {
		return &Game{}, err
	}

	return game, nil
}

func (manager Manager) GetGamesCount() (int, error) {
	return manager.gameRepository.Count()
}

func (manager Manager) GetGamesCountForSeason(season seasons.Season) (int, error) {
	return manager.gameRepository.CountForSeason(season)
}
