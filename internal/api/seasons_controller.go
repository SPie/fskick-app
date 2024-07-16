package api

import (
	"github.com/gin-gonic/gin"
	"github.com/spie/fskick/internal/seasons"
)

type SeasonsController struct {
	seasonsManager seasons.Manager
}

func NewSeasonsController(seasonsManager seasons.Manager) SeasonsController {
	return SeasonsController{seasonsManager: seasonsManager}
}

func (controller SeasonsController) GetSeasons(c *gin.Context) {
	seasons, err := controller.seasonsManager.GetSeasons()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"seasons": seasons})
}
