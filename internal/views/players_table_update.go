package views

import (
	"context"
	"fmt"
	"io"

	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/templates/components"
)

type PlayersTableUpdate struct{}

func NewPlayersTableUpdate() PlayersTableUpdate {
	return PlayersTableUpdate{}
}

func (view PlayersTableUpdate) Render(
	playerStats []games.PlayerStats,
	gamesCount int,
	playerUuid string,
	ctx context.Context,
	w io.Writer,
) error {
	playersTableEndpoint := "/table/players"
	if playerUuid != "" {
		playersTableEndpoint = fmt.Sprintf("%s/%s", playersTableEndpoint, playerUuid)
	}
	options := components.TableHtmxOptions{Endpoint: playersTableEndpoint}

	return components.PlayerStatsTable(playerStats, gamesCount, options).Render(ctx, w)
}
