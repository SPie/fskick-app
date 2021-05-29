package cli

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/spie/fskick/players"
)

func Print(output string) {
	fmt.Println(fmt.Sprintf("%s\n", output))
}

func PrintTable(head []string, entries [][]string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	if len(head) > 0 {
		t.AppendHeader(transformHead(head))
	}
	t.AppendRows(transformEntries(entries))

	t.Render()
}

func transformHead(head []string) table.Row {
	headRow := table.Row{}
	for _, headString := range head {
		headRow = append(headRow, headString)
	}

	return headRow
}

func transformEntries(entries [][]string) []table.Row {
	rows := []table.Row{}
	for _, entry := range entries {
		row := table.Row{}
		for _, col := range entry {
			row = append(row, col)
		}

		rows = append(rows, row)
	}

	return rows
}

func CreateTableHead(gamesCount int, playerStats *[]players.PlayerStats) []string {
	return []string{
		fmt.Sprintf("Position (%d)", len(*playerStats)),
		"Name",
		"Points Ratio",
		"Points",
		"Wins",
		fmt.Sprintf("Games (%d)", gamesCount),
		"Win Ratio",
		"Games Ratio",
	}
}

func CreateTableEntries(gamesCount int, playerStats *[]players.PlayerStats) [][]string {
	tableEntries := make([][]string, len(*playerStats))
	for i, playerStats := range *playerStats {
		tableEntries[i] = []string{
			fmt.Sprint(playerStats.Position),
			playerStats.Name,
			fmt.Sprintf("%0.2f", playerStats.PointsRatio),
			fmt.Sprint(playerStats.Points),
			fmt.Sprint(playerStats.Wins),
			fmt.Sprint(playerStats.Games),
			fmt.Sprintf("%0.2f", (float32(playerStats.Wins) / float32(playerStats.Games))),
			fmt.Sprintf("%0.2f", (float32(playerStats.Games) / float32(gamesCount))),
		}
	}

	return tableEntries
}
