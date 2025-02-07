package team

type Team struct {
	Name       TeamName `json:"name"`
	Conference string   `json:"conference"`
	Division   string   `json:"division"`
}

type TeamName struct {
	Location string
	Name     string
}
