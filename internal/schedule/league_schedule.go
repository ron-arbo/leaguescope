package schedule

import (
	"encoding/json"
	"fmt"
	"log"
	"nfl-app/internal/entry"
	"nfl-app/internal/game"
	"nfl-app/internal/standings"
	"nfl-app/internal/stats"
	"nfl-app/internal/team"
	"os"
	"sort"
)

// LeagueSchedule is meant to be the big schedule object that we can update and use to create new standings
// Unlike TeamSchedules, there is only one object for each game. TeamSchedule has two objects, one for each competitor
type LeagueSchedule struct {
	Weeks []SharedWeek
}

type SharedWeek struct {
	WeekNumber int
	Games      []*game.Game
}

func NewLeagueSchedule() LeagueSchedule {
	// Create a new LeagueSchedule with 18 weeks
	weeks := make([]SharedWeek, 18)
	for i := 0; i < 18; i++ {
		weeks[i] = SharedWeek{
			WeekNumber: i + 1,
			Games:      []*game.Game{},
		}
	}
	return LeagueSchedule{
		Weeks: weeks,
	}
}

// TODO: object.Populate() standard practice? Better way?
func (ls *LeagueSchedule) Populate() {

	for _, team := range team.NFLTeams {
		schedule := GetESPNSchedule(team.ID)
		ls.Add(schedule)
	}

	// Remove duplicates
	for i := 0; i < 18; i++ {
		week := ls.Weeks[i]
		uniqueGames := make([]*game.Game, 0, len(week.Games))
		seen := make(map[string]bool)
		for _, game := range week.Games {
			key := fmt.Sprintf("%s@%s", game.HomeTeam, game.AwayTeam)
			if _, ok := seen[key]; !ok {
				uniqueGames = append(uniqueGames, game)
				seen[key] = true
			}
		}
		ls.Weeks[i].Games = uniqueGames
	}
}

// Populate2 is just so we don't have to make API calls every time we run tests. Use hardcoded JSON instead
func (ls *LeagueSchedule) Populate2() {
	// read from hardcoded JSON file
	fileBytes, _ := os.ReadFile("./sample-ss/week18-ss.json")
	err := json.Unmarshal(fileBytes, &ls)

	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}
}

// ToTeamSchedules converts a LeagueSchedule to a map of TeamSchedules
func (ls *LeagueSchedule) ToTeamSchedules() map[string]*TeamSchedule {
	teamSchedules := NewTeamSchedules()

	for _, week := range ls.Weeks {
		for _, game := range week.Games {
			// Create TeamWeek for home team
			homeTeamWeek := TeamWeek{
				WeekNumber: week.WeekNumber,
				Game:       game,
			}

			// Create TeamWeek for away team
			awayTeamWeek := TeamWeek{
				WeekNumber: week.WeekNumber,
				Game:       game,
			}

			// Add TeamWeeks to TeamSchedule. Bye week will contain a nil game
			teamSchedules[game.HomeTeam.Name.String()].Weeks[week.WeekNumber-1] = homeTeamWeek
			teamSchedules[game.AwayTeam.Name.String()].Weeks[week.WeekNumber-1] = awayTeamWeek
		}
	}

	return teamSchedules
}

// Add adds all the games from a team's ESPNSchedule to the LeagueSchedule
func (ls *LeagueSchedule) Add(espnSchedule ESPNSchedule) {
	// The espnSchedule contains all the games for one team
	// Loop through each of those games and add them to the shared schedule
	for _, event := range espnSchedule.Events {
		weekNum := event.Week.Number
		game := EventToGame(event)

		// Add the game to the correct week
		weeksGames := ls.Weeks[weekNum-1].Games
		ls.Weeks[weekNum-1].Games = append(weeksGames, game)
	}
}

