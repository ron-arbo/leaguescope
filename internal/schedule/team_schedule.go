package schedule

import (
	"fmt"
	"nfl-app/internal/entry"
	"nfl-app/internal/game"
	"nfl-app/internal/team"
	"slices"
)

// Holds all team schedules, will take multiple API calls to populate
type TeamSchedule struct {
	// TODO: Do we want to convert this to Team object?
	Team  string     // Team name
	Weeks []TeamWeek // Weeks 1-18, with index=week-1
}

type TeamWeek struct {
	WeekNumber int
	Game       *game.Game
}

func NewTeamSchedules() map[string]*TeamSchedule {
	teamSchedules := make(map[string]*TeamSchedule)

	for _, team := range team.NFLTeams {
		teamSchedules[team.Name.String()] = &TeamSchedule{
			Team:  team.Name.String(),
			Weeks: make([]TeamWeek, 18),
		}
	}

	return teamSchedules
}

// ToStandings converts a TeamSchedule to Entry for that team
func (ts *TeamSchedule) ToEntry() *entry.Entry {
	entry := entry.NewEntry(ts.Team)

	for _, week := range ts.Weeks {
		if week.Game == nil {
			// Bye week, continue
			continue
		}

		if week.Game.Completed {
			// Updates counting stats only
			entry.UpdateStats(week.Game)
		}
	}

	return entry
}

func (ts *TeamSchedule) OpponentMap() map[string][]game.Game {
	// Need a string of games because teams will play division opponents twice
	opponentMap := make(map[string][]game.Game)
	for _, week := range ts.Weeks {
		if week.Game == nil || !week.Game.Completed {
			continue
		}

		if week.Game.HomeTeamName() == ts.Team {
			opponentMap[week.Game.AwayTeamName()] = append(opponentMap[week.Game.AwayTeamName()], *week.Game)
		} else {
			opponentMap[week.Game.HomeTeamName()] = append(opponentMap[week.Game.HomeTeamName()], *week.Game)
		}
	}

	return opponentMap
}

func CommonOpponents(teams []string, teamSchedules map[string]*TeamSchedule) []string {
	if len(teams) < 2 {
		// TODO
		return nil
	}

	oppMaps := make([]map[string][]game.Game, len(teams))
	for i, team := range teams {
		oppMaps[i] = teamSchedules[team].OpponentMap()
	}

	// Get all opponents for first team
	firstOppMap := oppMaps[0]
	firstOpps := make([]string, 0)
	for opp := range firstOppMap {
		firstOpps = append(firstOpps, opp)
	}

	commonOpps := make([]string, 0)
	for _, opp := range firstOpps {
		// Check if all other teams have played this opponent
		played := true
		for _, opponentMap := range oppMaps[1:] {
			if _, ok := opponentMap[opp]; !ok {
				played = false
				break
			}
		}

		if played {
			// All teams have played this opponent
			commonOpps = append(commonOpps, opp)
		}
	}

	return commonOpps
}

func CommonGamesRecords(entries []entry.Entry, teamSchedules map[string]*TeamSchedule) map[string]*entry.Record {
	// Create slice of names
	var teamNames []string
	for _, entry := range entries {
		teamNames = append(teamNames, entry.TeamName())
	}

	commonOpps := CommonOpponents(teamNames, teamSchedules)
	commonGamesRecord := make(map[string]*entry.Record)

	for _, team := range teamNames {
		record := &entry.Record{}
		opponentMap := teamSchedules[team].OpponentMap()

		for _, opp := range commonOpps {
			games := opponentMap[opp]
			for _, game := range games {
				if game.Completed {
					if game.Tied() {
						record.AddTie()
					}
					if game.WonBy(team) {
						record.AddWin()
					} else {
						record.AddLoss()
					}
				}
			}
		}
		commonGamesRecord[team] = record
	}

	return commonGamesRecord
}

// HeadToHeadGames returns a slice containg all Games between the given teams
func HeadToHeadGames(teams []entry.Entry, teamSchedules map[string]*TeamSchedule) []*game.Game {
	// Create slice of team names
	teamNames := make([]string, len(teams))
	for i, team := range teams {
		teamNames[i] = team.Team.Name.String()
	}

	// Find all games between these teams
	h2hGamesSeen := make(map[string]bool)
	h2hGames := make([]*game.Game, 0)
	for _, team := range teamNames {
		schedule := teamSchedules[team]
		for _, week := range schedule.Weeks {
			if week.Game == nil || !week.Game.Completed {
				continue
			}

			// Check if we've already seen this game
			if _, ok := h2hGamesSeen[fmt.Sprintf("%s@%s", week.Game.HomeTeamName(), week.Game.AwayTeamName())]; ok {
				continue
			}

			// Check if the game is amongst these teams
			if slices.Contains(teamNames, week.Game.HomeTeamName()) &&
				slices.Contains(teamNames, week.Game.AwayTeamName()) {
				h2hGames = append(h2hGames, week.Game)
				// Mark the game as seen so we don't add it again when we loop through the other team's schedule
				// Note that we're relying on the fact that the game is unique by home and away team
				h2hGamesSeen[fmt.Sprintf("%s@%s", week.Game.HomeTeamName(), week.Game.AwayTeamName())] = true
			}
		}
	}

	return h2hGames
}

// HeadToHeadGames returns a map of team names to their record against one another
func HeadToHeadRecords(entries []entry.Entry, teamSchedules map[string]*TeamSchedule) map[string]*entry.Record {
	h2hGames := HeadToHeadGames(entries, teamSchedules)
	h2hRecords := make(map[string]*entry.Record)

	// Initialize record for each team
	for _, ent := range entries {
		h2hRecords[ent.Team.Name.String()] = &entry.Record{}
	}

	for _, game := range h2hGames {
		homeRecord := h2hRecords[game.HomeTeamName()]
		awayRecord := h2hRecords[game.AwayTeamName()]

		homeRecord.UpdateFor(game.HomeTeamName(), game)
		awayRecord.UpdateFor(game.AwayTeamName(), game)
	}

	return h2hRecords
}
