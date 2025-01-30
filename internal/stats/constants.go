package stats

// TODO: Can we make win/losses/ties --> Record struct? And pointsFor/pointsAgainst/pointDifferential --> Points struct?

const (
	// Overall
	StatWins       string = "wins"
	StatLosses     string = "losses"
	StatTies       string = "ties"
	StatWinPercent string = "winPercent"

	// Home
	StatHomeWins   string = "homeWins"
	StatHomeLosses string = "homeLosses"
	StatHomeTies   string = "homeTies"

	// Road
	StatRoadWins   string = "roadWins"
	StatRoadLosses string = "roadLosses"
	StatRoadTies   string = "roadTies"

	// Division
	StatDivisionWins   string = "divisionWins"
	StatDivisionLosses string = "divisionLosses"
	StatDivisionTies   string = "divisionTies"

	// Conference
	StatConferenceWins   string = "conferenceWins"
	StatConferenceLosses string = "conferenceLosses"
	StatConferenceTies   string = "conferenceTies"

	// Points
	StatPointsFor         string = "pointsFor"
	StatPointsAgainst     string = "pointsAgainst"
	StatPointDifferential string = "pointDifferential"

	// Conference Points
	StatConferencePointsFor         string = "conferencePointsFor"
	StatConferencePointsAgainst     string = "conferencePointsAgainst"
	StatConferencePointDifferential string = "conferencePointDifferential"

	// Strength Of
	StatStrengthOfVictory  string = "strengthOfVictory"
	StatStrengthOfSchedule string = "strengthOfSchedule"

	// Rank
	ConferenceRankPointsFor     string = "conferenceRankPointsFor"
	ConferenceRankPointsAgainst string = "conferenceRankPointsAgainst"
	LeagueRankPointsFor         string = "leagueRankPointsFor"
	LeagueRankPointsAgainst     string = "leagueRankPointsAgainst"

	StatStreak string = "streak"

	StatClincher    string = "clincher"
	StatPlayoffSeed string = "playoffSeed"
)

var (
	statNames = []string{
		StatWins,
		StatLosses,
		StatTies,
		StatWinPercent,
		StatHomeWins,
		StatHomeLosses,
		StatHomeTies,
		StatRoadWins,
		StatRoadLosses,
		StatRoadTies,
		StatDivisionWins,
		StatDivisionLosses,
		StatDivisionTies,
		StatConferenceWins,
		StatConferenceLosses,
		StatConferenceTies,
		StatPointsFor,
		StatPointsAgainst,
		StatPointDifferential,
		StatConferencePointsFor,
		StatConferencePointsAgainst,
		StatConferencePointDifferential,
		StatStrengthOfVictory,
		StatStrengthOfSchedule,
		ConferenceRankPointsFor,
		ConferenceRankPointsAgainst,
		LeagueRankPointsFor,
		LeagueRankPointsAgainst,
		StatStreak,
		StatClincher,
		StatPlayoffSeed,
	}
)
