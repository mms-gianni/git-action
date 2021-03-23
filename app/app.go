package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

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

func RunAction(c *clif.Command, out clif.Output, in clif.Input) github.Workflow {
	client := login(c)
	o = out
	var workflow *github.Workflow

	_, repo := GetGitdir()
	if repo != nil {
		repodetails := getRepodetails(repo)
		workflow = selectAction(client, in, c.Argument("action").String(), repodetails)

		_, err := client.Actions.CreateWorkflowDispatchEventByID(ctx, repodetails.owner, repodetails.name, workflow.GetID(), github.CreateWorkflowDispatchEventRequest{
			Ref: "master",
		})

		if err == nil {
			out.Printf("\n  <ok> started Action \"%s\"\n\n", workflow.GetName())
		}

		return *workflow
	}

	return *workflow

}

func FollowAction(c *clif.Command, out clif.Output, workflow github.Workflow) *github.WorkflowRun {

	client := login(c)
	time.Sleep(2 * time.Second)

	_, repo := GetGitdir()
	repodetails := getRepodetails(repo)

	workflowRuns, _, _ := client.Actions.ListWorkflowRunsByID(ctx, repodetails.owner, repodetails.name, workflow.GetID(), nil)

	for {
		for count, run := range workflowRuns.WorkflowRuns {
			if count >= 1 {
				break
			}

			runDetails, _, _ := client.Actions.GetWorkflowRunByID(ctx, repodetails.owner, repodetails.name, run.GetID())

			if runDetails.GetStatus() != "completed" {
				out.Printf("\r  <%s> %s \"%s\"", runDetails.GetStatus(), runDetails.GetStatus(), workflow.GetName())
			} else {
				var icon string
				if runDetails.GetConclusion() == "success" {
					icon = "ok"
				} else {
					icon = "failure"
				}
				out.Printf("\r  <%s> %s \"%s\"                           \n\n", icon, runDetails.GetConclusion(), workflow.GetName())
				return runDetails
			}

		}
		time.Sleep(2 * time.Second)
	}
}

func ShowLog(c *clif.Command, out clif.Output, run *github.WorkflowRun) {

	client := login(c)
	_, repo := GetGitdir()
	repodetails := getRepodetails(repo)
	jobs, _, _ := client.Actions.ListWorkflowJobs(ctx, repodetails.owner, repodetails.name, run.GetID(), nil)

	for _, job := range jobs.Jobs {
		logurl, _, _ := client.Actions.GetWorkflowJobLogs(ctx, repodetails.owner, repodetails.name, job.GetID(), true)

		req, _ := client.NewRequest("GET", logurl.String(), nil)
		var data string
		resp, _ := client.Do(ctx, req, data)

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Print(string(body))

	}
}
