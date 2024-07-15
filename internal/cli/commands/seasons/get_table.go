package seasons

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/spie/fskick/internal/cli"
	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/seasons"
)

type getTableCommand struct {
	cc             *cobra.Command
	gamesManager   games.Manager
	seasonsManager seasons.Manager
}

func newGetTableCommand(
	gamesManager games.Manager,
	seasonsManager seasons.Manager,
) *getTableCommand {
	getTableCommand := &getTableCommand{
		gamesManager: gamesManager,
		seasonsManager: seasonsManager,
	}

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

	playerStats, err := tableCommand.gamesManager.GetPlayerStatsForSeason(season, sortName)
	if err != nil {
		return err
	}

	gamesCount, err := tableCommand.gamesManager.GetGamesCountForSeason(season)
	if err != nil {
		return err
	}

	head := cli.CreateTableHead(gamesCount, len(playerStats))
	tableEntries := cli.CreateTableEntries(gamesCount, playerStats)

	cli.Print(fmt.Sprintf("Season: %s", season.Name))
	cli.PrintTable(head, tableEntries)

	return nil
}

func (tableCommand *getTableCommand) getSeason(args []string) (seasons.Season, error) {
	if len(args) > 0 {
		season, err := tableCommand.seasonsManager.GetSeasonByName(args[0])

		return season, err
	}

	season, err := tableCommand.seasonsManager.ActiveSeason()

	return season, err
}

func getSortName(cmd *cobra.Command) (string, error) {
	sortName, err := cmd.Flags().GetString("sort")
	if err != nil {
		return "", err
	}

	if sortName == "" {
		return "pointsRatio", nil
	}

	if sortName != "pointsRatio" && sortName != "wins" && sortName != "games" && sortName != "winRatio" {
		return "", errors.New("Sort flag has to be pointsRatio, games, wins or winRatio")
	}

	return sortName, nil
}
