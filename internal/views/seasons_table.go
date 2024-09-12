package views

import (
	"context"
	"io"

	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/seasons"
	"github.com/spie/fskick/internal/templates"
)

type SeasonsTable struct{}

func NewSeasonTable() SeasonsTable {
	return SeasonsTable{}
}

func (view SeasonsTable) Render(
	seasons []seasons.Season,
	activeSeason seasons.Season,
	playerStats []games.PlayerStats,
	gamesCount int,
	ctx context.Context,
	w io.Writer,
) error {
	return templates.SeasonsTable(seasons, activeSeason, playerStats, gamesCount).Render(ctx, w)
}
