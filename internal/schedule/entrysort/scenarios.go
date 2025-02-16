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

	// Patriots 1-1, 1-0 H2H
	// Jets     1-1, 0-1 H2H

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
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.NewYorkJets})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}

func Division2ClubsDivisionRecord() Scenario {
	sched := schedule.NewSchedule()

	// Pats 1-1, 1-0 in div
	// Jets 1-1, 0-1 in div

	// Week 1: Pats get their division win, Jets get their loss
	patsAtBills := game.Game{
		Winner: team.NewEnglandPatriots.Name,
		Loser:  team.BuffaloBills.Name,
		Home:   team.BuffaloBills.Name,
		Away:   team.NewEnglandPatriots.Name,
	}
	finsAtJets := game.Game{
		Winner: team.MiamiDolphins.Name,
		Loser:  team.NewYorkJets.Name,
		Home:   team.NewYorkJets.Name,
		Away:   team.MiamiDolphins.Name,
	}
	sched.AddGame(1, patsAtBills)
	sched.AddGame(1, finsAtJets)

	// Week 2: Pats lose and Jets win to even win pct. Both games out of division
	bengalsAtPats := game.Game{
		Winner: team.CincinnatiBengals.Name,
		Loser:  team.NewEnglandPatriots.Name,
		Home:   team.NewEnglandPatriots.Name,
		Away:   team.CincinnatiBengals.Name,
	}
	jetsAtSteelers := game.Game{
		Winner: team.NewYorkJets.Name,
		Loser:  team.PittsburghSteelers.Name,
		Home:   team.PittsburghSteelers.Name,
		Away:   team.NewYorkJets.Name,
	}
	sched.AddGame(2, bengalsAtPats)
	sched.AddGame(2, jetsAtSteelers)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.NewYorkJets})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}

func Division2ClubsCommonGames() Scenario {
	sched := schedule.NewSchedule()

	// Pats 1-1, 0-0 div, 1-0 v Texans
	// Jets 1-1, 0-0 div, 0-1 v Texans

	// Week 1: Pats get their Texans win. Jets also win
	patsAtTexans := game.Game{
		Winner: team.NewEnglandPatriots.Name,
		Loser:  team.HoustonTexans.Name,
		Home:   team.HoustonTexans.Name,
		Away:   team.NewEnglandPatriots.Name,
	}
	jetsAtSteelers := game.Game{
		Winner: team.NewYorkJets.Name,
		Loser:  team.PittsburghSteelers.Name,
		Home:   team.PittsburghSteelers.Name,
		Away:   team.NewYorkJets.Name,
	}
	sched.AddGame(1, patsAtTexans)
	sched.AddGame(1, jetsAtSteelers)

	// Week 2: Jets get their Texans loss. Pats also lose
	bengalsAtPats := game.Game{
		Winner: team.CincinnatiBengals.Name,
		Loser:  team.NewEnglandPatriots.Name,
		Home:   team.NewEnglandPatriots.Name,
		Away:   team.CincinnatiBengals.Name,
	}
	texansAtJets := game.Game{
		Winner: team.HoustonTexans.Name,
		Loser:  team.NewYorkJets.Name,
		Home:   team.NewYorkJets.Name,
		Away:   team.HoustonTexans.Name,
	}
	sched.AddGame(2, bengalsAtPats)
	sched.AddGame(2, texansAtJets)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.NewYorkJets})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}

