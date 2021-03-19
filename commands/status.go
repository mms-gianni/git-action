package commands

import (
	githubcommands "github.com/mms-gianni/git-action/app"
	"gopkg.in/ukautz/clif.v1"
)

func cmdStatus() *clif.Command {
	cb := func(c *clif.Command, out clif.Output) {
		githubcommands.GetStatus(c, out)
	}

	return clif.NewCommand("status", "List actions", cb)
}

func init() {
	Commands = append(Commands, cmdStatus)
}
