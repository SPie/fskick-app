package players

import (
	"github.com/gin-gonic/gin"
	p "github.com/spie/fskick/players"
)

type getPlayerAttendancesRequest struct {
	player string `binding:"required"`
}

func GetPlayerAttendances(manager p.PlayerStatsCalculator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request getPlayerAttendancesRequest
		err := c.ShouldBindUri(&request)
		if err != nil {
			c.Error(err)
			return
		}

		attendances, err := manager.GetPlayerAttendances(request.player)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, gin.H{"attendances": attendances})
	}
}
