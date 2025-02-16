package stats

type Points struct {
	For     int
	Against int
}

func (p *Points) Differential() int {
	return p.For - p.Against
}

func (p *Points) AddFor(other int) {
	p.For += other
}

func (p *Points) AddAgainst(other int) {
	p.Against += other
}
