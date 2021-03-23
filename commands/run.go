package commands

import (
	githubcommands "github.com/mms-gianni/git-action/app"
	"gopkg.in/ukautz/clif.v1"
)

func cmdRun() *clif.Command {
	cb := func(c *clif.Command, out clif.Output, in clif.Input) {
		workflow := githubcommands.RunAction(c, out, in)

		if c.Option("follow").Bool() {
			rundetails := githubcommands.FollowAction(c, out, workflow)

			githubcommands.ShowLog(c, out, rundetails)
		}
	}

	return clif.NewCommand("run", "Run a workflow", cb).
		NewArgument("action", "Name of the workflow to run", "", false, true).
		NewFlag("follow", "f", "Follow the triggered Workflow", false)
}

func init() {
	Commands = append(Commands, cmdRun)
}
