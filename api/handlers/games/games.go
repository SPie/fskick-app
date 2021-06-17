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
		season, err := getSeason(gamesManager, c.Param("season"))
		if err != nil {
			c.Error(err)
			return
		}

		playerStats, err := playersManager.GetPlayerStats(
			season,
			playersManager.GetSortFunction(c.DefaultQuery("sort", players.SortByPointsRatio)),
		)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, gin.H{"season": season, "playerStats": playerStats})
	}
}

func getSeason(gamesManager g.Manager, seasonUuid string) (g.Season, error) {
	if seasonUuid == "" {
		return gamesManager.ActiveSeason()
	}

	return gamesManager.GetSeasonByUuid(seasonUuid)
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
