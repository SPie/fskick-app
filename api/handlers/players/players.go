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

		playerName := c.Query("name")
		if playerName == "" {
			c.JSON(200, gin.H{"playerStats": playerStats})
			return
		}

		foundPlayers, err := playersManager.SearchPlayers(playerName)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, gin.H{"playerStats": filterPlayerStats(playerStats, foundPlayers)})
	}
}

func filterPlayerStats(playerStats *[]p.PlayerStats, foundPlayers *[]p.Player) *[]p.PlayerStats {
	playerStatsToShow := []p.PlayerStats{}
	for _, stats := range *playerStats {
		for _, player := range *foundPlayers {
			if player.UUID == stats.UUID {
				playerStatsToShow = append(playerStatsToShow, stats)
				break
			}
		}
	}

	return &playerStatsToShow
}
