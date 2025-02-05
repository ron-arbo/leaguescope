package stats

import (
	"fmt"
)

type Stat struct {
	Name string
	// Like ESPN API, we have a field for both number and string stats.
	// It's easiest to store both, even if one is unused in some cases
	value        float64
	displayValue string
}

// Make this a pointer so we can update the value
type Stats map[string]*Stat

func NewStats() Stats {
	stats := make(Stats)

	// Initialize all stats to 0 value
	for _, name := range statNames {
		stat := Stat{
			Name:         name,
			value:        0,
			displayValue: "",
		}
		stats[name] = &stat
	}

	return stats
}

// Use private fields to ensure that the value is only updated through the SetValue method,
// and thus the displayValue is always in sync with the value
func (s *Stat) SetValue(value float64) {
	s.value = value

	// Display win percent to 3 decimal places, otherwise display an integer
	if s.Name == StatWinPercent {
		s.displayValue = fmt.Sprintf("%.3f", value)
	} else {
		s.displayValue = fmt.Sprintf("%d", int(value))
	}
}

func (s *Stat) Increment() {
	s.SetValue(s.value + 1)
}

func (s *Stat) Value() float64 {
	return s.value
}

func (s *Stat) DisplayValue() string {
	return s.displayValue
}
