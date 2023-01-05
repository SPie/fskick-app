package games

import (
	"errors"
	"fmt"
	"time"

	"github.com/spie/fskick/db"
	"gorm.io/gorm"
)

type Manager interface {
	CreateSeason(name string) (*Season, error)
	GetSeasons() ([]Season, error)
	ActivateSeason(name string) (Season, error)
	CreateGame(playedAt time.Time) (*Game, error)
	ActiveSeason() (Season, error)
	GetSeasonByName(name string) (Season, error)
	GetSeasonByUuid(uuid string) (Season, error)
	GetGamesCount() (int, error)
}

type manager struct {
	gameRepository    GamesRepository
	seasonsRepository SeasonsRepository
}

func NewManager(gameRepository GamesRepository, seasonsRepository SeasonsRepository) Manager {
	return manager{
		gameRepository:    gameRepository,
		seasonsRepository: seasonsRepository,
	}
}

func (manager manager) CreateSeason(name string) (*Season, error) {
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

func (manager manager) GetSeasons() ([]Season, error) {
	return manager.seasonsRepository.GetAll()
}

func (manager manager) ActivateSeason(name string) (Season, error) {
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

func (manager manager) CreateGame(playedAt time.Time) (*Game, error) {
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

func (manager manager) ActiveSeason() (Season, error) {
	activeSeason, err := manager.seasonsRepository.FindActiveSeason()
	if err != nil {
		return Season{}, err
	}

	return activeSeason, nil
}

func (manager manager) GetSeasonByName(name string) (Season, error) {
	return manager.seasonsRepository.FindSeasonByName(name)
}

func (manager manager) GetGamesCount() (int, error) {
	return manager.gameRepository.Count()
}

func (manager manager) GetSeasonByUuid(uuid string) (Season, error) {
	return manager.seasonsRepository.FindSeasonByUuid(uuid)
}

type Game struct {
	db.Model
	PlayedAt time.Time `json:"playedAt"`
	SeasonID uint      `json:"-"`
	Season   *Season   `json:"season"`
}

type GamesRepository interface {
	db.Repository
	Count() (int, error)
	GetAll() (*[]Game, error)
	Save(game *Game) error
}

type gamesRepository struct {
	connectionHandler db.ConnectionHandler
}

func NewGamesRepository(connectionHandler db.ConnectionHandler) GamesRepository {
	return &gamesRepository{connectionHandler: connectionHandler}
}

func (repository *gamesRepository) Save(game *Game) error {
	if game.ID == 0 {
		return repository.connectionHandler.Create(game)
	}

	return repository.connectionHandler.Save(game)
}

func (repository *gamesRepository) Count() (int, error) {
	return repository.connectionHandler.Count(&Game{})
}

func (repository *gamesRepository) GetAll() (*[]Game, error) {
	games := &[]Game{}

	err := repository.connectionHandler.GetAll(games)
	if err != nil {
		return &[]Game{}, err
	}

	return games, nil
}

func (repository *gamesRepository) AutoMigrate() {
	repository.connectionHandler.AutoMigrate(&Game{})

	repository.connectionHandler.Exec("UPDATE games SET played_at = created_at WHERE played_at IS NULL")
}

type Season struct {
	db.Model
	Name   string  `gorm:"unique;not null" json:"name"`
	Active bool    `gorm:"default:false" json:"active"`
	Games  *[]Game `json:"games"`
}

type SeasonsRepository interface {
	db.Repository
	Save(season *Season) error
	Find(condition *Season) (*[]Season, error)
	FindActiveSeason() (Season, error)
	FindSeasonByName(name string) (Season, error)
	FindSeasonByUuid(uuid string) (Season, error)
	GetAll() ([]Season, error)
}

type seasonsRepository struct {
	connectionHandler db.ConnectionHandler
}

func NewSeasonsRepository(connectionHandler db.ConnectionHandler) SeasonsRepository {
	return &seasonsRepository{connectionHandler: connectionHandler}
}

func (repository *seasonsRepository) Save(season *Season) error {
	if season.ID == 0 {
		return repository.connectionHandler.Create(season)
	}

	return repository.connectionHandler.Save(season)
}

func (repository *seasonsRepository) FindSeasonByName(name string) (Season, error) {
	season := &Season{}

	err := repository.connectionHandler.Preload("Games").FindOne(season, &Season{Name: name})
	if err != nil {
		return Season{}, err
	}

	return *season, nil
}

func (repository *seasonsRepository) FindSeasonByUuid(uuid string) (Season, error) {
	season := &Season{}

	err := repository.connectionHandler.Preload("Games").FindOne(season, &Season{Model: db.Model{UUID: uuid}})
	if err != nil {
		return Season{}, err
	}

	return *season, nil
}

func (repository *seasonsRepository) GetAll() ([]Season, error) {
	seasons := &[]Season{}

	err := repository.connectionHandler.Preload("Games").GetAll(seasons)
	if err != nil {
		return []Season{}, err
	}

	return *seasons, nil
}

func (repository *seasonsRepository) FindActiveSeason() (Season, error) {
	season := &Season{}

	err := repository.connectionHandler.Preload("Games").FindOne(season, &Season{Active: true})
	if err != nil {
		return Season{}, err
	}

	return *season, nil
}

func (repository *seasonsRepository) Find(condition *Season) (*[]Season, error) {
	seasons := &[]Season{}

	err := repository.connectionHandler.Preload("Games").Find(seasons, condition)
	if err != nil {
		return &[]Season{}, err
	}

	return seasons, nil
}

func (repository *seasonsRepository) AutoMigrate() {
	repository.connectionHandler.AutoMigrate(&Season{})
}
