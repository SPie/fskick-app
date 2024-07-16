package api

import (
	"github.com/gin-gonic/gin"
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

type GetPlayersRequest struct {
	Player string `uri:"player"`
}

func (controller PlayersController) GetPlayers(c *gin.Context) {
	playerStats, err := controller.gamesManager.GetAllPlayerStats(c.DefaultQuery("sort", "pointsRatio"))
	if err != nil {
		c.Error(err)
		return
	}

	var request GetPlayersRequest
	err = c.ShouldBindUri(&request)
	if err != nil {
		c.Error(err)
		return
	}

	if request.Player == "" {
		c.JSON(200, gin.H{"playerStats": playerStats})
		return
	}

	c.JSON(200, gin.H{"playerStats": filterPlayersStatsForName(playerStats, request.Player)})
	return
}

func filterPlayersStatsForName(playersStats []games.PlayerStats, uuid string) []games.PlayerStats {
	for _, playerStats := range playersStats {
		if playerStats.Player.UUID == uuid {
			return []games.PlayerStats{playerStats}
		}
	}

	return []games.PlayerStats{}
}

type GetFavoriteTeamRequest struct {
	Player string `uri:"player" binding:"required"`
}

func (controller PlayersController) GetFavoriteTeam(c *gin.Context) {
	var request GetPlayersRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		c.Error(err)
		return
	}

	player, err := controller.playersManager.GetPlayerByUUID(request.Player)
	if err != nil {
		c.Error(err)
		return
	}

	teamPlayerStats, err := controller.gamesManager.GetFellowPlayerStats(
		player,
		c.DefaultQuery("sort", "getPointsRatio"),
	)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"playerStats": teamPlayerStats})
	return
}
