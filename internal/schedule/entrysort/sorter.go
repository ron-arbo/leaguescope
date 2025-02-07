package entrysort

import (
	"fmt"
	"nfl-app/internal/entry"
	"nfl-app/internal/schedule"
)

type Sorter struct {
	// Name is what we are sorting by, useful for debugging
	Name string

	// SortByMap is a function that will return a map of team names to a value that this Sorter will sort by
	SortByMap func([]entry.Entry, map[string]schedule.Schedule) map[string]float64

	// Tiebreaker is the next Sorter to use if there are ties
	Tiebreaker *Sorter

	TiebreakMethod string // "subgroup" | "elimination" |  "double elimination" | "triple elimination"

	// Maybe a Validate() function to ensure the entries passed are valid for this sorter?

	// Add a cache for recent results?
}

const (
	// Tiebreak methods
	subgroup          = "subgroup"
	elimination       = "elimination"
	doubleElimination = "double elimination"
	tripleElimination = "triple elimination"
)

// According to NFL tiebreaker rules (https://www.nfl.com/standings/tie-breaking-procedures),
// there are X sorting groups to consider
// 2 clubs - Within division
// 3+ clubs - Within division
// 2 clubs - Outside division but within conference
// 3+ clubs - Outside division but within conference
// 2 clubs - Across conferences (for draft)
// 3+ clubs - Across conferences (for draft)

func GetSorterFor(entries []entry.Entry) (*Sorter, error) {
	if len(entries) < 2 {
		return nil, fmt.Errorf("not enough entries to define sort group: %v", entries)
	}

	conferenceSeen := make(map[string]bool)
	divisionSeen := make(map[string]bool)

	for _, entry := range entries {
		conferenceSeen[entry.Team.Conference] = true
		divisionSeen[entry.Team.Division] = true
	}

	// Validate
	if len(conferenceSeen) == 0 || len(conferenceSeen) > 2 || len(divisionSeen) == 0 || len(divisionSeen) > 8 {
		return nil, fmt.Errorf("invalid number of conferences/divisions in entries: %v", entries)
	}

	if len(conferenceSeen) == 2 {
		if len(entries) == 2 {
			return WithinLeagueTwoClubsSorter(), nil
		}
		return WithinLeagueThreeClubsSorter(), nil
	}

	if len(divisionSeen) >= 2 {
		if len(entries) == 2 {
			return WithinConferenceTwoClubsSorter(), nil
		}
		return WithinConferenceThreeClubsSorter(), nil
	}

	if len(entries) == 2 {
		return WithinDivisionTwoClubsSorter(), nil
	}
	return WithinDivisionThreeClubsSorter(), nil
}

