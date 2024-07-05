package games

import (
	"time"

	"github.com/gin-gonic/gin"
	g "github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
)

type seasonResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Active bool `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	GamesCount int `json:"gamesCount"`
}

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
) (*[]players.PlayerStats, seasonResponse, error) {
	season, err := getSeason(gamesManager, seasonUuid)
	if err != nil {
		return &[]players.PlayerStats{}, seasonResponse{}, err
	}

	playerStats, err := playersManager.GetPlayersStats(season)
	if err != nil {
		return &[]players.PlayerStats{}, seasonResponse{}, err
	}

	playersManager.GetSortFunction(sort)(playerStats)

	seasonRes := getSeasonResponse(season)

	gamesCount, err := gamesManager.GetGamesCountForSeason(season)
	if err != nil {
		return &[]players.PlayerStats{}, seasonResponse{}, err
	}
	seasonRes.GamesCount = gamesCount

	return playerStats, seasonRes, nil
}

func getSeason(gamesManager g.Manager, seasonUuid string) (g.Season, error) {
	if seasonUuid != "" {
		return gamesManager.GetSeasonByUuid(seasonUuid)
	}

	return gamesManager.ActiveSeason()
}

func getSeasonResponse(season g.Season) seasonResponse {
	return seasonResponse{
		UUID: season.UUID,
		Name: season.Name,
		Active: season.Active,
		CreatedAt: season.CreatedAt,
		UpdatedAt: season.UpdatedAt,
	}
}
