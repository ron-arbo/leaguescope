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

func TestSortEntries(t *testing.T) {
	tests := []struct {
		name    string
		args    Scenario
		want    []team.Team // Just use the team instead of the entire entry
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Division 2 Clubs Head to Head",
			args:    Division2ClubsHeadToHead(),
			want:    []team.Team{team.BuffaloBills, team.NewEnglandPatriots, team.NewYorkJets, team.MiamiDolphins},
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
