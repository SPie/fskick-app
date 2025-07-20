package players

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/spie/fskick/internal/db"
)

type Player struct {
	db.Model
	Name string
}

var (
	ErrPlayerNotFound = db.ErrNotFound
)

type PlayerRepository struct {
	conn db.Connection
}

func NewPlayerRepository(conn db.Connection) PlayerRepository {
	return PlayerRepository{
		conn: conn,
	}
}

func (repository PlayerRepository) CreatePlayer(player *Player) error {
	err := player.CreateUUID()
	if err != nil {
		return fmt.Errorf("create uuid for insert player: %w", err)
	}

	player.CreatedAt = time.Now()
	player.UpdatedAt = time.Now()

	row := repository.conn.QueryRow(
		`INSERT INTO players (uuid, name, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
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
	row := repository.conn.QueryRow(
		fmt.Sprintf(
			`SELECT %s
			FROM players
			WHERE uuid = $1`,
			getPlayerColumns(),
		),
		uuid,
	)

	player, err := scanPlayer(row)
	if err != nil {
		return Player{}, fmt.Errorf("query player by uuid: %w", err)
	}

	return player, nil
}

func (repository PlayerRepository) FindPlayerByName(name string) (Player, error) {
	row := repository.conn.QueryRow(
		fmt.Sprintf(
			`SELECT %s
			FROM players
			WHERE name = $1`,
			getPlayerColumns(),
		),
		name,
	)

	player, err := scanPlayer(row)
	if err != nil {
		return Player{}, fmt.Errorf("query player by name: %w", err)
	}

	return player, nil
}

func (repository PlayerRepository) FindPlayersByNames(names []string) ([]Player, error) {
	rows, err := repository.conn.Query(
		fmt.Sprintf(
			`SELECT
			%s
			FROM players
			WHERE name IN (%s)`,
			getPlayerColumns(),
			getInPlaceholders(names),
		),
		nameStringsToParameters(names)...,
	)
	if err != nil {
		return []Player{}, fmt.Errorf("Query players by names: %w", err)
	}
	defer rows.Close()

	players, err := scanPlayers(rows)
	if err != nil {
		return []Player{}, fmt.Errorf("Scan player rows: %w", err)
	}

	return players, nil
}

func getPlayerColumns() string {
	return `
		id,
		uuid,
		name,
		created_at,
		updated_at
	`
}

func scanPlayer(row *sql.Row) (Player, error) {
	var player Player
	err := row.Scan(
		&player.ID,
		&player.UUID,
		&player.Name,
		&player.CreatedAt,
		&player.UpdatedAt,
	)
	if err != nil {
		return Player{}, err
	}

	return player, nil
}

func scanPlayers(rows *sql.Rows) ([]Player, error) {
	players := []Player{}
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
			return []Player{}, err
		}

		players = append(players, player)
	}

	return players, nil
}

func getInPlaceholders(names []string) string {
	placeholders := make([]string, len(names))
	for i := range names {
		placeholders[i] = fmt.Sprintf("$%d", i)
	}

	return strings.Join(placeholders, ",")
}

func nameStringsToParameters(names []string) []any {
	parameters := make([]any, len(names))
	for i, name := range names {
		parameters[i] = name
	}

	return parameters
}
