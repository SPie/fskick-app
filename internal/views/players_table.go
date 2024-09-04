package views

import (
    "context"
    "io"

    "github.com/spie/fskick/internal/games"
    "github.com/spie/fskick/internal/templates"
)

type PlayersTable struct {}

func NewPlayersTable() PlayersTable {
    return PlayersTable{}
}

func (view PlayersTable) Render(
    playerStats []games.PlayerStats,
    gamesCount int,
    ctx context.Context,
    w io.Writer,
) error {
    return templates.PlayersTable(playerStats, gamesCount).Render(ctx, w)
}
