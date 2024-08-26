package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
)

type PlayersController struct {
        playersManager players.Manager
        gamesManager games.Manager
}

func NewPlayersController(
        playersManager players.Manager,
        gamesManager games.Manager,
) PlayersController {
        return PlayersController{
                playersManager: playersManager,
                gamesManager: gamesManager,
        }
}

func (controller PlayersController) GetPlayers(res http.ResponseWriter, req *http.Request) {
        sort := req.URL.Query().Get("sort")
        if sort == "" {
                sort = "pointsRation"
        }

        playerStats, err := controller.gamesManager.GetAllPlayerStats(sort)
        if err != nil {
                handleInternalServerError(res, err)
                return
        }

        playerName := req.PathValue("player")
        if playerName != "" {
                playerStats = filterPlayersStatsForName(playerStats, playerName)
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

func filterPlayersStatsForName(playersStats []games.PlayerStats, uuid string) []games.PlayerStats {
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
