package main

import (
	"fmt"
	"log"

	"github.com/spie/fskick/internal/config"
	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
	"github.com/spie/fskick/internal/seasons"
	"github.com/spie/fskick/internal/server"
	"github.com/spie/fskick/internal/views"
	"github.com/spie/fskick/migrations"
	"github.com/spie/fskick/static"
)

func main() {
	cfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbHandler, err := db.OpenDbHandler(cfg.DbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer dbHandler.Close()

	err = dbHandler.MigrateFS(migrations.FS, ".")
	if err != nil {
		log.Fatal(err)
	}

	seasonsRepository := seasons.NewSeasonsRepository(dbHandler)
	seasonManager := seasons.NewManager(seasonsRepository)

	gamesRepository := games.NewGamesRepository(dbHandler)
	attendanceRepository := games.NewAttendanceRepository(dbHandler)
	gamesManager := games.NewManager(gamesRepository, attendanceRepository, seasonManager)

	playersRepository := players.NewPlayerRepository(dbHandler)
	playersManager := players.NewManager(playersRepository)

	gamesViews := server.NewGamesViews()
	gamesViews.SeasonsTable = views.NewSeasonTable()
	gamesViews.SeasonsTableUpdate = views.NewSeasonsTableUpdate()
	gamesViews.PlayersTableUpdate = views.NewPlayersTableUpdate()
	gamesViews.FavoriteTeamUpdate = views.NewFavoriteTeamUpdate()
	gamesController := server.NewGamesController(gamesManager, seasonManager, playersManager, gamesViews)

	seasonsController := server.NewSeasonsController(seasonManager)

	playersViews := server.NewPlayersViews()
	playersViews.PlayersTable = views.NewPlayersTable()
	playersViews.PlayerInfo = views.NewPlayerInfo()
	playersController := server.NewPlayersController(playersManager, gamesManager, playersViews)

	imprintView := views.NewImprintView()
	imprintController := server.NewImprintController(cfg.ImprintText, imprintView)

	s := server.New(cfg.ApiHost)

	s.Get("/", gamesController.TablePage)
	s.Get("/players", playersController.GetPlayersTable)
	s.Get("/players/{player}", playersController.GetPlayerInfo)

	s.Get("/api/seasons", seasonsController.GetSeasons)
	s.Get("/api/seasons/table", gamesController.GetTable)
	s.Get("/api/seasons/table/{season}", gamesController.GetTable)

	s.Get("/api/players", playersController.GetPlayers)
	s.Get("/api/players/{player}/team", playersController.GetFavoriteTeam)
	s.Get("/api/players/{player}", playersController.GetPlayers)

	s.Get("/api/games/count", gamesController.GetGamesCount)

	s.Get("/imprint", imprintController.Imprint)

	s.Get("/table/seasons", gamesController.SeasonsTableUpdate)
	s.Get("/table/players", gamesController.PlayersTableUpdate)
	s.Get("/table/players/{player}", gamesController.PlayersTableUpdate)
	s.Get("/table/players/{player}/team", gamesController.GetFavoriteTeam)

	s.HandleStatic(static.Dir)

	fmt.Printf("Starting the server on %s...\n", cfg.ApiHost)

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
