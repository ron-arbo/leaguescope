package entry

import (
	"nfl-app/internal/game"
	"nfl-app/internal/stats"
	"nfl-app/internal/team"
	"sort"
	"strings"
)

type Entry struct {
	Team  team.Team
	Stats stats.Stats
}

func ConferenceEntries(entries []Entry, conference string) []Entry {
	out := make([]Entry, 0)
	for _, entry := range entries {
		if entry.Team.Conference == conference {
			out = append(out, entry)
		}
	}
	return out
}

func (e *Entry) TeamName() string {
	return e.Team.Name
}

func (e *Entry) AddGame(game game.Game) {
	// Records
	e.UpdateRecords(game)

	// Points
	e.UpdatePoints(game)

	// Streak
	e.UpdateStreak(game)

	// Certain stats (Strength of Victory, Strength of Schedule, Clincher, etc) cannot be calculated from one game,
	// they are dependent on the entire schedule. Calculate these later
}

func (e *Entry) UpdateRecords(game game.Game) {
	teamname := e.Team.Name

	switch {
	case game.Winner == teamname:
		// Overall
		e.Stats.Record.AddWin()

		// Home/Away
		if game.Home == teamname {
			e.Stats.HomeRecord.AddWin()
		} else {
			e.Stats.AwayRecord.AddWin()
		}

		// Division
		if team.SameDivision(game.Winner, game.Loser) {
			e.Stats.DivisionRecord.AddWin()
		}

		// Conference
		if team.SameConference(game.Winner, game.Loser) {
			e.Stats.ConferenceRecord.AddWin()
		}
	case game.Loser == teamname:
		// Overall
		e.Stats.Record.AddLoss()

		// Home/Away
		if game.Home == teamname {
			e.Stats.HomeRecord.AddLoss()
		} else {
			e.Stats.AwayRecord.AddLoss()
		}

		// Division
		if team.SameDivision(game.Winner, game.Loser) {
			e.Stats.DivisionRecord.AddLoss()
		}

		// Conference
		if team.SameConference(game.Winner, game.Loser) {
			e.Stats.ConferenceRecord.AddLoss()
		}
	}
}

func (e *Entry) UpdatePoints(game game.Game) {
	switch {
	case game.Winner == e.Team.Name:
		// Overall
		e.Stats.Points.AddFor(game.PtsWin)
		e.Stats.Points.AddAgainst(game.PtsLose)

		// Conference
		if team.SameConference(game.Winner, game.Loser) {
			e.Stats.ConferencePoints.AddFor(game.PtsWin)
			e.Stats.ConferencePoints.AddAgainst(game.PtsLose)
		}
	case game.Loser == e.Team.Name:
		// Overall
		e.Stats.Points.AddFor(game.PtsLose)
		e.Stats.Points.AddAgainst(game.PtsWin)

		// Conference
		if team.SameConference(game.Winner, game.Loser) {
			e.Stats.ConferencePoints.AddFor(game.PtsLose)
			e.Stats.ConferencePoints.AddAgainst(game.PtsWin)
		}
	}
}

func (e *Entry) UpdateStreak(game game.Game) {
	switch {
	case game.Winner == e.Team.Name:
		if e.Stats.Streak > 0 {
			e.Stats.Streak++
		} else {
			e.Stats.Streak = 1
		}
	case game.Loser == e.Team.Name:
		if e.Stats.Streak < 0 {
			e.Stats.Streak--
		} else {
			e.Stats.Streak = -1
		}
	default:
		e.Stats.Streak = 0
	}
}

func NewEntry(teamname string) *Entry {
	return &Entry{
		Team:  team.DisplayNameToTeam(teamname),
		Stats: stats.NewStats(),
	}
}

func GroupByDivision(entries []Entry) map[string][]Entry {
	out := make(map[string][]Entry)
	for _, entry := range entries {
		out[entry.Team.Division] = append(out[entry.Team.Division], entry)
	}
	return out
}

func GroupByConference(entries []Entry) map[string][]Entry {
	out := make(map[string][]Entry)
	for _, entry := range entries {
		out[entry.Team.Conference] = append(out[entry.Team.Conference], entry)
	}
	return out
}

// GroupEntries will split the given slice into multiple slices, each containing a group of tied entries
// For example {1, 1, 2, 3, 3, 3, 4, 5} --> {{1, 1}, {2}, {3, 3, 3}, {4}, {5}}
func GroupEntries(entries []Entry, sortBy map[string]float64) [][]Entry {
	out := make([][]Entry, 0)

	if len(entries) == 0 {
		return out
	}

	// Sort the entries in descending order using sortBy
	sort.Slice(entries, func(i, j int) bool {
		return sortBy[entries[i].Team.Name] > sortBy[entries[j].Team.Name]
	})

	// Initialize with the first entry
	curValue := sortBy[entries[0].Team.Name]
	curSlice := []Entry{entries[0]}

	for i := 1; i < len(entries); i++ {
		// If the current value is the same as the previous, add the entry to the current slice
		if sortBy[entries[i].Team.Name] == curValue {
			curSlice = append(curSlice, entries[i])
		} else {
			// New value hit, add the current slice to the output and create a new one with the current entry
			out = append(out, curSlice)

			// Start a new slice
			curValue = sortBy[entries[i].Team.Name]
			curSlice = []Entry{entries[i]}
		}
	}

	// Add the last slice, since in the loop we only add when we hit a new value
	out = append(out, curSlice)

	return out
}

func Teams(entries []Entry) string {
	out := make([]string, 0)
	for _, entry := range entries {
		out = append(out, entry.Team.Name)
	}
	return strings.Join(out, ", ")
}

// SplitAround will split the given slice of entries into two slices
// The first slice will contain all entries up to the given index
// The second slice will contain all entries after and including the given index
func SplitAround(entries []Entry, index int) ([]Entry, []Entry) {
	return entries[:index], entries[index:]
}
