package commands

import (
	githubcommands "github.com/mms-gianni/git-action/app"
	"gopkg.in/ukautz/clif.v1"
)

func cmdEnable() *clif.Command {
	cb := func(c *clif.Command, out clif.Output, in clif.Input) {
		githubcommands.Enable(c, out, in, true)
	}

	return clif.NewCommand("enable", "Enable actions", cb).
		NewArgument("action", "Name of the workflow to run", "", false, true)
}

func cmdDisable() *clif.Command {
	cb := func(c *clif.Command, out clif.Output, in clif.Input) {
		githubcommands.Enable(c, out, in, false)
	}

	return clif.NewCommand("disable", "Disable actions", cb).
		NewArgument("action", "Name of the workflow to run", "", false, true)
}

func init() {
	Commands = append(Commands, cmdEnable, cmdDisable)
}
