package players

import (
	"github.com/gin-gonic/gin"
	"github.com/spie/fskick/games"
	p "github.com/spie/fskick/players"
)

type GetPlayersRequest struct {
	Player string `uri:"player"`
}

func GetPlayers(playersManager p.PlayerStatsCalculator) gin.HandlerFunc {
	return func(c *gin.Context) {
		playerStats, err := playersManager.GetPlayersStats(games.Season{})
		if err != nil {
			c.Error(err)
			return
		}

		playersManager.GetSortFunction(c.DefaultQuery("sort", p.SortByPointsRatio))(playerStats)

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

func filterPlayersStatsForName(playersStats *[]p.PlayerStats, uuid string) *[]p.PlayerStats {
	for _, playerStats := range *playersStats {
		if playerStats.Player.UUID == uuid {
			return &[]p.PlayerStats{playerStats}
		}
	}

	return &[]p.PlayerStats{}
}
