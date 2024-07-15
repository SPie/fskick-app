package games

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
	"github.com/spie/fskick/internal/seasons"
)

type seasonResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Active bool `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	GamesCount int `json:"gamesCount"`
}

func GetTable(
	playersManager players.Manager,
	gamesManager games.Manager,
	seasonsManager seasons.Manager,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		playerStats, season, err := getTable(
			playersManager,
			gamesManager,
			seasonsManager,
			"",
			c.DefaultQuery("sort", "pointsRatio"),
		)
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

func GetTableForSeason(
	playersManager players.Manager,
	gamesManager games.Manager,
	seasonsManager seasons.Manager,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request getTableForSeasonRequest
		if err := c.ShouldBindUri(&request); err != nil {
			c.Error(err)
			return
		}

		playerStats, season, err := getTable(
			playersManager,
			gamesManager,
			seasonsManager,
			request.Season,
			c.DefaultQuery("sort", "pointsRatio"),
		)
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


func getTable(
	playersManager players.Manager,
	gamesManager games.Manager,
	seasonsManager seasons.Manager,
	seasonUuid string,
	sort string,
) ([]games.PlayerStats, seasonResponse, error) {
	season, err := getSeason(seasonsManager, seasonUuid)
	if err != nil {
		return nil, seasonResponse{}, err
	}

	playerStats, err := gamesManager.GetPlayerStatsForSeason(season, sort)
	if err != nil {
		return nil, seasonResponse{}, err
	}

	seasonRes := getSeasonResponse(season)

	gamesCount, err := gamesManager.GetGamesCountForSeason(season)
	if err != nil {
		return nil, seasonResponse{}, err
	}
	seasonRes.GamesCount = gamesCount

	return playerStats, seasonRes, nil
}

func getSeason(seasonsManager seasons.Manager, seasonUuid string) (seasons.Season, error) {
	if seasonUuid != "" {
		return seasonsManager.GetSeasonByUuid(seasonUuid)
	}

	return seasonsManager.ActiveSeason()
}

func getSeasonResponse(season seasons.Season) seasonResponse {
	return seasonResponse{
		UUID: season.UUID,
		Name: season.Name,
		Active: season.Active,
		CreatedAt: season.CreatedAt,
		UpdatedAt: season.UpdatedAt,
	}
}
