package scraper

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
)

// ScrapedRow holds all the available data in a single scraped row
type ScrapedRow struct {
	Week string

	// Time
	DayOfWeek string
	Date      string
	Gametime  string

	Winner       string
	GameLocation string
	Loser        string
	PtsWin       string
	PtsLose      string
	YardsWin     string
	YardsLose    string
	ToWin        string
	ToLose       string
}

// Note: ScrapeYear will include games that have not been played yet
// It is expected that the caller of this function will handle that accordingly
func ScrapeYear(year string) ([]ScrapedRow, error) {
	// Create a new collector
	c := colly.NewCollector()

	// Define the URL (replace YEAR with the desired season)
	url := fmt.Sprintf("https://www.pro-football-reference.com/years/%s/games.htm", year)

	scrapedRows := make([]ScrapedRow, 0)

	// Define what to do when visiting a row in the games table
	c.OnHTML("table#games tbody tr", func(e *colly.HTMLElement) {
		week := e.ChildText("th")
		date := e.ChildText("td[data-stat='game_date']")

		// In between each week, there is a header row
		// After the regular season, there is a header row for the playoffs
		// Make sure to skip these rows
		if week == "WeekDayDateTimeWinner/tieLoser/tiePtsWPtsLYdsWTOWYdsLTOL" ||
			date == "Playoffs" {
			return
		}

		// For now, ignore playoff games
		if week == "WildCard" || week == "Division" ||
			week == "ConfChamp" || week == "SuperBowl" {
			return
		}

		// Create a new scrapedRow instance
		row := ScrapedRow{
			Week:         week,
			DayOfWeek:    e.ChildText("td[data-stat='game_day_of_week']"),
			Date:         date,
			Gametime:     e.ChildText("td[data-stat='gametime']"),
			Winner:       e.ChildText("td[data-stat='winner']"),
			GameLocation: e.ChildText("td[data-stat='game_location']"),
			Loser:        e.ChildText("td[data-stat='loser']"),
			PtsWin:       e.ChildText("td[data-stat='pts_win']"),
			PtsLose:      e.ChildText("td[data-stat='pts_lose']"),
			YardsWin:     e.ChildText("td[data-stat='yards_win']"),
			YardsLose:    e.ChildText("td[data-stat='yards_lose']"),
			ToWin:        e.ChildText("td[data-stat='to_win']"),
			ToLose:       e.ChildText("td[data-stat='to_lose']"),
		}

		// Append the row to the slice
		scrapedRows = append(scrapedRows, row)
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		// TODO
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start scraping
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return scrapedRows, nil
}

func ParseTime(row ScrapedRow) time.Time {
	layout := "Mon 2006-01-02 3:04PM"

	// Match the layout using row fields
	input := fmt.Sprintf("%s %s %s", row.DayOfWeek, row.Date, row.Gametime)

	t, _ := time.Parse(layout, input) // TODO: handle error
	return t
}
