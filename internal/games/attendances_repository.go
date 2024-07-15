package games

import (
	"fmt"

	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/players"
	"github.com/spie/fskick/internal/seasons"
)

type Attendance struct {
	db.Model
	Win      bool        `json:"win"`
	PlayerID uint        `json:"-"`
	GameID   uint        `json:"-"`
}

type PlayerAttendance struct {
	players.Player
	Wins int
	Games int
}

type AttendanceRepository struct {
	dbHandler db.Handler
}

func NewAttendanceRepository(dbHandler db.Handler) AttendanceRepository {
	return AttendanceRepository{dbHandler: dbHandler}
}

func (repository AttendanceRepository) CollectPlayerAttendancesForSeason(
	season seasons.Season,
) ([]PlayerAttendance, error) {
	rows, err := repository.dbHandler.Query(
		`SELECT
			p.id,
			p.uuid,
			p.name,
			p.created_at,
			p.updated_at,
			COUNT(a.id) AS games_played,
			SUM(CASE WHEN a.win THEN 1 ELSE 0 END) as wins
		FROM players p
		JOIN attendances a ON p.id = a.player_id
		JOIN games g ON g.id = a.game_id
		WHERE g.season_id = $1
		GROUP BY p.id
		`,
		season.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("collect player attendances for season: %w", err)
	}
	defer rows.Close()

	var playerAttendances []PlayerAttendance
	for rows.Next() {
		var playerAttendance PlayerAttendance
		err = rows.Scan(
			&playerAttendance.ID,
			&playerAttendance.UUID,
			&playerAttendance.Name,
			&playerAttendance.CreatedAt,
			&playerAttendance.UpdatedAt,
			&playerAttendance.Games,
			&playerAttendance.Wins,
		)
		if err != nil {
			return nil, fmt.Errorf("scan row for collect player attendances for season: %w", err)
		}

		playerAttendances = append(playerAttendances, playerAttendance)
	}

	return playerAttendances, nil
}

func (repository AttendanceRepository) CollectAllPlayerAttendances() ([]PlayerAttendance, error) {
	rows, err := repository.dbHandler.Query(
		`SELECT
			p.id,
			p.uuid,
			p.name,
			p.created_at,
			p.updated_at,
			COUNT(a.id) AS games_played,
			SUM(CASE WHEN a.win THEN 1 ELSE 0 END) as wins
		FROM players p
		JOIN attendances a ON p.id = a.player_id
		JOIN games g ON g.id = a.game_id
		GROUP BY p.id
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("collect all player attendances: %w", err)
	}
	defer rows.Close()

	var playerAttendances []PlayerAttendance
	for rows.Next() {
		var playerAttendance PlayerAttendance
		err = rows.Scan(
			&playerAttendance.ID,
			&playerAttendance.UUID,
			&playerAttendance.Name,
			&playerAttendance.CreatedAt,
			&playerAttendance.UpdatedAt,
			&playerAttendance.Games,
			&playerAttendance.Wins,
		)
		if err != nil {
			return nil, fmt.Errorf("scan row for collect all player attendances: %w", err)
		}

		playerAttendances = append(playerAttendances, playerAttendance)
	}

	return playerAttendances, nil
}

func (repository AttendanceRepository) CollectFellowPlayerAttendances(
	player players.Player,
) ([]PlayerAttendance, error) {
	rows, err := repository.dbHandler.Query(
		`WITH player_games AS (
			SELECT g.id AS game_id, a.win
			FROM attendances a
			JOIN games g ON g.id = a.game_id
			WHERE a.player_id = $1
		)
		SELECT
			p.id,
			p.uuid,
			p.name,
			p.created_at,
			p.updated_at,
			COUNT(a.id) AS games_played,
			SUM(CASE WHEN a.win THEN 1 ELSE 0 END) as wins
		FROM players p
		JOIN attendances a ON p.id = a.player_id
		JOIN games g ON g.id = a.game_id
		JOIN player_games pg ON g.id = pg.game_id AND a.win = pg.win
		WHERE p.id != $1
		GROUP BY p.id
		`,
		player.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("collect fellow player attendances: %w", err)
	}
	defer rows.Close()

	var playerAttendances []PlayerAttendance
	for rows.Next() {
		var playerAttendance PlayerAttendance
		err = rows.Scan(
			&playerAttendance.ID,
			&playerAttendance.UUID,
			&playerAttendance.Name,
			&playerAttendance.CreatedAt,
			&playerAttendance.UpdatedAt,
			&playerAttendance.Games,
			&playerAttendance.Wins,
		)
		if err != nil {
			return nil, fmt.Errorf("scan row for collect fellow player attendances: %w", err)
		}

		playerAttendances = append(playerAttendances, playerAttendance)
	}

	return playerAttendances, nil
}
