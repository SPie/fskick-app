package seasons

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/spie/fskick/internal/cli"
	"github.com/spie/fskick/internal/games"
)

type createSeasonCommand struct {
	cc           *cobra.Command
	gamesManager games.Manager
}

func newCreateSeasonCommand(gamesManager games.Manager) *createSeasonCommand {
	createSeasonCommand := &createSeasonCommand{gamesManager: gamesManager}

	cc := &cobra.Command{
		Use:   "new [name]",
		Short: "Create a new season",
		Long:  "Create a new season with the given name. Will return an error if the name is already taken by another season.",
		Args:  cobra.MinimumNArgs(1),
		RunE:  createSeasonCommand.createSeason,
	}

	createSeasonCommand.cc = cc

	return createSeasonCommand
}

func (createScreateSeasonCommand *createSeasonCommand) createSeason(cmd *cobra.Command, args []string) error {
	season, err := createScreateSeasonCommand.gamesManager.CreateSeason(args[0])
	if err != nil {
		return err
	}

	cli.Print(fmt.Sprintf("Season %s created", season.Name))

	return nil
}
