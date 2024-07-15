package players

import (
	"fmt"
	"strings"
	"time"

	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/uuid"
)

var (
	ErrPlayerNotFound = db.ErrNotFound
)

type PlayerRepository struct {
	connectionHandler *db.ConnectionHandler
	dbHandler db.Handler
}

func NewPlayerRepository(
	connectionHandler *db.ConnectionHandler,
	dbHandler db.Handler,
) PlayerRepository {
	return PlayerRepository{
		connectionHandler: connectionHandler,
		dbHandler: dbHandler,
	}
}

func (repository PlayerRepository) CreatePlayer(player *Player) error {
	uuid, err := uuid.GenerateUuidString()
	if err != nil {
		return fmt.Errorf("create uuid for insert player: %w", err)
	}

	player.UUID = uuid
	player.CreatedAt = time.Now()
	player.UpdatedAt = time.Now()

	row := repository.dbHandler.QueryRow(
		`INSERT INTO players (uuid, name, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5)`,
		player.UUID,
		player.Name,
		player.CreatedAt,
		player.UpdatedAt,
		nil,
	)
	err = row.Scan(&player.ID)
	if err != nil {
		return fmt.Errorf("insert player: %w", err)
	}

	return nil
}

func (repository PlayerRepository) FindPlayerByUUID(uuid string) (Player, error) {
	player := Player{}
	row := repository.dbHandler.QueryRow(
		`SELECT id, uuid, created_at, updated_at, name
		FROM players
		WHERE uuid = $1`,
		uuid,
	)

	err := row.Scan(
		&player.ID,
		&player.UUID,
		&player.CreatedAt,
		&player.UpdatedAt,
		&player.Name,
	)
	if err != nil {
		return Player{}, fmt.Errorf("query player by uuid: %w", err)
	}

	return player, nil
}

func (repository PlayerRepository) FindPlayerByName(name string) (Player, error) {
	player := Player{}
	row := repository.dbHandler.QueryRow(
		`SELECT id, uuid, created_at, updated_at, name
		FROM players
		WHERE name = $1`,
		name,
	)

	err := row.Scan(
		&player.ID,
		&player.UUID,
		&player.CreatedAt,
		&player.UpdatedAt,
		&player.Name,
	)
	if err != nil {
		return Player{}, fmt.Errorf("query player by name: %w", err)
	}

	return player, nil
}

func (repository PlayerRepository) FindPlayersByNames(names []string) ([]Player, error) {
	players := []Player{}

	rows, err := repository.dbHandler.Query(
		`SELECT id, uuid, name, created_at, updated_at FROM players WHERE name IN ($1)`,
		strings.Join(names, ","),
	)
	if err != nil {
		return []Player{}, fmt.Errorf("Query players by names: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var player Player
		err := rows.Scan(
			&player.ID,
			&player.UUID,
			&player.Name,
			&player.CreatedAt,
			&player.UpdatedAt,
		)
		if err != nil {
			return []Player{}, fmt.Errorf("Scan player rows: %w", err)
		}

		players = append(players, player)
	}

	return players, nil
}
