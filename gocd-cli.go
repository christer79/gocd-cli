package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"

	"github.com/christer79/gocd-api"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("example")

//Selection is the basis of whcih pipeline/stage and Job to operate on
type Selection struct {
	Pipeline string
	Stage    string
	Job      string
	ID       string
}

func main() {
	var argHost string
	flag.StringVar(&argHost, "host", "http://localhost:8153/go/", "Set hostname")

	var argPipelineGroup string
	flag.StringVar(&argPipelineGroup, "pipeline-group", "", "Specify name of a Pipeline group")

	var argPipeline string
	flag.StringVar(&argPipeline, "pipeline", "", "Specify name of Pipeline")

	var argPipelineCount int
	flag.IntVar(&argPipelineCount, "pipeline-count", 0, "Specify name of Pipeline")

	var argStage string
	flag.StringVar(&argStage, "stage", "", "Specify name of Stage")

	var argStageCount int
	flag.IntVar(&argStageCount, "stage-count", 0, "Specify name of Stage")

	var argJob string
	flag.StringVar(&argJob, "job", "", "Specify name of Job")

	var argUserName string
	flag.StringVar(&argUserName, "username", "", "Specify a username for authentication")

	var argPassword string
	flag.StringVar(&argPassword, "password", "", "Specify a password for authentication")

	var argPath string
	flag.StringVar(&argPath, "path", "", "Specify a path to a file when fetching artifacts")

	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to configuration file")

	flag.Parse()

	coordinate := gocdapi.Coordinate{PipelineGroup: argPipelineGroup, PipelineName: argPipeline, PipelineCount: argPipelineCount, StageName: argStage, StageCount: argStageCount, JobName: argJob, FilePath: argPath}

	fmt.Println("host: " + argHost)

	client := gocdapi.NewClient(&http.Client{})
	client.BaseURL, _ = url.Parse(argHost)

	if argUserName != "" && argPassword != "" {
		client.UserName = &argUserName
		client.Password = &argPassword
	}

	switch flag.Arg(0) {
	case "group":
		pipelinegroups, _, _ := client.PipelineGroup.Get()
		fmt.Println(pipelinegroups)
	case "agent":
		agents, _, _ := client.Agents.GetAll()
		fmt.Println(agents)
		/*
			agent, _, _ := client.Agents.Get(agents.Embedded.Agents[1].UUID)
			fmt.Println(agent)
			fmt.Println(agent.Resources)
		*/
	case "user":
		users, _, _ := client.Users.GetAll()
		fmt.Println(users)
	case "artifacts":
		artifacts, _, _ := client.Artifacts.GetAllArtifacts(coordinate)
		fmt.Println(artifacts)
		/*
			artifact, _, _ := client.Artifacts.GetArtifactFile(coordinate)
			fmt.Println(artifact)
			coordinate.FilePath = "cruise-output"
			directory, _, _ := client.Artifacts.GetArtifactDirectory(coordinate)
			fmt.Println(directory)
		*/
	case "pipeline":
		switch flag.Arg(1) {
		case "instance":
			pipelineinstance, _, _ := client.Pipelines.GetInstance(coordinate)
			fmt.Println(pipelineinstance)
		case "history":

			pipelinehistory, _, _ := client.Pipelines.GetHistory(coordinate)
			fmt.Println(pipelinehistory)
		case "status":
			pipelinestatus, _, _ := client.Pipelines.GetStatus(coordinate)
			fmt.Println(pipelinestatus)
		}

	}

}
