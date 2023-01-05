package seasons

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/spie/fskick/cli"
	"github.com/spie/fskick/games"
	"github.com/spie/fskick/players"
)

type getTableCommand struct {
	cc             *cobra.Command
	gamesManager   games.Manager
	playersManager players.PlayerStatsCalculator
}

func newGetTableCommand(gamesManager games.Manager, playersManager players.PlayerStatsCalculator) *getTableCommand {
	getTableCommand := &getTableCommand{gamesManager: gamesManager, playersManager: playersManager}

	cc := &cobra.Command{
		Use:   "table",
		Short: "Get the seasons table",
		Long:  "Get the seasons table. If no season name is provided, the active season will be used.",
		RunE:  getTableCommand.getTable,
	}

	cc.Flags().StringP("sort", "s", "", "Table sort by")

	getTableCommand.cc = cc

	return getTableCommand
}

func (tableCommand *getTableCommand) getTable(cmd *cobra.Command, args []string) error {
	season, err := tableCommand.getSeason(args)
	if err != nil {
		return err
	}

	sortName, err := getSortName(cmd)
	if err != nil {
		return err
	}

	playerStats, err := tableCommand.playersManager.GetPlayersStats(season)
	if err != nil {
		return err
	}

	tableCommand.playersManager.GetSortFunction(sortName)(playerStats)

	gamesCount := len(*season.Games)

	head := cli.CreateTableHead(gamesCount, playerStats)
	tableEntries := cli.CreateTableEntries(gamesCount, playerStats)

	cli.Print(fmt.Sprintf("Season: %s", season.Name))
	cli.PrintTable(head, tableEntries)

	return nil
}

func (tableCommand *getTableCommand) getSeason(args []string) (games.Season, error) {
	if len(args) > 0 {
		season, err := tableCommand.gamesManager.GetSeasonByName(args[0])

		return season, err
	}

	season, err := tableCommand.gamesManager.ActiveSeason()

	return season, err
}

func getSortName(cmd *cobra.Command) (string, error) {
	sortName, err := cmd.Flags().GetString("sort")
	if err != nil {
		return "", err
	}

	if sortName == "" {
		return players.SortByPointsRatio, nil
	}

	if sortName != players.SortByPointsRatio && sortName != players.SortByWins && sortName != players.SortByGames && sortName != players.SortByWinRatio {
		return "", errors.New("Sort flag has to be pointsRatio, games, wins or winRatio")
	}

	return sortName, nil
}
