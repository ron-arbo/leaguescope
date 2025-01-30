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

	if g.HomeTeam.Name.String() == team {
		return g.HomeScore > g.AwayScore
	}

	return g.AwayScore > g.HomeScore
}
