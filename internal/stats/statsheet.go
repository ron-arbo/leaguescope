package stats

type StatSheet struct {
	// Records
	Record           Record
	HomeRecord       Record
	AwayRecord       Record
	DivisionRecord   Record
	ConferenceRecord Record

	// Points
	Points           Points
	ConferencePoints Points

	// Strength Of
	StrengthOfVictory  float64
	StrengthOfSchedule float64

	// Rank
	ConferenceRankPointsFor     int
	ConferenceRankPointsAgainst int
	LeagueRankPointsFor         int
	LeagueRankPointsAgainst     int

	Streak int

	Seed int

	// TODO: Clincher
}

func NewStatSheet() StatSheet {
	return StatSheet{}
}
