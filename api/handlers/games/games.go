package games

import (
	"github.com/gin-gonic/gin"
	g "github.com/spie/fskick/games"
	"github.com/spie/fskick/players"
)

func GetSeasons(gamesManager g.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		seasons, err := gamesManager.GetSeasons()
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, gin.H{"seasons": seasons})
	}
}

func GetTable(playersManager players.Manager, gamesManager g.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		season, err := gamesManager.ActiveSeason()
		if err != nil {
			c.Error(err)
			return
		}

		playerStats, err := playersManager.GetPlayerStats(season, playersManager.GetSortFunction(""))
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, gin.H{"season": season, "playerStats": playerStats})
	}
}
