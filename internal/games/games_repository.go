package games

import (
	"fmt"
	"time"

	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/uuid"
)

type Game struct {
	db.Model
	PlayedAt time.Time `json:"playedAt"`
	SeasonID uint      `json:"-"`
	Season   *Season   `json:"season"`
}

type GamesRepository struct {
	connectionHandler *db.ConnectionHandler
	dbHandler db.Handler
	uuidGenerator uuid.Generator
}

func NewGamesRepository(
	connectionHandler *db.ConnectionHandler,
	dbHandler db.Handler,
	uuidGenerator uuid.Generator,
) GamesRepository {
	return GamesRepository{
		connectionHandler: connectionHandler,
		dbHandler: dbHandler,
		uuidGenerator: uuidGenerator,
	}
}

func (repository GamesRepository) Save(game *Game) error {
	if game.ID == 0 {
		return repository.connectionHandler.Create(game)
	}

	return repository.connectionHandler.Save(game)
}

func (repository GamesRepository) Count() (int, error) {
	return repository.connectionHandler.Count(&Game{})
}

func (repository GamesRepository) GetAll() (*[]Game, error) {
	games := &[]Game{}

	err := repository.connectionHandler.GetAll(games)
	if err != nil {
		return &[]Game{}, err
	}

	return games, nil
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
