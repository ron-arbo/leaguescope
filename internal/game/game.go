package game

import (
	"nfl-app/internal/team"
)

// Common struct for both TeamSchedule and LeagueSchedule
type Game struct {
	HomeTeam  team.Team
	AwayTeam  team.Team
	HomeScore int
	AwayScore int
	Completed bool
}

func (g *Game) Tied() bool {
	return g.HomeScore == g.AwayScore && g.Completed
}

func (g *Game) WonBy(team string) bool {
	if !g.Completed {
		return false
	}

	if g.HomeTeamName() == team {
		return g.HomeScore > g.AwayScore
	}

	return g.AwayScore > g.HomeScore
}

func (g* Game) HomeTeamName() string {
	return g.HomeTeam.Name.String()
}

func (g* Game) AwayTeamName() string {
	return g.AwayTeam.Name.String()
}
