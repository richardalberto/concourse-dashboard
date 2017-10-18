package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/richardalberto/concourse-dashboard/pkg/concourse"
	"github.com/richardalberto/concourse-dashboard/pkg/config"
)

const (
	StatusSucceeded = "succeeded"
	StatusFailed    = "failed"
	StatusErrored   = "errored"
)

type team struct {
	Name      string
	Pipelines []pipelineStatus
}

type pipelineStatus struct {
	Name   string
	Status string
}

var (
	client concourse.Client
	conf   *config.Config
)

func main() {
	conf = config.Load("config")

	log.SetLevel(log.DebugLevel)

	client = concourse.NewClient(conf.Concourse.URL, conf.Concourse.APIPath)

	r := gin.Default()
	r.LoadHTMLGlob("views/*")
	r.Static("/static", "resources")
	r.GET("/", func(c *gin.Context) {
		teams, err := getTeams()
		if err != nil {
			log.Errorf("Couldn't retrieve teams information, %s", err)
		}

		c.HTML(http.StatusOK, "overview.tmpl", gin.H{
			"Title": "Concourse Pipelines Dashboard",
			"Teams": teams,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func getTeams() ([]team, error) {
	teams := make([]team, 0)

	if len(conf.Concourse.Teams) == 0 {
		return nil, fmt.Errorf("No teams found on the configuration")
	}

	for _, t := range conf.Concourse.Teams {
		log.Infof("Fetching pipelines status for team=%s", t.Name)

		log.Infof("Authenticating with username=%s password=%s", t.Username, t.Password)
		token, err := client.GetToken(t.Name, t.Username, t.Password)
		if err != nil {
			log.Errorf("Couldn't get token for team=%s, %s", t.Name, err)
			continue
		}
		log.Debugf("Received token=%s", token)

		pipelines, err := getPipelines(t.Name, token)
		if err != nil {
			return nil, err
		}

		team := team{
			Name:      t.Name,
			Pipelines: pipelines,
		}

		teams = append(teams, team)
	}

	return teams, nil
}

func getPipelines(team, token string) ([]pipelineStatus, error) {
	pipelines, err := client.GetPipelines(team, token)
	if err != nil {
		return nil, err
	}

	statuses := make([]pipelineStatus, 0)
	for _, p := range pipelines {
		jobs, err := client.GetJobs(team, p.Name, token)
		if err != nil {
			log.Errorf("An error ocurred when trying to retrieve jobs for pipeline=%s on team=%s", p.Name, team)
			continue
		}

		status := pipelineStatus{
			Name:   p.Name,
			Status: statusFromJobs(jobs),
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

func statusFromJobs(jobs []concourse.Job) string {
	for _, j := range jobs {
		if j.FinishedBuild != nil {
			if j.FinishedBuild.Status == StatusErrored || j.FinishedBuild.Status == StatusFailed {
				return j.FinishedBuild.Status
			}
		}
	}

	return StatusSucceeded
}
