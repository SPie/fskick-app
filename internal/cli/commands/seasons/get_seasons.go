package seasons

import (
	"github.com/spf13/cobra"

	"github.com/spie/fskick/internal/cli"
	"github.com/spie/fskick/internal/seasons"
)

type getSeasonsCommand struct {
	cc           *cobra.Command
	seasonsManager seasons.Manager
}

func newGetSeasonsCommand(seasonsManager seasons.Manager) *getSeasonsCommand {
	getSeasonsCommand := &getSeasonsCommand{seasonsManager: seasonsManager}

	cc := &cobra.Command{
		Use:   "list",
		Short: "List all seasons",
		Long:  "List all seasons with name and status",
		RunE:  getSeasonsCommand.getSeasons,
	}

	getSeasonsCommand.cc = cc

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
