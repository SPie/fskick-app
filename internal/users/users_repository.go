package users

import (
	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/players"
)

type User struct {
    players.Player
    Email string `json:"email"`
    Password string `json:"-"`
}

type UsersRepository struct {
    dbHandler db.Handler
}

func NewUsersRepository(dbHandler db.Handler) UsersRepository {
    return UsersRepository{dbHandler: dbHandler}
}
