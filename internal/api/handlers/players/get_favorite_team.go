package players

import (
	"github.com/gin-gonic/gin"
	"github.com/spie/fskick/internal/games"
	p "github.com/spie/fskick/internal/players"
)

type GetFavoriteTeamRequest struct {
	Player string `uri:"player" binding:"required"`
}

func GetFavoriteTeam(playersManager p.Manager, gamesManager games.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request GetPlayersRequest
		err := c.ShouldBindUri(&request)
		if err != nil {
			c.Error(err)
			return
		}

		player, err := playersManager.GetPlayerByUUID(request.Player)
		if err != nil {
			c.Error(err)
			return
		}

		teamPlayerStats, err := gamesManager.GetFellowPlayerStats(
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
}
