package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

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

			limit := c.Option("limit").Int()
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
				out.Printf("    #%d <%s> %s, created:%s, duration:%dsec %s\n", runDetails.GetRunNumber(), icon, runDetails.GetEvent(), runDetails.GetCreatedAt(), runUsage.GetRunDurationMS()/1000, runDetails.GetHeadCommit().GetMessage())

			}
			out.Printf("<reset>\n")
		}
	}
}

func selectAction(client *github.Client, in clif.Input, preselectedWorkflow string, repodetails *repoDetails) *github.Workflow {

	choices := make(map[string]string)

	repoWorkflows, _, _ := client.Actions.ListWorkflows(ctx, repodetails.owner, repodetails.name, nil)
	for key, workflow := range repoWorkflows.Workflows {

		yamlFile, err := ioutil.ReadFile(workflow.GetPath())
		if err != nil {
			panic(err)
		}
		s := string(yamlFile)

		if strings.Contains(s, "workflow_dispatch") {
			choices[strconv.Itoa(key+1)] = workflow.GetName()
			if workflow.GetName() == preselectedWorkflow {
				return workflow
			}
		}
	}

	selectedNr, _ := strconv.Atoi(in.Choose("Where do you want to add a task?", choices))
	return repoWorkflows.Workflows[selectedNr-1]
}

func RunAction(c *clif.Command, out clif.Output, in clif.Input) {
	client := login(c)
	o = out

	_, repo := GetGitdir()

	if repo != nil {
		repodetails := getRepodetails(repo)
		workflow := selectAction(client, in, "", repodetails)

		fmt.Println(workflow.GetName())

		client.Actions.CreateWorkflowDispatchEventByID(ctx, repodetails.owner, repodetails.name, workflow.GetID(), github.CreateWorkflowDispatchEventRequest{
			Ref: "master",
		})
	}
	return

}
