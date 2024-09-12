package views

import (
	"context"
	"fmt"
	"io"

	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/templates/components"
)

type FavoriteTeamUpdate struct {}

func NewFavoriteTeamUpdate() FavoriteTeamUpdate {
    return FavoriteTeamUpdate{}
}

func (view FavoriteTeamUpdate) Render(
    playerStats []games.PlayerStats,
    gamesCount int,
    playerUuid string,
    ctx context.Context,
    w io.Writer,
) error {
    options := components.TableHtmxOptions{
	Endpoint: fmt.Sprintf("/table/players/%s/team", playerUuid),
    }

    return components.PlayerStatsTable(playerStats[:5], gamesCount, options).Render(ctx, w)
}
