package players

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/spie/fskick/internal/cli"
	"github.com/spie/fskick/internal/games"
)

type getPlayersCommand struct {
	cc             *cobra.Command
	gamesManager   games.Manager
}

func newGetPlayersCommand(gamesManager games.Manager) *getPlayersCommand {
	getPlayersCommand := getPlayersCommand{
		gamesManager: gamesManager,
	}

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

	playersStats, err := getPlayersCommand.gamesManager.GetAllPlayerStats(sortName)
	if err != nil {
		return err
	}

	if len(args) > 0 {
		playersStats = filterPlayerStatsByName(args[0], playersStats)
	}

	head := cli.CreateTableHead(gamesCount, len(playersStats))
	tableEntries := cli.CreateTableEntries(gamesCount, playersStats)

	cli.PrintTable(head, tableEntries)

	return nil
}

func filterPlayerStatsByName(name string, playersStats []games.PlayerStats) []games.PlayerStats {
	for _, playerStats := range playersStats {
		if playerStats.Name == name {
			return []games.PlayerStats{playerStats}
		}
	}

	return []games.PlayerStats{}
}

func getSortName(cmd *cobra.Command) (string, error) {
	sortName, err := cmd.Flags().GetString("sort")
	if err != nil {
		return "", err
	}

	if sortName == "" {
		return "pointsRatio", nil
	}

	if sortName != "pointsRatio" && sortName != "wint" && sortName != "games" && sortName != "winRatio" {
		return "", errors.New("Sort flag has to be pointsRatio, games, wins or winRatio")
	}

	return sortName, nil
}
