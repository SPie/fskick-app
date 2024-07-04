package games

import (
	"fmt"
	"time"

	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/uuid"
)

type Season struct {
	db.Model
	Name   string  `gorm:"unique;not null" json:"name"`
	Active bool    `gorm:"default:false" json:"active"`
	Games  []Game `json:"games"`
}

type SeasonsRepository struct {
	connectionHandler *db.ConnectionHandler
	dbHandler db.Handler
	uuidGenerator uuid.Generator
}

func NewSeasonsRepository(
	connectionHandler *db.ConnectionHandler,
	dbHandler db.Handler,
	uuidGenerator uuid.Generator,
) SeasonsRepository {
	return SeasonsRepository{
		connectionHandler: connectionHandler,
		dbHandler: dbHandler,
		uuidGenerator: uuidGenerator,
	}
}

func (repository SeasonsRepository) CreateSeason(season *Season) error {
	uuid, err := repository.uuidGenerator.GenerateUuidString()
	if err != nil {
		return fmt.Errorf("create uuid for insert season: %w", err)
	}

	season.UUID = uuid
	season.CreatedAt = time.Now()
	season.UpdatedAt = time.Now()

	row := repository.dbHandler.QueryRow(
		`INSERT INTO seasons (uuid, name, created_at, updated_at, deleted_at, active)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		season.UUID,
		season.Name,
		season.CreatedAt,
		season.UpdatedAt,
		nil,
		season.Active,
	)
	err = row.Scan(season.ID)
	if err != nil {
		return fmt.Errorf("insert season: %w", err)
	}

	return nil
}

func (repository SeasonsRepository) FindSeasonByName(name string) (Season, error) {
	var season Season

	row := repository.dbHandler.QueryRow(
		`SELECT s.id, s.uuid, s.created_at, s.updated_at, s.name, s.active
		FROM seasons s
		WHERE s.name = $1`,
		name,
	)
	err := row.Scan(
		&season.ID,
		&season.UUID,
		&season.CreatedAt,
		&season.UpdatedAt,
		&season.Name,
		&season.Active,
	)
	if err != nil {
		return Season{}, fmt.Errorf("query season by name: %w", err)
	}

	return season, nil
}

func (repository SeasonsRepository) FindSeasonByUuid(uuid string) (Season, error) {
	var season Season

	row := repository.dbHandler.QueryRow(
		`SELECT s.id, s.uuid, s.created_at, s.updated_at, s.name, s.active
		FROM seasons s
		WHERE s.uuid = $1`,
		uuid,
	)
	err := row.Scan(
		&season.ID,
		&season.UUID,
		&season.CreatedAt,
		&season.UpdatedAt,
		&season.Name,
		&season.Active,
	)
	if err != nil {
		return Season{}, fmt.Errorf("query season by uuid: %w", err)
	}

	rows, err := repository.dbHandler.Query(
		`SELECT g.id, g.uuid, g.created_at, g.updated_at, g.deleted_at, g.played_at
		FROM games g
		WHERE season_id = $1`,
		season.ID,
	)
	if err != nil {
		return Season{}, fmt.Errorf("query season games: %w", err)
	}
	defer rows.Close()

	games := []Game{}
	for rows.Next() {
		var game Game

		err = rows.Scan(
			&game.ID,
			&game.UUID,
			&game.CreatedAt,
			&game.UpdatedAt,
			&game.DeletedAt,
			&game.PlayedAt,
		)
		if err != nil {
			return Season{}, fmt.Errorf("scan row in query season games: %w", err)
		}

		games = append(games, game)
	}

	season.Games = games

	return season, nil
}

func (repository SeasonsRepository) GetAll() ([]Season, error) {
	rows, err := repository.dbHandler.Query(`
		SELECT s.id, s.uuid, s.created_at, s.updated_at, s.name, s.active
		FROM seasons s
	`)
	if err != nil {
		return []Season{}, fmt.Errorf("query all seasons: %w", err)
	}
	defer rows.Close()

	seasons := []Season{}
	for rows.Next() {
		var season Season

		err = rows.Scan(
			&season.ID,
			&season.UUID,
			&season.CreatedAt,
			&season.UpdatedAt,
			&season.Name,
			&season.Active,
		)
		if err != nil {
			return []Season{}, fmt.Errorf("scan row in query all seasons: %w", err)
		}

		seasons = append(seasons, season)
	}

	return seasons, nil
}

func (repository SeasonsRepository) FindActiveSeason() (Season, error) {
	var season Season

	row := repository.dbHandler.QueryRow(
		`SELECT s.id, s.uuid, s.created_at, s.updated_at, s.name, s.active
		FROM seasons s
		WHERE s.active = true`,
	)

	err := row.Scan(
		&season.ID,
		&season.UUID,
		&season.CreatedAt,
		&season.UpdatedAt,
		&season.Name,
		&season.Active,
	)
	if err != nil {
		return Season{}, fmt.Errorf("scan row in query active season: %w", err)
	}

	return season, nil
}

func (repository SeasonsRepository) ActivateSeason(season *Season) error {
	tx, err := repository.dbHandler.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction in activate season: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec("UPDATE seasons SET active = false WHERE active = true")
	if err != nil {
		return fmt.Errorf("deactivate active seasons in activate season: %w", err)
	}

	_, err = tx.Exec("UPDATE seasons SET active = true WHERE id = $1", season.ID)
	if err != nil {
		return fmt.Errorf("activate season: %w", err)
	}

	tx.Commit()

	season.Active = true

	return nil
}