func WithinDivisionTwoClubsSorter() *Sorter {
	div2ClubsSorter := &Sorter{
		Name:           "Division 2 Clubs",
		SortByMap:      WinPercentageMap,
		Tiebreaker:     nil,
		TiebreakMethod: subgroup,
	}

	h2hSorter := &Sorter{
		Name:           "Head to Head",
		SortByMap:      HeadToHeadMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	divisionRecordSorter := &Sorter{
		Name:           "Division Record",
		SortByMap:      DivisionMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	commonRecordSorter := &Sorter{
		Name:           "Common Games",
		SortByMap:      CommonGamesMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	conferenceRecordSorter := &Sorter{
		Name:           "Conference Record",
		SortByMap:      ConferenceMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	strengthOfVictorySorter := &Sorter{
		Name:           "Strength of Victory",
		SortByMap:      StrengthOfVictoryMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	strengthOfScheduleSorter := &Sorter{
		Name:           "Strength of Schedule",
		SortByMap:      StrengthOfScheduleMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	combinedRankConferenceSorter := &Sorter{
		Name:           "Combined Rank Conference",
		SortByMap:      CombinedRankingConferenceMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	combinedRankAllSorter := &Sorter{
		Name:           "Combined Rank All",
		SortByMap:      CombinedRankingConferenceMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	netPointsCommonSorter := &Sorter{
		Name:           "Net points in common games",
		SortByMap:      NetPointsCommonMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	netPointsSorter := &Sorter{
		Name:           "Net Points",
		SortByMap:      NetPointsMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	// TODO: touchdowns

	coinTossSorter := &Sorter{
		Name:           "Coin Toss",
		SortByMap:      CoinTossMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	// Set tiebreakers
	div2ClubsSorter.Tiebreaker = h2hSorter
	h2hSorter.Tiebreaker = divisionRecordSorter
	divisionRecordSorter.Tiebreaker = commonRecordSorter
	commonRecordSorter.Tiebreaker = conferenceRecordSorter
	conferenceRecordSorter.Tiebreaker = strengthOfVictorySorter
	strengthOfVictorySorter.Tiebreaker = strengthOfScheduleSorter
	strengthOfScheduleSorter.Tiebreaker = combinedRankConferenceSorter
	combinedRankConferenceSorter.Tiebreaker = combinedRankAllSorter
	combinedRankAllSorter.Tiebreaker = netPointsCommonSorter
	netPointsCommonSorter.Tiebreaker = netPointsSorter
	netPointsSorter.Tiebreaker = coinTossSorter

	return div2ClubsSorter
}

func WithinDivisionThreeClubsSorter() *Sorter {
	div3ClubsSorter := &Sorter{
		Name:           "Division 3 Clubs",
		SortByMap:      WinPercentageMap,
		Tiebreaker:     nil,
		TiebreakMethod: subgroup,
	}

	h2hSorter := &Sorter{
		Name:       "Head to Head",
		SortByMap:  HeadToHeadMap,
		Tiebreaker: nil,
	}

	divisionRecordSorter := &Sorter{
		Name:       "Division Record",
		SortByMap:  DivisionMap,
		Tiebreaker: nil,
	}

	commonRecordSorter := &Sorter{
		Name:       "Common Games",
		SortByMap:  CommonGamesMap,
		Tiebreaker: nil,
	}

	conferenceRecordSorter := &Sorter{
		Name:       "Conference Record",
		SortByMap:  ConferenceMap,
		Tiebreaker: nil,
	}

	strengthOfVictorySorter := &Sorter{
		Name:       "Strength of Victory",
		SortByMap:  StrengthOfVictoryMap,
		Tiebreaker: nil,
	}

	strengthOfScheduleSorter := &Sorter{
		Name:       "Strength of Schedule",
		SortByMap:  StrengthOfScheduleMap,
		Tiebreaker: nil,
	}

	combinedRankConferenceSorter := &Sorter{
		Name:       "Combined Rank Conference",
		SortByMap:  CombinedRankingConferenceMap,
		Tiebreaker: nil,
	}

	combinedRankAllSorter := &Sorter{
		Name:       "Combined Rank All",
		SortByMap:  CombinedRankingLeagueMap,
		Tiebreaker: nil,
	}

	netPointsCommonSorter := &Sorter{
		Name:       "Net points in common games",
		SortByMap:  NetPointsCommonMap,
		Tiebreaker: nil,
	}

	netPointsSorter := &Sorter{
		Name:       "Net Points",
		SortByMap:  NetPointsMap,
		Tiebreaker: nil,
	}

	// TODO: touchdowns

	coinTossSorter := &Sorter{
		Name:       "Coin Toss",
		SortByMap:  CoinTossMap,
		Tiebreaker: nil,
	}

	// Set tiebreakers
	div3ClubsSorter.Tiebreaker = h2hSorter
	h2hSorter.Tiebreaker = divisionRecordSorter
	divisionRecordSorter.Tiebreaker = commonRecordSorter
	commonRecordSorter.Tiebreaker = conferenceRecordSorter
	conferenceRecordSorter.Tiebreaker = strengthOfVictorySorter
	strengthOfVictorySorter.Tiebreaker = strengthOfScheduleSorter
	strengthOfScheduleSorter.Tiebreaker = combinedRankConferenceSorter
	combinedRankConferenceSorter.Tiebreaker = combinedRankAllSorter
	combinedRankAllSorter.Tiebreaker = netPointsCommonSorter
	netPointsCommonSorter.Tiebreaker = netPointsSorter
	// TODO: touchdowns
	netPointsSorter.Tiebreaker = coinTossSorter

	return div3ClubsSorter
}

func WithinConferenceTwoClubsSorter() *Sorter {
	conf2ClubsSorter := &Sorter{
		Name:           "Conference 2 Clubs",
		SortByMap:      WinPercentageMap,
		Tiebreaker:     nil,
		TiebreakMethod: subgroup,
	}

	h2hSorter := &Sorter{
		Name:           "Head to Head",
		SortByMap:      HeadToHeadMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	conferenceRecordSorter := &Sorter{
		Name:           "Conference Record",
		SortByMap:      ConferenceMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	commonRecordSorter := &Sorter{
		Name:           "Common Games",
		SortByMap:      CommonGamesMin4Map, // Minimum of 4 common games
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	strengthOfVictorySorter := &Sorter{
		Name:           "Strength of Victory",
		SortByMap:      StrengthOfVictoryMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	strengthOfScheduleSorter := &Sorter{
		Name:           "Strength of Schedule",
		SortByMap:      StrengthOfScheduleMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	combinedRankConferenceSorter := &Sorter{
		Name:           "Best combined ranking among conference teams in points scored and points allowed in all games",
		SortByMap:      CombinedRankingConferenceMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	combinedRankAllSorter := &Sorter{
		Name:           "Best combined ranking among all teams in points scored and points allowed in all games",
		SortByMap:      CombinedRankingLeagueMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	netPointsConferenceSorter := &Sorter{
		Name:           "Net points in conference games",
		SortByMap:      NetPointsConferenceMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	netPointsSorter := &Sorter{
		Name:           "Net Points",
		SortByMap:      NetPointsMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	// TODO: touchdowns

	coinTossSorter := &Sorter{
		Name:           "Coin Toss",
		SortByMap:      CoinTossMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	// Set tiebreakers
	conf2ClubsSorter.Tiebreaker = h2hSorter
	h2hSorter.Tiebreaker = conferenceRecordSorter
	conferenceRecordSorter.Tiebreaker = commonRecordSorter
	commonRecordSorter.Tiebreaker = strengthOfVictorySorter
	strengthOfVictorySorter.Tiebreaker = strengthOfScheduleSorter
	strengthOfScheduleSorter.Tiebreaker = combinedRankConferenceSorter
	combinedRankConferenceSorter.Tiebreaker = combinedRankAllSorter
	combinedRankAllSorter.Tiebreaker = netPointsConferenceSorter
	netPointsConferenceSorter.Tiebreaker = netPointsSorter
	// TODO: Touchdowns
	netPointsSorter.Tiebreaker = coinTossSorter

	return conf2ClubsSorter
}

func WithinConferenceThreeClubsSorter() *Sorter {
	conf3ClubsSorter := &Sorter{
		Name:           "Conference 3+ Clubs",
		SortByMap:      WinPercentageMap,
		Tiebreaker:     nil,
		TiebreakMethod: subgroup,
	}

	// Create Sorter for each tiebreaker and link accordingly
	h2hSweepSorter := &Sorter{
		Name:           "Head to Head sweep",
		SortByMap:      HeadToHeadSweepMap,
		Tiebreaker:     nil,
		TiebreakMethod: doubleElimination,
	}

	conferenceRecordSorter := &Sorter{
		Name:           "Conference Record",
		SortByMap:      ConferenceMap,
		Tiebreaker:     nil,
		TiebreakMethod: doubleElimination,
	}

	commonRecordSorter := &Sorter{
		Name:           "Common Games",
		SortByMap:      CommonGamesMin4Map, // Minimum of 4 common games
		Tiebreaker:     nil,
		TiebreakMethod: doubleElimination,
	}

	strengthOfVictorySorter := &Sorter{
		Name:           "Strength of Victory",
		SortByMap:      StrengthOfVictoryMap,
		Tiebreaker:     nil,
		TiebreakMethod: doubleElimination,
	}

	strengthOfScheduleSorter := &Sorter{
		Name:           "Strength of Schedule",
		SortByMap:      StrengthOfScheduleMap,
		Tiebreaker:     nil,
		TiebreakMethod: doubleElimination,
	}

	combinedRankConferenceSorter := &Sorter{
		Name:           "Combined Rank Conference",
		SortByMap:      CombinedRankingConferenceMap,
		Tiebreaker:     nil,
		TiebreakMethod: doubleElimination,
	}

	combinedRankAllSorter := &Sorter{
		Name:           "Combined Rank All",
		SortByMap:      CombinedRankingLeagueMap,
		Tiebreaker:     nil,
		TiebreakMethod: doubleElimination,
	}

	netPointsConferenceSorter := &Sorter{
		Name:           "Net points in conference games",
		SortByMap:      NetPointsConferenceMap,
		Tiebreaker:     nil,
		TiebreakMethod: doubleElimination,
	}

	netPointsSorter := &Sorter{
		Name:           "Net Points",
		SortByMap:      NetPointsMap,
		Tiebreaker:     nil,
		TiebreakMethod: doubleElimination,
	}

	// TODO: touchdowns

	coinTossSorter := &Sorter{
		Name:           "Coin Toss",
		SortByMap:      CoinTossMap,
		Tiebreaker:     nil,
		TiebreakMethod: doubleElimination,
	}

	// Set tiebreakers
	conf3ClubsSorter.Tiebreaker = h2hSweepSorter
	h2hSweepSorter.Tiebreaker = conferenceRecordSorter
	conferenceRecordSorter.Tiebreaker = commonRecordSorter
	commonRecordSorter.Tiebreaker = strengthOfVictorySorter
	strengthOfVictorySorter.Tiebreaker = strengthOfScheduleSorter
	strengthOfScheduleSorter.Tiebreaker = combinedRankConferenceSorter
	combinedRankConferenceSorter.Tiebreaker = combinedRankAllSorter
	combinedRankAllSorter.Tiebreaker = netPointsConferenceSorter
	netPointsConferenceSorter.Tiebreaker = netPointsSorter
	// TODO: touchdowns
	netPointsSorter.Tiebreaker = coinTossSorter

	return conf3ClubsSorter

}

func WithinLeagueTwoClubsSorter() *Sorter {
	league2ClubsSorter := &Sorter{
		Name:           "League 2 Clubs",
		SortByMap:      WinPercentageMap,
		Tiebreaker:     nil,
		TiebreakMethod: subgroup,
	}

	// Create Sorter for each tiebreaker and link accordingly
	h2hSorter := &Sorter{
		Name:           "Head to Head",
		SortByMap:      HeadToHeadMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	commonRecordSorter := &Sorter{
		Name:           "Common Games",
		SortByMap:      CommonGamesMin4Map, // Minimum of 4 common games
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	strengthOfVictorySorter := &Sorter{
		Name:           "Strength of Victory",
		SortByMap:      StrengthOfVictoryMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	combinedRankAllSorter := &Sorter{
		Name:           "Combined Rank All",
		SortByMap:      CombinedRankingLeagueMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	netPointsSorter := &Sorter{
		Name:           "Net Points",
		SortByMap:      NetPointsMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	// TODO: touchdowns

	coinTossSorter := &Sorter{
		Name:           "Coin Toss",
		SortByMap:      CoinTossMap,
		Tiebreaker:     nil,
		TiebreakMethod: elimination,
	}

	// Set tiebreakers
	league2ClubsSorter.Tiebreaker = h2hSorter
	h2hSorter.Tiebreaker = commonRecordSorter
	commonRecordSorter.Tiebreaker = strengthOfVictorySorter
	strengthOfVictorySorter.Tiebreaker = combinedRankAllSorter
	combinedRankAllSorter.Tiebreaker = netPointsSorter
	netPointsSorter.Tiebreaker = coinTossSorter

	return league2ClubsSorter
}

func WithinLeagueThreeClubsSorter() *Sorter {
	league3ClubsSorter := &Sorter{
		Name:           "League 3 Clubs",
		SortByMap:      WinPercentageMap,
		Tiebreaker:     nil,
		TiebreakMethod: tripleElimination,
	}

	return league3ClubsSorter
}
