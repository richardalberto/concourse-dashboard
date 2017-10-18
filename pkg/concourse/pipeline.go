package concourse

// Pipeline represents a concourse pipeline
type Pipeline struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Paused   bool   `json:"paused"`
	Public   bool   `json:"public"`
	TeamName string `json:"team_name"`
}
