package stats

// ESPN represents a statistic for a team in regard to the standings, as returned by the ESPN API
type ESPNStat struct {
	Name             string  `json:"name"`
	DisplayName      string  `json:"displayName,omitempty"`
	ShortDisplayName string  `json:"shortDisplayName,omitempty"`
	Description      string  `json:"description,omitempty"`
	Abbreviation     string  `json:"abbreviation,omitempty"`
	Type             string  `json:"type"`
	Value            float64 `json:"value,omitempty"`
	DisplayValue     string  `json:"displayValue"`
	ID               string  `json:"id,omitempty"`
}
