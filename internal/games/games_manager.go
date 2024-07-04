package games

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Manager struct {
	gameRepository    GamesRepository
	seasonsRepository SeasonsRepository
}

func NewManager(gameRepository GamesRepository, seasonsRepository SeasonsRepository) Manager {
	return Manager{
		gameRepository:    gameRepository,
		seasonsRepository: seasonsRepository,
	}
}

func (manager Manager) CreateSeason(name string) (Season, error) {
	_, err := manager.seasonsRepository.FindSeasonByName(name)
	if err == nil {
		return Season{}, errors.New(fmt.Sprintf("Season with name %s exists", name))
	}
	if err != gorm.ErrRecordNotFound {
		return Season{}, err
	}

	season := Season{Name: name, Active: false}

	err = manager.seasonsRepository.CreateSeason(&season)
	if err != nil {
		return Season{}, err
	}

	return season, nil
}

func (manager Manager) GetSeasons() ([]Season, error) {
	return manager.seasonsRepository.GetAll()
}

func (manager Manager) ActivateSeason(name string) (Season, error) {
	season, err := manager.seasonsRepository.FindSeasonByName(name)
	if err != nil {
		return Season{}, err
	}

	manager.seasonsRepository.ActivateSeason(&season)

	return season, nil
}

func (manager Manager) CreateGame(playedAt time.Time) (*Game, error) {
	activeSeason, err := manager.seasonsRepository.FindActiveSeason()
	if err != nil {
		return &Game{}, err
	}

	if playedAt.IsZero() {
		playedAt = time.Now()
	}

	game := &Game{Season: &activeSeason, PlayedAt: playedAt}

	err = manager.gameRepository.Save(game)
	if err != nil {
		return &Game{}, err
	}

	return game, nil
}

func (manager Manager) ActiveSeason() (Season, error) {
	activeSeason, err := manager.seasonsRepository.FindActiveSeason()
	if err != nil {
		return Season{}, err
	}

	return activeSeason, nil
}

func (manager Manager) GetSeasonByName(name string) (Season, error) {
	return manager.seasonsRepository.FindSeasonByName(name)
}

func (manager Manager) GetGamesCount() (int, error) {
	return manager.gameRepository.Count()
}

func (manager Manager) GetSeasonByUuid(uuid string) (Season, error) {
	return manager.seasonsRepository.FindSeasonByUuid(uuid)
}
