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

	rootCommand := createCommands(seasonManager, gamesManager, playersManager)

	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}

func createCommands(seasonsManager seasons.Manager, gamesManager games.Manager, playersManager players.Manager) commands.Command {
	createPlayer := commands.NewCreatePlayerCommand(playersManager)
	getPlayers := commands.NewGetPlayersCommand(gamesManager)
	playersCommand := commands.NewPlayersCommand()
	playersCommand.AddCommand(createPlayer)
	playersCommand.AddCommand(getPlayers)

	createSeason := commands.NewCreateSeasonCommand(seasonsManager)
	getSeason := commands.NewGetSeasonsCommand(seasonsManager)
	activateSeason := commands.NewActivateSeasonCommand(seasonsManager)
	tableCommand := commands.NewGetTableCommand(gamesManager, seasonsManager)
	seasonsCommand := commands.NewSeasonsCommand()
	seasonsCommand.AddCommand(createSeason)
	seasonsCommand.AddCommand(getSeason)
	seasonsCommand.AddCommand(activateSeason)
	seasonsCommand.AddCommand(tableCommand)

	createGame := commands.NewCreateGameCommand(gamesManager, playersManager)
	gamesCommands := commands.NewGamesCommand()
	gamesCommands.AddCommand(createGame)

	rootCommand := commands.NewRootCommand()
	rootCommand.AddCommand(playersCommand)
	rootCommand.AddCommand(seasonsCommand)
	rootCommand.AddCommand(gamesCommands)

	return rootCommand
}
