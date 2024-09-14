package views

import (
	"context"
	"fmt"
	"io"

	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/templates/components"
)

type FavoriteTeamUpdate struct{}

func NewFavoriteTeamUpdate() FavoriteTeamUpdate {
	return FavoriteTeamUpdate{}
}

func (view FavoriteTeamUpdate) Render(
	playerStats []games.PlayerStats,
	gamesCount int,
	playerUuid string,
	sort string,
	ctx context.Context,
	w io.Writer,
) error {
	options := components.TableHtmxOptions{
		Endpoint: fmt.Sprintf("/table/players/%s/team", playerUuid),
	}

	return components.PlayerStatsTable(
		getFavoriteTeamOf5(playerStats),
		gamesCount,
		sort,
		options,
	).Render(ctx, w)
}

func getFavoriteTeamOf5(playerStats []games.PlayerStats) []games.PlayerStats {
	if len(playerStats) <= 5 {
		return playerStats
	}

	return playerStats[:5]
}
