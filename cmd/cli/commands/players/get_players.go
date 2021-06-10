package players

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/spie/fskick/cli"
	"github.com/spie/fskick/games"
	p "github.com/spie/fskick/players"
)

type getPlayersCommand struct {
	cc             *cobra.Command
	playersManager p.Manager
	gamesManager   games.Manager
}

func newGetPlayersCommand(playersManager p.Manager, gamesManager games.Manager) *getPlayersCommand {
	getPlayersCommand := getPlayersCommand{playersManager: playersManager, gamesManager: gamesManager}

	cc := &cobra.Command{
		Use:   "get [name]",
		Short: "Get all players with stats",
		RunE:  getPlayersCommand.getPlayers,
	}

	cc.Flags().StringP("sort", "s", "", "Table sort by")

	getPlayersCommand.cc = cc

	return &getPlayersCommand
}

func (getPlayersCommand *getPlayersCommand) getPlayers(cmd *cobra.Command, args []string) error {
	gamesCount, err := getPlayersCommand.gamesManager.GetGamesCount()
	if err != nil {
		return err
	}

	sortName, err := getSortName(cmd)
	if err != nil {
		return err
	}

	playersStats, err := getPlayersCommand.playersManager.GetPlayerStats(games.Season{}, getPlayersCommand.playersManager.GetSortFunction(sortName))
	if err != nil {
		return err
	}

	if len(args) > 0 {
		playersStats = filterPlayerStatsByName(args[0], playersStats)
	}

	head := cli.CreateTableHead(gamesCount, playersStats)
	tableEntries := cli.CreateTableEntries(gamesCount, playersStats)

	cli.PrintTable(head, tableEntries)

	return nil
}

func filterPlayerStatsByName(name string, playersStats *[]p.PlayerStats) *[]p.PlayerStats {
	for _, playerStats := range *playersStats {
		if playerStats.Name == name {
			return &[]p.PlayerStats{playerStats}
		}
	}

	return &[]p.PlayerStats{}
}

func getSortName(cmd *cobra.Command) (string, error) {
	sortName, err := cmd.Flags().GetString("sort")
	if err != nil {
		return "", err
	}

	if sortName == "" {
		return p.SortByPointsRatio, nil
	}

	if sortName != p.SortByPointsRatio && sortName != p.SortByWins && sortName != p.SortByGames && sortName != p.SortByWinRatio {
		return "", errors.New("Sort flag has to be pointsRatio, games, wins or winRatio")
	}

	return sortName, nil
}
