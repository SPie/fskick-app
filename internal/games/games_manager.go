package games

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/spie/fskick/internal/players"
	"github.com/spie/fskick/internal/seasons"
)

type PlayerStats struct {
	PlayerAttendance
	PointsRatio float64
	Points      int
	WinRatio    float64
	GamesRatio  float64
	Position    int
}

type Manager struct {
	gameRepository       GamesRepository
	attendanceRepository AttendanceRepository
	seasonsManager       seasons.Manager
}

func NewManager(
	gameRepository GamesRepository,
	attendanceRepository AttendanceRepository,
	seasonsManager seasons.Manager,
) Manager {
	return Manager{
		gameRepository:       gameRepository,
		attendanceRepository: attendanceRepository,
		seasonsManager:       seasonsManager,
	}
}

func (manager Manager) CreateGame(
	playedAt time.Time,
	winners players.Team,
	losers players.Team,
) (*Game, error) {
	activeSeason, err := manager.seasonsManager.ActiveSeason()
	if err != nil {
		return &Game{}, err
	}

	if playedAt.IsZero() {
		playedAt = time.Now()
	}

	game := &Game{Season: &activeSeason, PlayedAt: playedAt}
	attendances := append(
		createAttendances(winners, true),
		createAttendances(losers, false)...,
	)

	err = manager.gameRepository.CreateGame(game, attendances)
	if err != nil {
		return &Game{}, err
	}

	return game, nil
}

func createAttendances(team players.Team, win bool) []Attendance {
	attendances := make([]Attendance, len(team))

	for i, player := range team {
		attendances[i] = Attendance{
			Win:      win,
			PlayerID: player.ID,
		}
	}

	return attendances
}

func (manager Manager) GetGamesCount() (int, error) {
	return manager.gameRepository.Count()
}

func (manager Manager) GetGamesCountForSeason(season seasons.Season) (int, error) {
	return manager.gameRepository.CountForSeason(season)
}

func (manager Manager) GetGamesCountForPlayer(player players.Player) (int, error) {
	return manager.gameRepository.CountForPlayer(player)
}

func (manager Manager) GetPlayerStatsForSeason(season seasons.Season, sort string) ([]PlayerStats, error) {
	playerAttendances, err := manager.attendanceRepository.CollectPlayerAttendancesForSeason(season)
	if err != nil {
		return nil, fmt.Errorf("get player stats: %w", err)
	}

	gamesCount, err := manager.gameRepository.CountForSeason(season)
	if err != nil {
		return nil, fmt.Errorf("get player stats: %w", err)
	}

	maxGamesCount, err := manager.gameRepository.MaxGamesForSeason(season)
	if err != nil {
		return nil, fmt.Errorf("get player stats: %w", err)
	}

	playerStats := createPlayerStats(playerAttendances, gamesCount, maxGamesCount)

	sortPlayerStats(playerStats, sort)

	return playerStats, nil
}

func (manager Manager) GetAllPlayerStats(sort string) ([]PlayerStats, error) {
	playerAttendances, err := manager.attendanceRepository.CollectAllPlayerAttendances()
	if err != nil {
		return nil, fmt.Errorf("get player stats: %w", err)
	}

	gamesCount, err := manager.gameRepository.Count()
	if err != nil {
		return nil, fmt.Errorf("get player stats: %w", err)
	}

	maxGamesCount, err := manager.gameRepository.MaxGames()
	if err != nil {
		return nil, fmt.Errorf("get player stats: %w", err)
	}

	playerStats := createPlayerStats(playerAttendances, gamesCount, maxGamesCount)

	sortPlayerStats(playerStats, sort)

	return playerStats, nil
}

func (manager Manager) GetFellowPlayerStats(player players.Player, sort string) ([]PlayerStats, error) {
	playerAttendances, err := manager.attendanceRepository.CollectFellowPlayerAttendances(player)
	if err != nil {
		return nil, fmt.Errorf("get fellow player stats: %w", err)
	}

	gamesCount, err := manager.gameRepository.CountForPlayer(player)
	if err != nil {
		return nil, fmt.Errorf("get player stats: %w", err)
	}

	maxGamesCount, err := manager.gameRepository.MaxGamesForPlayer(player)
	if err != nil {
		return nil, fmt.Errorf("get player stats: %w", err)
	}

	playerStats := createPlayerStats(playerAttendances, gamesCount, maxGamesCount)

	sortPlayerStats(playerStats, sort)

	return playerStats, nil
}

func (manager Manager) GetAttendancesForPlayer(player players.Player) ([]Attendance, error) {
	return manager.attendanceRepository.GetAttendancesForPlayer(player)
}

func createPlayerStats(playerAttendances []PlayerAttendance, gamesCount int, maxGamesCount int) []PlayerStats {
	playerStats := make([]PlayerStats, len(playerAttendances))
	for i, playerAttendance := range playerAttendances {
		stats := PlayerStats{PlayerAttendance: playerAttendance}
		stats.WinRatio = float64(stats.Wins) / float64(stats.Games)
		stats.GamesRatio = float64(stats.Games) / float64(gamesCount)
		stats.Points = stats.Wins * 3
		stats.PointsRatio = float64(stats.Points) /
			math.Max(float64(stats.Games), float64(maxGamesCount/2))
		playerStats[i] = stats
	}

	return playerStats
}

func sortPlayerStats(playerStats []PlayerStats, sortName string) {
	sortFunc, positionFunc := getSortAndPositionFunc(playerStats, sortName)

	sort.Slice(playerStats, sortFunc)

	currentValue := positionFunc((playerStats)[0])
	(playerStats)[0].Position = 1
	position := 1
	playersCount := 0
	for i, stats := range playerStats {
		playersCount++
		if positionFunc(stats) < currentValue {
			position = playersCount
			currentValue = positionFunc(stats)
		}

		(playerStats)[i].Position = position
	}
}

func getSortAndPositionFunc(
	playerStats []PlayerStats,
	sortName string,
) (func(p, q int) bool, func(p PlayerStats) float64) {
	switch sortName {
	case "wins":
		return func(p, q int) bool {
				if (playerStats)[p].Wins == (playerStats)[q].Wins {
					return (playerStats)[p].Games > (playerStats)[q].Games
				}

				return (playerStats)[p].Wins > (playerStats)[q].Wins
			},
			func(p PlayerStats) float64 {
				return float64(p.Wins)
			}
	case "games":
		return func(p, q int) bool {
				if (playerStats)[p].Games == (playerStats)[q].Games {
					return (playerStats)[p].Wins > (playerStats)[q].Wins
				}

				return (playerStats)[p].Games > (playerStats)[q].Games
			},
			func(p PlayerStats) float64 {
				return float64(p.Games)
			}
	case "winRatio":
		return func(p, q int) bool {
				if (playerStats)[p].WinRatio == (playerStats)[q].WinRatio {
					return (playerStats)[p].Games > (playerStats)[q].Games
				}

				return (playerStats)[p].WinRatio > (playerStats)[q].WinRatio
			},
			func(p PlayerStats) float64 {
				return p.WinRatio
			}
	default:
		return func(p, q int) bool {
				if (playerStats)[p].PointsRatio == (playerStats)[q].PointsRatio {
					return (playerStats)[p].Games > (playerStats)[q].Games
				}

				return (playerStats)[p].PointsRatio > (playerStats)[q].PointsRatio
			},
			func(p PlayerStats) float64 {
				return p.PointsRatio
			}
	}
}
