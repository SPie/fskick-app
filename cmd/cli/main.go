package main

import (
	"log"

	"github.com/spie/fskick/internal/cli/commands"
	"github.com/spie/fskick/internal/config"
	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
	"github.com/spie/fskick/internal/seasons"
	"github.com/spie/fskick/migrations"
)

func main() {
	cfg, err := config.LoadCliConfig()
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

	if err := commands.NewRootCommand(playersManager, gamesManager, seasonManager).Execute(); err != nil {
		log.Fatal(err)
	}
}
