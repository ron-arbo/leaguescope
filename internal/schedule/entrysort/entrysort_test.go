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

	for i, _ := range entries {
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
		// 	args:    Division2ClubsHeadToHead(),
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
		// 	args:    Division2ClubsConferenceRecord(),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.NewYorkJets},
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Division 2 Clubs Strength of Victory",
		// 	args:    Division2ClubsStrengthOfVictory(),
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
		// 	name:    "Best combined ranking among conference teams in points scored and points allowed in all games",
		// 	args:    ConferenceRank(),
		// 	want:    []team.Team{team.NewEnglandPatriots, team.NewYorkJets},
		// 	wantErr: false,
		// },
		{
			name:    "Best combined ranking among all teams in points scored and points allowed in all games",
			args:    LeagueRank(),
			want:    []team.Team{team.NewEnglandPatriots, team.NewYorkJets},
			wantErr: false,
		},
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
