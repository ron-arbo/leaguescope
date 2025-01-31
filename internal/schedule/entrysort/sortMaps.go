package entrysort

import (
	"nfl-app/internal/entry"
	"nfl-app/internal/schedule"
	"nfl-app/internal/stats"
)

// This file contains function to get the sortBy maps needed for sorting entries
// Many are attainable from the entry itself, but some require the teamSchedules map (common games, h2h, etc)

func WinPercentageMap(entries []entry.Entry, teamSchedules map[string]*schedule.TeamSchedule) map[string]float64 {
	sortBy := make(map[string]float64)
	for _, entry := range entries {
		sortBy[entry.TeamName()] = entry.OverallRecord().WinPercentage()
	}

	return sortBy
}

func HeadToHeadMap(entries []entry.Entry, teamSchedules map[string]*schedule.TeamSchedule) map[string]float64 {
	headToHeadRecords := schedule.HeadToHeadRecords(entries, teamSchedules)

	sortBy := make(map[string]float64)
	for team, record := range headToHeadRecords {
		sortBy[team] = record.WinPercentage()
	}

	return sortBy
}

// "Sweep": Applicable only if one club has defeated each of the others or if one club has lost to each of the others.
func HeadToHeadSweepMap(entries []entry.Entry, teamSchedules map[string]*schedule.TeamSchedule) map[string]float64 {
	headToHeadRecords := schedule.HeadToHeadRecords(entries, teamSchedules)

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

func DivisionMap(entries []entry.Entry, teamSchedules map[string]*schedule.TeamSchedule) map[string]float64 {
	divisionRecords := entry.DivisionRecords(entries)

	sortBy := make(map[string]float64)
	for team, record := range divisionRecords {
		sortBy[team] = record.WinPercentage()
	}

	return sortBy
}

func ConferenceMap(entries []entry.Entry, teamSchedules map[string]*schedule.TeamSchedule) map[string]float64 {
	conferenceRecords := entry.ConferenceRecords(entries)

	sortBy := make(map[string]float64)
	for team, record := range conferenceRecords {
		sortBy[team] = record.WinPercentage()
	}

	return sortBy
}

func CommonGamesMap(entries []entry.Entry, teamSchedules map[string]*schedule.TeamSchedule) map[string]float64 {
	commonGamesRecords := schedule.CommonGamesRecords(entries, teamSchedules)

	sortBy := make(map[string]float64)
	for team, record := range commonGamesRecords {
		sortBy[team] = record.WinPercentage()
	}

	return sortBy
}

func CommonGamesMin4Map(entries []entry.Entry, teamSchedules map[string]*schedule.TeamSchedule) map[string]float64 {
	commonGamesRecords := schedule.CommonGamesRecords(entries, teamSchedules)

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

// Just return as is for now
func CoinTossMap(entries []entry.Entry, teamSchedules map[string]*schedule.TeamSchedule) map[string]float64 {
	coinTossMap := make(map[string]float64)
	for i, entry := range entries {
		coinTossMap[entry.TeamName()] = float64(i)
	}
	return coinTossMap
}

func NetPointsCommonMap(entries []entry.Entry, teamSchedules map[string]*schedule.TeamSchedule) map[string]float64 {
	var teams []string
	for _, entry := range entries {
		teams = append(teams, entry.TeamName())
	}

	commonOpponents := schedule.CommonOpponents(teams, teamSchedules)
	commonGamesNet := make(map[string]float64)

	for _, team := range teams {
		teamNet := 0
		opponentMap := teamSchedules[team].OpponentMap()

		for _, opp := range commonOpponents {
			games := opponentMap[opp]
			for _, game := range games {
				if game.Completed {
					if game.HomeTeamName() == team {
						teamNet += game.HomeScore - game.AwayScore
					} else {
						teamNet += game.AwayScore - game.HomeScore
					}
				}
			}
		}
		commonGamesNet[team] = float64(teamNet)
	}

	return commonGamesNet
}

func NetPointsMap(entries []entry.Entry, teamSchedules map[string]*schedule.TeamSchedule) map[string]float64 {
	netPointsMap := make(map[string]float64)
	for _, entry := range entries {
		netPointsMap[entry.TeamName()] = entry.Stats[stats.StatPointDifferential].Value()
	}

	return netPointsMap
}

func NetPointsConferenceMap(entries []entry.Entry, teamSchedules map[string]*schedule.TeamSchedule) map[string]float64 {
	netPointsConfMap := make(map[string]float64)
	for _, entry := range entries {
		netPointsConfMap[entry.TeamName()] = entry.Stats[stats.StatConferencePointDifferential].Value()
	}

	return netPointsConfMap
}

func StrengthOfVictoryMap(entries []entry.Entry, teamSchedules map[string]*schedule.TeamSchedule) map[string]float64 {
	sovMap := make(map[string]float64)
	for _, entry := range entries {
		sovMap[entry.TeamName()] = entry.Stats[stats.StatStrengthOfVictory].Value()
	}

	return sovMap
}

func StrengthOfScheduleMap(entries []entry.Entry, teamSchedules map[string]*schedule.TeamSchedule) map[string]float64 {
	sosMap := make(map[string]float64)
	for _, entry := range entries {
		sosMap[entry.TeamName()] = entry.Stats[stats.StatStrengthOfSchedule].Value()
	}

	return sosMap
}

func CombinedRankingConferenceMap(entries []entry.Entry, teamSchedules map[string]*schedule.TeamSchedule) map[string]float64 {
	sortBy := make(map[string]float64)

	for _, entry := range entries {
		// Calculate the combined ranking for the entry
		pointsForRanking := entry.Stats[stats.ConferenceRankPointsFor].Value()
		pointsAgainstRanking := entry.Stats[stats.ConferenceRankPointsAgainst].Value()

		combinedRanking := pointsForRanking + pointsAgainstRanking

		// We want to reward lower combined rankings, so multiply by -1 before adding to sortBy
		sortBy[entry.TeamName()] = (combinedRanking * -1)
	}

	return sortBy
}

func CombinedRankingLeagueMap(entries []entry.Entry, teamSchedules map[string]*schedule.TeamSchedule) map[string]float64 {
	sortBy := make(map[string]float64)

	for _, entry := range entries {
		// Calculate the combined ranking for the entry
		pointsForRanking := entry.Stats[stats.LeagueRankPointsFor].Value()
		pointsAgainstRanking := entry.Stats[stats.LeagueRankPointsAgainst].Value()

		combinedRanking := pointsForRanking + pointsAgainstRanking

		// We want to reward lower combined rankings, so multiply by -1 before adding to sortBy
		sortBy[entry.TeamName()] = (combinedRanking * -1)
	}

	return sortBy
}
