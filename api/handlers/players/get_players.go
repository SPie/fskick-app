package players

import (
	"github.com/gin-gonic/gin"
	"github.com/spie/fskick/games"
	p "github.com/spie/fskick/players"
)

func GetPlayers(playersManager p.PlayerStatsCalculator) gin.HandlerFunc {
	return func(c *gin.Context) {
		playerStats, err := playersManager.GetPlayersStats(games.Season{})
		if err != nil {
			c.Error(err)
			return
		}

		playersManager.GetSortFunction(c.DefaultQuery("sort", p.SortByPointsRatio))(playerStats)

		name := c.DefaultQuery("name", "")
		if name == "" {
			c.JSON(200, gin.H{"playerStats": playerStats})
			return
		}

		c.JSON(200, gin.H{"playerStats": filterPlayersStatsForName(playerStats, name)})
		return
	}
}

func filterPlayersStatsForName(playersStats *[]p.PlayerStats, name string) *[]p.PlayerStats {
	for _, playerStats := range *playersStats {
		if playerStats.Player.Name == name {
			return &[]p.PlayerStats{playerStats}
		}
	}

	return &[]p.PlayerStats{}
}
