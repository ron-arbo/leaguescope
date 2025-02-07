package schedule

import (
	"fmt"
	"nfl-app/internal/entry"
	"nfl-app/internal/game"
	"nfl-app/internal/scraper"
	"nfl-app/internal/team"
	"strconv"
)

type Schedule struct {
	Weeks []Week
}

type Week struct {
	Number int
	Games  []game.Game
}

// NewEmptySchedule creates a new empty schedule with 18 weeks
func NewEmptySchedule() Schedule {
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
		teamSchedules[team.Name.String()] = NewEmptySchedule()
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

// OpponentMapFor returns a map of opponent names to games played against that opponent
// for a given team
func (s *Schedule) OpponentMapFor(teamName string) map[string][]game.Game {
	// Need a string of games because teams will play division opponents twice
	opponentMap := make(map[string][]game.Game)
	for _, week := range s.Weeks {
		for _, game := range week.Games {
			if game.Home == teamName {
				opponentMap[game.Away] = append(opponentMap[game.Away], game)
			} else if game.Away == teamName {
				opponentMap[game.Home] = append(opponentMap[game.Home], game)
			}
		}
	}

	return opponentMap
}

func CreateEntries(schedule Schedule) []entry.Entry {
	// Create a map to track entries for each team
	entryMap := make(map[string]entry.Entry)
	for _, team := range team.NFLTeams {
		teamname := team.Name.String()
		entryMap[teamname] = *entry.NewEntry(teamname)
	}

	for _, week := range schedule.Weeks {
		for _, game := range week.Games {
			// Update the entries for each team
			homeEntry := entryMap[game.Home]
			awayEntry := entryMap[game.Away]

			homeEntry.AddGame(game)
			awayEntry.AddGame(game)

			// Need to reassign, could make pointers instead?
			entryMap[game.Home] = homeEntry
			entryMap[game.Away] = awayEntry
		}
	}

	// Now that the basic counting stats are in place,
	// calculate the stats that depend on other teams

	// Get slice of entries
	entries := make([]entry.Entry, 0)
	for _, entry := range entryMap {
		entries = append(entries, entry)
	}

	teamSchedules := schedule.SplitToTeams()

	// SOV and SOS
	for i, entry := range entries {
		sov := StrengthOfVictory(entry.TeamName(), entries, teamSchedules)
		sos := StrengthOfSchedule(entry.TeamName(), entries, teamSchedules)
		entry.Stats.StrengthOfVictory = sov
		entry.Stats.StrengthOfSchedule = sos

		// TODO: Use pointer to avoid this?
		entries[i] = entry
	}

	// Combined ranking among conference teams in points scored and points allowed
	leagueRankPointsFor := Ranking(entries, pointsFor)
	leagueRankPointsAgainst := Ranking(entries, pointsAgainst)

	for _, conference := range team.Conferences {
		conferenceEntries := entry.ConferenceEntries(entries, conference)
		conferenceRankPointsFor := Ranking(conferenceEntries, pointsFor)
		conferenceRankPointsAgainst := Ranking(conferenceEntries, pointsAgainst)

		for _, entry := range conferenceEntries {
			entry.Stats.LeagueRankPointsFor = leagueRankPointsFor[entry.TeamName()]
			entry.Stats.LeagueRankPointsAgainst = leagueRankPointsAgainst[entry.TeamName()]

			entry.Stats.ConferenceRankPointsFor = conferenceRankPointsFor[entry.TeamName()]
			entry.Stats.ConferenceRankPointsAgainst = conferenceRankPointsAgainst[entry.TeamName()]
		}
	}

	return entries
}
