package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spie/fskick/internal/cli/commands"
	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
)

func main() {
	godotenv.Load()

	connectionHandler, err := db.NewConnectionHandler(
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_LOG") != "false",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer connectionHandler.Close()

	if os.Getenv("DB_DEBUG") == "true" {
		connectionHandler.SetDebug()
	}

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