func Division2ClubsConferenceRecord() Scenario {
	sched := schedule.NewSchedule()

	// Pats 1-1, 1-0 in conf
	// Jets 1-1, 0-1 in conf

	// Week 1: Pats get their conference win, Jets get their loss (both non-divisional)
	patsAtTexans := game.Game{
		Winner: team.NewEnglandPatriots.Name,
		Loser:  team.HoustonTexans.Name,
		Home:   team.HoustonTexans.Name,
		Away:   team.NewEnglandPatriots.Name,
	}
	jetsAtSteelers := game.Game{
		Winner: team.PittsburghSteelers.Name,
		Loser:  team.NewYorkJets.Name,
		Home:   team.PittsburghSteelers.Name,
		Away:   team.NewYorkJets.Name,
	}
	sched.AddGame(1, patsAtTexans)
	sched.AddGame(1, jetsAtSteelers)

	// Week 2: Pats lose and Jets win to even win pct. Both games out of conference
	cardinalsAtPats := game.Game{
		Winner: team.ArizonaCardinals.Name,
		Loser:  team.NewEnglandPatriots.Name,
		Home:   team.NewEnglandPatriots.Name,
		Away:   team.ArizonaCardinals.Name,
	}
	jetsAtGiants := game.Game{
		Winner: team.NewYorkJets.Name,
		Loser:  team.NewYorkGiants.Name,
		Home:   team.NewYorkJets.Name,
		Away:   team.NewYorkGiants.Name,
	}
	sched.AddGame(2, cardinalsAtPats)
	sched.AddGame(2, jetsAtGiants)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.NewYorkJets})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}

func Division2ClubsStrengthOfVictory() Scenario {
	sched := schedule.NewSchedule()

	// Pats 2-0, Opponents beaten 2-2
	// Jets 2-0, Opponents beaten 0-4

	// Week 1: Pats and Jets both win.
	// Seahawks (Pats future opponent) win, Rams (Jets future opponent) lose
	patsAtCards := game.Game{
		Winner: team.NewEnglandPatriots.Name,
		Loser:  team.ArizonaCardinals.Name,
		Home:   team.ArizonaCardinals.Name,
		Away:   team.NewEnglandPatriots.Name,
	}
	jetsAtGiants := game.Game{
		Winner: team.NewYorkJets.Name,
		Loser:  team.NewYorkGiants.Name,
		Home:   team.NewYorkGiants.Name,
		Away:   team.NewYorkJets.Name,
	}
	seahawksAtRams := game.Game{
		Winner: team.SeattleSeahawks.Name,
		Loser:  team.LosAngelesRams.Name,
		Home:   team.LosAngelesRams.Name,
		Away:   team.SeattleSeahawks.Name,
	}
	sched.AddGame(1, patsAtCards)
	sched.AddGame(1, jetsAtGiants)
	sched.AddGame(1, seahawksAtRams)

	// Week 2: Pats and Jets both win
	// Cardinals (Pats previous opponent) win, Giants (Jets previous opponent) lose
	seahawksAtPats := game.Game{
		Winner: team.NewEnglandPatriots.Name,
		Loser:  team.SeattleSeahawks.Name,
		Home:   team.NewEnglandPatriots.Name,
		Away:   team.SeattleSeahawks.Name,
	}
	ramsAtJets := game.Game{
		Winner: team.NewYorkJets.Name,
		Loser:  team.LosAngelesRams.Name,
		Home:   team.NewYorkJets.Name,
		Away:   team.LosAngelesRams.Name,
	}
	cardinalsAtGiants := game.Game{
		Winner: team.ArizonaCardinals.Name,
		Loser:  team.NewYorkGiants.Name,
		Home:   team.NewYorkGiants.Name,
		Away:   team.ArizonaCardinals.Name,
	}
	sched.AddGame(2, seahawksAtPats)
	sched.AddGame(2, ramsAtJets)
	sched.AddGame(2, cardinalsAtGiants)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.NewYorkJets})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}

