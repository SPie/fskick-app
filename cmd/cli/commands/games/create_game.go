package games

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/spie/fskick/cli"
	g "github.com/spie/fskick/games"
	"github.com/spie/fskick/players"
)

type createGameCommand struct {
	cc             *cobra.Command
	gamesManager   g.Manager
	playersManager players.Manager
}

func newCreateGame(gamesManager g.Manager, playersManager players.Manager) *createGameCommand {
	createGameCommand := createGameCommand{gamesManager: gamesManager, playersManager: playersManager}

	cc := &cobra.Command{
		Use:  "new",
		RunE: createGameCommand.CreateGame,
	}

	cc.Flags().StringP("winners", "w", "", "comma seperated names of winners")
	cc.Flags().StringP("losers", "l", "", "comma seperated names of losers")

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

	game, err := createGameCommand.gamesManager.CreateGame()
	if err != nil {
		return err
	}

	winnersTeam, losersTeam, err := createGameCommand.playersManager.CreateAttendances(game, winners, losers)
	if err != nil {
		return err
	}

	cli.PrintTable(
		[]string{},
		[][]string{
			{"Winners", getPlayerNamesAsString(winnersTeam)},
			{"Losers", getPlayerNamesAsString(losersTeam)},
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

func getPlayerNamesAsString(team players.Team) string {
	names := make([]string, len(*team))
	for i := 0; i < len(*team); i++ {
		names[i] = (*team)[i].Name
	}

	return strings.Join(names, ",")
}
