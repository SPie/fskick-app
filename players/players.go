package players

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/spie/fskick/db"
	"github.com/spie/fskick/games"
	"gorm.io/gorm"
)

type Player struct {
	db.Model
	Name        string        `gorm:"unique;not null" json:"name"`
	Attendances *[]Attendance `json:"attendances"`
}

type PlayerStats struct {
	Player
	Position    int     `json:"position"`
	PointsRatio float32 `json:"pointsRatio"`
	Points      int     `json:"points"`
	Wins        int     `json:"wins"`
	Games       int     `json:"games"`
}

type Team *[]Player

type Manager interface {
	CreatePlayer(name string) (*Player, error)
	CreateAttendances(game *games.Game, winnerNames []string, loserNames []string) (Team, Team, error)
	GetPlayerStats(season games.Season, sortFunction sortFunction) (*[]PlayerStats, error)
	GetSortFunction(sortName string) sortFunction
	SearchPlayers(query string) (*[]Player, error)
}

type manager struct {
	playerRepository      PlayerRepository
	attendancesRepository AttendancesRepository
}

func NewManager(playerRepository PlayerRepository, attendancesRepository AttendancesRepository) Manager {
	return manager{playerRepository: playerRepository, attendancesRepository: attendancesRepository}
}

func (manager manager) CreatePlayer(name string) (*Player, error) {
	_, err := manager.playerRepository.FindPlayerByName(name)
	if err == nil {
		return &Player{}, errors.New(fmt.Sprintf("Player with name %s exists", name))
	}
	if err != gorm.ErrRecordNotFound {
		return &Player{}, err
	}

	player := &Player{Name: name}

	err = manager.playerRepository.Save(player)
	if err != nil {
		return &Player{}, err
	}

	return player, nil
}

func (manager manager) CreateAttendances(game *games.Game, winnerNames []string, loserNames []string) (Team, Team, error) {
	winners, losers, err := manager.getTeams(winnerNames, loserNames)
	if err != nil {
		return &[]Player{}, &[]Player{}, err
	}

	attendances := append(*createAttendances(game, winners, true), *createAttendances(game, losers, false)...)
	if len(attendances) < 1 {
		return &[]Player{}, &[]Player{}, errors.New("No attendances for game")
	}

	err = manager.attendancesRepository.Create(&attendances)
	if err != nil {
		return &[]Player{}, &[]Player{}, err
	}

	return winners, losers, err
}

func (manager manager) getTeams(winnerNames []string, loserNames []string) (Team, Team, error) {
	winners, err := manager.getTeam(winnerNames)
	if err != nil {
		return &[]Player{}, &[]Player{}, err
	}

	losers, err := manager.getTeam(loserNames)
	if err != nil {
		return &[]Player{}, &[]Player{}, err
	}

	return winners, losers, nil
}

func (manager manager) getTeam(names []string) (Team, error) {
	if len(names) < 1 {
		return &[]Player{}, nil
	}

	players, err := manager.playerRepository.FindPlayersByNames(names)
	if err != nil {
		return &[]Player{}, err
	}

	if len(*players) != len(names) {
		return &[]Player{}, errors.New(fmt.Sprintf("Players not found: %s", strings.Join(getIncorrectPlayerNames(names, players), ",")))
	}

	return players, nil
}

func getIncorrectPlayerNames(names []string, players *[]Player) []string {
	incorrectNames := []string{}
	playerNames := map[string]string{}
	for _, player := range *players {
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
	attendances := make([]Attendance, len(*team))
	for i, player := range *team {
		attendancePlayer := player
		attendances[i] = Attendance{Game: game, Player: &attendancePlayer, Win: winning}
	}

	return &attendances
}

func (manager manager) GetPlayerStats(season games.Season, sortFunction sortFunction) (*[]PlayerStats, error) {
	players, err := manager.getPlayersForPlayerStats(season)
	if err != nil {
		return &[]PlayerStats{}, err
	}

	playersStats, maxGames := initializePlayerStatsAndGetMaxGames(players)
	calculateAndSetPointsRatio(playersStats, maxGames)
	sortFunction(playersStats)

	return playersStats, nil
}

func (manager manager) getPlayersForPlayerStats(season games.Season) (*[]Player, error) {
	if season.ID != 0 {
		return manager.playerRepository.FindPlayersPlayedInSeason(season)
	}

	return manager.playerRepository.AllPlayersWithAttendances()
}

func initializePlayerStatsAndGetMaxGames(players *[]Player) (*[]PlayerStats, int) {
	maxGames := 0
	playersStats := make([]PlayerStats, len(*players))
	for i, player := range *players {
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
			Player: player,
			Games:  games,
			Wins:   wins,
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

func (manager manager) GetSortFunction(sortName string) sortFunction {
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

func (manager manager) SearchPlayers(query string) (*[]Player, error) {
	return manager.playerRepository.SearchPlayers(query)
}

type PlayerRepository interface {
	db.Repository
	Save(player *Player) error
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
	err := repository.connectionHandler.Preload("Attendances").Find(players)
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
	err := repository.connectionHandler.Joins("Player").Joins("Game").Find(attendances, &Attendance{Game: &games.Game{Season: season}})
	if err != nil {
		return &[]Attendance{}, err
	}

	return attendances, nil
}

func (repository *attendancesRepository) Create(attendances *[]Attendance) error {
	return repository.connectionHandler.Create(attendances)
}

func (repository *attendancesRepository) AutoMigrate() {
	repository.connectionHandler.AutoMigrate(&Attendance{})
}
