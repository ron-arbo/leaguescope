package entrysort

import (
	"fmt"
	"nfl-app/internal/entry"
	"nfl-app/internal/schedule"
	"sort"
)

// SortEntries will sort the given entries by win percentage, default to tiebreakers specified
// in https://www.nfl.com/standings/tie-breaking-procedures
func SortEntries(entries []entry.Entry, teamSchedules map[string]schedule.Schedule) ([]entry.Entry, error) {
	if len(entries) < 2 {
		return entries, nil
	}

	sorter, err := GetSorterFor(entries)
	if err != nil {
		return nil, err
	}
	return sorter.Sort(entries, teamSchedules), nil
}

// One general sort function for any sorter to use
func (s *Sorter) Sort(entries []entry.Entry, teamSchedules map[string]schedule.Schedule) []entry.Entry {
	// Validate
	if len(entries) < 2 {
		return entries
	}
	fmt.Printf("Sorting %s by %s\n", entry.Teams(entries), s.Name)

	// Get sortBy
	sortBy := s.SortByMap(entries, teamSchedules)
	for sortByKey, sortByValue := range sortBy {
		fmt.Printf("%s has %f\n", sortByKey, sortByValue)
	}

	// Check is sorted already?

	// Sort the entries in descending order using sortBy
	sort.Slice(entries, func(i, j int) bool {
		return sortBy[entries[i].TeamName()] > sortBy[entries[j].TeamName()]
	})

	// Now, figure out how we want to handle ties
	// Two methods:
	// 1. Elimination method: Eliminate one team at a time, either the best or the worst. If a team is eliminated, sort the remaining teams again from the start
	// 						  This is done for most tiebreakers, as defined in the NFL rulebook
	// 2. Subgroup method: Split the tied teams into subgroups based on their sortBy value. Sort each subgroup using the tiebreaker. Combine the sorted subgroups
	// 					   This is done for win percentage, because we don't want to sort the entire group again after identifying the top team, we know their win
	//					   percentage doesn't change based on the teams included (while something like head-to-head might)
	// 3. Double elimination method: Similar to elimination method, but only consider the top teams from each division before moving to tiebreakers. From there, find
	//							     a best team and sort the rest (including the original eliminated teams) using the root sort. Cannot find a worst team because
	//							     it may not actually be the worst team. The worst team may have been eliminated when finding the top division teams.
	switch s.TiebreakMethod {
	case subgroup:
		return s.SubgroupSort(entries, teamSchedules, sortBy)
	case elimination:
		return s.EliminationSort(entries, teamSchedules, sortBy)
	case doubleElimination:
		return s.DoubleEliminationSort(entries, teamSchedules, sortBy)
	case tripleElimination:
		return TripleEliminationSort(entries, teamSchedules)
	}

	panic(fmt.Sprintf("Unknown tiebreak method %s for sorter %s", s.TiebreakMethod, s.Name))
}

// SubgroupSort is used specifically at the top level for win percentage.
// It will split the entries into subgroups based on the sortBy value and sort each subgroup using the tiebreaker
func (s *Sorter) SubgroupSort(entries []entry.Entry, teamSchedules map[string]schedule.Schedule, sortBy map[string]float64) []entry.Entry {
	var sortedEntries []entry.Entry
	var sortedSubgroup []entry.Entry

	// Split entries into subgroups based on sortBy value
	subgroups := entry.GroupEntries(entries, sortBy)

	for _, subgroup := range subgroups {
		// Sort the subgroup based on the number of entries
		if len(subgroup) == 1 {
			sortedSubgroup = subgroup
		} else {
			// Sort the subgroup using the tiebreaker
			sortedSubgroup = s.Tiebreaker.Sort(subgroup, teamSchedules)
		}

		// Append the sorted subgroup to the sortedEntries slice
		sortedEntries = append(sortedEntries, sortedSubgroup...)
	}

	return sortedEntries
}

