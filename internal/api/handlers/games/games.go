package games

import (
	"github.com/gin-gonic/gin"
	g "github.com/spie/fskick/internal/games"
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

func GetGamesCount(gamesManager g.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		gamesCount, err := gamesManager.GetGamesCount()
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, gin.H{"gamesCount": gamesCount})
	}
}
