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

		if week.Game.HomeTeam.Name.String() == ts.Team {
			opponentMap[week.Game.AwayTeam.Name.String()] = append(opponentMap[week.Game.AwayTeam.Name.String()], *week.Game)
		} else {
			opponentMap[week.Game.HomeTeam.Name.String()] = append(opponentMap[week.Game.HomeTeam.Name.String()], *week.Game)
		}
	}

	return opponentMap
}

func CommonOpponents(teams []string, teamSchedules map[string]*TeamSchedule) []string {
	if len(teams) < 2 {
		// TODO
		return nil
	}

	opponentMaps := make([]map[string][]game.Game, len(teams))
	for i, team := range teams {
		opponentMaps[i] = teamSchedules[team].OpponentMap()
	}

	// Get all opponents for first team
	firstOpponentMap := opponentMaps[0]
	firstOpponents := make([]string, len(firstOpponentMap))
	i := 0
	for k := range firstOpponentMap {
		firstOpponents[i] = k
		i++
	}

	commonOpponents := make([]string, 0)
	for _, opp := range firstOpponents {
		// Check if all other teams have played this opponent
		played := true
		for _, opponentMap := range opponentMaps[1:] {
			if _, ok := opponentMap[opp]; !ok {
				played = false
				break
			}
		}

		if played {
			// All teams have played this opponent
			commonOpponents = append(commonOpponents, opp)
		}
	}

	return commonOpponents
}

func CommonGamesRecords(entries []entry.Entry, teamSchedules map[string]*TeamSchedule) map[string]*entry.Record {
	// Create slice of names
	var teamNames []string
	for _, entry := range entries {
		teamNames = append(teamNames, entry.TeamName())
	}

	commonOpponents := CommonOpponents(teamNames, teamSchedules)
	commonGamesRecord := make(map[string]*entry.Record)

	for _, team := range teamNames {
		record := &entry.Record{}
		opponentMap := teamSchedules[team].OpponentMap()

		for _, opp := range commonOpponents {
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
	headToHeadGamesSeen := make(map[string]bool)
	headToHeadGames := make([]*game.Game, 0)
	for _, team := range teamNames {
		schedule := teamSchedules[team]
		for _, week := range schedule.Weeks {
			if week.Game == nil || !week.Game.Completed {
				continue
			}

			// Check if we've already seen this game
			if _, ok := headToHeadGamesSeen[fmt.Sprintf("%s@%s", week.Game.HomeTeam.Name.String(), week.Game.AwayTeam.Name.String())]; ok {
				continue
			}

			// Check if the game is amongst these teams
			if slices.Contains(teamNames, week.Game.HomeTeam.Name.String()) &&
				slices.Contains(teamNames, week.Game.AwayTeam.Name.String()) {
				headToHeadGames = append(headToHeadGames, week.Game)
				// Mark the game as seen so we don't add it again when we loop through the other team's schedule
				// Note that we're relying on the fact that the game is unique by home and away team
				headToHeadGamesSeen[fmt.Sprintf("%s@%s", week.Game.HomeTeam.Name.String(), week.Game.AwayTeam.Name.String())] = true
			}
		}
	}

	return headToHeadGames
}

// HeadToHeadGames returns a map of team names to their record against one another
func HeadToHeadRecords(teams []entry.Entry, teamSchedules map[string]*TeamSchedule) map[string]*entry.Record {
	headToHeadGames := HeadToHeadGames(teams, teamSchedules)
	headToHeadRecords := make(map[string]*entry.Record)

	// Initialize record for each team
	for _, team := range teams {
		headToHeadRecords[team.Team.Name.String()] = &entry.Record{}
	}

	for _, game := range headToHeadGames {
		switch {
		case game.HomeScore > game.AwayScore:
			headToHeadRecords[game.HomeTeam.Name.String()].AddWin()
			headToHeadRecords[game.AwayTeam.Name.String()].AddLoss()
		case game.HomeScore < game.AwayScore:
			headToHeadRecords[game.HomeTeam.Name.String()].AddLoss()
			headToHeadRecords[game.AwayTeam.Name.String()].AddWin()
		default:
			headToHeadRecords[game.HomeTeam.Name.String()].AddTie()
			headToHeadRecords[game.AwayTeam.Name.String()].AddTie()
		}
	}

	return headToHeadRecords
}
