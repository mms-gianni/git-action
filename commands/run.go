package commands

import (
	githubcommands "github.com/mms-gianni/git-action/app"
	"gopkg.in/ukautz/clif.v1"
)

func cmdRun() *clif.Command {
	cb := func(c *clif.Command, out clif.Output, in clif.Input) {
		githubcommands.RunAction(c, out, in)
	}

	return clif.NewCommand("run", "Run actions", cb).
		NewArgument("action", "Name of the action to run", "", false, true)
}

func init() {
	Commands = append(Commands, cmdRun)
}
