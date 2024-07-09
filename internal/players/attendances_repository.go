package players

import (
	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/games"
)

type Attendance struct {
	db.Model
	Win      bool        `json:"win"`
	PlayerID uint        `json:"-"`
	Player   *Player     `json:"-"`
	GameID   uint        `json:"-"`
	Game     *games.Game `json:"game"`
}

type AttendancesRepository struct {
	connectionHandler *db.ConnectionHandler
	dbHandler db.Handler
}

func NewAttendancesRepository(
	connectionHandler *db.ConnectionHandler,
	dbHandler db.Handler,
) AttendancesRepository {
	return AttendancesRepository{
		connectionHandler: connectionHandler,
		dbHandler: dbHandler,
	}
}

func (repository AttendancesRepository) Create(attendances *[]Attendance) error {
	return repository.connectionHandler.Create(attendances)
}

func (repository AttendancesRepository) FindFellowAttendancesForPlayer(player Player) ([]Attendance, error) {
	fellowWinnerAttendances, err := repository.findFellowAttendancesForPlayerByWin(player, true)
	if err != nil {
		return []Attendance{}, err
	}

	fellowLoserAttendances, err := repository.findFellowAttendancesForPlayerByWin(player, false)
	if err != nil {
		return []Attendance{}, err
	}

	fellowAttendances := append(fellowWinnerAttendances, fellowLoserAttendances...)

	return fellowAttendances, nil
}

func (repository AttendancesRepository) findFellowAttendancesForPlayerByWin(player Player, win bool) ([]Attendance, error) {
	attendances := &[]Attendance{}
	err := repository.connectionHandler.
		Where("win = ? AND player_id = ?", win, player.ID).
		Find(attendances)
	if err != nil {
		return []Attendance{}, err
	}

	fellowAttendances := []Attendance{}
	err = repository.connectionHandler.
		Preload("Player").
		Where("win = ? AND game_id IN ? AND player_id != ?", win, getGameIdsFromAttendances(attendances), player.ID).
		Find(&fellowAttendances)
	if err != nil {
		return []Attendance{}, err
	}

	return fellowAttendances, nil
}

func getGameIdsFromAttendances(attendances *[]Attendance) []uint {
	gameIds := make([]uint, len(*attendances))
	for i, attendance := range *attendances {
		gameIds[i] = attendance.GameID
	}

	return gameIds
}
