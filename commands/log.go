package commands

import (
	githubcommands "github.com/mms-gianni/git-action/app"
	"gopkg.in/ukautz/clif.v1"
)

func cmdLog() *clif.Command {
	cb := func(c *clif.Command, out clif.Output) {
		githubcommands.GetLog(c, out)
	}

	return clif.NewCommand("log", "Follow a log", cb)
}

func init() {
	Commands = append(Commands, cmdLog)
}