// EliminationSort is the most commonly used tiebreaker method.
// When evaluating a group of teams, EliminationSort will attempt to isolate either the sole best or sole worst team in the group.
// If one of those teams is identified, its spot will be locked in and all of the remaining teams will be sorted again, starting from the beginning.
// If the sole worst/best is not found, the top group of tied entries will be sorted using the tiebreaker to find the sole best.
func (s *Sorter) EliminationSort(entries []entry.Entry, teamSchedules map[string]schedule.Schedule, sortBy map[string]float64) []entry.Entry {
	subgroups := entry.GroupEntries(entries, sortBy)
	topGroup := subgroups[0]
	bottomGroup := subgroups[len(subgroups)-1]

	if len(topGroup) == 1 {
		// There is a best team. Award it first position and sort the rest using the root sort
		restSorted, _ := SortEntries(entries[1:], teamSchedules) // TODO: Handle err
		return append(topGroup, restSorted...)
	}

	if len(bottomGroup) == 1 {
		// There is a worst team. Award it last position and sort the rest using the root sort
		restSorted, _ := SortEntries(entries[:len(entries)-1], teamSchedules) // TODO: Handle err
		return append(restSorted, bottomGroup...)
	}

	// We could not find a best/worst team using this Sorter
	var topGroupSorted []entry.Entry
	if len(topGroup) == len(entries) {
		// If the top group is the same as the original entries, revert to the tiebreaker Sorter
		topGroupSorted = s.Tiebreaker.Sort(topGroup, teamSchedules)
	} else {
		// If the top group is a subset of the original entries, it may require a different Sorter. Use the general SortEntries
		topGroupSorted, _ = SortEntries(topGroup, teamSchedules) // TODO: Handle err
	}

	// Separate the top entry from the others
	topEntry := topGroupSorted[:1]
	otherEntries := append(topGroupSorted[1:], entries[len(topGroup):]...)

	// Sort the remaining teams using the root sort and combine with the top entry
	otherEntriesSorted, _ := SortEntries(otherEntries, teamSchedules) // TODO: Handle err
	return append(topEntry, otherEntriesSorted...)
}

// DoubleEliminationSort is EliminationSort, however we make sure to eliminate any team that is not the top ranked team in their division
// before determining the sole best. Because of this, we neglect determining the sole worst team.
func (s *Sorter) DoubleEliminationSort(entries []entry.Entry, teamSchedules map[string]schedule.Schedule, sortBy map[string]float64) []entry.Entry {
	// Before eliminating the sole best or worst, we need to eliminate
	// any team that is not the top ranked team in their division
	divTop, divBottom := FindDivisionTopTeams(entries, teamSchedules)

	fmt.Printf("Top teams: %s\n", entry.Teams(divTop))
	fmt.Printf("Eliminated teams: %s\n\n", entry.Teams(divBottom))

	// Now basically just run EliminationSort on the top teams
	subgroups := entry.GroupEntries(divTop, sortBy)
	topGroup := subgroups[0]
	bottomGroup := subgroups[len(subgroups)-1]

	if len(topGroup) == 1 {
		// There is a best team. Award it first position and sort the rest using the root sort
		rest := append(divTop[1:], divBottom...)
		restSorted, _ := SortEntries(rest, teamSchedules) // TODO: Handle err
		return append(topGroup, restSorted...)
	}

	// Only eliminate the worst team if we didn't eliminate any teams when finding the top division teams,
	// otherwise we might eliminate a team that is not actually the worst
	if len(bottomGroup) == 1 && len(divBottom) == 0 {
		// There is a worst team. Award it last position and sort the rest using the root sort
		rest := divTop[:len(divTop)-1]
		restSorted, _ := SortEntries(rest, teamSchedules) // TODO: Handle err
		return append(restSorted, bottomGroup...)
	}

	// TODO: We probably don't need to do the whole sort below, since we sort all but the top team
	// in the next round anyway. We just need to find a way to get the top team, maybe we could be faster

	// We could not find a best/worst team using this Sorter. Sort the top tied group to find one
	var topGroupSorted []entry.Entry
	if len(topGroup) == len(entries) {
		// The top group is the same as the original entries, revert to the tiebreaker Sorter
		topGroupSorted = s.Tiebreaker.Sort(topGroup, teamSchedules)
	} else {
		// The top group is a subset of the original entries, it may require a different Sorter
		topGroupSorted, _ = SortEntries(topGroup, teamSchedules) // TODO: Handle err
	}

	// Separate the top entry from the others
	topEntry := topGroupSorted[:1]
	restOfDivTop := append(topGroupSorted[1:], divTop[len(topGroup):]...)
	restOfEntries := append(restOfDivTop, divBottom...)

	// Sort the remaining teams using the root sort and combine with the top entry
	restOfEntriesSorted, _ := SortEntries(restOfEntries, teamSchedules) // TODO: Handle err
	return append(topEntry, restOfEntriesSorted...)
}

