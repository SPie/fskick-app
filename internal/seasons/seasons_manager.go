package seasons

import (
	"errors"
	"fmt"
)

type Manager struct {
	seasonsRepository SeasonsRepository
}

func NewManager(seasonRepository SeasonsRepository) Manager {
	return Manager{seasonsRepository: seasonRepository}
}

func (manager Manager) CreateSeason(name string) (Season, error) {
	_, err := manager.seasonsRepository.FindSeasonByName(name)
	if err == nil {
		return Season{}, errors.New(fmt.Sprintf("Season with name %s exists", name))
	}
	if !errors.Is(err, ErrSeasonNotFound) {
		return Season{}, fmt.Errorf("Check for season with name in CreateSeason: %w", err)
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

func (manager Manager) ActiveSeason() (Season, error) {
	activeSeason, err := manager.seasonsRepository.FindActiveSeason()
	if err != nil {
		return Season{}, err
	}

	return activeSeason, nil
}

func (manager Manager) GetSeasonByUuid(uuid string) (Season, error) {
	return manager.seasonsRepository.FindSeasonByUuid(uuid)
}

func (manager Manager) GetSeasonByName(name string) (Season, error) {
	return manager.seasonsRepository.FindSeasonByName(name)
}
