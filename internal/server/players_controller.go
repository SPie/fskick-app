package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
	"github.com/spie/fskick/internal/views"
)

type PlayersViews struct {
        PlayersTable views.PlayersTable
        PlayerInfo views.PlayerInfo
}

func NewPlayersViews() PlayersViews {
        return PlayersViews{}
}

type PlayersController struct {
        playersManager players.Manager
        gamesManager games.Manager
        views PlayersViews
}

func NewPlayersController(
        playersManager players.Manager,
        gamesManager games.Manager,
        views PlayersViews,
) PlayersController {
        return PlayersController{
                playersManager: playersManager,
                gamesManager: gamesManager,
                views: views,
        }
}

func (controller PlayersController) GetPlayersTable(res http.ResponseWriter, req *http.Request) {
        playerStats, err := controller.gamesManager.GetAllPlayerStats("pointsRatio")
        if err != nil {
                handleInternalServerError(res, err)
                return
        }

        gamesCount, err := controller.gamesManager.GetGamesCount()
        if err != nil {
                handleInternalServerError(res, err)
                return
        }

        err = controller.views.PlayersTable.Render(playerStats, gamesCount, req.Context(), res)
        if err != nil {
                handleInternalServerError(res, err)
                return
        }
}

func (controller PlayersController) GetPlayerInfo(res http.ResponseWriter, req *http.Request) {
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

        playerStats, err := controller.gamesManager.GetAllPlayerStats("pointsRatio")
        if err != nil {
                handleInternalServerError(res, err)
                return
        }

        playerStats = filterPlayersStatsForUuid(playerStats, player.UUID)
        if len(playerStats) != 1 {
                http.Error(res, fmt.Sprintf("Player %s not found", playerUuid), http.StatusUnprocessableEntity)
                return
        }

        gamesCount, err := controller.gamesManager.GetGamesCount()
        if err != nil {
                handleInternalServerError(res, err)
                return
        }

        attendances, err := controller.gamesManager.GetAttendancesForPlayer(player)
        if err != nil {
                handleInternalServerError(res, err)
                return
        }

        teamPlayerStats, err := controller.gamesManager.GetFellowPlayerStats(player, "pointsRatio")
        if err != nil {
                handleInternalServerError(res, err)
                return
        }

        err = controller.views.PlayerInfo.Render(
                playerStats[0],
                gamesCount,
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

func (controller PlayersController) GetPlayers(res http.ResponseWriter, req *http.Request) {
        sort := req.URL.Query().Get("sort")
        if sort == "" {
                sort = "pointsRatio"
        }

        playerStats, err := controller.gamesManager.GetAllPlayerStats(sort)
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

func filterPlayersStatsForUuid(playersStats []games.PlayerStats, uuid string) []games.PlayerStats {
        for _, playerStats := range playersStats {
                if playerStats.Player.UUID == uuid {
                        return []games.PlayerStats{playerStats}
                }
        }

        return []games.PlayerStats{}
}

func (controller PlayersController) GetFavoriteTeam(res http.ResponseWriter, req *http.Request) {
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

        err = writeJsonResponse(
                res,
                map[string][]playerStatsResponse{"playerStats": newPlayerStatsResponsesFromPlayerStats(teamPlayerStats)},
        )
        if err != nil {
                handleInternalServerError(res, err)
                return
        }
}
