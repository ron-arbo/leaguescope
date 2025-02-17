package entrysort

import (
	"nfl-app/internal/entry"
	"nfl-app/internal/game"
	"nfl-app/internal/schedule"
	"nfl-app/internal/team"
)

// This file contains a series of "scenarios" (a combination of entries and team schedules)
//
// Each scenario is meant to test a specific tiebreaker scenario
// from https://www.nfl.com/standings/tie-breaking-procedures

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
	// 1. Patriots  (AFC) - 40
	// 2. Cardinals (NFC) - 30
	// 3. Giants    (NFC) - 25
	// 4. Jets      (AFC) - 20
	// 5. Bears     (NFC) - 16

	// Points Against:
	// 1. Bears     (NFC) - 15
	// 1. Jets      (AFC) - 17
	// 2. Giants    (NFC) - 20
	// 3. Patriots  (AFC) - 30
	// 4. Cardinals (NFC) - 55

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

func Conference2ClubsHeadtoHead() Scenario {
	sched := schedule.NewSchedule()

	// Patriots 1-1, 1-0 H2H
	// Raiders  1-1, 0-1 H2H

	// Week 1: Pats beat Raiders
	raidersAtPats := game.Game{
		Winner: team.NewEnglandPatriots.Name,
		Loser:  team.LasVegasRaiders.Name,
		Home:   team.NewEnglandPatriots.Name,
		Away:   team.LasVegasRaiders.Name,
	}
	billsAtFins := game.Game{
		Winner: team.BuffaloBills.Name,
		Loser:  team.MiamiDolphins.Name,
		Home:   team.MiamiDolphins.Name,
		Away:   team.BuffaloBills.Name,
	}
	sched.AddGame(1, raidersAtPats)
	sched.AddGame(2, billsAtFins)

	// Week 2: Pats lose and Raiders win to even record
	patsAtBills := game.Game{
		Winner: team.BuffaloBills.Name,
		Loser:  team.NewEnglandPatriots.Name,
		Home:   team.BuffaloBills.Name,
		Away:   team.NewEnglandPatriots.Name,
	}
	finsAtRaiders := game.Game{
		Winner: team.LasVegasRaiders.Name,
		Loser:  team.MiamiDolphins.Name,
		Home:   team.LasVegasRaiders.Name,
		Away:   team.MiamiDolphins.Name,
	}
	sched.AddGame(2, patsAtBills)
	sched.AddGame(2, finsAtRaiders)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.LasVegasRaiders})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}

func Conference2ClubsConferenceRecord() Scenario {
	sched := schedule.NewSchedule()

	// Pats    1-1, 1-0 in conf
	// Raiders 1-1, 0-1 in conf

	// Week 1: Pats get their conference win, Raiders get their loss
	patsAtTexans := game.Game{
		Winner: team.NewEnglandPatriots.Name,
		Loser:  team.HoustonTexans.Name,
		Home:   team.HoustonTexans.Name,
		Away:   team.NewEnglandPatriots.Name,
	}
	raidersAtSteelers := game.Game{
		Winner: team.PittsburghSteelers.Name,
		Loser:  team.LasVegasRaiders.Name,
		Home:   team.PittsburghSteelers.Name,
		Away:   team.LasVegasRaiders.Name,
	}
	sched.AddGame(1, patsAtTexans)
	sched.AddGame(1, raidersAtSteelers)

	// Week 2: Pats lose and Raiders win to even win pct. Both games out of conference
	cardinalsAtPats := game.Game{
		Winner: team.ArizonaCardinals.Name,
		Loser:  team.NewEnglandPatriots.Name,
		Home:   team.NewEnglandPatriots.Name,
		Away:   team.ArizonaCardinals.Name,
	}
	raidersAtGiants := game.Game{
		Winner: team.LasVegasRaiders.Name,
		Loser:  team.NewYorkGiants.Name,
		Home:   team.LasVegasRaiders.Name,
		Away:   team.NewYorkGiants.Name,
	}
	sched.AddGame(2, cardinalsAtPats)
	sched.AddGame(2, raidersAtGiants)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.LasVegasRaiders})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}

