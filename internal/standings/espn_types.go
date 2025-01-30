package standings

import (
	"nfl-app/internal/stats"
	"nfl-app/internal/team"
)

type ESPNStandings struct {
	UID          string           `json:"uid"`
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Abbreviation string           `json:"abbreviation"`
	ShortName    string           `json:"shortName"`
	Conferences  []ESPNConference `json:"children"`
	Links        []ESPNLink       `json:"links"`
	Seasons      []ESPNSeason     `json:"seasons"`
}

type ESPNConference struct {
	UID          string                  `json:"uid"`
	ID           string                  `json:"id"`
	Name         string                  `json:"name"`
	Abbreviation string                  `json:"abbreviation"`
	Standings    ESPNConferenceStandings `json:"standings"`
}

// ESPNLink represents a link to a resource (unused for us)
type ESPNLink struct {
	Language   string   `json:"language"`
	Rel        []string `json:"rel"`
	Href       string   `json:"href"`
	Text       string   `json:"text"`
	ShortText  string   `json:"shortText"`
	IsExternal bool     `json:"isExternal"`
	IsPremium  bool     `json:"isPremium"`
}

// ESPNSeason represents a overall season in the standings, including all subseasons
type ESPNSeason struct {
	Year        int             `json:"year"`
	StartDate   string          `json:"startDate"`
	EndDate     string          `json:"endDate"`
	DisplayName string          `json:"displayName"`
	Types       []ESPNSubSeason `json:"types"`
}

type ESPNConferenceStandings struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"` // "overall"
	DisplayName       string      `json:"displayName"`
	Links             []ESPNLink  `json:"links"`
	Season            int         `json:"season"`
	SeasonType        int         `json:"seasonType"`
	SeasonDisplayName string      `json:"seasonDisplayName"`
	Entries           []ESPNEntry `json:"entries"`
}

// ESPNEntry represents one entry in the standings, with a team and their associated stats
type ESPNEntry struct {
	Team  team.ESPNTeam    `json:"team"`
	Stats []stats.ESPNStat `json:"stats"`
}

// ESPNSubSeason represents a subseason within an overall season (preseason, regular season, postseason, offseason)
type ESPNSubSeason struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	HasStandings bool   `json:"hasStandings"`
}
