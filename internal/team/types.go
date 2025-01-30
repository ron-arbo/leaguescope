package team

type Team struct {
	ID         string   `json:"id"`
	Name       TeamName `json:"name"`
	Conference string   `json:"conference"`
	Division   string   `json:"division"`
}

type TeamName struct {
	Location string
	Name     string
}

// ------------ ESPN Types ------------

// ESPNTeam represents a team in the standings. It contains metadata about the team
type ESPNTeam struct {
	ID               string `json:"id"`
	UID              string `json:"uid"`
	Location         string `json:"location"`
	Name             string `json:"name"`
	Abbreviation     string `json:"abbreviation"`
	DisplayName      string `json:"displayName"`
	ShortDisplayName string `json:"shortDisplayName"`
	IsActive         bool   `json:"isActive"`
	Logos            []ESPNLogo
}

// ESPNLogo describes a team logo
type ESPNLogo struct {
	Href        string   `json:"href"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Alt         string   `json:"alt"`
	Rel         []string `json:"rel"`
	LastUpdated string   `json:"lastUpdated"`
}
