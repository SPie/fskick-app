package commands

import (
	"github.com/spf13/cobra"

	"github.com/spie/fskick/cmd/cli/commands/games"
	"github.com/spie/fskick/cmd/cli/commands/players"
	"github.com/spie/fskick/cmd/cli/commands/seasons"
	g "github.com/spie/fskick/games"
	p "github.com/spie/fskick/players"
)

type RootCommand interface {
	Execute() error
}

type rootCommand struct {
	cc *cobra.Command
}

func NewRootCommand(playersManager p.Manager, gamesManager g.Manager) RootCommand {
	rootCommand := &rootCommand{cc: &cobra.Command{
		Use:   "fskick",
		Short: "FSKick CLI app",
		Long:  "CLI app for FSKick to create new players, seasons, games and show results and statistics",
	}}

	playersCommand := players.NewPlayersCommand(playersManager, gamesManager)
	rootCommand.cc.AddCommand(playersCommand.GetPlayersCommand())
	seasonsCommand := seasons.NewSeasonsCommand(gamesManager, playersManager)
	rootCommand.cc.AddCommand(seasonsCommand.GetSeasonsCommand())
	gamesCommand := games.NewGamesCommand(gamesManager, playersManager)
	rootCommand.cc.AddCommand(gamesCommand.GetGamesCommand())

	return rootCommand
}

func (rootCommand *rootCommand) Execute() error {
	return rootCommand.cc.Execute()
}
