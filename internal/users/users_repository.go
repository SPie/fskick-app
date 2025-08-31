package users

import (
	"fmt"
	"time"

	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/players"
)

type User struct {
	players.Player
	Email    string `json:"email"`
	Password string `json:"-"`
}

type UsersRepository struct {
	conn db.Connection
}

func NewUsersRepository(conn db.Connection) UsersRepository {
	return UsersRepository{conn: conn}
}

func (repo UsersRepository) CreateUser(user *User) error {
	_, err := repo.conn.Exec(
		"UPDATE players SET email = ?, password = ?, updated_at = ? WHERE id = ?",
		user.Email,
		user.Password,
		time.Now(),
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("Error update player to create user: %w", err)
	}

	return nil
}
