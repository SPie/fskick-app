package players

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spie/fskick/internal/db"
)

type Player struct {
	db.Model
	Name string `gorm:"unique;not null" json:"name"`
}

type PlayerAttendance struct {
	Player
	Wins  int `json:"wins"`
	Games int `json:"games"`
}

type PlayerStats struct {
	PlayerAttendance
	Position    int     `json:"position"`
	PointsRatio float32 `json:"pointsRatio"`
	Points      int     `json:"points"`
}

type Team []Player

type PlayerCreator interface {
	CreatePlayer(name string) (Player, error)
}

type AttendanceCreator interface {
	GetTeamsByNames(winnerNames []string, loserNames []string) (Team, Team, error)
}

type Manager struct {
	playerRepository      PlayerRepository
}

func NewManager(playerRepository PlayerRepository) Manager {
	return Manager{playerRepository: playerRepository}
}

func (manager Manager) CreatePlayer(name string) (Player, error) {
	_, err := manager.playerRepository.FindPlayerByName(name)
	if err == nil {
		return Player{}, errors.New(fmt.Sprintf("Player with name %s exists", name))
	}
	if !errors.Is(err, ErrPlayerNotFound) {
		return Player{}, fmt.Errorf("Check for player with name in CreatePlayer: %w", err)
	}

	player := Player{Name: name}

	err = manager.playerRepository.CreatePlayer(&player)
	if err != nil {
		return Player{}, err
	}

	return player, nil
}

func (manager Manager) GetPlayerByUUID(uuid string) (Player, error) {
	player, err := manager.playerRepository.FindPlayerByUUID(uuid)
	if err != nil {
		return Player{}, fmt.Errorf("get player by uuid: %w", err)
	}

	return player, nil
}

func (manager Manager) GetTeamsByNames(winnerNames []string, loserNames []string) (Team, Team, error) {
	winners, err := manager.getTeamByNames(winnerNames)
	if err != nil {
		return []Player{}, []Player{}, err
	}

	losers, err := manager.getTeamByNames(loserNames)
	if err != nil {
		return []Player{}, []Player{}, err
	}

	return winners, losers, nil
}

func (manager Manager) getTeamByNames(names []string) (Team, error) {
	if len(names) < 1 {
		return []Player{}, nil
	}

	players, err := manager.playerRepository.FindPlayersByNames(names)
	if err != nil {
		return []Player{}, err
	}

	if len(players) != len(names) {
		return []Player{}, errors.New(fmt.Sprintf("Players not found: %s", strings.Join(getIncorrectPlayerNames(names, players), ",")))
	}

	return players, nil
}

func getIncorrectPlayerNames(names []string, players []Player) []string {
	incorrectNames := []string{}
	playerNames := map[string]string{}
	for _, player := range players {
		playerNames[player.Name] = player.Name
	}

	for _, name := range names {
		if _, ok := playerNames[name]; !ok {
			incorrectNames = append(incorrectNames, name)
		}
	}

	return incorrectNames
}
