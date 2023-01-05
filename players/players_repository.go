package players

import (
	"fmt"

	"github.com/spie/fskick/db"
	"github.com/spie/fskick/games"
)

type PlayerRepository interface {
	db.Repository
	Save(player *Player) error
	FindPlayerByUUID(uuid string) (*Player, error)
	FindPlayerByName(name string) (*Player, error)
	FindPlayersByNames(names []string) (*[]Player, error)
	FindPlayersPlayedInSeason(season games.Season) (*[]Player, error)
	AllPlayersWithAttendances() (*[]Player, error)
	SearchPlayers(query string) (*[]Player, error)
}

type playerRepository struct {
	connectionHandler db.ConnectionHandler
}

func NewPlayerRepository(connectionHandler db.ConnectionHandler) PlayerRepository {
	return &playerRepository{connectionHandler: connectionHandler}
}

func (repository *playerRepository) AutoMigrate() {
	repository.connectionHandler.AutoMigrate(&Player{})
}

func (repository *playerRepository) Save(player *Player) error {
	if player.ID == 0 {
		return repository.connectionHandler.Create(player)
	}

	return repository.connectionHandler.Save(player)
}

func (repository *playerRepository) FindPlayerByUUID(uuid string) (*Player, error) {
	player := &Player{}
	err := repository.connectionHandler.FindOne(player, &Player{Model: db.Model{UUID: uuid}})
	if err != nil {
		return &Player{}, err
	}

	return player, nil
}

func (repository *playerRepository) FindPlayerByName(name string) (*Player, error) {
	player := &Player{}
	err := repository.connectionHandler.FindOne(player, &Player{Name: name})
	if err != nil {
		return &Player{}, err
	}

	return player, nil
}

func (repository *playerRepository) FindPlayersByNames(names []string) (*[]Player, error) {
	players := &[]Player{}

	err := repository.connectionHandler.Find(players, "name IN ?", names)
	if err != nil {
		return &[]Player{}, err
	}

	return players, nil
}

func (repository *playerRepository) FindPlayersPlayedInSeason(season games.Season) (*[]Player, error) {
	players := &[]Player{}

	err := repository.connectionHandler.Preload("Attendances.Game.Season").Find(players)
	if err != nil {
		return &[]Player{}, err
	}

	return getPlayersForSeason(season, players), nil
}

func getPlayersForSeason(season games.Season, players *[]Player) *[]Player {
	playersPlayed := []Player{}
	for _, player := range *players {
		attendancesInSeason := getAttendancesForSeason(season, player.Attendances)
		if len(*attendancesInSeason) > 0 {
			player.Attendances = attendancesInSeason
			playersPlayed = append(playersPlayed, player)
		}
	}

	return &playersPlayed
}

func getAttendancesForSeason(season games.Season, attendances *[]Attendance) *[]Attendance {
	attendancesInSeason := []Attendance{}
	for _, attendance := range *attendances {
		if attendance.Game.Season.ID == season.ID {
			attendancesInSeason = append(attendancesInSeason, attendance)
		}
	}

	return &attendancesInSeason
}

func (repository *playerRepository) AllPlayersWithAttendances() (*[]Player, error) {
	players := &[]Player{}
	err := repository.connectionHandler.Preload("Attendances.Game").Find(players)
	if err != nil {
		return &[]Player{}, err
	}

	return players, nil
}

func (repository *playerRepository) SearchPlayers(query string) (*[]Player, error) {
	players := &[]Player{}

	err := repository.connectionHandler.Find(players, "name LIKE ?", fmt.Sprintf("%%%s%%", query))
	if err != nil {
		return &[]Player{}, err
	}
	fmt.Println(len(*players))

	return players, nil
}
