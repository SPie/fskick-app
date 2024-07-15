package players

import (
	"github.com/gin-gonic/gin"
	"github.com/spie/fskick/internal/games"
)

type GetPlayersRequest struct {
	Player string `uri:"player"`
}

func GetPlayers(gamesManager games.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		playerStats, err := gamesManager.GetAllPlayerStats(c.DefaultQuery("sort", "pointsRatio"))
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
}

func filterPlayersStatsForName(playersStats []games.PlayerStats, uuid string) []games.PlayerStats {
	for _, playerStats := range playersStats {
		if playerStats.Player.UUID == uuid {
			return []games.PlayerStats{playerStats}
		}
	}

	return []games.PlayerStats{}
}
