package seasons

import (
	// "fmt"

	"fmt"

	"github.com/spf13/cobra"

	"github.com/spie/fskick/cli"
	"github.com/spie/fskick/games"
)

type activateSeasonCommand struct {
	cc           *cobra.Command
	gamesManager games.Manager
}

func newActivateSeasonCommand(gamesManager games.Manager) *activateSeasonCommand {
	activateSeasonCommand := activateSeasonCommand{gamesManager: gamesManager}

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
	season, err := activateSeasonCommand.gamesManager.ActivateSeason(args[0])
	if err != nil {
		return err
	}

	cli.Print(fmt.Sprintf("Season %s activated", season.Name))

	return nil
}
