package app

import (
	"context"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
	"gopkg.in/ukautz/clif.v1"
)

var ctx = context.Background()

var o clif.Output
var i clif.Input

func login(c *clif.Command) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Option("githubtoken").String()},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

func GetStatus(c *clif.Command, out clif.Output) {
	client := login(c)
	o = out

	_, repo := GetGitdir()

	if repo != nil {
		repodetails := getRepodetails(repo)
		repoWorkflows, _, _ := client.Actions.ListWorkflows(ctx, repodetails.owner, repodetails.name, nil)

		for _, workflow := range repoWorkflows.Workflows {

			out.Printf(" <%s> %s\n", workflow.GetState(), workflow.GetName())

			workflowRuns, _, _ := client.Actions.ListWorkflowRunsByID(ctx, repodetails.owner, repodetails.name, workflow.GetID(), nil)

			limit := 5
			for count, run := range workflowRuns.WorkflowRuns {
				if count >= limit {
					break
				}

				runDetails, _, _ := client.Actions.GetWorkflowRunByID(ctx, repodetails.owner, repodetails.name, run.GetID())
				runUsage, _, _ := client.Actions.GetWorkflowRunUsageByID(ctx, repodetails.owner, repodetails.name, run.GetID())
				// possible status : queued, in_progress, completed
				// possible Conclusions: success, failure
				var icon string
				switch status := runDetails.GetStatus(); status {
				case "queued":
					icon = "queued"
				case "in_progress":
					icon = "in_progress"
				case "completed":
					if runDetails.GetConclusion() == "success" {
						icon = "ok"
					} else {
						icon = "failure"
					}
				}
				out.Printf("   #%d <%s> %s, created:%s, duration:%dsec \n", runDetails.GetRunNumber(), icon, runDetails.GetEvent(), runDetails.GetCreatedAt(), runUsage.GetRunDurationMS()/1000)
			}
			out.Printf("<reset>\n")
		}
	}
}

func printRunners(runners *github.Runners, title string) {

	if runners == nil {
		return
	}
	o.Printf("<important>%s<reset>\n", title)
	statusicon := ""
	for _, runner := range runners.Runners {

		if runner.GetStatus() == "offline" {
			statusicon = "offline"
		} else {
			if runner.GetBusy() {
				statusicon = "busy"
			} else {
				statusicon = "online"
			}
		}
		o.Printf("   <%s> %s\n", statusicon, runner.GetName())
	}
}
