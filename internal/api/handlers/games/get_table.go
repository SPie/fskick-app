package games

import (
	"github.com/gin-gonic/gin"
	g "github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
)

func GetTable(playersManager players.Manager, gamesManager g.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		playerStats, season, err := getTable(playersManager, gamesManager, "", c.DefaultQuery("sort", players.SortByPointsRatio))
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, gin.H{
			"season":      season,
			"playerStats": playerStats,
		})
	}
}

type getTableForSeasonRequest struct {
	Season string `uri:"season" binding:"required"`
}

func GetTableForSeason(playersManager players.Manager, gamesManager g.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request getTableForSeasonRequest
		if err := c.ShouldBindUri(&request); err != nil {
			c.Error(err)
			return
		}

		playerStats, season, err := getTable(playersManager, gamesManager, request.Season, c.DefaultQuery("sort", players.SortByPointsRatio))
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, gin.H{
			"season":      season,
			"playerStats": playerStats,
		})
	}
}

type getTableForPlayerRequest struct {
	Season string `uri:"season" binding:"required"`
	Player string `uri:"player" binding:"required"`
}

func GetTableForPlayer(playersManager players.Manager, gamesManager g.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request getTableForPlayerRequest
		err := c.ShouldBindUri(&request)
		if err != nil {
			c.Error(err)
			return
		}

		playersStats, season, err := getTable(playersManager, gamesManager, request.Season, c.DefaultQuery("sort", players.SortByPointsRatio))
		if err != nil {
			c.Error(err)
			return
		}

		for _, playerStats := range *playersStats {
			if playerStats.Player.Name == request.Player {
				c.JSON(200, gin.H{
					"season":      season,
					"playerStats": playerStats,
				})
				return
			}
		}

		c.JSON(200, gin.H{
			"season":      season,
			"playerStats": players.PlayerStats{},
		})
	}
}

func getTable(
	playersManager players.Manager,
	gamesManager g.Manager,
	seasonUuid string,
	sort string,
) (*[]players.PlayerStats, g.Season, error) {
	season, err := getSeason(gamesManager, seasonUuid)
	if err != nil {
		return &[]players.PlayerStats{}, g.Season{}, err
	}

	playerStats, err := playersManager.GetPlayersStats(season)
	if err != nil {
		return &[]players.PlayerStats{}, g.Season{}, err
	}

	playersManager.GetSortFunction(sort)(playerStats)

	return playerStats, season, nil
}

func getSeason(gamesManager g.Manager, seasonUuid string) (g.Season, error) {
	if seasonUuid != "" {
		return gamesManager.GetSeasonByUuid(seasonUuid)
	}

	return gamesManager.ActiveSeason()
}