// TODO: This is just an inverse of the process used to determine draft order.
// It's not entirely accurate because the draft order adds an initial tiebreaker to all ties:
// 3. If ties exist in any grouping, such ties shall be broken by figuring the aggregate won-lost-tied percentage of each involved club's regular-season opponents and awarding preferential selection order to the club that faced the schedule of teams with the lowest aggregate won-lost-tied percentage.
func TripleEliminationSort(entries []entry.Entry, teamSchedules map[string]schedule.Schedule) []entry.Entry {
	// (ii) Ties involving THREE-OR-MORE clubs from different conferences will be broken by applying
	// (a) divisional tiebreakers to determine the lowest-ranked team in a division
	divisionGroups := entry.GroupByDivision(entries)
	var remaining []entry.Entry
	worstInDivision := make([]entry.Entry, 0)
	for _, teams := range divisionGroups {
		if len(teams) == 1 {
			worstInDivision = append(worstInDivision, teams[0])
			continue
		}

		// Sort the teams and take the worst one
		sortedTeams, _ := SortEntries(teams, teamSchedules) // TODO: Handle err

		// Split the slice into the worst team and the remaining teams
		remainingTeams, worstTeam := entry.SplitAround(sortedTeams, len(sortedTeams)-1)

		worstInDivision = append(worstInDivision, worstTeam...)
		remaining = append(remaining, remainingTeams...)
	}

	// (b) conference tiebreakers to determine the lowest-ranked team within a conference, and
	conferenceGroups := entry.GroupByConference(worstInDivision)
	worstInConference := make([]entry.Entry, 0)
	for _, teams := range conferenceGroups {
		if len(teams) == 1 {
			worstInConference = append(worstInConference, teams[0])
			continue
		}

		// Sort the teams and take the worst one
		sortedTeams, _ := SortEntries(teams, teamSchedules) // TODO: Handle err

		// Split the slice into the worst team and the remaining teams
		remainingTeams, worstTeam := entry.SplitAround(sortedTeams, len(sortedTeams)-1)

		worstInConference = append(worstInConference, worstTeam...)
		remaining = append(remaining, remainingTeams...)
	}

	// (c) interconference tiebreakers to determine the lowest ranked team in the league
	// Sort the teams and take the worst one
	sortedTeams, _ := SortEntries(worstInConference, teamSchedules) // TODO: Handle err

	// Split the slice into the worst team and the remaining teams
	remainingTeams, worstTeam := entry.SplitAround(sortedTeams, len(sortedTeams)-1)

	worstInLeague := worstTeam
	remaining = append(remaining, remainingTeams...)

	remainingSorted, _ := SortEntries(remaining, teamSchedules) // TODO: Handle error
	return append(remainingSorted, worstInLeague...)
}

// FindDivisionTopTeams will return a subset of the original group of entries containing all the highest ranked teams
// in each division, as well as a subset of the original group of entries containing all the other teams in each division
func FindDivisionTopTeams(entries []entry.Entry, teamSchedules map[string]schedule.Schedule) ([]entry.Entry, []entry.Entry) {
	top := make([]entry.Entry, 0)
	other := make([]entry.Entry, 0)

	divisionGroups := entry.GroupByDivision(entries)
	for _, group := range divisionGroups {
		// If there is only one team in the division, they are the top team
		if len(group) == 1 {
			top = append(top, group[0])
			continue
		}

		// Otherwise, sort the teams and take the top one
		sortedGroup, _ := SortEntries(group, teamSchedules) // TODO: Handle err

		topTeam, otherTeams := entry.SplitAround(sortedGroup, 1)
		top = append(top, topTeam...)
		other = append(other, otherTeams...)
	}

	return top, other
}

// Elimination vs Subgroup sort:

// We must start from the beginning because the elimination of a team may change the result of previously used tiebreakers, such as H2H/comon games.
// For example, consider Team A, Team B, and Team C. They're equal in win percentage, and they're each 1-1 against the others.
// The next tiebreaker for wildcard spots is conference record. Let's say Team A has the best, Teams B and C are tied for second.
// We identify team Team A as the team to "eliminate". We do so and sort B and C from the beginning, using H2H as the tiebreaker.
// Well Team B's loss came to team A and Team C's win came to Team A. Now Team B has the better H2H record and is ranked ahead of Team C.
// But consider if we did not eliminate and start from the beginning and instead just sorted B and C by the next tiebreaker because they were
// tie for second. Then we find out Team C has a better record in common games (the next tiebreaker) than Team B. Team C is ranked ahead of Team B in
// this scenario, whereas Team B is ranked ahead of Team C in the other.

// SeedEntries will sort the entries as they would be seeded in the playoffs
// This means that the top team in each division is seeded 1-4, and the rest are seeded 5+
func SeedEntries(entries []entry.Entry, ts map[string]schedule.Schedule) ([]entry.Entry, error) {
	// Sort entries by win percentage
	sortedEntries, err := SortEntries(entries, ts)
	if err != nil {
		return nil, err
	}

	// TODO: For now assume they are in the same conference
	divSeed := 1
	wildSeed := 5

	divisionsSeen := make(map[string]bool)
	for _, entry := range sortedEntries {
		// If we haven't seen the division yet, we know this is a division winner,
		// so it assign it to the highest available seed in 1-4. Otherwise assign it to 5+
		if _, ok := divisionsSeen[entry.Team.Division]; !ok {
			entry.StatSheet.Seed = divSeed
			divSeed++
			divisionsSeen[entry.Team.Division] = true
		} else {
			entry.StatSheet.Seed = wildSeed
			wildSeed++
		}
	}

	return sortedEntries, nil
}
