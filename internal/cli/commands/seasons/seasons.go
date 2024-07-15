package seasons

import (
	"github.com/spf13/cobra"

	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/seasons"
)

type SeasonsCommand interface {
	GetSeasonsCommand() *cobra.Command
}

type seasonsCommand struct {
	cc *cobra.Command
}

func NewSeasonsCommand(
	gamesManager games.Manager,
	seasonsManager seasons.Manager,
) SeasonsCommand {
	seasonsCommand := seasonsCommand{cc: &cobra.Command{
		Use:   "seasons",
		Short: "Commands to handle seasons",
		Long:  "All commands to handle seasons like creating new seasons, switch active seasons, show tables...",
	}}

	createSeasonCommand := newCreateSeasonCommand(seasonsManager)
	seasonsCommand.cc.AddCommand(createSeasonCommand.cc)
	getSeasonsCommand := newGetSeasonsCommand(seasonsManager)
	seasonsCommand.cc.AddCommand(getSeasonsCommand.cc)
	activateSeasonComand := newActivateSeasonCommand(seasonsManager)
	seasonsCommand.cc.AddCommand(activateSeasonComand.cc)
	tableCommand := newGetTableCommand(gamesManager, seasonsManager)
	seasonsCommand.cc.AddCommand(tableCommand.cc)

	return &seasonsCommand
}

func (seasonsCommand *seasonsCommand) GetSeasonsCommand() *cobra.Command {
	return seasonsCommand.cc
}
