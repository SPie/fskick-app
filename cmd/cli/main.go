package main

import (
	"log"

	"github.com/spie/fskick/internal/cli/commands"
	"github.com/spie/fskick/internal/config"
	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
	"github.com/spie/fskick/migrations"
)

func main() {
	cfg, err := config.LoadCliConfig()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.OpenDbConnection(cfg.DbConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = db.MigrateFS(conn, migrations.FS, ".")
	if err != nil {
		log.Fatal(err)
	}

	connectionHandler, err := db.NewConnectionHandler(cfg.DbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer connectionHandler.Close()

	gamesRepository := games.NewGamesRepository(connectionHandler)
	seasonsRepository := games.NewSeasonsRepository(connectionHandler)
	gamesManager := games.NewManager(gamesRepository, seasonsRepository)

	playersRepository := players.NewPlayerRepository(connectionHandler)
	attentanceRepository := players.NewAttendancesRepository(connectionHandler)
	playersManager := players.NewManager(playersRepository, attentanceRepository)

	if err := commands.NewRootCommand(playersManager, gamesManager).Execute(); err != nil {
		log.Fatal(err)
	}
}
