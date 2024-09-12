package views

import (
	"context"
	"io"

	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/templates/components"
)

type SeasonsTableUpdate struct{}

func NewSeasonsTableUpdate() SeasonsTableUpdate {
	return SeasonsTableUpdate{}
}

func (view SeasonsTableUpdate) Render(
	playerStats []games.PlayerStats,
	gamesCount int,
	ctx context.Context,
	w io.Writer,
) error {
	options := components.TableHtmxOptions{
		Endpoint: "/table/seasons",
		Include:  "select[name='season']",
	}

	return components.PlayerStatsTable(playerStats, gamesCount, options).Render(ctx, w)
}
