package entrysort

import (
	"nfl-app/internal/entry"
	"nfl-app/internal/game"
	"nfl-app/internal/schedule"
	"nfl-app/internal/team"
)

// This file contains a series of "scenarios" (a combination of entries and team schedules)
//
// Each scenario is meant to test a specific tiebreaker scenario from https://www.nfl.com/standings/tie-breaking-procedures

type Scenario struct {
	entries   []entry.Entry
	schedules map[string]schedule.Schedule
}

// Two clubs, within division
// Head to Head
func Division2ClubsHeadToHead() Scenario {
	sched := schedule.NewSchedule()
	// Bills    2-0
	// Patriots 1-1
	// Jets     1-1
	// Dolphins 0-2
	// TODO: Should we only include the entries for patriots and jets?

	jetsAtPats := game.Game{
		Winner: team.NewEnglandPatriots.Name,
		Loser:  team.NewYorkJets.Name,
		Home:   team.NewEnglandPatriots.Name,
		Away:   team.NewYorkJets.Name,
	}
	billsAtFins := game.Game{
		Winner: team.BuffaloBills.Name,
		Loser:  team.MiamiDolphins.Name,
		Home:   team.MiamiDolphins.Name,
		Away:   team.BuffaloBills.Name,
	}
	sched.AddGame(1, jetsAtPats)
	sched.AddGame(2, billsAtFins)

	patsAtBills := game.Game{
		Winner: team.BuffaloBills.Name,
		Loser:  team.NewEnglandPatriots.Name,
		Home:   team.BuffaloBills.Name,
		Away:   team.NewEnglandPatriots.Name,
	}
	finsAtJets := game.Game{
		Winner: team.NewYorkJets.Name,
		Loser:  team.MiamiDolphins.Name,
		Home:   team.NewYorkJets.Name,
		Away:   team.MiamiDolphins.Name,
	}
	sched.AddGame(2, patsAtBills)
	sched.AddGame(2, finsAtJets)

	entries := schedule.CreateEntries(sched)
	return Scenario{
		entries:   entries,
		schedules: sched.SplitToTeams(),
	}
}
