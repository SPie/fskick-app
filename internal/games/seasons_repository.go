package games

import (
	"fmt"
	"time"

	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/uuid"
)

type Season struct {
	db.Model
	Name   string  `json:"name"`
	Active bool    `json:"active"`
	Games  []Game `json:"games"`
}

type SeasonsRepository struct {
	dbHandler db.Handler
	uuidGenerator uuid.Generator
}

func NewSeasonsRepository(
	dbHandler db.Handler,
	uuidGenerator uuid.Generator,
) SeasonsRepository {
	return SeasonsRepository{
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
	season, err := repository.selectSeason("name = $1", name)
	if err != nil {
		return Season{}, fmt.Errorf("query season by name: %w", err)
	}

	return season, nil
}

func (repository SeasonsRepository) FindSeasonByUuid(uuid string) (Season, error) {
	season, err := repository.selectSeason("uuid = $1", uuid)
	if err != nil {
		return Season{}, fmt.Errorf("query season by uuid: %w", err)
	}

	rows, err := repository.dbHandler.Query(
		fmt.Sprintf(
			`SELECT %s
			FROM games
			WHERE season_id = $1`,
			getGamesColumns(),
		),
		season.ID,
	)
	if err != nil {
		return Season{}, fmt.Errorf("query season games: %w", err)
	}
	defer rows.Close()

	games, err := scanGames(rows)
	if err != nil {
		return Season{}, fmt.Errorf("scan row in query season games: %w", err)
	}

	season.Games = games

	return season, nil
}

func (repository SeasonsRepository) GetAll() ([]Season, error) {
	rows, err := repository.dbHandler.Query(fmt.Sprintf(`SELECT %s FROM seasons`, getSeasonsColumns()))
	if err != nil {
		return []Season{}, fmt.Errorf("query all seasons: %w", err)
	}
	defer rows.Close()

	seasons, err := scanSeasons(rows)
	if err != nil {
		return []Season{}, fmt.Errorf("scan row in query all seasons: %w", err)
	}

	return seasons, nil
}

func (repository SeasonsRepository) FindActiveSeason() (Season, error) {
	season, err := repository.selectSeason("active = true")
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

func getSeasonsColumns() string {
	return "id, uuid, created_at, updated_at, name, active"
}

func (repository SeasonsRepository) selectSeason(whereQuery string, args ...any) (Season, error) {
	var season Season
	row := repository.dbHandler.QueryRow(
		fmt.Sprintf(
			`SELECT %s
			FROM seasons
			WHERE %s`,
			getSeasonsColumns(),
			whereQuery,
		),
		args...,
	)

	err := row.Scan(row, &season)
	if err != nil {
		return Season{}, fmt.Errorf("scan row in query active season: %w", err)
	}

	return season, nil
}

func scanSeason(row db.Row, season *Season) error {
	return row.Scan(
		season.ID,
		season.UUID,
		season.CreatedAt,
		season.UpdatedAt,
		season.Name,
		season.Active,
	)
}

func scanSeasons(rows db.Rows) ([]Season, error) {
	var seasons []Season
	for rows.Next() {
		var season Season

		err := rows.Scan(
			&season.ID,
			&season.UUID,
			&season.CreatedAt,
			&season.UpdatedAt,
			&season.Name,
			&season.Active,
		)
		if err != nil {
			return []Season{}, fmt.Errorf("scan season rows: %w", err)
		}

		seasons = append(seasons, season)
	}

	return seasons, nil
}
