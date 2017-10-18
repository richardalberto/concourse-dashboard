package concourse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Client interface {
	GetToken(team, username, password string) (string, error)
	GetPipelines(team, token string) ([]Pipeline, error)
	GetJobs(team, pipeline, token string) ([]Job, error)
}

type client struct {
	url  string
	path string

	client *http.Client
}

// NewClient creates a new concourse client
func NewClient(url, path string) Client {
	return &client{
		url:    url,
		path:   path,
		client: &http.Client{},
	}
}

func (c *client) GetToken(team, username, password string) (string, error) {
	url := fmt.Sprintf("%s/auth/token", c.teamURL(team))

	log.Debugf("Requesting token from %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(username, password)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Server replied with status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var t token
	if err := json.Unmarshal(body, &t); err != nil {
		return "", err
	}

	return t.Value, nil
}

func (c *client) GetPipelines(team, token string) ([]Pipeline, error) {
	url := fmt.Sprintf("%s/pipelines", c.teamURL(team))

	body, err := c.get(url, token)
	if err != nil {
		return nil, err
	}

	var pipelines []Pipeline
	if err := json.Unmarshal(body, &pipelines); err != nil {
		return nil, err
	}

	return pipelines, nil
}

func (c *client) GetJobs(team, pipeline, token string) ([]Job, error) {
	url := fmt.Sprintf("%s/pipelines/%s/jobs", c.teamURL(team), pipeline)

	body, err := c.get(url, token)
	if err != nil {
		return nil, err
	}

	var jobs []Job
	if err := json.Unmarshal(body, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (c *client) get(url, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("cookie", fmt.Sprintf("ATC-Authorization=Bearer %s", token))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Server replied with status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return body, nil
}

func (c *client) teamURL(team string) string {
	url := fmt.Sprintf("%s%s/teams/%s", c.url, c.path, team)
	return strings.TrimRight(url, "/")
}
