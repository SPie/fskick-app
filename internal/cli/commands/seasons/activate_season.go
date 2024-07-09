package seasons

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/spie/fskick/internal/cli"
	"github.com/spie/fskick/internal/seasons"
)

type activateSeasonCommand struct {
	cc *cobra.Command
	seasonsManager seasons.Manager
}

func newActivateSeasonCommand(seasonsManager seasons.Manager) *activateSeasonCommand {
	activateSeasonCommand := activateSeasonCommand{seasonsManager: seasonsManager}

	cc := &cobra.Command{
		Use:   "activate [name]",
		Short: "Activates an inactive season",
		Long:  "Activates the given inactive season",
		Args:  cobra.MinimumNArgs(1),
		RunE:  activateSeasonCommand.activateSeason,
	}

	activateSeasonCommand.cc = cc

	return &activateSeasonCommand
}

func (activateSeasonCommand *activateSeasonCommand) activateSeason(cmd *cobra.Command, args []string) error {
	season, err := activateSeasonCommand.seasonsManager.ActivateSeason(args[0])
	if err != nil {
		return err
	}

	cli.Print(fmt.Sprintf("Season %s activated", season.Name))

	return nil
}