func Conference2ClubsCommonGames() Scenario {
	sched := schedule.NewSchedule()

	// Need 4 common games. They will each play
	// Texans, Steelers, Cardinals, Rams

	// Pats    3-2, 3-1 v common, 0-1 v other
	// Raiders 3-2, 2-2 v common, 1-0 v other

	// Week 1: Pats lose to other, Raiders win against other
	patsAtBroncos := game.Game{
		Winner: team.DenverBroncos.Name,
		Loser:  team.NewEnglandPatriots.Name,
		Home:   team.DenverBroncos.Name,
		Away:   team.NewEnglandPatriots.Name,
	}
	raidersAtJets := game.Game{
		Winner: team.LasVegasRaiders.Name,
		Loser:  team.NewYorkJets.Name,
		Home:   team.NewYorkJets.Name,
		Away:   team.LasVegasRaiders.Name,
	}
	sched.AddGame(1, patsAtBroncos)
	sched.AddGame(1, raidersAtJets)

	// Week 2: Pats win v common, Raiders lose v common
	texansAtPats := game.Game{
		Winner: team.NewEnglandPatriots.Name,
		Loser:  team.HoustonTexans.Name,
		Home:   team.NewEnglandPatriots.Name,
		Away:   team.HoustonTexans.Name,
	}
	steelersAtRaiders := game.Game{
		Winner: team.PittsburghSteelers.Name,
		Loser:  team.LasVegasRaiders.Name,
		Home:   team.LasVegasRaiders.Name,
		Away:   team.PittsburghSteelers.Name,
	}
	sched.AddGame(2, texansAtPats)
	sched.AddGame(2, steelersAtRaiders)

	// Weeks 3-5: Pats and Raiders both win v remaining common
	steelersAtPats := game.Game{
		Winner: team.NewEnglandPatriots.Name,
		Loser:  team.PittsburghSteelers.Name,
		Home:   team.NewEnglandPatriots.Name,
		Away:   team.PittsburghSteelers.Name,
	}
	texansAtRaiders := game.Game{
		Winner: team.LasVegasRaiders.Name,
		Loser:  team.HoustonTexans.Name,
		Home:   team.LasVegasRaiders.Name,
		Away:   team.HoustonTexans.Name,
	}
	sched.AddGame(3, steelersAtPats)
	sched.AddGame(3, texansAtRaiders)

	patsAtCards := game.Game{
		Winner: team.NewEnglandPatriots.Name,
		Loser:  team.ArizonaCardinals.Name,
		Home:   team.ArizonaCardinals.Name,
		Away:   team.NewEnglandPatriots.Name,
	}
	raidersAtRams := game.Game{
		Winner: team.LasVegasRaiders.Name,
		Loser:  team.LosAngelesRams.Name,
		Home:   team.LosAngelesRams.Name,
		Away:   team.LasVegasRaiders.Name,
	}
	sched.AddGame(4, patsAtCards)
	sched.AddGame(4, raidersAtRams)

	ramsAtPats := game.Game{
		Winner: team.NewEnglandPatriots.Name,
		Loser:  team.LosAngelesRams.Name,
		Home:   team.NewEnglandPatriots.Name,
		Away:   team.LosAngelesRams.Name,
	}
	cardsAtRaiders := game.Game{
		Winner: team.LasVegasRaiders.Name,
		Loser:  team.ArizonaCardinals.Name,
		Home:   team.LasVegasRaiders.Name,
		Away:   team.ArizonaCardinals.Name,
	}
	sched.AddGame(5, ramsAtPats)
	sched.AddGame(5, cardsAtRaiders)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.LasVegasRaiders})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}

