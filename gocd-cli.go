package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	logging "github.com/op/go-logging"
	"github.com/parnurzeal/gorequest"
)

var log = logging.MustGetLogger("example")

HOST="http://localhost:8153"

//Selection is the basis of whcih pipeline/stage and Job to operate on
type Selection struct {
	Pipeline string
	Stage    string
	Job      string
	ID       string
}

func makeRequest(uri string) []byte {
	log.Debug("URI:" + uri)

	request := gorequest.New()
	request.SetBasicAuth("<name>", "<pass>")
	response, body, err := request.Get(uri).End()

	if err != nil {
		panic(err)
	}

	log.Debug(string(response.Status))

	data := []byte(body)
	return data
}

//Agent represents a Agent index as in the response when listing agents
type Agent struct {
	Os           string   `json:"os"`
	Environments []string `json:"environments"`
	UUID         string   `json:"uuid"`
	AgentName    string   `json:"agent_name"`
	FreeSpace    string   `json:"free_space"`
	Resources    []string `json:"resources"`
	Sandbox      string   `json:"sandbox"`
	Status       string   `json:"status"`
	BuildLocator string   `json:"build_locator"`
	IPAddress    string   `json:"ip_address"`
}

func listAgents() {
	requestString := HOST+"/go/api/agents"
	response := makeRequest(requestString)
	var agents []Agent
	json.Unmarshal(response, &agents)

	log.Debug("Number of agents:  %d", len(agents))

	fmt.Printf("| %-45s | %-11s | %-11s | \n", "NAME", "STATUS", "RESOURCES")
	for _, agent := range agents {
		fmt.Printf("| %-45s | %-11s | %-11s | \n", agent.AgentName, agent.Status, strings.Join(agent.Resources, ", "))
	}
}

//Stage holds informaton about a Stage of a Pipeline
type Stage struct {
	Name string `json:"name"`
}

//Material represents information about material such as Git or Mercurial repo
type Material struct {
	Fingerprint string `json:"fingerprint"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

//Pipeline represents information about a single pipeline
type Pipeline struct {
	Name      string     `json:"name"`
	Label     string     `json:"label"`
	Materials []Material `json:"materials"`
	Stages    []Stage    `json:"stages"`
}

//PipelineGroup holds informationa bout a single pipeline group
type PipelineGroup struct {
	Name      string     `json:"name"`
	Pipelines []Pipeline `json:"pipelines"`
}

func listGroups() {
	requestString := HOST + "/go/api/config/pipeline_groups"
	response := makeRequest(requestString)

	var groups []PipelineGroup
	json.Unmarshal(response, &groups)
	log.Debug("Number of groups:  %d", len(groups))

	for _, group := range groups {
		fmt.Printf("G: %-20s \n", group.Name)
		for _, pipeline := range group.Pipelines {
			fmt.Printf("  P: %-20s | \n", pipeline.Name)
			for _, material := range pipeline.Materials {
				fmt.Printf("   |M: %-4s | %-70s \n", material.Type, material.Description)

			}
			for _, stage := range pipeline.Stages {
				fmt.Printf("   |S: %-20s \n", stage.Name)
			}
		}

	}
}

//Job holds information about a Job as presented in JobHistory
type Job struct {
	AgentUUID           string   `json:"agent_uuid"`
	Name                string   `json:"name"`
	JobStateTransitions []string `json:"job_state_transitions"`
	ScheduledDate       string   `json:"scheduled_date"`
	OriginalJobID       string   `json:"original_job_id"`
	PipelineCounter     string   `json:"pipeline_counter"`
	Rerun               string   `json:"rerun"`
	PipelineName        string   `json:"pipeline_name"`
	Result              string   `json:"result"`
	State               string   `json:"state"`
	ID                  string   `json:"id"`
	StageCounter        string   `json:"stage_counter"`
	StageName           string   `json:"stage_name"`
}

// JobHistory hold a list of jobs as present when querying ajobs history
type JobHistory struct {
	Jobs []Job `json:"jobs"`
}

func getJobHistory(selection Selection) {

	requestString := fmt.Sprintf("%s/go/api/jobs/%s/%s/%s/history)", HOST, selection.Pipeline, selection.Stage, selection.Job)

	response := makeRequest(requestString)
	var jobHistory JobHistory
	log.Debug(string(response))
	json.Unmarshal(response, &jobHistory)

	log.Debug("Number of jobs in history:  %d", len(jobHistory.Jobs))

	fmt.Printf("| %-14s -> %-11s -> %-11s \n", selection.Pipeline, selection.Stage, selection.Job)

	fmt.Printf("| %-14s | %-11s | %-11s | %-20s | \n", "DATE", "STATE", "RESULT", "AGENT UUID")
	for _, job := range jobHistory.Jobs {
		fmt.Printf("| %-14s | %-11s | %-11s | %-20s | \n", job.ScheduledDate, job.State, job.Result, job.AgentUUID)
	}
}

func main() {
	var argHost string
	flag.StringVar(&argHost, "host", "localhost:8153", "Set hostname")

	var argPipeline string
	flag.StringVar(&argPipeline, "pipeline", "", "Specify name of Pipeline")

	var argStage string
	flag.StringVar(&argStage, "stage", "", "Specify name of Stage")

	var argJob string
	flag.StringVar(&argJob, "job", "", "Specify name of Job")

	flag.Parse()
	selection := Selection{Pipeline: argPipeline, Stage: argStage, Job: argJob}
	//listAgents()
	listGroups()
	log.Debug("Pipeline %s:%s", argPipeline, selection.Pipeline)
	if selection.Pipeline != "" {
		getJobHistory(selection)
	}

}
