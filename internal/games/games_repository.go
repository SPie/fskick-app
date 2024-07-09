package games

import (
	"fmt"
	"time"

	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/seasons"
	"github.com/spie/fskick/internal/uuid"
)

type Game struct {
	db.Model
	PlayedAt time.Time `json:"playedAt"`
	SeasonID uint      `json:"-"`
	Season   *seasons.Season   `json:"season"`
}

type GamesRepository struct {
	dbHandler db.Handler
}

func NewGamesRepository(dbHandler db.Handler) GamesRepository {
	return GamesRepository{dbHandler: dbHandler}
}

func (repository GamesRepository) CreateGame(game *Game) error {
	uuid, err := uuid.GenerateUuidString()
	if err != nil {
		return fmt.Errorf("create uuid for insert game: %w", err)
	}

	game.UUID = uuid
	game.CreatedAt = time.Now()
	game.UpdatedAt = time.Now()

	row := repository.dbHandler.QueryRow(
		`INSERT INTO games (uuid, played_at, season_id, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		game.UUID,
		game.PlayedAt,
		game.Season.ID,
		game.CreatedAt,
		game.UpdatedAt,
		nil,
	)
	err = row.Scan(&game.ID)
	if err != nil {
		return fmt.Errorf("insert game: %w", err)
	}

	return nil
}

func (repository GamesRepository) Count() (int, error) {
	var count int
	err := repository.dbHandler.
		QueryRow("SELECT COUNT(*) FROM games").
		Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count games: %w", err)
	}

	return count, nil
}

func (repository GamesRepository) CountForSeason(season seasons.Season) (int, error) {
	var count int

	err := repository.dbHandler.
		QueryRow("SELECT COUNT(*) FROM games WHERE season_id = $1", season.ID).
		Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count games: %w", err)
	}

	return count, nil
}

func getGamesColumns() string {
	return "id, uuid, created_at, updated_at, deleted_at, played_at"
}

func scanGames(rows db.Rows) ([]Game, error) {
	var games []Game
	for rows.Next() {
		var game Game

		err := rows.Scan(
			&game.ID,
			&game.UUID,
			&game.CreatedAt,
			&game.UpdatedAt,
			&game.DeletedAt,
			&game.PlayedAt,
		)
		if err != nil {
			return []Game{}, fmt.Errorf("scan games: %w", err)
		}

		games = append(games, game)
	}

	return games, nil
}
