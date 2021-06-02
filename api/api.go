package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	gh "github.com/spie/fskick/api/handlers/games"
	"github.com/spie/fskick/games"
	"github.com/spie/fskick/players"
)

func SetUp(playersManager players.Manager, gamesManager games.Manager) *gin.Engine {
	engine := gin.Default()

	engine.Use(setUpCors())

	api := engine.Group("/api")
	{
		api.GET("/seasons", gh.GetSeasons(gamesManager))
		api.GET("/seasons/table", gh.GetTable(playersManager, gamesManager))
	}

	return engine
}

func setUpCors() gin.HandlerFunc {
	return cors.Default()
}
