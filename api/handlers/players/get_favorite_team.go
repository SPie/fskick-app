package players

import (
	"github.com/gin-gonic/gin"
	p "github.com/spie/fskick/players"
)

type GetFavoriteTeamRequest struct {
	Player string `uri:"player" binding:"required"`
}

func GetFavoriteTeam(playersManager p.PlayerStatsCalculator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request GetPlayersRequest
		err := c.ShouldBindUri(&request)
		if err != nil {
			c.Error(err)
			return
		}

		teamPlayerStats, err := playersManager.GetFavoriteTeam(request.Player)
		if err != nil {
			c.Error(err)
			return
		}

		playersManager.GetSortFunction(c.DefaultQuery("sort", p.SortByPointsRatio))(teamPlayerStats)

		c.JSON(200, gin.H{"playerStats": teamPlayerStats})
		return
	}
}
