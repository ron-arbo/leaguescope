package entrysort

import (
	"nfl-app/internal/entry"
	"nfl-app/internal/schedule"
)

// This file contains function to get the sortBy maps needed for sorting entries
// Many are attainable from the entry itself, but some require the teamSchedules map (common games, h2h, etc)

func WinPercentageMap(entries []entry.Entry, ts map[string]schedule.Schedule) map[string]float64 {
	sortBy := make(map[string]float64)
	for _, entry := range entries {
		sortBy[entry.Team.Name.String()] = entry.StatSheet.Record.WinPercentage()
	}

	return sortBy
}

func HeadToHeadMap(entries []entry.Entry, ts map[string]schedule.Schedule) map[string]float64 {
	headToHeadRecords := schedule.HeadToHeadRecords(entries, ts)

	sortBy := make(map[string]float64)
	for team, record := range headToHeadRecords {
		sortBy[team] = record.WinPercentage()
	}

	return sortBy
}

func HeadToHeadSweepMap(entries []entry.Entry, ts map[string]schedule.Schedule) map[string]float64 {
	headToHeadRecords := schedule.HeadToHeadRecords(entries, ts)

	sortBy := make(map[string]float64)
	for team, record := range headToHeadRecords {
		// We should only be considering non-divisional teams here, consider step 1 of the 3+ club conference tiebreaker
		// So, we can assume that in order for a team to have played every other team, they must have len(entries)-1 games played
		// If you haven't played every other team, you can't have swept. So continue in that case
		if record.GamesPlayed() != len(entries)-1 {
			sortBy[team] = 0.5
			continue
		}
		// Give 1.000 to the team that has beaten *every* other team
		// Give 0.000 to the team that has lost to *every* other team
		// Give .500 to every other tea
		switch record.WinPercentage() {
		case 1.0:
			sortBy[team] = 1.0
		case 0.0:
			sortBy[team] = 0.0
		default:
			sortBy[team] = 0.5
		}
	}

	return sortBy
}

func DivisionMap(entries []entry.Entry, ts map[string]schedule.Schedule) map[string]float64 {
	sortBy := make(map[string]float64)
	for _, entry := range entries {
		sortBy[entry.Team.Name.String()] = entry.StatSheet.Record.WinPercentage()
	}
	return sortBy
}

func ConferenceMap(entries []entry.Entry, ts map[string]schedule.Schedule) map[string]float64 {
	sortBy := make(map[string]float64)
	for _, entry := range entries {
		sortBy[entry.Team.Name.String()] = entry.StatSheet.ConferenceRecord.WinPercentage()
	}
	return sortBy
}

func CommonGamesMap(entries []entry.Entry, ts map[string]schedule.Schedule) map[string]float64 {
	commonGamesRecords := schedule.CommonGamesRecords(entries, ts)

	sortBy := make(map[string]float64)
	for team, record := range commonGamesRecords {
		sortBy[team] = record.WinPercentage()
	}

	return sortBy
}

func CommonGamesMin4Map(entries []entry.Entry, ts map[string]schedule.Schedule) map[string]float64 {
	commonGamesRecords := schedule.CommonGamesRecords(entries, ts)

	sortBy := make(map[string]float64)
	var minGamesPlayed = true
	for team, record := range commonGamesRecords {
		if record.GamesPlayed() < 4 {
			// Does not meet minimum games played, mark
			minGamesPlayed = false
			break
		}
		sortBy[team] = record.WinPercentage()
	}

	// If not enough games played, return a sortBy that will trigger the next tiebreaker
	if !minGamesPlayed {
		for team := range commonGamesRecords {
			sortBy[team] = -1 // Use negative to indicate not enough games played
		}
	}

	return sortBy
}

func CoinTossMap(entries []entry.Entry, ts map[string]schedule.Schedule) map[string]float64 {
	coinTossMap := make(map[string]float64)
	for i, entry := range entries {
		coinTossMap[entry.Team.Name.String()] = float64(i)
	}
	return coinTossMap
}

func NetPointsCommonMap(entries []entry.Entry, ts map[string]schedule.Schedule) map[string]float64 {
	var teams []string
	for _, entry := range entries {
		teams = append(teams, entry.Team.Name.String())
	}

	commonOpponents := schedule.CommonOpponents(teams, ts)
	commonGamesNet := make(map[string]float64)

	for _, team := range teams {
		teamNet := 0
		teamSchedule := ts[team]
		opponentMap := teamSchedule.OpponentMapFor(team)

		for _, opp := range commonOpponents {
			games := opponentMap[opp]
			for _, game := range games {
				if game.Winner == team {
					teamNet += game.PtsWin - game.PtsLose
				} else {
					teamNet += game.PtsLose - game.PtsWin
				}
			}
		}
		commonGamesNet[team] = float64(teamNet)
	}

	return commonGamesNet
}

func NetPointsMap(entries []entry.Entry, ts map[string]schedule.Schedule) map[string]float64 {
	sortBy := make(map[string]float64)
	for _, entry := range entries {
		sortBy[entry.Team.Name.String()] = float64(entry.StatSheet.Points.Differential())
	}

	return sortBy
}

func NetPointsConferenceMap(entries []entry.Entry, ts map[string]schedule.Schedule) map[string]float64 {
	sortBy := make(map[string]float64)
	for _, entry := range entries {
		sortBy[entry.Team.Name.String()] = float64(entry.StatSheet.ConferencePoints.Differential())
	}

	return sortBy
}

func StrengthOfVictoryMap(entries []entry.Entry, ts map[string]schedule.Schedule) map[string]float64 {
	sovMap := make(map[string]float64)
	for _, entry := range entries {
		sovMap[entry.Team.Name.String()] = entry.StatSheet.StrengthOfVictory
	}

	return sovMap
}

func StrengthOfScheduleMap(entries []entry.Entry, ts map[string]schedule.Schedule) map[string]float64 {
	sosMap := make(map[string]float64)
	for _, entry := range entries {
		sosMap[entry.Team.Name.String()] = entry.StatSheet.StrengthOfSchedule
	}

	return sosMap
}

func CombinedRankingConferenceMap(entries []entry.Entry, ts map[string]schedule.Schedule) map[string]float64 {
	sortBy := make(map[string]float64)

	for _, entry := range entries {
		// Calculate the combined ranking for the entry
		pfRank := entry.StatSheet.ConferenceRankPointsFor
		paRank := entry.StatSheet.ConferenceRankPointsAgainst

		combinedRanking := pfRank + paRank

		// We want to reward lower combined rankings, so multiply by -1 before adding to sortBy
		sortBy[entry.TeamName()] = float64(combinedRanking * -1)
	}

	return sortBy
}

func CombinedRankingLeagueMap(entries []entry.Entry, ts map[string]schedule.Schedule) map[string]float64 {
	sortBy := make(map[string]float64)

	for _, entry := range entries {
		// Calculate the combined ranking for the entry
		pfRank := entry.StatSheet.LeagueRankPointsFor
		paRank := entry.StatSheet.LeagueRankPointsAgainst

		combinedRanking := pfRank + paRank

		// We want to reward lower combined rankings, so multiply by -1 before adding to sortBy
		sortBy[entry.TeamName()] = float64(combinedRanking * -1)
	}

	return sortBy
}
