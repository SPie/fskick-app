package players

import (
	"github.com/gin-gonic/gin"
	"github.com/spie/fskick/games"
	p "github.com/spie/fskick/players"
)

func GetPlayers(playersManager p.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		playerStats, err := playersManager.GetPlayerStats(
			games.Season{},
			playersManager.GetSortFunction(c.DefaultQuery("sort", p.SortByPointsRatio)),
		)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, gin.H{"playerStats": playerStats})
	}
}
