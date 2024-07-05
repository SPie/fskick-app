package main

import (
	"log"

	"github.com/spie/fskick/internal/cli/commands"
	"github.com/spie/fskick/internal/config"
	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
	"github.com/spie/fskick/internal/uuid"
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

	connectionHandler, err := db.NewConnectionHandler(cfg.DbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer connectionHandler.Close()

	uuidGenerator := uuid.NewGenerator()

	gamesRepository := games.NewGamesRepository(connectionHandler, dbHandler, uuidGenerator)
	seasonsRepository := games.NewSeasonsRepository(dbHandler, uuidGenerator)
	gamesManager := games.NewManager(gamesRepository, seasonsRepository)

	playersRepository := players.NewPlayerRepository(connectionHandler, dbHandler, uuidGenerator)
	attentanceRepository := players.NewAttendancesRepository(connectionHandler, dbHandler, uuidGenerator)
	playersManager := players.NewManager(playersRepository, attentanceRepository)

	if err := commands.NewRootCommand(playersManager, gamesManager).Execute(); err != nil {
		log.Fatal(err)
	}
}
