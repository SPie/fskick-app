package games

import (
	"fmt"
	"time"

	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/players"
	"github.com/spie/fskick/internal/seasons"
)

type Game struct {
	db.Model
	PlayedAt time.Time `json:"playedAt"`
	SeasonID uint      `json:"-"`
	Season   *seasons.Season   `json:"season"`
	Attendances []Attendance
}

type GamesRepository struct {
	dbHandler db.Handler
}

func NewGamesRepository(dbHandler db.Handler) GamesRepository {
	return GamesRepository{dbHandler: dbHandler}
}

func (repository GamesRepository) CreateGame(game *Game, attendances []Attendance) error {
	err := game.CreateUUID()
	if err != nil {
		return fmt.Errorf("create uuid for insert game: %w", err)
	}

	now := time.Now()

	game.CreatedAt = now
	game.UpdatedAt = now

	tx, err := repository.dbHandler.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction for insert game: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	row := tx.QueryRow(
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

	createdAttendancees := make([]Attendance, len(attendances))
	for i, attendance := range attendances {
		err = attendance.CreateUUID()
		if err != nil {
			return fmt.Errorf("create uuid for insert attendance: %w", err)
		}

		attendance.CreatedAt = now
		attendance.UpdatedAt = now
		attendance.GameID = game.ID

		row = tx.QueryRow(
			`INSERT INTO attendances (uuid, win, player_id, game_id, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id`,
			attendance.UUID,
			attendance.Win,
			attendance.PlayerID,
			attendance.GameID,
			attendance.CreatedAt,
			attendance.UpdatedAt,
			nil,
		)
		err = row.Scan(&attendance.ID)
		if err != nil {
			return fmt.Errorf("insert attendance: %w", err)
		}

		createdAttendancees[i] = attendance
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("insert game: %w", err)
	}

	game.Attendances = createdAttendancees

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

func (repository GamesRepository) CountForPlayer(player players.Player) (int, error) {
	var count int
	err := repository.dbHandler.
		QueryRow(
			`SELECT COUNT(*)
			FROM games g
			JOIN attendances a
			WHERE a.player_id = $1`,
			player.ID,
		).
		Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count games: %w", err)
	}

	return count, nil
}

func (repository GamesRepository) MaxGamesForSeason(season seasons.Season) (int, error) {
	var maxGames int
	row := repository.dbHandler.QueryRow(
		`SELECT MAX(games_played) as max_games_played
		FROM (
			SELECT COUNT(a.id) as games_played
			FROM players p
			JOIN attendances a ON p.id = a.player_id
			JOIN games g ON g.id = a.game_id
			WHERE g.season_id = $1
			GROUP BY p.id
		)`,
		season.ID,
	)

	err := row.Scan(&maxGames)
	if err != nil {
		return 0, fmt.Errorf("query max games for seasons: %w", err)
	}

	return maxGames, nil
}

func (repository GamesRepository) MaxGames() (int, error) {
	var maxGames int
	row := repository.dbHandler.QueryRow(
		`SELECT MAX(games_played) as max_games_played
		FROM (
			SELECT COUNT(a.id) as games_played
			FROM players p
			JOIN attendances a ON p.id = a.player_id
			JOIN games g ON g.id = a.game_id
			GROUP BY p.id
		)`,
	)

	err := row.Scan(&maxGames)
	if err != nil {
		return 0, fmt.Errorf("query max games: %w", err)
	}

	return maxGames, nil
}

func (repository GamesRepository) MaxGamesForPlayer(player players.Player) (int, error) {
	var maxGames int
	row := repository.dbHandler.QueryRow(
		`WITH player_games AS (
			SELECT g.id AS game_id, a.win
			FROM attendances a
			JOIn games g ON a.game_id = g.id
			WHERE a.player_id = $1
		)
		SELECT MAX(games_played) AS max_games_played
		FROM (
			SELECT p.id, COUNT(a.id) as games_played
			FROM players p
			JOIN attendances a ON p.id = a.player_id
			JOIN games g ON g.id = a.game_id
			JOIN player_games pg ON g.id = pg.game_id AND a.win = pg.win
			WHERE p.id != $1
			GROUP BY p.id
		) subquery`,
		player.ID,
	)

	err := row.Scan(&maxGames)
	if err != nil {
		return 0, fmt.Errorf("query max games for player: %w", err)
	}

	return maxGames, nil
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
