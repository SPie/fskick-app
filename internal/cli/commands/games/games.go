package games

import (
	"github.com/spf13/cobra"

	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
)

type GamesCommand interface {
	GetGamesCommand() *cobra.Command
}

type gamesCommand struct {
	cc *cobra.Command
}

func NewGamesCommand(gamesManager games.Manager, playersManager players.Manager) GamesCommand {
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
