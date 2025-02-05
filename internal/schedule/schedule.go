package schedule

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"nfl-app/internal/entry"
	"nfl-app/internal/game"
	"nfl-app/internal/scraper"
	"nfl-app/internal/team"
	"strconv"
)

func GetESPNSchedule(teamId string) ESPNSchedule {
	url := fmt.Sprintf("http://site.api.espn.com/apis/site/v2/sports/football/nfl/teams/%s/schedule", teamId)

	// Send a GET request
	resp, err := http.Get(url)
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
	var schedule ESPNSchedule
	if err := json.Unmarshal(body, &schedule); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	return schedule
}

func EventToGame(event ESPNEvent) *game.Game {
	competition := event.Competitions[0] // Usually only one competition per event
	return &game.Game{
		// TODO: Confirm home team is always Competitors[0]
		HomeTeam:  team.DisplayNameToTeam(competition.Competitors[0].Team.DisplayName),
		AwayTeam:  team.DisplayNameToTeam(competition.Competitors[1].Team.DisplayName),
		HomeScore: int(competition.Competitors[0].Score.Value),
		AwayScore: int(competition.Competitors[1].Score.Value),
		Completed: competition.Status.Type.Completed,
	}
}

// -------------------- New STUFF --------------------
type Schedule struct {
	Weeks []Week
}

type Week struct {
	Number int
	Games  []game.Game2
}

// NewEmptySchedule creates a new empty schedule with 18 weeks
func NewEmptySchedule() Schedule {
	weeks := make([]Week, 18)
	for i := 0; i < 18; i++ {
		weeks[i] = Week{
			Number: i + 1,
			Games:  make([]game.Game2, 0),
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
		g := game.Game2{
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
			week.Games = make([]game.Game2, 0)
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

func CreateEntries(schedule Schedule) []entry.Entry2 {
	// Create a map to track entries for each team
	entries := make(map[string]entry.Entry2)
	for _, team := range team.NFLTeams {
		teamname := team.Name.String()
		entries[teamname] = *entry.NewEntry2(teamname)
	}

	for _, week := range schedule.Weeks {
		for _, game := range week.Games {
			// Update the entries for each team
			homeEntry := entries[game.Home]
			awayEntry := entries[game.Away]

			homeEntry.AddGame(game)
			awayEntry.AddGame(game)

			// Need to reassign, could make pointers instead?
			entries[game.Home] = homeEntry
			entries[game.Away] = awayEntry
		}
	}

	out := make([]entry.Entry2, 0)
	for _, entry := range entries {
		out = append(out, entry)
	}
	return out
}
