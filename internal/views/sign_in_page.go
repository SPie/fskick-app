package views

import (
	"context"
	"io"

	"github.com/spie/fskick/internal/templates"
)

type SignInPage struct{}

func NewSignInPage() SignInPage {
	return SignInPage{}
}

func (view SignInPage) Render(email string, ctx context.Context, w io.Writer) error {
	return templates.SignInPage(email).Render(ctx, w)
}
