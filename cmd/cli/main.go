package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spie/fskick/cmd/cli/commands"
	"github.com/spie/fskick/db"
	"github.com/spie/fskick/games"
	"github.com/spie/fskick/html"
	"github.com/spie/fskick/players"
)

func main() {
	godotenv.Load()

	connectionHandler, err := db.NewConnectionHandler(
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DRIVER"),
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

	htmlTableFile, err := os.OpenFile(os.Getenv("SEASON_TABLE_HTML"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer htmlTableFile.Close()

	if err := commands.NewRootCommand(playersManager, gamesManager, html.NewHtmlWriter(htmlTableFile)).Execute(); err != nil {
		log.Fatal(err)
	}
}
