package games

import (
	"github.com/gin-gonic/gin"
	"github.com/spie/fskick/internal/games"
)

func GetGamesCount(gamesManager games.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		gamesCount, err := gamesManager.GetGamesCount()
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, gin.H{"gamesCount": gamesCount})
	}
}
