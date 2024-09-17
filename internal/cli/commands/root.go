package commands

import (
	"github.com/spf13/cobra"
)

type Command interface {
	Execute() error
	AddCommand(command Command)
	getCommand() *cobra.Command
}

type command struct {
	cc *cobra.Command
}

func newCommand(cc *cobra.Command) command {
	return command{cc: cc}
}

func (cmd *command) Execute() error {
	return cmd.cc.Execute()
}

func (cmd *command) AddCommand(c Command) {
	cmd.cc.AddCommand(c.getCommand())
}

func (cmd *command) getCommand() *cobra.Command {
	return cmd.cc
}

type rootCommand struct {
	command
}

func NewRootCommand() *rootCommand {
	rootCommand := rootCommand{command: newCommand(&cobra.Command{
		Use:   "fskick",
		Short: "FSKick CLI app",
		Long:  "CLI app for FSKick to create new players, seasons, games and show results and statistics",
	})}

	return &rootCommand
}
