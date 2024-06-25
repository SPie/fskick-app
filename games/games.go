package games

import (
	"errors"
	"fmt"
	"time"

	"github.com/spie/fskick/db"
	"gorm.io/gorm"
)

type Manager struct {
	gameRepository    *GamesRepository
	seasonsRepository *SeasonsRepository
}

func NewManager(gameRepository *GamesRepository, seasonsRepository *SeasonsRepository) Manager {
	return Manager{
		gameRepository:    gameRepository,
		seasonsRepository: seasonsRepository,
	}
}

func (manager Manager) CreateSeason(name string) (*Season, error) {
	_, err := manager.seasonsRepository.FindSeasonByName(name)
	if err == nil {
		return &Season{}, errors.New(fmt.Sprintf("Season with name %s exists", name))
	}
	if err != gorm.ErrRecordNotFound {
		return &Season{}, err
	}

	season := &Season{Name: name}

	err = manager.seasonsRepository.Save(season)
	if err != nil {
		return &Season{}, err
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

	activeSeason, err := manager.seasonsRepository.FindActiveSeason()
	if err != nil && err != gorm.ErrRecordNotFound {
		return Season{}, err
	}

	if err != gorm.ErrRecordNotFound {
		activeSeason.Active = false
		manager.seasonsRepository.Save(&activeSeason)
	}

	season.Active = true
	manager.seasonsRepository.Save(&season)

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

type Game struct {
	db.Model
	PlayedAt time.Time `json:"playedAt"`
	SeasonID uint      `json:"-"`
	Season   *Season   `json:"season"`
}

type GamesRepository struct {
	connectionHandler *db.ConnectionHandler
}

func NewGamesRepository(connectionHandler *db.ConnectionHandler) *GamesRepository {
	return &GamesRepository{connectionHandler: connectionHandler}
}

func (repository *GamesRepository) Save(game *Game) error {
	if game.ID == 0 {
		return repository.connectionHandler.Create(game)
	}

	return repository.connectionHandler.Save(game)
}

func (repository *GamesRepository) Count() (int, error) {
	return repository.connectionHandler.Count(&Game{})
}

func (repository *GamesRepository) GetAll() (*[]Game, error) {
	games := &[]Game{}

	err := repository.connectionHandler.GetAll(games)
	if err != nil {
		return &[]Game{}, err
	}

	return games, nil
}

func (repository *GamesRepository) AutoMigrate() {
	repository.connectionHandler.AutoMigrate(&Game{})

	repository.connectionHandler.Exec("UPDATE games SET played_at = created_at WHERE played_at IS NULL")
}

type Season struct {
	db.Model
	Name   string  `gorm:"unique;not null" json:"name"`
	Active bool    `gorm:"default:false" json:"active"`
	Games  *[]Game `json:"games"`
}

type SeasonsRepository struct {
	connectionHandler *db.ConnectionHandler
}

func NewSeasonsRepository(connectionHandler *db.ConnectionHandler) *SeasonsRepository {
	return &SeasonsRepository{connectionHandler: connectionHandler}
}

func (repository *SeasonsRepository) Save(season *Season) error {
	if season.ID == 0 {
		return repository.connectionHandler.Create(season)
	}

	return repository.connectionHandler.Save(season)
}

func (repository *SeasonsRepository) FindSeasonByName(name string) (Season, error) {
	season := &Season{}

	err := repository.connectionHandler.Preload("Games").FindOne(season, &Season{Name: name})
	if err != nil {
		return Season{}, err
	}

	return *season, nil
}

func (repository *SeasonsRepository) FindSeasonByUuid(uuid string) (Season, error) {
	season := &Season{}

	err := repository.connectionHandler.Preload("Games").FindOne(season, &Season{Model: db.Model{UUID: uuid}})
	if err != nil {
		return Season{}, err
	}

	return *season, nil
}

func (repository *SeasonsRepository) GetAll() ([]Season, error) {
	seasons := &[]Season{}

	err := repository.connectionHandler.Preload("Games").GetAll(seasons)
	if err != nil {
		return []Season{}, err
	}

	return *seasons, nil
}

func (repository *SeasonsRepository) FindActiveSeason() (Season, error) {
	season := &Season{}

	err := repository.connectionHandler.Preload("Games").FindOne(season, &Season{Active: true})
	if err != nil {
		return Season{}, err
	}

	return *season, nil
}

func (repository *SeasonsRepository) Find(condition *Season) (*[]Season, error) {
	seasons := &[]Season{}

	err := repository.connectionHandler.Preload("Games").Find(seasons, condition)
	if err != nil {
		return &[]Season{}, err
	}

	return seasons, nil
}

func (repository *SeasonsRepository) AutoMigrate() {
	repository.connectionHandler.AutoMigrate(&Season{})
}
