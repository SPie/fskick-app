package commands

import (
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/spie/fskick/internal/cli"
	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
)

type gamesCommand struct {
	cc *cobra.Command
}

func NewGamesCommand() *gamesCommand {
	gamesCommand := gamesCommand{cc: &cobra.Command{
		Use:   "games",
		Short: "Commands to handle games",
	}}

	return &gamesCommand
}

func (command *gamesCommand) AddCommand(c Command) {
	command.cc.AddCommand(c.getCommand())
}

func (command *gamesCommand) getCommand() *cobra.Command {
	return command.cc
}

type createGameCommand struct {
	cc             *cobra.Command
	gamesManager   games.Manager
	playersManager players.Manager
}

func NewCreateGameCommand(gamesManager games.Manager, playersManager players.Manager) *createGameCommand {
	createGameCommand := createGameCommand{gamesManager: gamesManager, playersManager: playersManager}

	cc := &cobra.Command{
		Use:  "new",
		RunE: createGameCommand.CreateGame,
	}

	cc.Flags().StringP("winners", "w", "", "comma seperated names of winners")
	cc.Flags().StringP("losers", "l", "", "comma seperated names of losers")
	cc.Flags().StringP("playedAt", "p", "", "Date and time of the game")

	createGameCommand.cc = cc

	return &createGameCommand
}

func (createGameCommand *createGameCommand) CreateGame(cmd *cobra.Command, args []string) error {
	winnerNames, _ := cmd.Flags().GetString("winners")
	loserNames, _ := cmd.Flags().GetString("losers")

	winners, losers, err := createGameCommand.playersManager.GetTeamsByNames(
		getPlayerNamesFromFlag(winnerNames),
		getPlayerNamesFromFlag(loserNames),
	)
	if err != nil {
		return err
	}

	playedAt, err := createGameCommand.getPlayedAt()
	if err != nil {
		return err
	}

	_, err = createGameCommand.gamesManager.CreateGame(playedAt, winners, losers)
	if err != nil {
		return err
	}

	cli.PrintTable(
		[]string{},
		[][]string{
			{"Winners", winnerNames},
			{"Losers", loserNames},
		},
	)

	return nil
}

func getPlayerNamesFromFlag(names string) []string {
	if names == "" {
		return []string{}
	}

	return strings.Split(names, ",")
}

func (command *createGameCommand) getPlayedAt() (time.Time, error) {
	playedAtFlag, _ := command.cc.Flags().GetString("playedAt")
	if playedAtFlag == "" {
		return time.Time{}, nil
	}

	playedAt, err := time.Parse("2006-01-02", playedAtFlag)
	if err != nil {
		return time.Time{}, err
	}

	return playedAt, nil
}

func (command *createGameCommand) AddCommand(c Command) {
	command.cc.AddCommand(c.getCommand())
}

func (command *createGameCommand) getCommand() *cobra.Command {
	return command.cc
}
