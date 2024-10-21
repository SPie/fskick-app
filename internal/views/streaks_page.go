package views

import (
	"context"
	"io"

	"github.com/spie/fskick/internal/streaks"
	"github.com/spie/fskick/internal/templates"
)

type StreaksPage struct{}

func NewStreaksPage() StreaksPage {
	return StreaksPage{}
}

func (view StreaksPage) Render(
	longestWiningStreak streaks.Streak,
	longestLosingStreak streaks.Streak,
	currentStreaks []streaks.Streak,
	ctx context.Context,
	w io.Writer,
) error {
	return templates.StreaksPage(
		longestWiningStreak,
		longestLosingStreak,
		currentStreaks[:10],
	).Render(ctx, w)
}
