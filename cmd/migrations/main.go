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
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

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

	seasonsRepository := games.NewSeasonsRepository(connectionHandler)
	seasonsRepository.AutoMigrate()

	gamesRepository := games.NewGamesRepository(connectionHandler)
	gamesRepository.AutoMigrate()

	playersRepository := players.NewPlayerRepository(connectionHandler)
	playersRepository.AutoMigrate()

	attendancesRepository := players.NewAttendancesRepository(connectionHandler)
	attendancesRepository.AutoMigrate()
}
