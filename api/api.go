package api

import (
	"github.com/gin-gonic/gin"

	gh "github.com/spie/fskick/api/handlers/games"
	"github.com/spie/fskick/games"
	"github.com/spie/fskick/players"
)

func SetUp(playersManager players.Manager, gamesManager games.Manager) *gin.Engine {
	engine := gin.Default()

	api := engine.Group("/api")
	{
		api.GET("/seasons", gh.GetSeasons(gamesManager))
		api.GET("/seasons/table", gh.GetTable(playersManager, gamesManager))
	}

	return engine
}
