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

var HOST = "http://localhost:8153"

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

	//coordinate := gocdapi.Coordinate{PipelineName: "AlaskaMultiNodeLassiBlue", PipelineCount: 4, StageName: "BootStrap", StageCount: 1, JobName: "BootStrap"}

	coordinate := gocdapi.Coordinate{PipelineName: argPipeline, PipelineCount: argPipelineCount, StageName: argStage, StageCount: argStageCount, JobName: argJob, FilePath: argPath}

	fmt.Println("host: " + argHost)

	client := gocdapi.NewClient(&http.Client{})
	client.BaseURL, _ = url.Parse(argHost)
	if argUserName != "" && argPassword != "" {
		client.UserName = &argUserName
		client.Password = &argPassword
	}

	fmt.Println("Agent: GETALL")
	agents, _, _ := client.Agents.GetAll()
	fmt.Println(agents)
	fmt.Println("Agent: GET UUID: " + agents.Embedded.Agents[1].UUID)
	agent, _, _ := client.Agents.Get(agents.Embedded.Agents[1].UUID)
	fmt.Println(agent)
	fmt.Println(agent.Resources)
	fmt.Println("PipelineGroup: GET")
	pipelinegroups, _, _ := client.PipelineGroup.Get()
	fmt.Println(pipelinegroups)

	fmt.Println("Users: GETALL")
	users, _, _ := client.Users.GetAll()
	fmt.Println(users)
	/*
		userName := "patrik.berglund"
		fmt.Println("Users: GET: " + userName)
		user, _, _ := client.Users.Get(userName)
		fmt.Println(user)
	*/

	fmt.Println("Artifacts: GETALL")
	//curl 'http://localhost:8153/go/files/AlaskaMultiNodeLassiBlue/4/BootStrap/1/BootStrap.json'       -u 'christer.eriksson:p4rqBE24'

	artifacts, _, _ := client.Artifacts.GetAllArtifacts(coordinate)
	fmt.Println(artifacts)

	artifact, _, _ := client.Artifacts.GetArtifactFile(coordinate)
	fmt.Println(artifact)

	coordinate.FilePath = "cruise-output"
	directory, _, _ := client.Artifacts.GetArtifactDirectory(coordinate)
	fmt.Println(directory)

	pipelineinstance, _, _ := client.Pipelines.GetInstance(coordinate)
	fmt.Println(pipelineinstance)

	pipelinehistory, _, _ := client.Pipelines.GetHistory(coordinate)
	fmt.Println(pipelinehistory)

	pipelinestatus, _, _ := client.Pipelines.GetStatus(coordinate)
	fmt.Println(pipelinestatus)
}
