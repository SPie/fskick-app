package players

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/seasons"
)

type Player struct {
	db.Model
	Name        string        `gorm:"unique;not null" json:"name"`
	Attendances *[]Attendance `json:"attendances"`
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

type PlayerStatsCalculator interface {
	GetPlayersStats(season seasons.Season) (*[]PlayerStats, error)
	GetSortFunction(sortName string) sortFunction
	GetFavoriteTeam(playerUuid string) (*[]PlayerStats, error)
}

type AttendanceCreator interface {
	GetTeamsByNames(winnerNames []string, loserNames []string) (Team, Team, error)
	CreateAttendances(game *games.Game, winners Team, losers Team) (Team, Team, error)
}

type Manager struct {
	playerRepository      PlayerRepository
	attendancesRepository AttendancesRepository
}

func NewManager(playerRepository PlayerRepository, attendancesRepository AttendancesRepository) Manager {
	return Manager{playerRepository: playerRepository, attendancesRepository: attendancesRepository}
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

func (manager Manager) CreateAttendances(game *games.Game, winners Team, losers Team) (Team, Team, error) {
	attendances := append(*createAttendances(game, winners, true), *createAttendances(game, losers, false)...)
	if len(attendances) < 1 {
		return []Player{}, []Player{}, errors.New("No attendances for game")
	}

	err := manager.attendancesRepository.Create(&attendances)
	if err != nil {
		return []Player{}, []Player{}, err
	}

	return winners, losers, err
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

func createAttendances(game *games.Game, team Team, winning bool) *[]Attendance {
	attendances := make([]Attendance, len(team))
	for i, player := range team {
		attendancePlayer := player
		attendances[i] = Attendance{Game: game, Player: &attendancePlayer, Win: winning}
	}

	return &attendances
}

func (manager Manager) GetPlayersStats(season seasons.Season) (*[]PlayerStats, error) {
	players, err := manager.getPlayersForPlayerStats(season)
	if err != nil {
		return &[]PlayerStats{}, err
	}

	playersStats, maxGames := initializePlayerStatsAndGetMaxGames(players)
	calculateAndSetPointsRatio(playersStats, maxGames)

	return playersStats, nil
}

func (manager Manager) getPlayersForPlayerStats(season seasons.Season) ([]Player, error) {
	if season.ID != 0 {
		return manager.playerRepository.FindPlayersPlayedInSeason(season)
	}

	return manager.playerRepository.AllPlayersWithAttendances()
}

func initializePlayerStatsAndGetMaxGames(players []Player) (*[]PlayerStats, int) {
	maxGames := 0
	playersStats := make([]PlayerStats, len(players))
	for i, player := range players {
		wins := 0
		for _, attendance := range *player.Attendances {
			if attendance.Win {
				wins++
			}
		}

		games := len(*player.Attendances)
		if games > maxGames {
			maxGames = games
		}

		playersStats[i] = PlayerStats{
			PlayerAttendance: PlayerAttendance{
				Player: player,
				Games:  games,
				Wins:   wins,
			},
			Points: wins * 3,
		}
	}

	return &playersStats, maxGames
}

func calculateAndSetPointsRatio(playersStats *[]PlayerStats, maxGames int) {
	minDivider := int(maxGames / 2)
	for i, playerStats := range *playersStats {
		divider := playerStats.Games
		if divider < minDivider {
			divider = minDivider
		}

		playerStats.PointsRatio = float32(playerStats.Points) / float32(divider)
		(*playersStats)[i] = playerStats
	}
}

func (manager Manager) GetFavoriteTeam(playerUuid string) (*[]PlayerStats, error) {
	player, err := manager.playerRepository.FindPlayerByUUID(playerUuid)
	if err != nil {
		return &[]PlayerStats{}, err
	}

	attendances, err := manager.attendancesRepository.FindFellowAttendancesForPlayer(player)
	if err != nil {
		return &[]PlayerStats{}, err
	}

	playerMap := map[uint]*Player{}

	for _, attendance := range attendances {
		if _, ok := playerMap[attendance.PlayerID]; !ok {
			fellowPlayer := attendance.Player
			fellowPlayer.Attendances = &[]Attendance{}
			playerMap[attendance.PlayerID] = fellowPlayer
		}

		playerAttendances := *playerMap[attendance.PlayerID].Attendances
		playerAttendances = append(playerAttendances, attendance)

		playerMap[attendance.PlayerID].Attendances = &playerAttendances
	}

	players := []Player{}
	for _, fellowPlayer := range playerMap {
		players = append(players, *fellowPlayer)
	}

	playerStats, maxGames := initializePlayerStatsAndGetMaxGames(players)
	calculateAndSetPointsRatio(playerStats, maxGames)

	return playerStats, nil
}

func (manager Manager) GetSortFunction(sortName string) sortFunction {
	switch sortName {
	case SortByWins:
		return sortByWins
	case SortByGames:
		return sortByGames
	case SortByWinRatio:
		return sortByWinRatio
	}

	return sortByPointsRatio
}

type sortFunction func(playersStats *[]PlayerStats)

var (
	SortByPointsRatio = "pointsRatio"
	SortByWins        = "wins"
	SortByWinRatio    = "winRatio"
	SortByGames       = "games"
)

var sortByPointsRatio = func(playersStats *[]PlayerStats) {
	sort.Slice(*playersStats, func(p, q int) bool {
		if (*playersStats)[p].PointsRatio > (*playersStats)[q].PointsRatio {
			return true
		}
		if (*playersStats)[p].PointsRatio < (*playersStats)[q].PointsRatio {
			return false
		}
		return (*playersStats)[p].Games > (*playersStats)[q].Games
	})

	currentPositionPointsRatio := (*playersStats)[0].PointsRatio
	(*playersStats)[0].Position = 1
	position := 1
	playersCount := 0
	for i, playerStats := range *playersStats {
		playersCount++
		if playerStats.PointsRatio < currentPositionPointsRatio {
			position = playersCount
			currentPositionPointsRatio = playerStats.PointsRatio
		}

		(*playersStats)[i].Position = position
	}
}

var sortByWins = func(playersStats *[]PlayerStats) {
	sort.Slice(*playersStats, func(p, q int) bool {
		if (*playersStats)[p].Wins > (*playersStats)[q].Wins {
			return true
		}
		if (*playersStats)[p].Wins < (*playersStats)[q].Wins {
			return false
		}
		return (*playersStats)[p].Games < (*playersStats)[q].Games
	})

	currentPositionWins := (*playersStats)[0].Wins
	(*playersStats)[0].Position = 1
	position := 1
	playersCount := 0
	for i, playerStats := range *playersStats {
		playersCount++
		if playerStats.Wins < currentPositionWins {
			position = playersCount
			currentPositionWins = playerStats.Wins
		}

		(*playersStats)[i].Position = position
	}
}

var sortByWinRatio = func(playersStats *[]PlayerStats) {
	sort.Slice(*playersStats, func(p, q int) bool {
		winRatioP := float32((*playersStats)[p].Wins) / float32((*playersStats)[p].Games)
		winRatioQ := float32((*playersStats)[q].Wins) / float32((*playersStats)[q].Games)

		if winRatioP > winRatioQ {
			return true
		}
		if winRatioP < winRatioQ {
			return false
		}
		return (*playersStats)[p].Games > (*playersStats)[q].Games
	})

	currentPositionWinRatio := float32((*playersStats)[0].Wins) / float32((*playersStats)[0].Games)
	(*playersStats)[0].Position = 1
	position := 1
	playersCount := 0
	for i, playerStats := range *playersStats {
		playersCount++
		playerWinRatio := float32(playerStats.Wins) / float32(playerStats.Games)
		if playerWinRatio < currentPositionWinRatio {
			position = playersCount
			currentPositionWinRatio = playerWinRatio
		}

		(*playersStats)[i].Position = position
	}
}

var sortByGames = func(playersStats *[]PlayerStats) {
	sort.Slice(*playersStats, func(p, q int) bool {
		if (*playersStats)[p].Games > (*playersStats)[q].Games {
			return true
		}
		if (*playersStats)[p].Games < (*playersStats)[q].Games {
			return false
		}

		return (*playersStats)[p].Wins > (*playersStats)[q].Wins
	})

	currentPositionGames := (*playersStats)[0].Games
	(*playersStats)[0].Position = 1
	position := 1
	playersCount := 0
	for i, playerStats := range *playersStats {
		playersCount++
		if playerStats.Games < currentPositionGames {
			position = playersCount
			currentPositionGames = playerStats.Games
		}

		(*playersStats)[i].Position = position
	}
}
