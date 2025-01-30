package main

import (
	"fmt"
	"nfl-app/internal/entry"
	"nfl-app/internal/schedule"
	"nfl-app/internal/schedule/entrysort"
)

func main() {

	scheduleObj := schedule.NewLeagueSchedule()
	scheduleObj.Populate2()

	teamSchedules := scheduleObj.ToTeamSchedules()

	// standings, err := scheduleObj.ToStandings()
	// if err != nil {
	// 	panic(err)
	// }

	standingsObj := scheduleObj.ToStandings()

	// Find the entry for the panthers and saints
	var entries []entry.Entry
	for _, entry := range standingsObj.AllEntries() {
		if entry.TeamName() == "Carolina Panthers" || entry.TeamName() == "New Orleans Saints" || entry.TeamName() == "Chicago Bears" {
			entries = append(entries, entry)
		}
	}

	// Sort the entries
	sortedEntries, err := entrysort.SortEntries(entries, teamSchedules)
	if err != nil {
		panic(err)
	}
	fmt.Println(entry.Teams(sortedEntries))

}
