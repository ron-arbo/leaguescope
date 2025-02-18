package entrysort

import (
	"nfl-app/internal/entry"
	"nfl-app/internal/team"
	"testing"
)

func orderMatched(entries []entry.Entry, teams []team.Team) bool {
	if len(entries) != len(teams) {
		return false
	}

	for i := range entries {
		if entries[i].Team.Name != teams[i].Name {
			return false
		}
	}

	return true
}

// TODO: If possible, should we test both ways in case a group of entries is already sorted by default?
// Or maybe just ensure it isn't when we give it as input?
func TestSortEntries(t *testing.T) {
	tests := []struct {
		name    string
		args    Scenario
		want    []team.Team // Just use the team instead of the entire entry
		wantErr bool
	}{
		// {
		// 	name:    "Division 2 Clubs Head to Head",
		// 	args:    HeadToHead(team.NewEnglandPatriots, team.NewYorkJets),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.NewYorkJets},
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Division 2 Clubs Division Record",
		// 	args:    Division2ClubsDivisionRecord(),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.NewYorkJets},
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Division 2 Clubs Common Games",
		// 	args:    Division2ClubsCommonGames(),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.NewYorkJets},
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Division 2 Clubs Conference Record",
		// 	args:    ConferenceRecord(team.NewEnglandPatriots, team.NewYorkJets),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.NewYorkJets},
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Division 2 Clubs Strength of Victory",
		// 	args:    StrengthOfVictory(team.NewEnglandPatriots, team.NewYorkJets),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.NewYorkJets},
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Division 2 Clubs Strength of Schedule",
		// 	args:    Division2ClubsStrengthOfSchedule(),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.NewYorkJets},
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Division 2 Clubs best combined ranking among conference teams in points scored and points allowed in all games",
		// 	args:    Division2ClubsConferenceRank(),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.NewYorkJets},
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Division 2 Clubs best combined ranking among all teams in points scored and points allowed in all games",
		// 	args:    Division2ClubsLeagueRank(),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.NewYorkJets},
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Division 2 Clubs Net points common games",
		// 	args:    Division2ClubsNetPointsCommonGames(),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.NewYorkJets},
		// 	wantErr: false,
		// },
		// Cannot test coin toss as it is random
		// {
		// 	name:    "Division 2 Clubs Net points",
		// 	args:    Division2ClubsNetPoints(),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.NewYorkJets},
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Conference 2 Clubs Head to Head",
		// 	args:    HeadToHead(team.NewEnglandPatriots, team.LasVegasRaiders),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.LasVegasRaiders},
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Conference 2 Clubs Conference Record",
		// 	args:    ConferenceRecord(team.NewEnglandPatriots, team.LasVegasRaiders),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.LasVegasRaiders},
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Conference 2 Clubs Common Games",
		// 	args:    Conference2ClubsCommonGames(),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.LasVegasRaiders},
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Conference 2 Clubs Strength of Victory",
		// 	args:    StrengthOfVictory(team.NewEnglandPatriots, team.LasVegasRaiders),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.LasVegasRaiders},
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Conference 2 Clubs Strength of Schedule",
		// 	args:    Conference2ClubsStrengthOfSchedule(),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.LasVegasRaiders},
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Conference 2 Clubs best combined ranking among conference teams in points scored and points allowed in all games",
		// 	args:    Conference2ClubsConferenceRank(),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.LasVegasRaiders},
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Conference 2 Clubs best combined ranking among all teams in points scored and points allowed in all games",
		// 	args:    Conference2ClubsLeagueRank(),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.LasVegasRaiders},
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Conference 2 Clubs net points conference",
		// 	args:    Conference2ClubsNetPointsConference(),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.LasVegasRaiders},
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Conference 2 Clubs net points",
		// 	args:    Conference2ClubsNetPoints(),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.LasVegasRaiders},
		// 	wantErr: false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SortEntries(tt.args.entries, tt.args.schedules)
			if (err != nil) != tt.wantErr {
				t.Errorf("SortEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("Length mismatch")
			}
			if !orderMatched(got, tt.want) {
				t.Errorf("Order mismatch\nGot: %s\nWant: %s\n", entry.Names(got), team.Names(tt.want))
			}
		})
	}
}
