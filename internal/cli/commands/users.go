package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spie/fskick/internal/cli"
	"github.com/spie/fskick/internal/users"
)

type usersCommand struct {
	command
}

func NewUsersCommand() *usersCommand {
	return &usersCommand{command: newCommand(&cobra.Command{
		Use:   "users",
		Short: "Commands to handle users",
		Long:  "All commands handling users, like creating new users",
	})}
}

type createUserFromPlayerCommand struct {
	command
	usersManager users.Manager
}

func NewCreateUserFromPlayerCommand(userManager users.Manager) *createUserFromPlayerCommand {
	createUserFromPlayerCommand := &createUserFromPlayerCommand{usersManager: userManager}

	cc := &cobra.Command{
		Use:   "player [name] [email] [password]",
		Short: "Creates a new user for an existing player",
		Long:  "Creates a new user with a given email and password for a player with the given name",
		Args:  cobra.ExactArgs(3),
		RunE:  createUserFromPlayerCommand.createUser,
	}

	createUserFromPlayerCommand.command = newCommand(cc)

	return createUserFromPlayerCommand
}

func (createUserFromPlayerCommand createUserFromPlayerCommand) createUser(
	cmd *cobra.Command,
	args []string,
) error {
	user, err := createUserFromPlayerCommand.usersManager.CreateUserFromPlayer(
		args[0],
		args[1],
		args[2],
	)
	if err != nil {
		return err
	}

	cli.Print(fmt.Sprintf("User %s for Player %s created\n", user.Email, args[0]))

	return nil
}