func Division2ClubsStrengthOfSchedule() Scenario {
	sched := schedule.NewSchedule()

	// Pats 1-1, Opponents 2-2
	// Jets 1-1, Opponents 1-3

	// Week 1: Pats and Jets both win.
	// Seahawks (Pats future opponent) win, Rams (Jets future opponent) lose
	patsAtCards := game.Game{
		Winner: team.NewEnglandPatriots.Name,
		Loser:  team.ArizonaCardinals.Name,
		Home:   team.ArizonaCardinals.Name,
		Away:   team.NewEnglandPatriots.Name,
	}
	jetsAtGiants := game.Game{
		Winner: team.NewYorkJets.Name,
		Loser:  team.NewYorkGiants.Name,
		Home:   team.NewYorkGiants.Name,
		Away:   team.NewYorkJets.Name,
	}
	seahawksAtRams := game.Game{
		Winner: team.SeattleSeahawks.Name,
		Loser:  team.LosAngelesRams.Name,
		Home:   team.LosAngelesRams.Name,
		Away:   team.SeattleSeahawks.Name,
	}
	sched.AddGame(1, patsAtCards)
	sched.AddGame(1, jetsAtGiants)
	sched.AddGame(1, seahawksAtRams)

	// Week 2: Pats and Jets both lose. Pats to a better team
	// Cardinals (Pats previous opponent) and Giants (Jets previous opponent) both lose to keep SOV the same
	seahawksAtPats := game.Game{
		Winner: team.SeattleSeahawks.Name,
		Loser:  team.NewEnglandPatriots.Name,
		Home:   team.NewEnglandPatriots.Name,
		Away:   team.SeattleSeahawks.Name,
	}
	ramsAtJets := game.Game{
		Winner: team.LosAngelesRams.Name,
		Loser:  team.NewYorkJets.Name,
		Home:   team.NewYorkJets.Name,
		Away:   team.LosAngelesRams.Name,
	}
	cardinalsAtPanthers := game.Game{
		Winner: team.CarolinaPanthers.Name,
		Loser:  team.ArizonaCardinals.Name,
		Home:   team.CarolinaPanthers.Name,
		Away:   team.ArizonaCardinals.Name,
	}
	giantsAtBears := game.Game{
		Winner: team.ChicagoBears.Name,
		Loser:  team.NewYorkGiants.Name,
		Home:   team.ChicagoBears.Name,
		Away:   team.NewYorkGiants.Name,
	}
	sched.AddGame(2, seahawksAtPats)
	sched.AddGame(2, ramsAtJets)
	sched.AddGame(2, cardinalsAtPanthers)
	sched.AddGame(2, giantsAtBears)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.NewYorkJets})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}

func Division2ClubsConferenceRank() Scenario {
	sched := schedule.NewSchedule()

	// Pats 1-0, 30 points for, 20 points against
	// Jets 1-0, 29 points for, 21 points against

	// Week 1: Pats and Jets both win. Pats score more and allow less
	patsAtCards := game.Game{
		Winner:  team.NewEnglandPatriots.Name,
		Loser:   team.ArizonaCardinals.Name,
		Home:    team.ArizonaCardinals.Name,
		Away:    team.NewEnglandPatriots.Name,
		PtsWin:  30,
		PtsLose: 20,
	}
	jetsAtGiants := game.Game{
		Winner:  team.NewYorkJets.Name,
		Loser:   team.NewYorkGiants.Name,
		Home:    team.NewYorkGiants.Name,
		Away:    team.NewYorkJets.Name,
		PtsWin:  29,
		PtsLose: 21,
	}
	sched.AddGame(1, patsAtCards)
	sched.AddGame(1, jetsAtGiants)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.NewYorkJets})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}

