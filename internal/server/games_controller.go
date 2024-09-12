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

func (controller GamesController) TablePage(res http.ResponseWriter, req *http.Request) {
	seasons, err := controller.seasonsManager.GetSeasons()
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	season, err := controller.getSeason("")
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	playerStats, err := controller.gamesManager.GetPlayerStatsForSeason(season, "pointsRatio")
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	gamesCount, err := controller.gamesManager.GetGamesCountForSeason(season)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	if err = controller.views.SeasonsTable.Render(
		seasons,
		season,
		playerStats,
		gamesCount,
		req.Context(),
		res,
	); err != nil {
		handleInternalServerError(res, err)
		return
	}
}

func (controller GamesController) SeasonsTableUpdate(res http.ResponseWriter, req *http.Request) {
	sort := req.URL.Query().Get("sort")
	if sort == "" {
		sort = "pointsRatio"
	}

	season, err := controller.getSeason(req.URL.Query().Get("season"))
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	playerStats, err := controller.gamesManager.GetPlayerStatsForSeason(season, sort)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	gamesCount, err := controller.gamesManager.GetGamesCountForSeason(season)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	if err = controller.views.SeasonsTableUpdate.Render(
		playerStats,
		gamesCount,
		req.Context(),
		res,
	); err != nil {
		handleInternalServerError(res, err)
		return
	}
}

func (controller GamesController) PlayersTableUpdate(res http.ResponseWriter, req *http.Request) {
	sort := req.URL.Query().Get("sort")
	if sort == "" {
		sort = "pointsRatio"
	}

	playerStats, err := controller.gamesManager.GetAllPlayerStats(sort)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	gamesCount, err := controller.gamesManager.GetGamesCount()
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	playerUuid := req.PathValue("player")
	if playerUuid != "" {
		playerStats = filterPlayersStatsForUuid(playerStats, playerUuid)
	}

	if err = controller.views.PlayersTableUpdate.Render(
		playerStats,
		gamesCount,
		playerUuid,
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

func (controller GamesController) GetFavoriteTeam(res http.ResponseWriter, req *http.Request) {
	playerUuid := req.PathValue("player")

	player, err := controller.playersManager.GetPlayerByUUID(playerUuid)
	if errors.Is(err, players.ErrPlayerNotFound) {
		http.Error(res, fmt.Sprintf("Player %s not found", playerUuid), http.StatusUnprocessableEntity)
		return
	}
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	sort := req.URL.Query().Get("sort")
	if sort == "" {
		sort = "pointsRatio"
	}

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
		req.Context(),
		res,
	); err != nil {
		handleInternalServerError(res, err)
		return
	}
}

func (controller GamesController) GetTable(res http.ResponseWriter, req *http.Request) {
	sort := req.URL.Query().Get("sort")
	if sort == "" {
		sort = "pointsRatio"
	}

	tableRes, err := controller.getTable(req.PathValue("season"), sort)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	err = writeJsonResponse(res, tableRes)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}
}

func (controller GamesController) getTable(
	seasonUuid string,
	sort string,
) (tableResponse, error) {
	season, err := controller.getSeason(seasonUuid)
	if err != nil {
		return tableResponse{}, err
	}

	playerStats, err := controller.gamesManager.GetPlayerStatsForSeason(season, sort)
	if err != nil {
		return tableResponse{}, err
	}

	gamesCount, err := controller.gamesManager.GetGamesCountForSeason(season)
	if err != nil {
		return tableResponse{}, err
	}

	return newTableResponse(season, gamesCount, playerStats), nil
}

func (controller GamesController) getSeason(seasonUuid string) (seasons.Season, error) {
	if seasonUuid != "" {
		return controller.seasonsManager.GetSeasonByUuid(seasonUuid)
	}

	return controller.seasonsManager.ActiveSeason()
}
