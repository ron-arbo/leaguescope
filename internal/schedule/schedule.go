package schedule

import (
	"fmt"
	"nfl-app/internal/entry"
	"nfl-app/internal/game"
	"nfl-app/internal/scraper"
	"nfl-app/internal/team"
	"strconv"
)

// Schedule represents an NFL schedule for one or more teams
type Schedule struct {
	Weeks []Week
}

// Week represent an NFL week from the prespective of one or more teams
// Games will contain multiple games if from the prespective of multiple teams
type Week struct {
	Number int
	Games  []game.Game
}

// NewSchedule creates a new empty schedule with 18 weeks
func NewSchedule() Schedule {
	weeks := make([]Week, 18)
	for i := 0; i < 18; i++ {
		weeks[i] = Week{
			Number: i + 1,
			Games:  make([]game.Game, 0),
		}
	}

	return Schedule{
		Weeks: weeks,
	}
}

// AddGame will "append" a game to a schedule, adding it to
// the first week with an empty Games slice
func (s *Schedule) AddGame(week int, game game.Game) {
	s.Weeks[week-1].Games = append(s.Weeks[week-1].Games, game)
}

func CreateSchedule(rows []scraper.ScrapedRow) Schedule {
	weeks := make([]Week, 18)

	for _, row := range rows {
		// Parse the week number
		weekNum, _ := strconv.Atoi(row.Week) // TODO: handle error
		week := &weeks[weekNum-1]

		// Parse the time
		t := scraper.ParseTime(row)

		// Determine home/away. Winner is listed first, @ is used optionally
		var home, away string
		if row.GameLocation == "@" {
			home = row.Loser
			away = row.Winner
		} else {
			home = row.Winner
			away = row.Loser
		}

		// Convert int fields. TODO: handle error
		ptsWin, _ := strconv.Atoi(row.PtsWin)
		ptsLose, _ := strconv.Atoi(row.PtsLose)
		yardsWin, _ := strconv.Atoi(row.YardsWin)
		yardsLose, _ := strconv.Atoi(row.YardsLose)
		toWin, _ := strconv.Atoi(row.ToWin)
		toLose, _ := strconv.Atoi(row.ToLose)

		// Create the game
		g := game.Game{
			Time:      t,
			Winner:    row.Winner,
			Loser:     row.Loser,
			Home:      home,
			Away:      away,
			PtsWin:    ptsWin,
			PtsLose:   ptsLose,
			YardsWin:  yardsWin,
			YardsLose: yardsLose,
			ToWin:     toWin,
			ToLose:    toLose,
		}

		// Add the game to the week
		if week.Games == nil {
			week.Games = make([]game.Game, 0)
		}
		week.Games = append(week.Games, g)
		week.Number = weekNum // TODO: Don't need to do this every loop, but it works
	}

	return Schedule{Weeks: weeks}
}

func (s *Schedule) Print() {
	for _, week := range s.Weeks {
		fmt.Printf("Week %d\n", week.Number)
		for _, game := range week.Games {
			fmt.Printf("%s vs %s\n", game.Home, game.Away)
		}
	}
}

// SplitToTeams splits the schedule into a map of team schedules
func (s *Schedule) SplitToTeams() map[string]Schedule {
	// Initialize the map
	teamSchedules := make(map[string]Schedule)
	for _, team := range team.NFLTeams {
		teamSchedules[team.Name] = NewSchedule()
	}

	for i, week := range s.Weeks {
		for _, game := range week.Games {
			// Add the game to both team's schedules
			homeSchedule := teamSchedules[game.Home]
			homeSchedule.Weeks[i].Games = append(homeSchedule.Weeks[i].Games, game)

			awaySchedule := teamSchedules[game.Away]
			awaySchedule.Weeks[i].Games = append(awaySchedule.Weeks[i].Games, game)
		}
	}

	return teamSchedules
}

// OpponentMapFor returns a map of opponent names to games played against that opponent for a given team
func (s *Schedule) OpponentMapFor(team string) map[string][]game.Game {
	// Need a string of games because teams will play division opponents twice
	oppMap := make(map[string][]game.Game)
	for _, week := range s.Weeks {
		for _, game := range week.Games {
			if game.Home == team {
				oppMap[game.Away] = append(oppMap[game.Away], game)
			} else if game.Away == team {
				oppMap[game.Home] = append(oppMap[game.Home], game)
			}
		}
	}

	return oppMap
}

// CreateEntries will create a slice of entries representing the given schedule
// Note that this slice is not guaranteed to contain an entry for every team, only
// those teams in involved in the given schedule
func CreateEntries(schedule Schedule) []entry.Entry {
	// Create a map to track entries for each team
	// Use pointer to entry so we can update in place without re-assigning to map
	entryMap := make(map[string]*entry.Entry)

	for _, week := range schedule.Weeks {
		for _, game := range week.Games {
			// Intialize entry in map if needed
			if entryMap[game.Home] == nil {
				entryMap[game.Home] = entry.NewEntry(game.Home)
			}
			if entryMap[game.Away] == nil {
				entryMap[game.Away] = entry.NewEntry(game.Away)
			}

			// Update the entries for each team
			entryMap[game.Home].AddGame(game)
			entryMap[game.Away].AddGame(game)
		}
	}

	// Now that the basic counting stats are in place, calculate the stats that depend on other teams
	// To do this, split the schedule into individual team schedules
	teamSchedules := schedule.SplitToTeams()

	// Get slice of entries
	entries := make([]entry.Entry, 0)
	for _, entry := range entryMap {
		entries = append(entries, *entry)
	}

	// SOV and SOS
	for i, entry := range entries {
		sov := StrengthOfVictory(entry.Team.Name, entries, teamSchedules)
		sos := StrengthOfSchedule(entry.Team.Name, entries, teamSchedules)

		entry.Stats.StrengthOfVictory = sov
		entry.Stats.StrengthOfSchedule = sos

		// TODO: Use pointer to avoid this?
		entries[i] = entry
	}

	// Combined ranking among conference teams in points scored and points allowed
	leagueRankPf := Ranking(entries, pointsFor)
	leagueRankPa := Ranking(entries, pointsAgainst)

	for _, conf := range team.Conferences {
		confEntries := entry.ConferenceEntries(entries, conf)
		confRankPf := Ranking(confEntries, pointsFor)
		confRankPa := Ranking(confEntries, pointsAgainst)

		for _, entry := range confEntries {
			entry.Stats.LeagueRankPointsFor = leagueRankPf[entry.Team.Name]
			entry.Stats.LeagueRankPointsAgainst = leagueRankPa[entry.Team.Name]

			entry.Stats.ConferenceRankPointsFor = confRankPf[entry.Team.Name]
			entry.Stats.ConferenceRankPointsAgainst = confRankPa[entry.Team.Name]
		}
	}

	return entries
}
