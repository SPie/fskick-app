package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	gamesHandlers "github.com/spie/fskick/api/handlers/games"
	playersHandlers "github.com/spie/fskick/api/handlers/players"
	"github.com/spie/fskick/games"
	"github.com/spie/fskick/players"
)

func SetUp(playersManager players.Manager, gamesManager games.Manager) *gin.Engine {
	engine := gin.Default()

	engine.Use(setUpCors())

	api := engine.Group("/api")
	{
		api.GET("/seasons", gamesHandlers.GetSeasons(gamesManager))
		api.GET("/seasons/table", gamesHandlers.GetTable(playersManager, gamesManager))
		api.GET("/seasons/table/:season", gamesHandlers.GetTable(playersManager, gamesManager))

		api.GET("/players", playersHandlers.GetPlayers(playersManager))
		api.GET("/players/:player/streak", playersHandlers.GetPlayerAttendances(playersManager))

		api.GET("/games/count", gamesHandlers.GetGamesCount(gamesManager))
	}

	return engine
}

func setUpCors() gin.HandlerFunc {
	return cors.Default()
}
