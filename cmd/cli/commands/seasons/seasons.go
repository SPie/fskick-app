package seasons

import (
	"github.com/spf13/cobra"

	"github.com/spie/fskick/games"
	"github.com/spie/fskick/html"
	"github.com/spie/fskick/players"
)

type SeasonsCommand interface {
	GetSeasonsCommand() *cobra.Command
}

type seasonsCommand struct {
	cc *cobra.Command
}

func NewSeasonsCommand(gamesManager games.Manager, playersManager players.Manager, htmlWriter html.HtmlWriter) SeasonsCommand {
	seasonsCommand := seasonsCommand{cc: &cobra.Command{
		Use:   "seasons",
		Short: "Commands to handle seasons",
		Long:  "All commands to handle seasons like creating new seasons, switch active seasons, show tables...",
	}}

	createSeasonCommand := newCreateSeasonCommand(gamesManager)
	seasonsCommand.cc.AddCommand(createSeasonCommand.cc)
	getSeasonsCommand := newGetSeasonsCommand(gamesManager)
	seasonsCommand.cc.AddCommand(getSeasonsCommand.cc)
	activateSeasonComand := newActivateSeasonCommand(gamesManager)
	seasonsCommand.cc.AddCommand(activateSeasonComand.cc)
	tableCommand := newGetTableCommand(gamesManager, playersManager, htmlWriter)
	seasonsCommand.cc.AddCommand(tableCommand.cc)

	return &seasonsCommand
}

func (seasonsCommand *seasonsCommand) GetSeasonsCommand() *cobra.Command {
	return seasonsCommand.cc
}