func (ls *LeagueSchedule) ToStandings() *standings.Standings {
	standings := standings.NewStandings()

	teamSchedules := ls.ToTeamSchedules()

	for _, teamSchedule := range teamSchedules {
		entry := teamSchedule.ToEntry()
		standings.UpsertEntry(*entry)
	}

	// Now that the basic counting stats are in place,
	// calculate the stats that depend on other teams
	allEntries := standings.AllEntries()

	// SOV and SOS
	for _, entry := range allEntries {
		sov := StrengthOfVictory(entry.TeamName(), allEntries, teamSchedules)
		sos := StrengthOfSchedule(entry.TeamName(), allEntries, teamSchedules)
		entry.Stats[stats.StatStrengthOfVictory].SetValue(sov)
		entry.Stats[stats.StatStrengthOfSchedule].SetValue(sos)
	}

	// Combined ranking among conference teams in points scored and points allowed
	leagueRankPointsFor := Ranking(allEntries, stats.StatPointsFor)
	leagueRankPointsAgainst := Ranking(allEntries, stats.StatPointsAgainst)

	for _, conference := range team.Conferences {
		conferenceEntries := standings.ConferenceStandings[conference].Entries
		conferenceRankPointsFor := Ranking(conferenceEntries, stats.StatPointsFor)
		conferenceRankPointsAgainst := Ranking(conferenceEntries, stats.StatPointsAgainst)

		for _, entry := range conferenceEntries {
			entry.Stats[stats.LeagueRankPointsFor].SetValue(float64(leagueRankPointsFor[entry.TeamName()]))
			entry.Stats[stats.LeagueRankPointsAgainst].SetValue(float64(leagueRankPointsAgainst[entry.TeamName()]))

			entry.Stats[stats.ConferenceRankPointsFor].SetValue(float64(conferenceRankPointsFor[entry.TeamName()]))
			entry.Stats[stats.ConferenceRankPointsAgainst].SetValue(float64(conferenceRankPointsAgainst[entry.TeamName()]))
		}
	}

	// Calculate playoffs seeds
	// standings.AssignPlayoffSeeds()

	// TODO: Calculate clincher

	return standings
}

const (
	strengthOfVictory  = "victory"
	strengthOfSchedule = "schedule"
)

func StrengthOfVictory(team string, entries []entry.Entry, teamSchedules map[string]*TeamSchedule) float64 {
	return StrengthOf(team, entries, teamSchedules, strengthOfVictory)
}

func StrengthOfSchedule(team string, entries []entry.Entry, teamSchedules map[string]*TeamSchedule) float64 {
	return StrengthOf(team, entries, teamSchedules, strengthOfSchedule)
}

func StrengthOf(team string, entries []entry.Entry, teamSchedules map[string]*TeamSchedule, attribute string) float64 {
	if attribute != strengthOfVictory && attribute != strengthOfSchedule {
		panic("unknown attribute")
	}

	combinedOpponentRecord := entry.NewRecord(0, 0, 0)

	// Create map of team names to Entry for convenience
	teamMap := make(map[string]entry.Entry)
	for _, entry := range entries {
		teamMap[entry.TeamName()] = entry
	}

	schedule := teamSchedules[team]
	for _, week := range schedule.Weeks {
		if week.Game == nil || !week.Game.Completed {
			continue
		}

		// Get the opponent's record
		var opp string
		if week.Game.HomeTeam.Name.String() == team {
			opp = week.Game.AwayTeam.Name.String()
		} else {
			opp = week.Game.HomeTeam.Name.String()
		}

		// For victory, only consider games that were won
		if attribute == strengthOfVictory && !week.Game.WonBy(team) {
			continue
		}

		// For schedule, consider all games
		oppRecord := entry.NewRecord(
			int(teamMap[opp].Stats[stats.StatWins].Value()),
			int(teamMap[opp].Stats[stats.StatLosses].Value()),
			int(teamMap[opp].Stats[stats.StatTies].Value()),
		)

		combinedOpponentRecord.Add(oppRecord)
	}

	return combinedOpponentRecord.WinPercentage()
}

// Ranking returns a map of team names to their rank amongst the given entries based on the given stat
// Tied teams will shared the rank
func Ranking(entries []entry.Entry, stat string) map[string]int {
	if stat != stats.StatPointsFor && stat != stats.StatPointsAgainst {
		panic("unsupported stat")
	}

	var sortFunc func(i, j int) bool
	switch stat {
	case stats.StatPointsFor:
		// Sort PointsFor in descending order
		sortFunc = func(i, j int) bool {
			return entries[i].Stats[stats.StatPointsFor].Value() > entries[j].Stats[stats.StatPointsFor].Value()
		}
	case stats.StatPointsAgainst:
		// Sort PointsAgainst in ascending order
		sortFunc = func(i, j int) bool {
			return entries[i].Stats[stats.StatPointsAgainst].Value() < entries[j].Stats[stats.StatPointsAgainst].Value()
		}
	}
	sort.Slice(entries, sortFunc)

	// Assign rankings. If tied, give the same rank
	ranking := make(map[string]int)
	rank := 1
	for i, entry := range entries {
		// If the given entry is different than the last, we can increase the rank
		if i > 0 && entry.Stats[stat].Value() != entries[i-1].Stats[stat].Value() {
			rank = i + 1
		}
		ranking[entry.TeamName()] = rank
	}

	return ranking
}
