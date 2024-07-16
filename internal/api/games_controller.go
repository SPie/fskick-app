package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
	"github.com/spie/fskick/internal/seasons"
)

type GamesController struct {
	gamesManager games.Manager
	seasonsManager seasons.Manager
	playersManager players.Manager
}

func NewGamesController(
	gamesManager games.Manager,
	seasonsManager seasons.Manager,
	playersManager players.Manager,
) GamesController {
	return GamesController{
		gamesManager: gamesManager,
		seasonsManager: seasonsManager,
		playersManager: playersManager,
	}
}

func (controller GamesController) GetGamesCount(c *gin.Context) {
	gamesCount, err := controller.gamesManager.GetGamesCount()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"gamesCount": gamesCount})
}

type seasonResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Active bool `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	GamesCount int `json:"gamesCount"`
}

func (controller GamesController) GetTable(c *gin.Context) {
	playerStats, season, err := controller.getTable("", c.DefaultQuery("sort", "pointsRatio"))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"season":      season,
		"playerStats": playerStats,
	})
}

type getTableForSeasonRequest struct {
	Season string `uri:"season" binding:"required"`
}

func (controller GamesController) GetTableForSeason(c *gin.Context) {
	var request getTableForSeasonRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.Error(err)
		return
	}

	playerStats, season, err := controller.getTable(request.Season, c.DefaultQuery("sort", "pointsRatio"))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"season":      season,
		"playerStats": playerStats,
	})
}

type getTableForPlayerRequest struct {
	Season string `uri:"season" binding:"required"`
	Player string `uri:"player" binding:"required"`
}


func (controller GamesController) getTable(
	seasonUuid string,
	sort string,
) ([]games.PlayerStats, seasonResponse, error) {
	season, err := controller.getSeason(seasonUuid)
	if err != nil {
		return nil, seasonResponse{}, err
	}

	playerStats, err := controller.gamesManager.GetPlayerStatsForSeason(season, sort)
	if err != nil {
		return nil, seasonResponse{}, err
	}

	seasonRes := getSeasonResponse(season)

	gamesCount, err := controller.gamesManager.GetGamesCountForSeason(season)
	if err != nil {
		return nil, seasonResponse{}, err
	}
	seasonRes.GamesCount = gamesCount

	return playerStats, seasonRes, nil
}

func (controller GamesController) getSeason(seasonUuid string) (seasons.Season, error) {
	if seasonUuid != "" {
		return controller.seasonsManager.GetSeasonByUuid(seasonUuid)
	}

	return controller.seasonsManager.ActiveSeason()
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
