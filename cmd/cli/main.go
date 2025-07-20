package main

import (
	"log"

	"github.com/spie/fskick/internal/cli/commands"
	"github.com/spie/fskick/internal/config"
	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/passwords"
	"github.com/spie/fskick/internal/players"
	"github.com/spie/fskick/internal/seasons"
	"github.com/spie/fskick/internal/users"
	"github.com/spie/fskick/migrations"
)

var version string = "development"

func main() {
	cfg, err := config.LoadCliConfig()
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

	passwordService := passwords.NewPasswordService()

	seasonsRepository := seasons.NewSeasonsRepository(conn)
	seasonManager := seasons.NewManager(seasonsRepository)

	gamesRepository := games.NewGamesRepository(conn)
	attendanceRepository := games.NewAttendanceRepository(conn)
	gamesManager := games.NewManager(gamesRepository, attendanceRepository, seasonManager)

	playersRepository := players.NewPlayerRepository(conn)
	playersManager := players.NewManager(playersRepository)

	usersRepository := users.NewUsersRepository(conn)
	usersManager := users.NewManager(usersRepository, playersManager, passwordService)

	rootCommand := createCommands(seasonManager, gamesManager, playersManager, usersManager)

	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}

func createCommands(
	seasonsManager seasons.Manager,
	gamesManager games.Manager,
	playersManager players.Manager,
	usersManager users.Manager,
) commands.Command {
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

	createUserFromPlayer := commands.NewCreateUserFromPlayerCommand(usersManager)
	usersCommand := commands.NewUsersCommand()
	usersCommand.AddCommand(createUserFromPlayer)

	versionCommand := commands.NewVersionCommand(version)

	rootCommand := commands.NewRootCommand()
	rootCommand.AddCommand(versionCommand)
	rootCommand.AddCommand(playersCommand)
	rootCommand.AddCommand(seasonsCommand)
	rootCommand.AddCommand(gamesCommands)
	rootCommand.AddCommand(usersCommand)

	return rootCommand
}
