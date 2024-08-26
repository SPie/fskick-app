package server

import (
	"net/http"

	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
	"github.com/spie/fskick/internal/seasons"
)

type GamesController struct {
    gamesManager games.Manager
    seasonsManager seasons.Manager
    playersManager players.Manager
}

func NewGamesController(
    gamesManager games.Manager,
    seasonsManager seasons.Manager,
    playersManager players.Manager,
) GamesController {
    return GamesController{
        gamesManager: gamesManager,
        seasonsManager: seasonsManager,
        playersManager: playersManager,
    }
}

func (controller GamesController) TablePage(res http.ResponseWriter, req *http.Request) {

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

func (controller GamesController) GetTable(res http.ResponseWriter, req *http.Request) {
    sort := req.URL.Query().Get("sort")
    if sort == "" {
        sort = "pointsRatio"
    }

    tableRes, err := controller.getTable("", sort)
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
