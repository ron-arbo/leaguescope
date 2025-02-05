package stats

type StatSheet struct {
	// Records
	Record Record
	HomeRecord Record
	AwayRecord Record
	DivisionRecord Record
	ConferenceRecord Record

	// Points
	Points Points
	ConferencePoints Points

	// Strength Of
	StrengthOfVictory float64
	strengthOfSchedule  float64

	// Rank
	ConferenceRankPoints int
	LeagueRankPoints     int

	Streak int

	Seed int

	// TODO: Clincher
}

func NewStatSheet() StatSheet{
	return StatSheet{}
}