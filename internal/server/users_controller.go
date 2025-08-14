package server

import (
	"net/http"

	// "github.com/spie/fskick/internal/auth"
	// "github.com/spie/fskick/internal/users"
	"github.com/spie/fskick/internal/views"
)

// type usersManager interface {
// 	GetUserByEmail(email string) (users.User, error)
// }

// type authManager interface {
// 	Authenticate(user users.User, password string) (auth.Session, error)
// }

type UsersViews struct {
	SignInPage views.SignInPage
}

func NewUsersViews() UsersViews {
	return UsersViews{}
}

type UsersController struct {
	// usersManager usersManager
	// authManager authManager
	views UsersViews
}

func NewUsersController(
	// usersManager usersManager,
	// authManager authManager,
	usersViews UsersViews,
) UsersController {
	return UsersController{
		// usersManager: usersManager,
		// authManager: authManager,
		views: usersViews,
	}
}

func (usersController UsersController) SignInPage(res http.ResponseWriter, req *http.Request) {
	err := usersController.views.SignInPage.Render(
		req.URL.Query().Get("email"),
		req.Context(),
		res,
	)
	if err != nil {
		handleInternalServerError(res, err)
	}
}

func (usersController UsersController) DoSignIn(res http.ResponseWriter, req *http.Request) {
	// TODO
}
