package commands

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/spie/fskick/internal/cli"
	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/seasons"
)

type seasonsCommand struct {
	command
}

func NewSeasonsCommand() *seasonsCommand {
	seasonsCommand := seasonsCommand{command: newCommand(&cobra.Command{
		Use:   "seasons",
		Short: "Commands to handle seasons",
		Long:  "All commands to handle seasons like creating new seasons, switch active seasons, show tables...",
	})}

	return &seasonsCommand
}

type createSeasonCommand struct {
	command
	seasonsManager seasons.Manager
}

func NewCreateSeasonCommand(seasonsManager seasons.Manager) *createSeasonCommand {
	createSeasonCommand := &createSeasonCommand{seasonsManager: seasonsManager}

	cc := &cobra.Command{
		Use:   "new [name]",
		Short: "Create a new season",
		Long:  "Create a new season with the given name. Will return an error if the name is already taken by another season.",
		Args:  cobra.MinimumNArgs(1),
		RunE:  createSeasonCommand.createSeason,
	}

	createSeasonCommand.command = newCommand(cc)

	return createSeasonCommand
}

func (createScreateSeasonCommand *createSeasonCommand) createSeason(cmd *cobra.Command, args []string) error {
	season, err := createScreateSeasonCommand.seasonsManager.CreateSeason(args[0])
	if err != nil {
		return err
	}

	cli.Print(fmt.Sprintf("Season %s created", season.Name))

	return nil
}

type getSeasonsCommand struct {
	command
	seasonsManager seasons.Manager
}

func NewGetSeasonsCommand(seasonsManager seasons.Manager) *getSeasonsCommand {
	getSeasonsCommand := &getSeasonsCommand{seasonsManager: seasonsManager}

	cc := &cobra.Command{
		Use:   "list",
		Short: "List all seasons",
		Long:  "List all seasons with name and status",
		RunE:  getSeasonsCommand.getSeasons,
	}

	getSeasonsCommand.command = newCommand(cc)

	return getSeasonsCommand
}

func (getSeasonsCommand *getSeasonsCommand) getSeasons(cmd *cobra.Command, args []string) error {
	seasons, err := getSeasonsCommand.seasonsManager.GetSeasons()
	if err != nil {
		return err
	}

	seasonsTable := [][]string{}
	for _, season := range seasons {
		active := ""
		if season.Active {
			active = "Active"
		}
		seasonsTable = append(seasonsTable, []string{season.Name, active})
	}

	cli.PrintTable([]string{"Name", "Active"}, seasonsTable)

	return nil
}

type activateSeasonCommand struct {
	command
	seasonsManager seasons.Manager
}

func NewActivateSeasonCommand(seasonsManager seasons.Manager) *activateSeasonCommand {
	activateSeasonCommand := activateSeasonCommand{seasonsManager: seasonsManager}

	cc := &cobra.Command{
		Use:   "activate [name]",
		Short: "Activates an inactive season",
		Long:  "Activates the given inactive season",
		Args:  cobra.MinimumNArgs(1),
		RunE:  activateSeasonCommand.activateSeason,
	}

	activateSeasonCommand.command = newCommand(cc)

	return &activateSeasonCommand
}

func (activateSeasonCommand *activateSeasonCommand) activateSeason(cmd *cobra.Command, args []string) error {
	season, err := activateSeasonCommand.seasonsManager.ActivateSeason(args[0])
	if err != nil {
		return err
	}

	cli.Print(fmt.Sprintf("Season %s activated", season.Name))

	return nil
}

type getTableCommand struct {
	command
	gamesManager   games.Manager
	seasonsManager seasons.Manager
}

func NewGetTableCommand(
	gamesManager games.Manager,
	seasonsManager seasons.Manager,
) *getTableCommand {
	getTableCommand := &getTableCommand{
		gamesManager:   gamesManager,
		seasonsManager: seasonsManager,
	}

	cc := &cobra.Command{
		Use:   "table",
		Short: "Get the seasons table",
		Long:  "Get the seasons table. If no season name is provided, the active season will be used.",
		RunE:  getTableCommand.getTable,
	}

	cc.Flags().StringP("sort", "s", "", "Table sort by")

	getTableCommand.command = newCommand(cc)

	return getTableCommand
}

func (tableCommand *getTableCommand) getTable(cmd *cobra.Command, args []string) error {
	season, err := tableCommand.getSeason(args)
	if err != nil {
		return err
	}

	sortName, err := tableCommand.getSortName()
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

func (command *getTableCommand) getSortName() (string, error) {
	sortName, err := command.cc.Flags().GetString("sort")
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
