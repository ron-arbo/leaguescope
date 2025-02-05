package entry

import (
	"fmt"
	"nfl-app/internal/game"
	"nfl-app/internal/stats"
	"nfl-app/internal/team"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

type Entry struct {
	Team  team.Team // Team Name
	Stats stats.Stats
}

type Entry2 struct {
	Team team.Team
	StatSheet stats.StatSheet
}

func (e *Entry2) AddGame(game game.Game2) {
	// Records
	e.UpdateRecords(game)

	// Points
	e.UpdatePoints(game)

	// Streak
	e.UpdateStreak(game)

	// Certain stats cannot be calculated from one game, they are dependent on the entire schedule
	// Strength of Victory, Strength of Schedule, Clincher, etc
	// Calculate these later
}

func (e *Entry2) UpdateRecords(game game.Game2) {
	teamname := e.Team.Name.String()

	switch {
	case game.Winner == teamname:
		// Overall
		e.StatSheet.Record.AddWin()

		// Home/Away
		if game.Home == teamname {
			e.StatSheet.HomeRecord.AddWin()
		} else {
			e.StatSheet.AwayRecord.AddWin()
		}

		// Division
		if team.SameDivision(game.Winner, game.Loser) {
			e.StatSheet.DivisionRecord.AddWin()
		}

		// Conference
		if team.SameConference(game.Winner, game.Loser) {
			e.StatSheet.ConferenceRecord.AddWin()
		}
	case game.Loser == teamname:
		// Overall
		e.StatSheet.Record.AddLoss()

		// Home/Away
		if game.Home == teamname {
			e.StatSheet.HomeRecord.AddLoss()
		} else {
			e.StatSheet.AwayRecord.AddLoss()
		}

		// Division
		if team.SameDivision(game.Winner, game.Loser) {
			e.StatSheet.DivisionRecord.AddLoss()
		}

		// Conference
		if team.SameConference(game.Winner, game.Loser) {
			e.StatSheet.ConferenceRecord.AddLoss()
		}
	}
}

func (e *Entry2) UpdatePoints(game game.Game2) {
	switch {
	case game.Winner == e.Team.Name.String():
		// Overall
		e.StatSheet.Points.AddFor(game.PtsWin)
		e.StatSheet.Points.AddAgainst(game.PtsLose)

		// Conference
		if team.SameConference(game.Winner, game.Loser) {
			e.StatSheet.ConferencePoints.AddFor(game.PtsWin)
			e.StatSheet.ConferencePoints.AddAgainst(game.PtsLose)
		}
	case game.Loser == e.Team.Name.String():
		// Overall
		e.StatSheet.Points.AddFor(game.PtsLose)
		e.StatSheet.Points.AddAgainst(game.PtsWin)

		// Conference
		if team.SameConference(game.Winner, game.Loser) {
			e.StatSheet.ConferencePoints.AddFor(game.PtsLose)
			e.StatSheet.ConferencePoints.AddAgainst(game.PtsWin)
		}
	}
}

func (e *Entry2) UpdateStreak(game game.Game2) {
	switch {
	case game.Winner == e.Team.Name.String():
		if e.StatSheet.Streak > 0 {
			e.StatSheet.Streak++
		} else {
			e.StatSheet.Streak = 1
		}
	case game.Loser == e.Team.Name.String():
		if e.StatSheet.Streak < 0 {
			e.StatSheet.Streak--
		} else {
			e.StatSheet.Streak = -1
		}
	default:
		e.StatSheet.Streak = 0
	}
}

func NewEntry(teamname string) *Entry {
	return &Entry{
		Team:  team.DisplayNameToTeam(teamname),
		Stats: stats.NewStats(),
	}
}

func NewEntry2(teamname string) *Entry2 {
	return &Entry2{
		Team: team.DisplayNameToTeam(teamname),
		StatSheet: stats.NewStatSheet(),
	}
}


// Stat getters
func (e *Entry) Wins() int {
	return int(e.Stats[stats.StatWins].Value())
}
func (e *Entry) Losses() int {
	return int(e.Stats[stats.StatLosses].Value())
}
func (e *Entry) Ties() int {
	return int(e.Stats[stats.StatTies].Value())
}
func (e *Entry) WinPercentage() float64 {
	return e.Stats[stats.StatWinPercent].Value()
}
func (e *Entry) Seed() int {
	return int(e.Stats[stats.StatPlayoffSeed].Value())
}


func (e *Entry) UpdateStreak(game *game.Game) {
	stat := e.Stats[stats.StatStreak]

	if game.Tied() {
		stat.SetValue(0)
		e.Stats[stats.StatStreak] = stat
		return
	}

	if e.TeamWon(game) {
		// Add if streak is positive, otherwise reset to 1
		if stat.Value() > 0 {
			stat.SetValue(stat.Value() + 1)
		} else {
			stat.SetValue(1)
		}
	} else {
		// Subtract if streak is negative, otherwise reset to -1
		if stat.Value() < 0 {
			stat.SetValue(stat.Value() - 1)
		} else {
			stat.SetValue(-1)
		}
	}
}

func (e *Entry) TeamWon(game *game.Game) bool {
	if game.HomeTeamName() == e.Team.Name.String() {
		return game.HomeScore > game.AwayScore
	}
	return game.AwayScore > game.HomeScore
}

func (e *Entry) TeamLost(game *game.Game) bool {
	if game.HomeTeamName() == e.Team.Name.String() {
		return game.HomeScore < game.AwayScore
	}
	return game.AwayScore < game.HomeScore
}

func (e *Entry) TeamTied(game *game.Game) bool {
	return game.Tied()
}

func (e *Entry) TeamScore(game *game.Game) int {
	if game.HomeTeamName() == e.Team.Name.String() {
		return game.HomeScore
	}
	if game.AwayTeamName() == e.Team.Name.String() {
		return game.AwayScore
	}
	return 0
}

func (e *Entry) OpponentScore(game *game.Game) int {
	if game.HomeTeamName() == e.Team.Name.String() {
		return game.AwayScore
	}
	if game.AwayTeamName() == e.Team.Name.String() {
		return game.HomeScore
	}
	return 0
}

func (e *Entry) HomeTeam(game *game.Game) bool {
	return game.HomeTeamName() == e.Team.Name.String()
}

func (e *Entry) RoadTeam(game *game.Game) bool {
	return game.AwayTeamName() == e.Team.Name.String()
}

func (e *Entry) DivisionOpponent(game *game.Game) bool {
	return game.HomeTeam.Division == game.AwayTeam.Division
}

func (e *Entry) ConferenceOpponent(game *game.Game) bool {
	return game.HomeTeam.Conference == game.AwayTeam.Conference
}

func (e *Entry) OverallRecord() *Record {
	return NewRecord(
		int(e.Stats[stats.StatWins].Value()),
		int(e.Stats[stats.StatLosses].Value()),
		int(e.Stats[stats.StatTies].Value()),
	)
}

func (e *Entry) DivisionRecord() *Record {
	return NewRecord(
		int(e.Stats[stats.StatDivisionWins].Value()),
		int(e.Stats[stats.StatDivisionLosses].Value()),
		int(e.Stats[stats.StatDivisionTies].Value()),
	)
}

func (e *Entry) ConferenceRecord() *Record {
	return NewRecord(
		int(e.Stats[stats.StatConferenceWins].Value()),
		int(e.Stats[stats.StatConferenceLosses].Value()),
		int(e.Stats[stats.StatConferenceTies].Value()),
	)
}

// TODO: Duplicate code, refactor
func (e *Entry) UpdateStats(game *game.Game) {
	switch {
	case e.TeamWon(game):
		// Overall
		e.Stats[stats.StatWins].Increment()

		// Home/Road
		if e.HomeTeam(game) {
			e.Stats[stats.StatHomeWins].Increment()
		} else {
			e.Stats[stats.StatRoadWins].Increment()
		}

		// Division
		if e.DivisionOpponent(game) {
			e.Stats[stats.StatDivisionWins].Increment()
		}

		// Conference
		if e.ConferenceOpponent(game) {
			e.Stats[stats.StatConferenceWins].Increment()
		}

	case e.TeamLost(game):
		// Overall
		e.Stats[stats.StatLosses].Increment()

		// Home/Road
		if e.HomeTeam(game) {
			e.Stats[stats.StatHomeLosses].Increment()
		} else {
			e.Stats[stats.StatRoadLosses].Increment()
		}

		// Division
		if e.DivisionOpponent(game) {
			e.Stats[stats.StatDivisionLosses].Increment()
		}

		// Conference
		if e.ConferenceOpponent(game) {
			e.Stats[stats.StatConferenceLosses].Increment()
		}

	case e.TeamTied(game):
		// Overall
		e.Stats[stats.StatTies].Increment()

		// Home/Road
		if e.HomeTeam(game) {
			e.Stats[stats.StatHomeTies].Increment()
		} else {
			e.Stats[stats.StatRoadTies].Increment()
		}

		// Division
		if e.DivisionOpponent(game) {
			e.Stats[stats.StatDivisionTies].Increment()
		}

		// Conference
		if e.ConferenceOpponent(game) {
			e.Stats[stats.StatConferenceTies].Increment()
		}
	default:
		panic("impossible!")
	}

	// Streak
	e.UpdateStreak(game)

	// Points
	e.Stats[stats.StatPointsFor].SetValue(e.Stats[stats.StatPointsFor].Value() + float64(e.TeamScore(game)))
	e.Stats[stats.StatPointsAgainst].SetValue(e.Stats[stats.StatPointsAgainst].Value() + float64(e.OpponentScore(game)))
	e.Stats[stats.StatPointDifferential].SetValue(e.Stats[stats.StatPointsFor].Value() - e.Stats[stats.StatPointsAgainst].Value())

	if e.ConferenceOpponent(game) {
		e.Stats[stats.StatConferencePointsFor].SetValue(e.Stats[stats.StatConferencePointsFor].Value() + float64(e.TeamScore(game)))
		e.Stats[stats.StatConferencePointsAgainst].SetValue(e.Stats[stats.StatConferencePointsAgainst].Value() + float64(e.OpponentScore(game)))
		e.Stats[stats.StatConferencePointDifferential].SetValue(e.Stats[stats.StatConferencePointsFor].Value() - e.Stats[stats.StatConferencePointsAgainst].Value())
	}

	// Win Percent
	winPct := e.OverallRecord().WinPercentage()
	e.Stats[stats.StatWinPercent].SetValue(winPct)

	// TODO: Clincher, Games Behind, Playoff Seed must be update in a LeagueSchedule
	// They cannot be determined from a single entry/game
}

func (e *Entry) TeamName() string {
	return e.Team.Name.String()
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

// OverallRecords returns a map of team names to their overall records
func OverallRecords(entries []Entry) map[string]*Record {
	overallRecords := make(map[string]*Record)

	// Initialize overall record for each team
	for _, entry := range entries {
		overallRecords[entry.TeamName()] = entry.OverallRecord()
	}

	return overallRecords
}

// DivisionRecords returns a map of team names to their division records
func DivisionRecords(entries []Entry) map[string]*Record {
	divisionRecord := make(map[string]*Record)

	// Initialize division record for each team
	for _, entry := range entries {
		divisionRecord[entry.TeamName()] = entry.DivisionRecord()
	}

	return divisionRecord
}

// ConferenceRecords returns a map of team names to their conference records
func ConferenceRecords(entries []Entry) map[string]*Record {
	conferenceRecord := make(map[string]*Record)

	// Initialize conference record for each team
	for _, entry := range entries {
		conferenceRecord[entry.TeamName()] = entry.ConferenceRecord()
	}

	return conferenceRecord
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
		return sortBy[entries[i].Team.Name.String()] > sortBy[entries[j].Team.Name.String()]
	})

	// Initialize with the first entry
	curValue := sortBy[entries[0].Team.Name.String()]
	curSlice := []Entry{entries[0]}

	for i := 1; i < len(entries); i++ {
		// If the current value is the same as the previous, add the entry to the current slice
		if sortBy[entries[i].Team.Name.String()] == curValue {
			curSlice = append(curSlice, entries[i])
		} else {
			// New value hit, add the current slice to the output and create a new one with the current entry
			out = append(out, curSlice)

			// Start a new slice
			curValue = sortBy[entries[i].Team.Name.String()]
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
		out = append(out, entry.TeamName())
	}
	return strings.Join(out, ", ")
}

// SplitAround will split the given slice of entries into two slices
// The first slice will contain all entries up to the given index
// The second slice will contain all entries after and including the given index
func SplitAround(entries []Entry, index int) ([]Entry, []Entry) {
	return entries[:index], entries[index:]
}

func Print(entries []Entry) {
    // Use tabwriter to format the output
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
    fmt.Fprintln(w, "Seed\tTeam Name\tWins\tLosses\tTies\tWin Percentage\t")

    // Print columns for team name, wins, losses, ties, win percentage
    for _, entry := range entries {
        fmt.Fprintf(w, "%d\t%s\t%d\t%d\t%d\t%.2f\t\n", entry.Seed(), entry.TeamName(), entry.Wins(), entry.Losses(), entry.Ties(), entry.WinPercentage())
    }

    w.Flush()
}