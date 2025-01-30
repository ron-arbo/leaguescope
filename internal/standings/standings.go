package standings

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"nfl-app/internal/entry"
	"nfl-app/internal/stats"
	"nfl-app/internal/team"
)

type Standings struct {
	// Map conference name to conference standings
	ConferenceStandings map[string]ConferenceStandings
}

type ConferenceStandings struct {
	Conference string // Conference name
	Entries    []entry.Entry
}

var standingsUrl = "https://site.api.espn.com/apis/v2/sports/football/nfl/standings"

// TODO: We'll cache this result at the beginning of the program at some point, but for now call it as needed
func GetESPNStandings() ESPNStandings {
	// Send a GET request
	resp, err := http.Get(standingsUrl)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Parse the JSON into the struct
	var standings ESPNStandings
	if err := json.Unmarshal(body, &standings); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	return standings
}

func GetConferenceStandings(conference string) ESPNConferenceStandings {
	standings := GetESPNStandings()
	if standings.Conferences == nil {
		log.Fatalf("No conferences found in standings")
	}

	if standings.Conferences[0].Name == conference {
		return standings.Conferences[0].Standings
	}

	if standings.Conferences[1].Name == conference {
		return standings.Conferences[1].Standings
	}

	log.Fatalf("Conference %s not found in standings", conference)
	return ESPNConferenceStandings{}
}

// NewStandings creates a new Standings object with empty ConferenceStandings
func NewStandings() *Standings {
	confStandings := make(map[string]ConferenceStandings)
	for _, conference := range team.Conferences {
		confStandings[conference] = ConferenceStandings{
			Conference: conference,
			Entries:    []entry.Entry{},
		}
	}

	// Create empty entry for each team
	for _, team := range team.NFLTeams {
		// Create team entry
		teamEntry := entry.Entry{
			Team:  team,
			Stats: stats.NewStats(),
		}

		// Add entry to conference standings
		conferenceStandings := confStandings[team.Conference]
		conferenceStandings.Entries = append(conferenceStandings.Entries, teamEntry)
		confStandings[team.Conference] = conferenceStandings
	}

	return &Standings{
		ConferenceStandings: confStandings,
	}
}

// TODO: More pointer types to avoid re-assigning?
func (s *Standings) UpsertEntry(entry entry.Entry) {
	// Find the entry in the standings
	for _, conf := range team.Conferences {
		confStandings := s.ConferenceStandings[conf]

		if entry.Team.Conference != conf {
			continue
		}

		// Find the entry and update it
		for i, e := range confStandings.Entries {
			if e.Team == entry.Team {
				confStandings.Entries[i] = entry
				s.ConferenceStandings[conf] = confStandings
				return
			}
		}

		// If the entry wasn't found, add it
		confStandings.Entries = append(confStandings.Entries, entry)
		s.ConferenceStandings[conf] = confStandings
	}
}

func (s *Standings) AllEntries() []entry.Entry {
	return append(s.AFCEntries(), s.NFCEntries()...)
}

func (s *Standings) AFCEntries() []entry.Entry {
	return s.ConferenceStandings[team.AFC].Entries
}

func (s *Standings) NFCEntries() []entry.Entry {
	return s.ConferenceStandings[team.NFC].Entries
}

// Calculate seed
// Sort each conference by win pct (We should take tiebreakers into account https://www.nfl.com/standings/tie-breaking-procedures)
// Loop through teams, if team is from division we haven't seen, assign it to highest available seed in 1-4
// If team is from division we HAVE seen, assign it to highest available seed in 5+
// func (s *Standings) AssignPlayoffSeeds(teamSchedules map[string]*schedule.TeamSchedule) {

// 	for _, conf := range team.Conferences {
// 		sortedEntries, err := entrysort.SortEntries(s.ConferenceStandings[conf].Entries, teamSchedules)
// 		if err != nil {
// 			// TODO
// 			log.Fatalf("Failed to sort entries: %v", err)
// 		}

// 		divisionSeedAvailable := 1
// 		otherSeedAvailable := 5

// 		divisionsSeen := make(map[string]bool)
// 		for _, entry := range sortedEntries {
// 			// If we haven't seen the division yet, we know this is a division winner,
// 			// so it assign it to the highest available seed in 1-4. Otherwise assign it to 5+
// 			if _, ok := divisionsSeen[entry.Team.Division]; !ok {
// 				entry.Stats[stats.StatPlayoffSeed].SetValue(float64(divisionSeedAvailable))
// 				divisionSeedAvailable++
// 				divisionsSeen[entry.Team.Division] = true
// 			} else {
// 				entry.Stats[stats.StatPlayoffSeed].SetValue(float64(otherSeedAvailable))
// 				otherSeedAvailable++
// 			}
// 		}
// 	}
// }
