package team

type Team struct {
	Name       string `json:"name"` // Could split this into struct containing city + team, but let's keep it simple for now
	Conference string `json:"conference"`
	Division   string `json:"division"`
}
