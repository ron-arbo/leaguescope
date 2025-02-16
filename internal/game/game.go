package game

import (
	"time"
)

type Game struct {
	Time time.Time

	Winner string
	Loser  string

	Home string
	Away string

	PtsWin  int
	PtsLose int

	YardsWin  int
	YardsLose int

	ToWin  int
	ToLose int
}

// TODO:
// game.Validate() make sure no negative values, two unique teams, etc
