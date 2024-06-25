package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spie/fskick/db"
	"github.com/spie/fskick/games"
	"github.com/spie/fskick/players"
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

	seasonsRepository := games.NewSeasonsRepository(connectionHandler)
	seasonsRepository.AutoMigrate()

	gamesRepository := games.NewGamesRepository(connectionHandler)
	gamesRepository.AutoMigrate()

	playersRepository := players.NewPlayerRepository(connectionHandler)
	playersRepository.AutoMigrate()

	attendancesRepository := players.NewAttendancesRepository(connectionHandler)
	attendancesRepository.AutoMigrate()
}