func Conference2ClubsStrengthOfVictory() Scenario {
	sched := schedule.NewSchedule()

	// Pats    2-0, Opponents beaten 2-2
	// Raiders 2-0, Opponents beaten 0-4

	// Week 1: Pats and Raiders both win.
	// Seahawks (Pats future opponent) win, Rams (Raiders future opponent) lose
	patsAtCards := game.Game{
		Winner: team.NewEnglandPatriots.Name,
		Loser:  team.ArizonaCardinals.Name,
		Home:   team.ArizonaCardinals.Name,
		Away:   team.NewEnglandPatriots.Name,
	}
	raidersAtGiants := game.Game{
		Winner: team.LasVegasRaiders.Name,
		Loser:  team.NewYorkGiants.Name,
		Home:   team.NewYorkGiants.Name,
		Away:   team.LasVegasRaiders.Name,
	}
	seahawksAtRams := game.Game{
		Winner: team.SeattleSeahawks.Name,
		Loser:  team.LosAngelesRams.Name,
		Home:   team.LosAngelesRams.Name,
		Away:   team.SeattleSeahawks.Name,
	}
	sched.AddGame(1, patsAtCards)
	sched.AddGame(1, raidersAtGiants)
	sched.AddGame(1, seahawksAtRams)

	// Week 2: Pats and Raiders both win
	// Cardinals (Pats previous opponent) win, Giants (Raiders previous opponent) lose
	seahawksAtPats := game.Game{
		Winner: team.NewEnglandPatriots.Name,
		Loser:  team.SeattleSeahawks.Name,
		Home:   team.NewEnglandPatriots.Name,
		Away:   team.SeattleSeahawks.Name,
	}
	ramsAtRaiders := game.Game{
		Winner: team.LasVegasRaiders.Name,
		Loser:  team.LosAngelesRams.Name,
		Home:   team.LasVegasRaiders.Name,
		Away:   team.LosAngelesRams.Name,
	}
	cardinalsAtGiants := game.Game{
		Winner: team.ArizonaCardinals.Name,
		Loser:  team.NewYorkGiants.Name,
		Home:   team.NewYorkGiants.Name,
		Away:   team.ArizonaCardinals.Name,
	}
	sched.AddGame(2, seahawksAtPats)
	sched.AddGame(2, ramsAtRaiders)
	sched.AddGame(2, cardinalsAtGiants)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.LasVegasRaiders})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}

func Conference2ClubsStrengthOfSchedule() Scenario {
	sched := schedule.NewSchedule()

	// Pats 1-1, Opponents 2-2
	// Raiders 1-1, Opponents 1-3

	// Week 1: Pats and Raiders both win.
	// Seahawks (Pats future opponent) win, Rams (Raiders future opponent) lose
	patsAtCards := game.Game{
		Winner: team.NewEnglandPatriots.Name,
		Loser:  team.ArizonaCardinals.Name,
		Home:   team.ArizonaCardinals.Name,
		Away:   team.NewEnglandPatriots.Name,
	}
	raidersAtGiants := game.Game{
		Winner: team.LasVegasRaiders.Name,
		Loser:  team.NewYorkGiants.Name,
		Home:   team.NewYorkGiants.Name,
		Away:   team.LasVegasRaiders.Name,
	}
	seahawksAtRams := game.Game{
		Winner: team.SeattleSeahawks.Name,
		Loser:  team.LosAngelesRams.Name,
		Home:   team.LosAngelesRams.Name,
		Away:   team.SeattleSeahawks.Name,
	}
	sched.AddGame(1, patsAtCards)
	sched.AddGame(1, raidersAtGiants)
	sched.AddGame(1, seahawksAtRams)

	// Week 2: Pats and Raiders both lose. Pats to a better team
	// Cardinals (Pats previous opponent) and Giants (Raiders previous opponent) both lose to keep SOV the same
	seahawksAtPats := game.Game{
		Winner: team.SeattleSeahawks.Name,
		Loser:  team.NewEnglandPatriots.Name,
		Home:   team.NewEnglandPatriots.Name,
		Away:   team.SeattleSeahawks.Name,
	}
	ramsAtRaiders := game.Game{
		Winner: team.LosAngelesRams.Name,
		Loser:  team.LasVegasRaiders.Name,
		Home:   team.LasVegasRaiders.Name,
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
	sched.AddGame(2, ramsAtRaiders)
	sched.AddGame(2, cardinalsAtPanthers)
	sched.AddGame(2, giantsAtBears)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.LasVegasRaiders})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}

