package players

import (
	"errors"
	"fmt"
	"strings"
)

type Team []Player

type playerRepository interface {
	CreatePlayer(player *Player) error
	FindPlayerByUUID(uuid string) (Player, error)
	FindPlayerByName(name string) (Player, error)
	FindPlayersByNames(names []string) ([]Player, error)
}

type Manager struct {
	playerRepository playerRepository
}

func NewManager(playerRepository playerRepository) Manager {
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

func (manager Manager) GetPlayerByName(name string) (Player, error) {
	player, err := manager.playerRepository.FindPlayerByName(name)
	if err != nil {
		return Player{}, fmt.Errorf("get player by name: %w", err)
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
