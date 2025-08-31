package games

import (
	"database/sql"
	"fmt"

	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/players"
	"github.com/spie/fskick/internal/seasons"
)

type Attendance struct {
	db.Model
	Win      bool
	PlayerID uint
	GameID   uint
}

type PlayerAttendance struct {
	players.Player
	Wins  int
	Games int
}

type PlayerWithAttendances struct {
	players.Player
	Attendances []Attendance
}

type AttendanceRepository struct {
	conn db.Connection
}

func NewAttendanceRepository(conn db.Connection) AttendanceRepository {
	return AttendanceRepository{conn: conn}
}

func (repository AttendanceRepository) CollectPlayerAttendancesForSeason(
	season seasons.Season,
) ([]PlayerAttendance, error) {
	rows, err := repository.conn.Query(
		fmt.Sprintf(
			`SELECT
			%s
			FROM players p
			JOIN attendances a ON p.id = a.player_id
			JOIN games g ON g.id = a.game_id
			WHERE g.season_id = $1
			GROUP BY p.id
			`,
			getPlayerAttendanceColumns(),
		),
		season.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("collect player attendances for season: %w", err)
	}
	defer rows.Close()

	playerAttendances, err := scanPlayerAttendances(rows)
	if err != nil {
		return nil, fmt.Errorf("scan row for collect player attendances for season: %w", err)
	}

	return playerAttendances, nil
}

func (repository AttendanceRepository) CollectAllPlayerAttendances() ([]PlayerAttendance, error) {
	rows, err := repository.conn.Query(
		fmt.Sprintf(
			`SELECT
			%s
			FROM players p
			JOIN attendances a ON p.id = a.player_id
			JOIN games g ON g.id = a.game_id
			GROUP BY p.id
			`,
			getPlayerAttendanceColumns(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("collect all player attendances: %w", err)
	}
	defer rows.Close()

	playerAttendances, err := scanPlayerAttendances(rows)
	if err != nil {
		return nil, fmt.Errorf("scan row for collect all player attendances: %w", err)
	}

	return playerAttendances, nil
}

func (repository AttendanceRepository) CollectFellowPlayerAttendances(
	player players.Player,
) ([]PlayerAttendance, error) {
	rows, err := repository.conn.Query(
		fmt.Sprintf(
			`WITH player_games AS (
				SELECT g.id AS game_id, a.win
				FROM attendances a
				JOIN games g ON g.id = a.game_id
				WHERE a.player_id = $1
			)
			SELECT
			%s
			FROM players p
			JOIN attendances a ON p.id = a.player_id
			JOIN games g ON g.id = a.game_id
			JOIN player_games pg ON g.id = pg.game_id AND a.win = pg.win
			WHERE p.id != $1
			GROUP BY p.id
			`,
			getPlayerAttendanceColumns(),
		),
		player.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("collect fellow player attendances: %w", err)
	}
	defer rows.Close()

	playerAttendances, err := scanPlayerAttendances(rows)
	if err != nil {
		return nil, fmt.Errorf("scan row for collect fellow player attendances: %w", err)
	}

	return playerAttendances, nil
}

func (repository AttendanceRepository) CollectOponentPlayerAttendances(
	player players.Player,
) ([]PlayerAttendance, error) {
	rows, err := repository.conn.Query(
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
			SUM(CASE WHEN a.win THEN 0 ELSE 1 END) as wins
		FROM players p
		JOIN attendances a ON p.id = a.player_id
		JOIN games g ON g.id = a.game_id
		JOIN player_games pg ON g.id = pg.game_id AND a.win != pg.win
		WHERE p.id != $1
		GROUP BY p.id
		`,
		player.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("collect oponent player attendances: %w", err)
	}
	defer rows.Close()

	playerAttendances, err := scanPlayerAttendances(rows)
	if err != nil {
		return nil, fmt.Errorf("scan row for collect oponent player attendances: %w", err)
	}

	return playerAttendances, nil
}

func (repository AttendanceRepository) GetAttendancesForPlayer(player players.Player) ([]Attendance, error) {
	rows, err := repository.conn.Query(
		`SELECT a.id, a.uuid, a.win, a.created_at
		FROM attendances a
		JOIN games g ON a.game_id = g.id
		WHERE a.player_id = $1
		ORDER BY g.played_at ASC`,
		player.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("get last attendances for player: %w", err)
	}
	defer rows.Close()

	attendances := []Attendance{}
	for rows.Next() {
		var attendance Attendance
		err = rows.Scan(
			&attendance.ID,
			&attendance.UUID,
			&attendance.Win,
			&attendance.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan attendance rows: %w", err)
		}

		attendances = append(attendances, attendance)
	}

	return attendances, nil
}

func (repository AttendanceRepository) GetAttendancesForAllPlayers() ([]PlayerWithAttendances, error) {
	rows, err := repository.conn.Query(
		`SELECT p.id, p.uuid, p.name, p.created_at, a.id, a.uuid, a.win, a.created_at
		FROM players p
		JOIN attendances a ON p.id = a.player_id
		JOIN games g ON g.id = a.game_id
		ORDER BY g.played_at ASC
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("get all attendances for all players: %w", err)
	}
	defer rows.Close()

	playersWithAttendances := map[uint]*PlayerWithAttendances{}
	for rows.Next() {
		var player players.Player
		var attendance Attendance
		err = rows.Scan(
			&player.ID,
			&player.UUID,
			&player.Name,
			&player.CreatedAt,
			&attendance.ID,
			&attendance.UUID,
			&attendance.Win,
			&attendance.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan player with attendances rows: %w", err)
		}

		if _, ok := playersWithAttendances[player.ID]; !ok {
			playersWithAttendances[player.ID] = &PlayerWithAttendances{Player: player}
		}

		playersWithAttendances[player.ID].Attendances = append(
			playersWithAttendances[player.ID].Attendances,
			attendance,
		)
	}

	return getPlayersWithAttendancesFromMap(playersWithAttendances), nil
}

func getPlayerAttendanceColumns() string {
	return `
		p.id,
		p.uuid,
		p.name,
		p.created_at,
		p.updated_at,
		COUNT(a.id) AS games_played,
		SUM(CASE WHEN a.win THEN 1 ELSE 0 END) as wins`
}

func scanPlayerAttendances(rows *sql.Rows) ([]PlayerAttendance, error) {
	var playerAttendances []PlayerAttendance
	for rows.Next() {
		var playerAttendance PlayerAttendance
		err := rows.Scan(
			&playerAttendance.ID,
			&playerAttendance.UUID,
			&playerAttendance.Name,
			&playerAttendance.CreatedAt,
			&playerAttendance.UpdatedAt,
			&playerAttendance.Games,
			&playerAttendance.Wins,
		)
		if err != nil {
			return nil, fmt.Errorf("scan player attendances rows: %w", err)
		}

		playerAttendances = append(playerAttendances, playerAttendance)
	}

	return playerAttendances, nil
}

func getPlayersWithAttendancesFromMap(
	playersWithAttendanceMap map[uint]*PlayerWithAttendances,
) []PlayerWithAttendances {
	playersWithAttendances := []PlayerWithAttendances{}
	for _, playerWithAttendance := range playersWithAttendanceMap {
		playersWithAttendances = append(playersWithAttendances, *playerWithAttendance)
	}

	return playersWithAttendances
}