func Conference2ClubsConferenceRank() Scenario {
	sched := schedule.NewSchedule()

	// Pats    1-0, 30 points for, 20 points against
	// Raiders 1-0, 29 points for, 21 points against

	// Week 1: Pats and Raiders both win. Pats score more and allow less
	patsAtCards := game.Game{
		Winner:  team.NewEnglandPatriots.Name,
		Loser:   team.ArizonaCardinals.Name,
		Home:    team.ArizonaCardinals.Name,
		Away:    team.NewEnglandPatriots.Name,
		PtsWin:  30,
		PtsLose: 20,
	}
	raidersAtGiants := game.Game{
		Winner:  team.LasVegasRaiders.Name,
		Loser:   team.NewYorkGiants.Name,
		Home:    team.NewYorkGiants.Name,
		Away:    team.LasVegasRaiders.Name,
		PtsWin:  29,
		PtsLose: 21,
	}
	sched.AddGame(1, patsAtCards)
	sched.AddGame(1, raidersAtGiants)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.LasVegasRaiders})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}

func Conference2ClubsLeagueRank() Scenario {
	sched := schedule.NewSchedule()

	// Points For:
	// 1. Patriots  (AFC) - 40
	// 2. Cardinals (NFC) - 30
	// 3. Giants    (NFC) - 25
	// 4. Raiders   (AFC) - 20
	// 5. Bears     (NFC) - 16

	// Points Against:
	// 1. Bears     (NFC) - 15
	// 1. Raiders   (AFC) - 17
	// 2. Giants    (NFC) - 20
	// 3. Patriots  (AFC) - 30
	// 4. Cardinals (NFC) - 55

	// Week 1:
	patsAtCards := game.Game{
		Winner:  team.NewEnglandPatriots.Name,
		Loser:   team.ArizonaCardinals.Name,
		Home:    team.ArizonaCardinals.Name,
		Away:    team.NewEnglandPatriots.Name,
		PtsWin:  40,
		PtsLose: 30,
	}
	raidersAtGiants := game.Game{
		Winner:  team.LasVegasRaiders.Name,
		Loser:   team.NewYorkGiants.Name,
		Home:    team.NewYorkGiants.Name,
		Away:    team.LasVegasRaiders.Name,
		PtsWin:  20,
		PtsLose: 17,
	}
	sched.AddGame(1, patsAtCards)
	sched.AddGame(1, raidersAtGiants)

	// Week 2: Giants overtake Raiders in PF, lose to keep SOV/SOS same
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
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.LasVegasRaiders})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}

func Conference2ClubsNetPointsConference() Scenario {
	sched := schedule.NewSchedule()

	// Pats    0-1, -5 net conference
	// Raiders 0-1, -10 net conference

	// Week 1: Pats and Raiders both lost, Pats better diff
	patsAtSteelers := game.Game{
		Winner:  team.PittsburghSteelers.Name,
		Loser:   team.NewEnglandPatriots.Name,
		Home:    team.PittsburghSteelers.Name,
		Away:    team.NewEnglandPatriots.Name,
		PtsWin:  30,
		PtsLose: 25,
	}
	raidersAtBengals := game.Game{
		Winner:  team.CincinnatiBengals.Name,
		Loser:   team.LasVegasRaiders.Name,
		Home:    team.CincinnatiBengals.Name,
		Away:    team.LasVegasRaiders.Name,
		PtsWin:  20,
		PtsLose: 10,
	}
	sched.AddGame(1, patsAtSteelers)
	sched.AddGame(1, raidersAtBengals)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.LasVegasRaiders})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}

func Conference2ClubsNetPoints() Scenario {
	sched := schedule.NewSchedule()

	// Pats   0-1, -5 net
	// Raiders 0-1, -10 net

	// Week 1: Pats and Raiders both lost (non-conference), Pats better diff
	patsAtCards := game.Game{
		Winner:  team.ArizonaCardinals.Name,
		Loser:   team.NewEnglandPatriots.Name,
		Home:    team.ArizonaCardinals.Name,
		Away:    team.NewEnglandPatriots.Name,
		PtsWin:  30,
		PtsLose: 25,
	}
	raidersAtRams := game.Game{
		Winner:  team.LosAngelesRams.Name,
		Loser:   team.LasVegasRaiders.Name,
		Home:    team.LosAngelesRams.Name,
		Away:    team.LasVegasRaiders.Name,
		PtsWin:  20,
		PtsLose: 10,
	}
	sched.AddGame(1, patsAtCards)
	sched.AddGame(1, raidersAtRams)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{team.NewEnglandPatriots, team.LasVegasRaiders})
	return Scenario{
		entries:   filtered,
		schedules: sched.SplitToTeams(),
	}
}
