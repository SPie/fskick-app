package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	gamesHandlers "github.com/spie/fskick/internal/api/handlers/games"
	playersHandlers "github.com/spie/fskick/internal/api/handlers/players"
	seasonsHandlers "github.com/spie/fskick/internal/api/handlers/seasons"
	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
	"github.com/spie/fskick/internal/seasons"
)

func SetUp(
	playersManager players.Manager,
	gamesManager games.Manager,
	seasonsManager seasons.Manager,
) *gin.Engine {
	engine := gin.Default()

	engine.Use(setUpCors())

	api := engine.Group("/api")
	{
		api.GET("/seasons", seasonsHandlers.GetSeasons(seasonsManager))
		api.GET("/seasons/table", gamesHandlers.GetTable(playersManager, gamesManager, seasonsManager))
		api.GET(
			"/seasons/table/:season",
			gamesHandlers.GetTableForSeason(playersManager, gamesManager, seasonsManager),
		)

		api.GET("/players", playersHandlers.GetPlayers(playersManager))
		api.GET(
			"/playsers/:player/seasons/:seasons",
			gamesHandlers.GetTableForPlayer(playersManager, gamesManager, seasonsManager),
		)
		api.GET("/players/:player/team", playersHandlers.GetFavoriteTeam(playersManager))
		api.GET("/players/:player", playersHandlers.GetPlayers(playersManager))

		api.GET("/games/count", gamesHandlers.GetGamesCount(gamesManager))
	}

	return engine
}

func setUpCors() gin.HandlerFunc {
	return cors.Default()
}
