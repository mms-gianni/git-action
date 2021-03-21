package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
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

func GetLog(c *clif.Command, out clif.Output) {
	client := login(c)
	o = out
	_, repo := GetGitdir()

	if repo != nil {
		repodetails := getRepodetails(repo)
		repoWorkflows, _, _ := client.Actions.ListWorkflows(ctx, repodetails.owner, repodetails.name, nil)

		for {
			for _, workflow := range repoWorkflows.Workflows {

				workflowRuns, _, _ := client.Actions.ListWorkflowRunsByID(ctx, repodetails.owner, repodetails.name, workflow.GetID(), nil)

				limit := 1
				for count, run := range workflowRuns.WorkflowRuns {
					if count >= limit {
						break
					}

					runDetails, _, _ := client.Actions.GetWorkflowRunByID(ctx, repodetails.owner, repodetails.name, run.GetID())

					jobs, _, _ := client.Actions.ListWorkflowJobs(ctx, repodetails.owner, repodetails.name, run.GetID(), nil)

					if runDetails.GetStatus() != "completed" {
						//out.Printf("following curl -i -u 'mms-gianni:%s' %s %s\n", os.Getenv("GITHUB_TOKEN"), run.GetLogsURL(), runDetails.GetHeadCommit().GetID())

						for _, job := range jobs.Jobs {

							/*
								url := "https://github.com/mms-gianni/git-action/commit/" + runDetails.GetHeadCommit().GetID() + "/checks/" + strconv.FormatInt(job.GetID(), 10) + "/live_logs"
								//url := "https://api.github.com/repos/mms-gianni/git-action/commit/" + runDetails.GetHeadCommit().GetID() + "/checks/" + strconv.FormatInt(job.GetID(), 10) + "/live_logs"

								out.Printf("%s\n", url)
										https://github.com/mms-gianni/git-action/commit/dc8bb9a93f3c759bd253238168a1541786ea7b05/checks/2157500334/live_logs
										https://github.com/mms-gianni/git-action/commit/dc8bb9a93f3c759bd253238168a1541786ea7b05/checks/2157500334/live_logs
									client := &http.Client{}
									req, _ := http.NewRequest("GET", url, nil)
									req.Header.Set("accept", "application/json")
									req.SetBasicAuth("mms-gianni", os.Getenv("GITHUB_TOKEN"))
									resp, _ := client.Do(req)

									defer resp.Body.Close()
									body, _ := ioutil.ReadAll(resp.Body)
									out.Printf(string(body))
									out.Printf("\n")

									u := fmt.Sprintf("repos/%v/%v/actions/jobs/%v", owner, repo, jobID)
							*/

							url := "https://api.github.com/repos/mms-gianni/git-action/actions/jobs/" + strconv.FormatInt(job.GetID(), 10) + "/logs"
							req, _ := client.NewRequest("GET", url, nil)
							resp, _ := client.Do(ctx, req, nil)

							defer resp.Body.Close()
							body, _ := ioutil.ReadAll(resp.Body)
							out.Printf(string(body))

							out.Printf("\n")

						}

					} else {
						//out.Printf("not running curl -i -u 'mms-gianni:%s' %s %s\n", os.Getenv("GITHUB_TOKEN"), run.GetLogsURL(), run.GetID())
						for _, job := range jobs.Jobs {
							/*
									out.Printf(job.GetCheckRunURL())
									out.Printf(job.GetHTMLURL())
									url := "https://github.com/mms-gianni/git-action/commit/" + runDetails.GetHeadCommit().GetID() + "/checks/" + strconv.FormatInt(job.GetID(), 10) + "/live_logs"
									//url := "https://api.github.com/repos/mms-gianni/git-action/commit/" + runDetails.GetHeadCommit().GetID() + "/checks/" + strconv.FormatInt(job.GetID(), 10) + "/live_logs"

									out.Printf("%s\n", url)

								url := job.GetCheckRunURL()
								url = "https://api.github.com/repos/mms-gianni/commits/" + runDetails.GetHeadCommit().GetID() + "/checks/" + strconv.FormatInt(job.GetID(), 10) + "/live_logs"
								url = "https://api.github.com/repos/mms-gianni/git-action/commit/" + runDetails.GetHeadCommit().GetID() + "/checks/" + strconv.FormatInt(job.GetID(), 10) + "/live_logs"
								url = "https://api.github.com/repos/mms-gianni/git-action/commits/" + runDetails.GetHeadCommit().GetID() + "/check-runs"
								//url = "https://github.com/mms-gianni/git-action/commit/" + runDetails.GetHeadCommit().GetID() + "/checks/" + strconv.FormatInt(job.GetID(), 10) + "/live_logs"

								//client := &http.Client{}
								///req, _ := http.NewRequest("GET", url, nil)
								req, _ := client.NewRequest("GET", url, nil)
								//req.Header.Set("accept", "application/json")
								//req.SetBasicAuth("mms-gianni", os.Getenv("GITHUB_TOKEN"))
								resp, _ := client.Do(ctx, req, nil)

								defer resp.Body.Close()
								body, _ := ioutil.ReadAll(resp.Body)
								out.Printf(string(body))

								out.Printf("\n")
							*/

							//logurl, _, _ := client.Actions.GetWorkflowJobLogs(ctx, repodetails.owner, repodetails.name, job.GetID(), true)
							//fmt.Println(logurl.String())
							fmt.Println(job.GetID())

							url := "https://api.github.com/repos/mms-gianni/git-action/actions/jobs/" + strconv.FormatInt(job.GetID(), 10) + "/logs"
							req, _ := client.NewRequest("GET", url, nil)
							resp, _ := client.Do(ctx, req, nil)

							defer resp.Body.Close()
							body, _ := ioutil.ReadAll(resp.Body)
							out.Printf(string(body))

							out.Printf("\n")
						}
					}
				}
			}
			time.Sleep(2 * time.Second)
		}
	}

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
				out.Printf("    #%d <%s> %s, created:%s, duration:%dsec %s\n", runDetails.GetRunNumber(), icon, runDetails.GetEvent(), runDetails.GetCreatedAt(), runUsage.GetRunDurationMS()/1000, runDetails.GetHeadCommit().GetMessage())

			}
			out.Printf("<reset>\n")
		}
	}
}
