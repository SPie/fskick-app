package seasons

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/spie/fskick/cli"
	"github.com/spie/fskick/games"
	"github.com/spie/fskick/html"
	"github.com/spie/fskick/players"
)

type getTableCommand struct {
	cc             *cobra.Command
	gamesManager   games.Manager
	playersManager players.Manager
	htmlWriter     html.HtmlWriter
}

func newGetTableCommand(gamesManager games.Manager, playersManager players.Manager, htmlWriter html.HtmlWriter) *getTableCommand {
	getTableCommand := &getTableCommand{gamesManager: gamesManager, playersManager: playersManager, htmlWriter: htmlWriter}

	cc := &cobra.Command{
		Use:   "table",
		Short: "Get the seasons table",
		Long:  "Get the seasons table. If no season name is provided, the active season will be used.",
		RunE:  getTableCommand.getTable,
	}

	cc.Flags().StringP("sort", "s", "", "Table sort by")
	cc.Flags().BoolP("html", "H", false, "Print table to HTML file")

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

	playerStats, err := tableCommand.playersManager.GetPlayerStats(season, tableCommand.playersManager.GetSortFunction(sortName))
	if err != nil {
		return err
	}

	gamesCount := len(*season.Games)

	head := cli.CreateTableHead(gamesCount, playerStats)
	tableEntries := cli.CreateTableEntries(gamesCount, playerStats)

	err = tableCommand.writeHtmlTable(cmd, season, head, tableEntries)
	if err != nil {
		return err
	}

	cli.Print(fmt.Sprintf("Season: %s", season.Name))
	cli.PrintTable(head, tableEntries)

	return nil
}

func (tableCommand *getTableCommand) getSeason(args []string) (games.Season, error) {
	if len(args) > 0 {
		season, err := tableCommand.gamesManager.GetSeasonByName(args[0])

		return *season, err
	}

	season, err := tableCommand.gamesManager.ActiveSeason()

	return *season, err
}

func getSortName(cmd *cobra.Command) (string, error) {
	sortName, err := cmd.Flags().GetString("sort")
	if err != nil {
		return "", err
	}

	if sortName == "" {
		return players.SortByPointsRatio, nil
	}

	if sortName != players.SortByPointsRatio && sortName != players.SortByWins && sortName != players.SortByGames {
		return "", errors.New("Sort flag has to be pointsRatio, games or wins")
	}

	return sortName, nil
}

func (getTableCommand *getTableCommand) writeHtmlTable(cmd *cobra.Command, season games.Season, head []string, tableEntries [][]string) error {
	withHtml, err := cmd.Flags().GetBool("html")
	if err != nil {
		return err
	}

	if withHtml {
		return getTableCommand.htmlWriter.WriteSeasonTable(season, head, tableEntries)
	}

	return nil
}
