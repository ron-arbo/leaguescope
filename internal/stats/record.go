package stats

type Record struct {
	wins   int
	losses int
	ties   int
}

func NewRecord(wins, losses, ties int) *Record {
	return &Record{
		wins:   wins,
		losses: losses,
		ties:   ties,
	}
}

func (r *Record) Wins() int {
	return r.wins
}

func (r *Record) Losses() int {
	return r.losses
}

func (r *Record) Ties() int {
	return r.ties
}

func (r *Record) AddWin() {
	r.wins++
}

func (r *Record) AddLoss() {
	r.losses++
}

func (r *Record) AddTie() {
	r.ties++
}

func (r *Record) Add(other *Record) {
	r.wins += other.wins
	r.losses += other.losses
	r.ties += other.ties
}

func (r *Record) WinPercentage() float64 {
	wins := float64(r.wins) + float64(r.ties)/2
	total := float64(r.wins + r.losses + r.ties)

	if total == 0 {
		return 0
	}

	return wins / total
}

func (r *Record) GamesPlayed() int {
	return r.wins + r.losses + r.ties
}
