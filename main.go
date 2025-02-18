package main

import (
	"fmt"
	"nfl-app/internal/team"
)

func main() {
	// scrapedRows, err := scraper.ScrapeYear("2024")
	// if err != nil {
	// 	panic(err)
	// }

	teams := team.GetRandomTeams(7, []team.Team{team.ArizonaCardinals, team.BaltimoreRavens},
		[]string{"NFC South", "NFC East"}, "American Football Conference")
	fmt.Println(team.Names(teams))

	// sched := schedule.CreateSchedule(scrapedRows)
	// ts := sched.SplitToTeams()

	// entries := schedule.CreateEntries(sched)

	// // Take only afc entries
	// var afcEntries []entry.Entry
	// for _, e := range entries {
	// 	if e.Team.Conference == team.NFC {
	// 		afcEntries = append(afcEntries, e)
	// 	}
	// }

	// fmt.Println("Before")
	// for _, e := range afcEntries {
	// 	fmt.Println(e.Team.Name)
	// }

	// sortedEntries, err := entrysort.SortEntries(afcEntries, ts)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("After")
	// for _, e := range sortedEntries {
	// 	fmt.Println(e.Team.Name)
	// }
}
