package commands

import (
	"github.com/spf13/cobra"
)

type Command interface {
	AddCommand(command Command)
	getCommand() *cobra.Command
}

type rootCommand struct {
	cc *cobra.Command
}

func NewRootCommand() *rootCommand {
	rootCommand := rootCommand{cc: &cobra.Command{
		Use:   "fskick",
		Short: "FSKick CLI app",
		Long:  "CLI app for FSKick to create new players, seasons, games and show results and statistics",
	}}

	return &rootCommand
}

func (command *rootCommand) Execute() error {
	return command.cc.Execute()
}

func (command *rootCommand) AddCommand(c Command) {
	command.cc.AddCommand(c.getCommand())
}

func (command *rootCommand) getCommand() *cobra.Command {
	return command.cc
}
