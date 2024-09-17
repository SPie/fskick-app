package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
	"github.com/spie/fskick/internal/seasons"
	"github.com/spie/fskick/internal/views"
)

type GamesViews struct {
	SeasonsTable       views.SeasonsTable
	PlayersTable       views.PlayersTable
	PlayerInfo         views.PlayerInfo
	SeasonsTableUpdate views.SeasonsTableUpdate
	PlayersTableUpdate views.PlayersTableUpdate
	FavoriteTeamUpdate views.FavoriteTeamUpdate
}

func NewGamesViews() GamesViews {
	return GamesViews{}
}

type GamesController struct {
	gamesManager   games.Manager
	seasonsManager seasons.Manager
	playersManager players.Manager
	views          GamesViews
}

func NewGamesController(
	gamesManager games.Manager,
	seasonsManager seasons.Manager,
	playersManager players.Manager,
	gamesViews GamesViews,
) GamesController {
	return GamesController{
		gamesManager:   gamesManager,
		seasonsManager: seasonsManager,
		playersManager: playersManager,
		views:          gamesViews,
	}
}

func (controller GamesController) SeasonsTable(res http.ResponseWriter, req *http.Request) {
	seasons, err := controller.seasonsManager.GetSeasons()
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	seasonTableData, err := controller.getSeasonsTableData("", getSort(req))
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	if err = controller.views.SeasonsTable.Render(
		seasons,
		seasonTableData.season,
		seasonTableData.playerStats,
		seasonTableData.gamesCount,
		req.Context(),
		res,
	); err != nil {
		handleInternalServerError(res, err)
		return
	}
}

func (controller GamesController) SeasonsTableUpdate(res http.ResponseWriter, req *http.Request) {
	sort := getSort(req)

	seasonTableData, err := controller.getSeasonsTableData(req.URL.Query().Get("season"), sort)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	if err = controller.views.SeasonsTableUpdate.Render(
		seasonTableData.playerStats,
		seasonTableData.gamesCount,
		sort,
		req.Context(),
		res,
	); err != nil {
		handleInternalServerError(res, err)
		return
	}
}

func (controller GamesController) PlayersTable(res http.ResponseWriter, req *http.Request) {
	playersTableData, err := controller.getPlayersTableData("", getSort(req))
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	err = controller.views.PlayersTable.Render(
		playersTableData.playerStats,
		playersTableData.gamesCount,
		req.Context(),
		res,
	)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}
}

func (controller GamesController) PlayersTableUpdate(res http.ResponseWriter, req *http.Request) {
	playerUuid := req.PathValue("player")
	sort := getSort(req)

	playersTableData, err := controller.getPlayersTableData(playerUuid, sort)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	if err = controller.views.PlayersTableUpdate.Render(
		playersTableData.playerStats,
		playersTableData.gamesCount,
		playerUuid,
		sort,
		req.Context(),
		res,
	); err != nil {
		handleInternalServerError(res, err)
		return
	}
}

func (controller GamesController) PlayerInfo(res http.ResponseWriter, req *http.Request) {
	playerUuid := req.PathValue("player")

	sort := getSort(req)

	playersTableData, err := controller.getPlayersTableData(playerUuid, sort)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}
	if len(playersTableData.playerStats) != 1 {
		http.Error(res, fmt.Sprintf("Player %s not found", playerUuid), http.StatusUnprocessableEntity)
		return
	}

	attendances, err := controller.gamesManager.GetAttendancesForPlayer(playersTableData.playerStats[0].Player)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	teamPlayerStats, err := controller.gamesManager.GetFellowPlayerStats(
		playersTableData.playerStats[0].Player,
		sort,
	)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	err = controller.views.PlayerInfo.Render(
		playersTableData.playerStats[0],
		playersTableData.gamesCount,
		attendances,
		teamPlayerStats,
		req.Context(),
		res,
	)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}
}

func (controller GamesController) FavoriteTeamUpdate(res http.ResponseWriter, req *http.Request) {
	playerUuid := req.PathValue("player")

	player, err := controller.getPlayer(playerUuid)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	sort := getSort(req)

	teamPlayerStats, err := controller.gamesManager.GetFellowPlayerStats(player, sort)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	gamesCount, err := controller.gamesManager.GetGamesCountForPlayer(player)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	if err = controller.views.FavoriteTeamUpdate.Render(
		teamPlayerStats,
		gamesCount,
		playerUuid,
		sort,
		req.Context(),
		res,
	); err != nil {
		handleInternalServerError(res, err)
		return
	}
}

