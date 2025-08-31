package main

import (
	"fmt"
	"log"

	"github.com/spie/fskick/cmd/server/static"
	"github.com/spie/fskick/internal/config"
	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
	"github.com/spie/fskick/internal/seasons"
	"github.com/spie/fskick/internal/server"
	"github.com/spie/fskick/internal/streaks"
	"github.com/spie/fskick/internal/views"
	"github.com/spie/fskick/migrations"
)

//go:generate templ generate -path ../../internal/templates
//go:generate npx tailwindcss build -c tailwind.config.js -i ./static/css/style.css -o ./static/css/tailwind.css -m

func main() {
	cfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.OpenDbConnection(cfg.DbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	err = db.MigrateFS(conn, migrations.FS, ".")
	if err != nil {
		log.Fatal(err)
	}

	seasonsRepository := seasons.NewSeasonsRepository(conn)
	seasonManager := seasons.NewManager(seasonsRepository)

	gamesRepository := games.NewGamesRepository(conn)
	attendanceRepository := games.NewAttendanceRepository(conn)
	gamesManager := games.NewManager(gamesRepository, attendanceRepository, seasonManager)

	playersRepository := players.NewPlayerRepository(conn)
	playersManager := players.NewManager(playersRepository)

	streaksManager := streaks.NewManager(attendanceRepository)

	gamesViews := server.NewGamesViews()
	gamesViews.SeasonsTable = views.NewSeasonTable()
	gamesViews.PlayersTable = views.NewPlayersTable()
	gamesViews.PlayerInfo = views.NewPlayerInfo()
	gamesViews.SeasonsTableUpdate = views.NewSeasonsTableUpdate()
	gamesViews.PlayersTableUpdate = views.NewPlayersTableUpdate()
	gamesViews.FavoriteTeamUpdate = views.NewFavoriteTeamUpdate()
	gamesController := server.NewGamesController(
		gamesManager,
		seasonManager,
		playersManager,
		streaksManager,
		gamesViews,
	)

	seasonsController := server.NewSeasonsController(seasonManager)

	streaksViews := server.NewStreaksViews()
	streaksViews.StreaksPage = views.NewStreaksPage()
	streaksController := server.NewStreaksController(streaksManager, streaksViews)

	imprintView := views.NewImprintView()
	imprintController := server.NewImprintController(cfg.ImprintText, imprintView)

	s := server.New(cfg.ApiHost)

	s.Get("/", gamesController.SeasonsTable)
	s.Get("/players", gamesController.PlayersTable)
	s.Get("/players/{player}", gamesController.PlayerInfo)
	s.Get("/streaks", streaksController.StreaksPage)
	s.Get("/imprint", imprintController.Imprint)

	s.Get("/table/seasons", gamesController.SeasonsTableUpdate)
	s.Get("/table/players", gamesController.PlayersTableUpdate)
	s.Get("/table/players/{player}", gamesController.PlayersTableUpdate)
	s.Get("/table/players/{player}/team", gamesController.FavoriteTeamUpdate)
	s.Get("/table/players/{player}/oponents", gamesController.FavoriteOponentsUpdate)
	s.Get("/streaks/current", streaksController.CurrentStreaks)

	s.Get("/api/seasons", seasonsController.GetSeasons)
	s.Get("/api/seasons/table", gamesController.GetSeasonsTable)
	s.Get("/api/seasons/table/{season}", gamesController.GetSeasonsTable)
	s.Get("/api/players", gamesController.GetPlayers)
	s.Get("/api/players/{player}/team", gamesController.GetFavoriteTeam)
	s.Get("/api/players/{player}", gamesController.GetPlayers)
	s.Get("/api/games/count", gamesController.GetGamesCount)

	s.HandleStatic(static.Dir)

	fmt.Printf("Starting the server on %s...\n", cfg.ApiHost)

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
