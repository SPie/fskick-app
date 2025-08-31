package views

import (
	"context"
	"io"

	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/streaks"
	"github.com/spie/fskick/internal/templates"
)

type PlayerInfo struct{}

func NewPlayerInfo() PlayerInfo {
	return PlayerInfo{}
}

func (view PlayerInfo) Render(
	playerStats games.PlayerStats,
	gamesCount int,
	lastAttendances []games.Attendance,
	longestWinningStreak streaks.Streak,
	longestLosingStreak streaks.Streak,
	favoriteTeam []games.PlayerStats,
	favoriteOponents []games.PlayerStats,
	ctx context.Context,
	w io.Writer,
) error {
	return templates.Player(
		playerStats,
		gamesCount,
		lastAttendances,
		longestWinningStreak,
		longestLosingStreak,
		getFavoriteTeamOf5(favoriteTeam),
		getFavoriteTeamOf5(favoriteOponents),
	).Render(ctx, w)
}
