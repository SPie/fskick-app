package players

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/spie/fskick/internal/cli"
	"github.com/spie/fskick/internal/players"
)

type createPlayerCommand struct {
	cc             *cobra.Command
	playersManager players.Manager
}

func newCreatePlayerCommand(playersManager players.Manager) *createPlayerCommand {
	createPlayerCommand := &createPlayerCommand{playersManager: playersManager}

	cc := &cobra.Command{
		Use:   "new [name]",
		Short: "Creates a new player",
		Long:  "Creates a new player with the given name. Will return an error if the name is already taken by another player",
		Args:  cobra.MinimumNArgs(1),
		RunE:  createPlayerCommand.createPlayer,
	}

	createPlayerCommand.cc = cc

	return createPlayerCommand
}

func (createPlayerCommand *createPlayerCommand) createPlayer(cmd *cobra.Command, args []string) error {
	player, err := createPlayerCommand.playersManager.CreatePlayer(args[0])
	if err != nil {
		return err
	}

	cli.Print(fmt.Sprintf("Player %s created\n", player.Name))

	return nil
}
