package views

import (
    "context"
    "io"

    "github.com/spie/fskick/internal/templates"
)

type ImprintView struct {}

func NewImprintView() ImprintView {
    return ImprintView{}
}

func (view ImprintView) Render(imprintText string, ctx context.Context, w io.Writer) error {
    return templates.Imprint(imprintText).Render(ctx, w)
}
