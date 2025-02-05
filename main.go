package main

import (
	"fmt"
	"nfl-app/internal/schedule"
	"nfl-app/internal/scraper"
)

func main() {
	scrapedRows, err := scraper.ScrapeYear("2024")
	if err != nil {
		panic(err)
	}

	sched := schedule.CreateSchedule(scrapedRows)

	entries := schedule.CreateEntries(sched)
	fmt.Print(entries[0])
}
