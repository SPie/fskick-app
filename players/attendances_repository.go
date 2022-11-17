package players

import (
	"github.com/spie/fskick/db"
	"github.com/spie/fskick/games"
)

type Attendance struct {
	db.Model
	Win      bool        `json:"win"`
	PlayerID uint        `json:"-"`
	Player   *Player     `json:"-"`
	GameID   uint        `json:"-"`
	Game     *games.Game `json:"game"`
}

type AttendancesRepository interface {
	db.Repository
	Save(attendance *Attendance) error
	FindAttendancesForSeason(season *games.Season) (*[]Attendance, error)
	FindAttendancesForPlayer(player *Player) (*[]Attendance, error)
	FindFellowAttendancesForPlayer(player *Player) (*[]Attendance, error)
	Create(attendances *[]Attendance) error
}

type attendancesRepository struct {
	connectionHandler db.ConnectionHandler
}

func NewAttendancesRepository(connectionHandler db.ConnectionHandler) AttendancesRepository {
	return &attendancesRepository{connectionHandler: connectionHandler}
}

func (repository *attendancesRepository) Save(attendance *Attendance) error {
	return repository.connectionHandler.Save(attendance)
}

func (repository *attendancesRepository) FindAttendancesForSeason(season *games.Season) (*[]Attendance, error) {
	attendances := &[]Attendance{}
	err := repository.connectionHandler.
		Joins("Player").
		Joins("Game").
		Find(attendances, &Attendance{Game: &games.Game{Season: season}})
	if err != nil {
		return &[]Attendance{}, err
	}

	return attendances, nil
}

func (repository *attendancesRepository) FindAttendancesForPlayer(player *Player) (*[]Attendance, error) {
	attendances := &[]Attendance{}
	err := repository.connectionHandler.
		Joins("Player").
		Preload("Game").
		Find(attendances, &Attendance{PlayerID: player.ID})
	if err != nil {
		return &[]Attendance{}, err
	}

	return attendances, nil
}

func (repository *attendancesRepository) FindFellowAttendancesForPlayer(player *Player) (*[]Attendance, error) {
	fellowWinnerAttendances, err := repository.findFellowAttendancesForPlayerByWin(player, true)
	if err != nil {
		return &[]Attendance{}, err
	}

	fellowLoserAttendances, err := repository.findFellowAttendancesForPlayerByWin(player, false)
	if err != nil {
		return &[]Attendance{}, err
	}

	fellowAttendances := append(*fellowWinnerAttendances, *fellowLoserAttendances...)

	return &fellowAttendances, nil
}

func (repository *attendancesRepository) findFellowAttendancesForPlayerByWin(player *Player, win bool) (*[]Attendance, error) {
	attendances := &[]Attendance{}
	err := repository.connectionHandler.
		Where("win = ? AND player_id = ?", win, player.ID).
		Find(attendances)
	if err != nil {
		return &[]Attendance{}, err
	}

	fellowAttendances := []Attendance{}
	err = repository.connectionHandler.
		Preload("Player").
		Where("win = ? AND game_id IN ? AND player_id != ?", win, getGameIdsFromAttendances(attendances), player.ID).
		Find(&fellowAttendances)
	if err != nil {
		return &[]Attendance{}, err
	}

	return &fellowAttendances, nil
}

func getGameIdsFromAttendances(attendances *[]Attendance) []uint {
	gameIds := make([]uint, len(*attendances))
	for i, attendance := range *attendances {
		gameIds[i] = attendance.GameID
	}

	return gameIds
}

func (repository *attendancesRepository) Create(attendances *[]Attendance) error {
	return repository.connectionHandler.Create(attendances)
}

func (repository *attendancesRepository) AutoMigrate() {
	repository.connectionHandler.AutoMigrate(&Attendance{})
}
