package users

import (
	"fmt"

	"github.com/spie/fskick/internal/players"
)

type playersManager interface {
    GetPlayerByName(name string) (players.Player, error)
}

type usersRepository interface {
    CreateUser(user *User) error
}

type passwordService interface {
    HashPassword(password []byte) []byte
}

type Manager struct {
    usersRepository usersRepository
    playersManager playersManager
    passwordService passwordService
}

func NewManager(
    usersRepository usersRepository,
    playersManager playersManager,
    paspasswordService passwordService,
) Manager {
    return Manager{
        usersRepository: usersRepository,
        playersManager: playersManager,
    }
}

func (manager Manager) CreateUserFromPlayer(
    playerName string,
    email string,
    plaintextPassword string,
) (User, error) {
    player, err := manager.playersManager.GetPlayerByName(playerName)
    if err != nil {
        return User{}, fmt.Errorf("get player for CreateUserFromPlayer: %w", err)
    }

    user := User{
        Player: player,
        Email: email,
        Password: string(manager.passwordService.HashPassword([]byte(plaintextPassword))),
    }

    err = manager.usersRepository.CreateUser(&user)
    if err != nil {
        return User{}, fmt.Errorf("store user for CreateUserFromPlayer: %w", err)
    }

    return user, nil
}
