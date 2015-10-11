package main

import (
	"flag"
	"fmt"
	"net/http"

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
	flag.StringVar(&argHost, "host", "localhost:8153", "Set hostname")

	var argPipeline string
	flag.StringVar(&argPipeline, "pipeline", "", "Specify name of Pipeline")

	var argStage string
	flag.StringVar(&argStage, "stage", "", "Specify name of Stage")

	var argJob string
	flag.StringVar(&argJob, "job", "", "Specify name of Job")

	flag.Parse()

	client := gocdapi.NewClient(&http.Client{})
	fmt.Println("Agent: GETALL")
	agents, _, _ := client.Agents.GetAll()
	fmt.Println(agents)
	fmt.Println(agents[1].Os)
	fmt.Println("Agent: GET UUID: " + agents[1].UUID)
	agent, _, _ := client.Agents.Get(agents[1].UUID)
	fmt.Println(agent)
	fmt.Println(agent.Resources)
	/*
		fmt.Println("Users: GETALL")
		users, _, _ := client.Users.GetAll()
		fmt.Println(users)

		fmt.Println("PipelineGroup: GET")
		pipelinegroups, _, _ := client.PipelineGroup.Get()
		fmt.Println(pipelinegroups)
	*/
}
