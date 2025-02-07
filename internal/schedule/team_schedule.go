package schedule

import (
	"fmt"
	"nfl-app/internal/entry"
	"nfl-app/internal/game"
	"slices"
	"sort"
)

const (
	pointsFor     = "pointsFor"
	pointsAgainst = "pointsAgainst"
)

func CommonOpponents(teams []string, ts map[string]Schedule) []string {
	if len(teams) < 2 {
		// TODO
		return nil
	}

	oppMaps := make([]map[string][]game.Game, len(teams))
	for i, team := range teams {
		teamSchedule := ts[team]
		oppMaps[i] = teamSchedule.OpponentMapFor(team)
	}

	// Get all opponents for first team
	firstOppMap := oppMaps[0]
	firstOpps := make([]string, 0)
	for opp := range firstOppMap {
		firstOpps = append(firstOpps, opp)
	}

	// TODO: Maybe import some set.Intersection code here
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

func CommonGamesRecords(entries []entry.Entry, ts map[string]Schedule) map[string]*entry.Record {
	// Create slice of names
	var teamNames []string
	for _, entry := range entries {
		teamNames = append(teamNames, entry.Team.Name.String())
	}

	commonOpps := CommonOpponents(teamNames, ts)
	commonGamesRecord := make(map[string]*entry.Record)

	for _, team := range teamNames {
		record := &entry.Record{}
		teamSchedule := ts[team]
		opponentMap := teamSchedule.OpponentMapFor(team)

		for _, opp := range commonOpps {
			games := opponentMap[opp]
			for _, game := range games {
				switch {
				case game.Winner == team:
					record.AddWin()
				case game.Loser == team:
					record.AddLoss()
				default:
					record.AddTie()
				}
			}
		}
		commonGamesRecord[team] = record
	}

	return commonGamesRecord
}

func HeadToHeadGames(entries []entry.Entry, ts map[string]Schedule) []game.Game {
	// Create slice of team names
	teamNames := make([]string, len(entries))
	for i, team := range entries {
		teamNames[i] = team.Team.Name.String()
	}

	// Find all games between these teams
	h2hGamesSeen := make(map[string]bool)
	h2hGames := make([]game.Game, 0)
	for _, team := range teamNames {
		schedule := ts[team]
		for _, week := range schedule.Weeks {
			if len(week.Games) == 0 {
				// Bye week, continue
				continue
			}

			// Since we're dealing with team schedules, assume there is only 1 game
			game := week.Games[0]

			// Check if we've already seen this game
			if _, ok := h2hGamesSeen[fmt.Sprintf("%s@%s", game.Away, game.Home)]; ok {
				continue
			}

			// Check if the game is amongst these teams
			if slices.Contains(teamNames, game.Home) && slices.Contains(teamNames, game.Away) {
				h2hGames = append(h2hGames, game)
				// Mark the game as seen so we don't add it again when we loop through the other team's schedule
				// Note that we're relying on the fact that the game is unique by home and away team
				h2hGamesSeen[fmt.Sprintf("%s@%s", game.Away, game.Home)] = true
			}
		}
	}

	return h2hGames
}

func HeadToHeadRecords(entries []entry.Entry, ts map[string]Schedule) map[string]*entry.Record {
	h2hGames := HeadToHeadGames(entries, ts)
	h2hRecords := make(map[string]*entry.Record)

	// Initialize record for each team
	for _, ent := range entries {
		h2hRecords[ent.Team.Name.String()] = &entry.Record{}
	}

	for _, game := range h2hGames {
		homeRecord := h2hRecords[game.Home]
		awayRecord := h2hRecords[game.Away]

		switch {
		case game.Winner == game.Home:
			homeRecord.AddWin()
			awayRecord.AddLoss()
		case game.Winner == game.Away:
			homeRecord.AddLoss()
			awayRecord.AddWin()
		default:
			homeRecord.AddTie()
			awayRecord.AddTie()
		}

	}

	return h2hRecords
}

const (
	strengthOfVictory  = "victory"
	strengthOfSchedule = "schedule"
)

func StrengthOfVictory(team string, entries []entry.Entry, ts map[string]Schedule) float64 {
	return StrengthOf(team, entries, ts, strengthOfVictory)
}

func StrengthOfSchedule(team string, entries []entry.Entry, ts map[string]Schedule) float64 {
	return StrengthOf(team, entries, ts, strengthOfSchedule)
}

func StrengthOf(team string, entries []entry.Entry, ts map[string]Schedule, attribute string) float64 {
	if attribute != strengthOfVictory && attribute != strengthOfSchedule {
		panic("unknown attribute")
	}

	combinedOpponentRecord := entry.NewRecord(0, 0, 0)

	// Create map of team names to Entry for convenience
	teamMap := make(map[string]entry.Entry)
	for _, entry := range entries {
		teamMap[entry.Team.Name.String()] = entry
	}

	schedule := ts[team]
	for _, week := range schedule.Weeks {
		if len(week.Games) == 0 {
			continue
		}
		game := week.Games[0]

		// Get the opponent's record
		var opp string
		if game.Home == team {
			opp = game.Away
		} else {
			opp = game.Home
		}

		// For victory, only consider games that were won
		if attribute == strengthOfVictory && game.Winner != team {
			continue
		}

		// For schedule, consider all games
		oppEntry := teamMap[opp]
		oppRecord := entry.NewRecord(
			oppEntry.StatSheet.Record.Wins(),
			oppEntry.StatSheet.Record.Losses(),
			oppEntry.StatSheet.Record.Ties(),
		)

		combinedOpponentRecord.Add(oppRecord)
	}

	return combinedOpponentRecord.WinPercentage()
}

// Ranking returns a map of team names to their rank amongst the given entries based on the given stat
// Tied teams will shared the rank
func Ranking(entries []entry.Entry, stat string) map[string]int {
	if stat != pointsFor && stat != pointsAgainst {
		panic("unsupported stat")
	}

	var sortFunc func(i, j int) bool
	switch stat {
	case pointsFor:
		// Sort PointsFor in descending order
		sortFunc = func(i, j int) bool {
			return entries[i].StatSheet.Points.For > entries[j].StatSheet.Points.For
		}
	case pointsAgainst:
		// Sort PointsAgainst in ascending order
		sortFunc = func(i, j int) bool {
			return entries[i].StatSheet.Points.Against < entries[j].StatSheet.Points.Against
		}
	}
	sort.Slice(entries, sortFunc)

	// Assign rankings. If tied, give the same rank
	ranking := make(map[string]int)
	rank := 1
	for i, entry := range entries {
		// If the given entry is different than the last, we can increase the rank
		if stat == pointsFor {
			if i > 0 && entry.StatSheet.Points.For != entries[i-1].StatSheet.Points.For {
				rank = i + 1
			}
		} else {
			if i > 0 && entry.StatSheet.Points.Against != entries[i-1].StatSheet.Points.Against {
				rank = i + 1
			}
		}
		ranking[entry.TeamName()] = rank
	}

	return ranking
}
