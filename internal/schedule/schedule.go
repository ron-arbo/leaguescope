package schedule

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"nfl-app/internal/game"
	"nfl-app/internal/team"
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
