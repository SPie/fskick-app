package players

import (
	"github.com/spf13/cobra"

	"github.com/spie/fskick/games"
	p "github.com/spie/fskick/players"
)

type PlayersCommand interface {
	GetPlayersCommand() *cobra.Command
}

type playersCommand struct {
	cc *cobra.Command
}

func NewPlayersCommand(playersManager p.Manager, gamesManager games.Manager) PlayersCommand {
	playersCommand := playersCommand{cc: &cobra.Command{
		Use:   "players",
		Short: "Commands to handle players",
		Long:  "All commands handling players like creating new players, show a specific player, list all players...",
	}}

	createPlayerCommand := newCreatePlayerCommand(playersManager)
	playersCommand.cc.AddCommand(createPlayerCommand.cc)
	getPlayersCommand := newGetPlayersCommand(playersManager, gamesManager)
	playersCommand.cc.AddCommand(getPlayersCommand.cc)

	return &playersCommand
}

func (playersCommand *playersCommand) GetPlayersCommand() *cobra.Command {
	return playersCommand.cc
}
