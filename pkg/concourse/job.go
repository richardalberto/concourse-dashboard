package concourse

type Job struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	URL           string `json:"url"`
	NextBuild     string `json:"_"`
	FinishedBuild *Build `json:"finished_build"`
}

type Build struct {
	ID       int    `json:"id"`
	TeamName string `json:"team_name"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}