func (controller GamesController) GetGamesCount(res http.ResponseWriter, _ *http.Request) {
	gamesCount, err := controller.gamesManager.GetGamesCount()
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	err = writeJsonResponse(res, map[string]int{"gamesCount": gamesCount})
	if err != nil {
		handleInternalServerError(res, err)
		return
	}
}

func (controller GamesController) GetSeasonsTable(res http.ResponseWriter, req *http.Request) {
	seasonTableData, err := controller.getSeasonsTableData(req.PathValue("season"), getSort(req))
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	err = writeJsonResponse(
		res,
		newTableResponse(seasonTableData.season, seasonTableData.gamesCount, seasonTableData.playerStats),
	)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}
}

func (controller GamesController) GetPlayers(res http.ResponseWriter, req *http.Request) {
	playerStats, err := controller.gamesManager.GetAllPlayerStats(getSort(req))
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	playerUuid := req.PathValue("player")
	if playerUuid != "" {
		playerStats = filterPlayersStatsForUuid(playerStats, playerUuid)
	}

	err = writeJsonResponse(
		res,
		map[string][]playerStatsResponse{"playerStats": newPlayerStatsResponsesFromPlayerStats(playerStats)},
	)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}
}

func (controller GamesController) GetFavoriteTeam(res http.ResponseWriter, req *http.Request) {
	playerUuid := req.PathValue("player")

	player, err := controller.getPlayer(playerUuid)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	teamPlayerStats, err := controller.gamesManager.GetFellowPlayerStats(player, getSort(req))
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	err = writeJsonResponse(
		res,
		map[string][]playerStatsResponse{"playerStats": newPlayerStatsResponsesFromPlayerStats(teamPlayerStats)},
	)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}
}

type seasonTableData struct {
	playerTableData
	season seasons.Season
}

func (controller GamesController) getSeasonsTableData(
	seasonUuid string,
	sort string,
) (seasonTableData, error) {
	season, err := controller.getSeason(seasonUuid)
	if err != nil {
		return seasonTableData{}, err
	}

	playerStats, err := controller.gamesManager.GetPlayerStatsForSeason(season, sort)
	if err != nil {
		return seasonTableData{}, err
	}

	gamesCount, err := controller.gamesManager.GetGamesCountForSeason(season)
	if err != nil {
		return seasonTableData{}, err
	}

	return seasonTableData{
		season: season,
		playerTableData: playerTableData{
			playerStats: playerStats,
			gamesCount:  gamesCount,
		},
	}, nil
}

type playerTableData struct {
	playerStats []games.PlayerStats
	gamesCount  int
}

func (controller GamesController) getPlayersTableData(playerUuid string, sort string) (playerTableData, error) {
	playerStats, err := controller.gamesManager.GetAllPlayerStats(sort)
	if err != nil {
		return playerTableData{}, err
	}
	if playerUuid != "" {
		playerStats = filterPlayersStatsForUuid(playerStats, playerUuid)
	}

	gamesCount, err := controller.gamesManager.GetGamesCount()
	if err != nil {
		return playerTableData{}, err
	}

	return playerTableData{playerStats: playerStats, gamesCount: gamesCount}, nil
}

func (controller GamesController) getSeason(seasonUuid string) (seasons.Season, error) {
	if seasonUuid != "" {
		return controller.seasonsManager.GetSeasonByUuid(seasonUuid)
	}

	return controller.seasonsManager.ActiveSeason()
}

func (controller GamesController) getPlayer(playerUuid string) (players.Player, error) {
	player, err := controller.playersManager.GetPlayerByUUID(playerUuid)
	if errors.Is(err, players.ErrPlayerNotFound) {
		return players.Player{}, errors.New(fmt.Sprintf("Player %s not found", playerUuid))
	}
	if err != nil {
		return players.Player{}, err
	}

	return player, nil
}

func filterPlayersStatsForUuid(playersStats []games.PlayerStats, uuid string) []games.PlayerStats {
	for _, playerStats := range playersStats {
		if playerStats.Player.UUID == uuid {
			return []games.PlayerStats{playerStats}
		}
	}

	return []games.PlayerStats{}
}

func getSort(req *http.Request) string {
	sort := req.URL.Query().Get("sort")
	if sort == "" {
		sort = "pointsRatio"
	}

	return sort
}
