package entrysort

import (
	"fmt"
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

// HeadToHead returns a general scenario in which t1 breaks
// the tie with t2 by head to head record
func HeadToHead(t1, t2 team.Team) Scenario {
	// Get two other teams for them to play
	others := team.GetRandomTeams(2, []team.Team{t1, t2}, nil, "")

	sched := schedule.NewSchedule()

	// Week 1: t1 beats t2
	t1Att2 := game.Game{
		Winner: t1.Name,
		Loser:  t2.Name,
		Home:   t2.Name,
		Away:   t1.Name,
	}
	sched.AddGame(1, t1Att2)

	// Week 2: t1 loses to an other, t2 beats an other. Records evened
	t1AtOther1 := game.Game{
		Winner: others[0].Name,
		Loser:  t1.Name,
		Home:   others[0].Name,
		Away:   t1.Name,
	}
	t2AtOther2 := game.Game{
		Winner: t2.Name,
		Loser:  others[1].Name,
		Home:   others[1].Name,
		Away:   t2.Name,
	}
	sched.AddGame(2, t1AtOther1)
	sched.AddGame(2, t2AtOther2)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{t1, t2})
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

func StrengthOfVictory(t1, t2 team.Team) Scenario {
	sched := schedule.NewSchedule()
	// Get two other teams for them to play
	//
	// If the two teams given are from the same conference, select two teams from the
	// other conference as opponents so we don't interferece with previous tiebreakers
	//
	// If the two teams given are from different conferences, select any other two teams as the
	// only league tiebreakers before strength of victory are head to head and common games
	var noConf string
	if t1.Conference == t2.Conference {
		noConf = t1.Conference
	}
	others := team.GetRandomTeams(4, []team.Team{t1, t2}, nil, noConf)

	t1Opps := []string{others[0].Name, others[1].Name} // Win all games not against t1 --> good SOV
	t2Opps := []string{others[2].Name, others[3].Name} // Lose all games (including t2) --> bad SOV

	// t1 | 2-0 | Opponents 2-2
	// t2 | 2-0 | Opponents 0-4

	// Week 1: t1 and t2 both win
	// t1's future opponent wins, t2's future opponent loses
	t1AtOpp := game.Game{
		Winner: t1.Name,
		Loser:  t1Opps[0],
		Home:   t1Opps[0],
		Away:   t1.Name,
	}
	t2AtOpp := game.Game{
		Winner: t2.Name,
		Loser:  t2Opps[0],
		Home:   t2Opps[0],
		Away:   t2.Name,
	}
	t1OppAtt2Opp := game.Game{
		Winner: t1Opps[1],
		Loser:  t2Opps[1],
		Home:   t2Opps[1],
		Away:   t1Opps[1],
	}
	sched.AddGame(1, t1AtOpp)
	sched.AddGame(1, t2AtOpp)
	sched.AddGame(1, t1OppAtt2Opp)

	// Week 1: t1 and t2 both win
	// t1's previous opponent wins, t2's previous opponent loses
	oppAtt1 := game.Game{
		Winner: t1.Name,
		Loser:  t1Opps[1],
		Home:   t1.Name,
		Away:   t1Opps[1],
	}
	oppAtt2 := game.Game{
		Winner: t2.Name,
		Loser:  t2Opps[1],
		Home:   t2.Name,
		Away:   t2Opps[1],
	}
	t2OppAtt1Opp := game.Game{
		Winner: t1Opps[0],
		Loser:  t2Opps[0],
		Home:   t1Opps[0],
		Away:   t2Opps[0],
	}
	sched.AddGame(1, oppAtt1)
	sched.AddGame(1, oppAtt2)
	sched.AddGame(1, t2OppAtt1Opp)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{t1, t2})
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

func ConferenceRecord(t1, t2 team.Team) Scenario {
	sched := schedule.NewSchedule()

	// Conference record is only used as a tiebreaker when the two teams are from the same conference
	if t1.Conference != t2.Conference {
		// TODO: Don't panic
		panic(fmt.Sprintf("Teams not from the same conference: %s, %s", t1.Name, t2.Name))
	}

	// If the teams are from the same division, make sure to pick conference teams NOT from
	// that division, otherwise the division record tiebreaker might kick in before this one
	//
	// Not an issue if the teams are not from the same division, as the wild card tiebreaker
	// does not use division record
	var noDiv []string
	if t1.Division == t2.Division {
		noDiv = append(noDiv, t1.Division)
	}

	// Find 4 other teams, 2 from the same conference and 2 from the other
	var otherConf string
	if t1.Conference == team.AFC {
		otherConf = team.NFC
	} else {
		otherConf = team.AFC
	}

	confOpps := team.GetRandomTeams(2, nil, noDiv, otherConf)
	nonConfOpps := team.GetRandomTeams(2, nil, nil, t1.Conference)

	// t1 | 1-1 | 1-0 v conf
	// t2 | 1-1 | 0-1 v conf

	// Week 1: Both teams win, and t1's game is in-conference
	t1AtConf := game.Game{
		Winner: t1.Name,
		Loser:  confOpps[0].Name,
		Home:   confOpps[0].Name,
		Away:   t1.Name,
	}
	t2AtNonConf := game.Game{
		Winner: t2.Name,
		Loser:  nonConfOpps[0].Name,
		Home:   nonConfOpps[0].Name,
		Away:   t2.Name,
	}
	sched.AddGame(1, t1AtConf)
	sched.AddGame(1, t2AtNonConf)

	// Week 2: Both teams lose, and t2's game is in-conference
	t1AtNonConf := game.Game{
		Winner: nonConfOpps[1].Name,
		Loser:  t1.Name,
		Home:   nonConfOpps[1].Name,
		Away:   t1.Name,
	}
	t2AtConf := game.Game{
		Winner: confOpps[1].Name,
		Loser:  t2.Name,
		Home:   confOpps[1].Name,
		Away:   t2.Name,
	}
	sched.AddGame(2, t1AtNonConf)
	sched.AddGame(2, t2AtConf)

	entries := schedule.CreateEntries(sched)
	// Filter so only the tied teams we care about are sorted
	filtered := entry.FilterEntries(entries, []team.Team{t1, t2})
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
