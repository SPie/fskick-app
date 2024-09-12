package commands

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/spie/fskick/internal/cli"
	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
)

type playersCommand struct {
	command
}

func NewPlayersCommand() *playersCommand {
	playersCommand := playersCommand{command: newCommand(&cobra.Command{
		Use:   "players",
		Short: "Commands to handle players",
		Long:  "All commands handling players like creating new players, show a specific player, list all players...",
	})}

	return &playersCommand
}

type createPlayerCommand struct {
	command
	playersManager players.Manager
}

func NewCreatePlayerCommand(playersManager players.Manager) *createPlayerCommand {
	createPlayerCommand := &createPlayerCommand{playersManager: playersManager}

	cc := &cobra.Command{
		Use:   "new [name]",
		Short: "Creates a new player",
		Long:  "Creates a new player with the given name. Will return an error if the name is already taken by another player",
		Args:  cobra.MinimumNArgs(1),
		RunE:  createPlayerCommand.createPlayer,
	}

	createPlayerCommand.command = newCommand(cc)

	return createPlayerCommand
}

func (createPlayerCommand *createPlayerCommand) createPlayer(cmd *cobra.Command, args []string) error {
	player, err := createPlayerCommand.playersManager.CreatePlayer(args[0])
	if err != nil {
		return err
	}

	cli.Print(fmt.Sprintf("Player %s created\n", player.Name))

	return nil
}

type getPlayersCommand struct {
	command
	gamesManager games.Manager
}

func NewGetPlayersCommand(gamesManager games.Manager) *getPlayersCommand {
	getPlayersCommand := getPlayersCommand{
		gamesManager: gamesManager,
	}

	cc := &cobra.Command{
		Use:   "get [name]",
		Short: "Get all players with stats",
		RunE:  getPlayersCommand.getPlayers,
	}

	cc.Flags().StringP("sort", "s", "", "Table sort by")

	getPlayersCommand.command = newCommand(cc)

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
