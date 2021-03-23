package main

import (
	"os"

	"github.com/mms-gianni/git-action/commands"
	"gopkg.in/ukautz/clif.v1"
)

func addDefaultOptions(cli *clif.Cli) {
	githubtoken := clif.NewOption("githubtoken", "t", "Private Github Token", "", true, false).
		SetEnv("GITHUB_TOKEN")

	githubusername := clif.NewOption("username", "u", "Github username", "", false, false).
		SetEnv("GITHUB_USERNAME")

	githubOrganisations := clif.NewOption("organisations", "o", "Github organisations (comma separated)", "", false, false).
		SetEnv("GITHUB_ORGANISATIONS")
	cli.AddDefaultOptions(githubtoken, githubusername, githubOrganisations)
}

func main() {
	cli := clif.New("git-action", "DEV-VERSION", "Manage your github actions with git cli")

	var OwnStyles = map[string]string{
		"error":             "\033[31;1m",
		"warn":              "\033[33m",
		"info":              "\033[0;97m",
		"success":           "\033[32m",
		"debug":             "\033[30;1m",
		"headline":          "\033[4;1m",
		"subline":           "\033[4m",
		"important":         "\033[47;30;1m",
		"query":             "\033[36m",
		"reset":             "\033[0m",
		"ok":                "\U00002705",
		"failure":           "\U0000274C",
		"in_progress":       "\U0001F525",
		"queued":            "\U0001F553",
		"disabled_manually": "\033[30;1m",
		"active":            "",
	}

	cli.SetOutput(clif.NewColorOutput(os.Stdout).SetFormatter(clif.NewDefaultFormatter(OwnStyles)))

	addDefaultOptions(cli)

	for _, cb := range commands.Commands {
		cli.Add(cb())
	}

	cli.Run()
}
