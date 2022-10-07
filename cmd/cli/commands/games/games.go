package games

import (
	"github.com/spf13/cobra"

	g "github.com/spie/fskick/games"
	"github.com/spie/fskick/players"
)

type GamesCommand interface {
	GetGamesCommand() *cobra.Command
}

type gamesCommand struct {
	cc *cobra.Command
}

func NewGamesCommand(gamesManager g.Manager, playersManager players.AttendanceCreator) GamesCommand {
	gamesCommand := gamesCommand{cc: &cobra.Command{
		Use:   "games",
		Short: "Commands to handle games",
	}}

	createGameCommand := newCreateGame(gamesManager, playersManager)
	gamesCommand.cc.AddCommand(createGameCommand.cc)

	return &gamesCommand
}

func (gamesCommand *gamesCommand) GetGamesCommand() *cobra.Command {
	return gamesCommand.cc
}
