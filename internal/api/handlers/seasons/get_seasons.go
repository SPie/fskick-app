package seasons

import (
	"github.com/gin-gonic/gin"
	"github.com/spie/fskick/internal/seasons"
)

func GetSeasons(seasonsManager seasons.Manager) gin.HandlerFunc {
    return func(c *gin.Context) {
        seasons, err := seasonsManager.GetSeasons()
        if err != nil {
            c.Error(err)
            return
        }

        c.JSON(200, gin.H{"seasons": seasons})
    }
}
