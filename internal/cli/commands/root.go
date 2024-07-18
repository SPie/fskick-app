package commands

import (
	"github.com/spf13/cobra"

	"github.com/spie/fskick/internal/cli/commands/games"
	"github.com/spie/fskick/internal/cli/commands/seasons"
	g "github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
	s "github.com/spie/fskick/internal/seasons"
)

type Command interface {
	AddCommand(command Command)
	getCommand() *cobra.Command
}

type rootCommand struct {
	cc *cobra.Command
}

func NewRootCommand(playersManager players.Manager, gamesManager g.Manager, seasonsManagers s.Manager) *rootCommand {
	rootCommand := rootCommand{cc: &cobra.Command{
		Use:   "fskick",
		Short: "FSKick CLI app",
		Long:  "CLI app for FSKick to create new players, seasons, games and show results and statistics",
	}}

	playersCommand := NewPlayersCommand(playersManager, gamesManager)
	rootCommand.AddCommand(playersCommand)
	seasonsCommand := seasons.NewSeasonsCommand(gamesManager, seasonsManagers)
	rootCommand.cc.AddCommand(seasonsCommand.GetSeasonsCommand())
	gamesCommand := games.NewGamesCommand(gamesManager, playersManager)
	rootCommand.cc.AddCommand(gamesCommand.GetGamesCommand())

	return &rootCommand
}

func (command *rootCommand) Execute() error {
	return command.cc.Execute()
}

func (command *rootCommand) AddCommand(c Command) {
	command.cc.AddCommand(c.getCommand())
}
