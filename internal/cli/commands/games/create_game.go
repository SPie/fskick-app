package games

import (
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/spie/fskick/internal/cli"
	g "github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
)

type createGameCommand struct {
	cc             *cobra.Command
	gamesManager   g.Manager
	playersManager players.AttendanceCreator
}

func newCreateGame(gamesManager g.Manager, playersManager players.AttendanceCreator) *createGameCommand {
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

	playedAt, err := getPlayedAt(cmd)
	if err != nil {
		return err
	}

	game, err := createGameCommand.gamesManager.CreateGame(playedAt)
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
	names := make([]string, len(team))
	for i := 0; i < len(team); i++ {
		names[i] = (team)[i].Name
	}

	return strings.Join(names, ",")
}

func getPlayedAt(cmd *cobra.Command) (time.Time, error) {
	playedAtFlag, _ := cmd.Flags().GetString("playedAt")
	if playedAtFlag == "" {
		return time.Time{}, nil
	}

	playedAt, err := time.Parse("2006-01-02", playedAtFlag)
	if err != nil {
		return time.Time{}, err
	}

	return playedAt, nil
}
