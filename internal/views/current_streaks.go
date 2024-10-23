package views

import (
	"context"
	"io"

	"github.com/spie/fskick/internal/streaks"
	"github.com/spie/fskick/internal/templates"
)

type CurrentStreaks struct{}

func NewCurrentStreaks() CurrentStreaks {
	return CurrentStreaks{}
}

func (view CurrentStreaks) Render(currentStreaks []streaks.Streak, ctx context.Context, w io.Writer) error {
	return templates.CurrentStreaks(currentStreaks[:10]).Render(ctx, w)
}
