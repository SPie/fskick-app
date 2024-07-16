package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetUp(
	gamesController GamesController,
	seasonsController SeasonsController,
	playersController PlayersController,
) *gin.Engine {
	engine := gin.Default()

	engine.Use(setUpCors())

	api := engine.Group("/api")
	{
		api.GET("/seasons", seasonsController.GetSeasons)
		api.GET("/seasons/table", gamesController.GetTable)
		api.GET("/seasons/table/:season", gamesController.GetTableForSeason)

		api.GET("/players", playersController.GetPlayers)
		api.GET("/players/:player/team", playersController.GetFavoriteTeam)
		api.GET("/players/:player", playersController.GetPlayers)

		api.GET("/games/count", gamesController.GetGamesCount)
	}

	return engine
}

func setUpCors() gin.HandlerFunc {
	return cors.Default()
}