func Division2ClubsLeagueRank() Scenario {
	sched := schedule.NewSchedule()

	// Points For:
	// 1. Patriots (AFC) - 40
	// 2. Cardinals  (NFC) - 30
	// 3. Giants (NFC) - 25
	// 4. Jets     (AFC) - 20
	// Bear - 16

	// Points Against:
	// Bear - 15
	// 1. Jets     (AFC) - 17
	// 2. Giants (NFC) - 20
	// 3. Patriots (AFC) - 30
	// 4. Cardinals  (NFC) - 55

	// Week 1:
	patsAtCards := game.Game{
		Winner:  team.NewEnglandPatriots.Name,
		Loser:   team.ArizonaCardinals.Name,
		Home:    team.ArizonaCardinals.Name,
		Away:    team.NewEnglandPatriots.Name,
		PtsWin:  40,
		PtsLose: 30,
	}
	jetsAtGiants := game.Game{
		Winner:  team.NewYorkJets.Name,
		Loser:   team.NewYorkGiants.Name,
		Home:    team.NewYorkGiants.Name,
		Away:    team.NewYorkJets.Name,
		PtsWin:  20,
		PtsLose: 17,
	}
	sched.AddGame(1, patsAtCards)
	sched.AddGame(1, jetsAtGiants)

	// Week 2: Giants overtake Jets in PF, lose to keep SOV/SOS same
	giantsAtBears := game.Game{
		Winner:  team.ChicagoBears.Name,
		Loser:   team.NewYorkGiants.Name,
		Home:    team.NewYorkGiants.Name,
		Away:    team.ChicagoBears.Name,
		PtsWin:  16,
		PtsLose: 15,
	}
	sched.AddGame(2, giantsAtBears)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.NewYorkJets})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}

func Division2ClubsNetPointsCommonGames() Scenario {
	sched := schedule.NewSchedule()

	// Pats 1-1, 0-0 div, 1-0 v Texans (+10 net)
	// Jets 1-1, 0-0 div, 1-0 v Texans (+5 net)

	// Week 1: Pats get their Texans win. Jets lose
	patsAtTexans := game.Game{
		Winner:  team.NewEnglandPatriots.Name,
		Loser:   team.HoustonTexans.Name,
		Home:    team.HoustonTexans.Name,
		Away:    team.NewEnglandPatriots.Name,
		PtsWin:  30,
		PtsLose: 20,
	}
	jetsAtSteelers := game.Game{
		Winner: team.PittsburghSteelers.Name,
		Loser:  team.NewYorkJets.Name,
		Home:   team.PittsburghSteelers.Name,
		Away:   team.NewYorkJets.Name,
	}
	sched.AddGame(1, patsAtTexans)
	sched.AddGame(1, jetsAtSteelers)

	// Week 2: Jets get their Texans win. Pats lose
	bengalsAtPats := game.Game{
		Winner: team.CincinnatiBengals.Name,
		Loser:  team.NewEnglandPatriots.Name,
		Home:   team.NewEnglandPatriots.Name,
		Away:   team.CincinnatiBengals.Name,
	}
	texansAtJets := game.Game{
		Winner:  team.NewYorkJets.Name,
		Loser:   team.HoustonTexans.Name,
		Home:    team.NewYorkJets.Name,
		Away:    team.HoustonTexans.Name,
		PtsWin:  20,
		PtsLose: 15,
	}
	sched.AddGame(2, bengalsAtPats)
	sched.AddGame(2, texansAtJets)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.NewYorkJets})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}

func Division2ClubsNetPoints() Scenario {
	sched := schedule.NewSchedule()

	// Pats 0-1, -5 net
	// Jets 0-1, -10 net

	// Week 1: Pats and Jets both lost, Pats better diff
	patsAtCards := game.Game{
		Winner:  team.ArizonaCardinals.Name,
		Loser:   team.NewEnglandPatriots.Name,
		Home:    team.ArizonaCardinals.Name,
		Away:    team.NewEnglandPatriots.Name,
		PtsWin:  30,
		PtsLose: 25,
	}
	jetsAtRams := game.Game{
		Winner:  team.LosAngelesRams.Name,
		Loser:   team.NewYorkJets.Name,
		Home:    team.LosAngelesRams.Name,
		Away:    team.NewYorkJets.Name,
		PtsWin:  20,
		PtsLose: 10,
	}
	sched.AddGame(1, patsAtCards)
	sched.AddGame(1, jetsAtRams)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.NewYorkJets})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}
