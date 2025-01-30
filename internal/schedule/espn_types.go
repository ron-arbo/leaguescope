package schedule

import "time"

type ESPNSchedule struct {
	Timestamp       time.Time           `json:"timestamp"`
	Status          string              `json:"status"`
	Season          ESPNSeason          `json:"season"`
	Team            ESPNTeam            `json:"team"`
	Events          []ESPNEvent         `json:"events"`
	RequestedSeason ESPNRequestedSeason `json:"requestedSeason"`
	ByeWeek         int                 `json:"byeWeek"`
}

type ESPNSeason struct {
	Year        int    `json:"year"`
	Type        int    `json:"type"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Half        int    `json:"half"`
}

type ESPNTeam struct {
	ID              string         `json:"id"`
	Abbreviation    string         `json:"abbreviation"`
	Location        string         `json:"location"`
	Name            string         `json:"name"`
	DisplayName     string         `json:"displayName"`
	Clubhouse       string         `json:"clubhouse"`
	Color           string         `json:"color"`
	Logo            string         `json:"logo"`
	RecordSummary   string         `json:"recordSummary"`
	SeasonSummary   string         `json:"seasonSummary"`
	StandingSummary string         `json:"standingSummary"`
	Groups          ESPNTeamGroups `json:"groups"`
}

type ESPNTeamGroups struct {
	ID           string          `json:"id"`
	Parent       ESPNGroupParent `json:"parent"`
	IsConference bool            `json:"isConference"`
}

type ESPNGroupParent struct {
	ID string `json:"id"`
}

type ESPNEvent struct {
	ID           string            `json:"id"`
	Date         string            `json:"date"`
	Name         string            `json:"name"`
	ShortName    string            `json:"shortName"`
	Season       ESPNEventSeason   `json:"season"`
	SeasonType   ESPNSeasonType    `json:"seasonType"`
	Week         ESPNWeek          `json:"week"`
	TimeValid    bool              `json:"timeValid"`
	Competitions []ESPNCompetition `json:"competitions"`
	Links        []ESPNLink        `json:"links"`
}

type ESPNEventSeason struct {
	Year        int    `json:"year"`
	DisplayName string `json:"displayName"`
}

type ESPNSeasonType struct {
	ID           string `json:"id"`
	Type         int    `json:"type"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}

type ESPNWeek struct {
	Number int    `json:"number"`
	Text   string `json:"text"`
}

type ESPNCompetition struct {
	ID                string                `json:"id"`
	Date              string                `json:"date"`
	Attendance        int                   `json:"attendance"`
	Type              ESPNCompetitionType   `json:"type"`
	TimeValid         bool                  `json:"timeValid"`
	NeutralSite       bool                  `json:"neutralSite"`
	BoxscoreAvailable bool                  `json:"boxscoreAvailable"`
	TicketsAvailable  bool                  `json:"ticketsAvailable"`
	Venue             ESPNVenue             `json:"venue"`
	Competitors       []ESPNCompetitor      `json:"competitors"`
	Notes             []any                 `json:"notes"`
	Broadcasts        []ESPNBroadcast       `json:"broadcasts"`
	Status            ESPNCompetitionStatus `json:"status"`
}

type ESPNCompetitionType struct {
	ID           string `json:"id"`
	Text         string `json:"text"`
	Abbreviation string `json:"abbreviation"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
}

type ESPNVenue struct {
	FullName string           `json:"fullName"`
	Address  ESPNVenueAddress `json:"address"`
}

type ESPNVenueAddress struct {
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipCode"`
}

type ESPNCompetitor struct {
	ID       string             `json:"id"`
	Type     string             `json:"type"`
	Order    int                `json:"order"`
	HomeAway string             `json:"homeAway"`
	Winner   bool               `json:"winner"`
	Team     ESPNCompetitorTeam `json:"team"`
	Score    ESPNScore          `json:"score"`
	Record   []ESPNRecord       `json:"record"`
	Leaders  []ESPNLeader       `json:"leaders,omitempty"`
}

type ESPNCompetitorTeam struct {
	ID               string     `json:"id"`
	Location         string     `json:"location"`
	Nickname         string     `json:"nickname"`
	Abbreviation     string     `json:"abbreviation"`
	DisplayName      string     `json:"displayName"`
	ShortDisplayName string     `json:"shortDisplayName"`
	Logos            []ESPNLogo `json:"logos"`
	Links            []ESPNLink `json:"links"`
}

type ESPNLogo struct {
	Href        string   `json:"href"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Alt         string   `json:"alt"`
	Rel         []string `json:"rel"`
	LastUpdated string   `json:"lastUpdated"`
}

type ESPNLink struct {
	Rel  []string `json:"rel"`
	Href string   `json:"href"`
	Text string   `json:"text"`
}

type ESPNScore struct {
	Value        float64 `json:"value"`
	DisplayValue string  `json:"displayValue"`
}

type ESPNRecord struct {
	ID               string `json:"id"`
	Abbreviation     string `json:"abbreviation,omitempty"`
	DisplayName      string `json:"displayName"`
	ShortDisplayName string `json:"shortDisplayName"`
	Description      string `json:"description"`
	Type             string `json:"type"`
	DisplayValue     string `json:"displayValue"`
}

type ESPNLeader struct {
	Name         string             `json:"name"`
	DisplayName  string             `json:"displayName"`
	Abbreviation string             `json:"abbreviation"`
	Leaders      []ESPNLeaderDetail `json:"leaders"`
}

type ESPNLeaderDetail struct {
	DisplayValue string      `json:"displayValue"`
	Value        float64     `json:"value"`
	Athlete      ESPNAthlete `json:"athlete"`
}

type ESPNAthlete struct {
	ID          string     `json:"id"`
	LastName    string     `json:"lastName"`
	DisplayName string     `json:"displayName"`
	ShortName   string     `json:"shortName"`
	Links       []ESPNLink `json:"links"`
}

type ESPNBroadcast struct {
	Type   ESPNBroadcastType `json:"type"`
	Market ESPNMarket        `json:"market"`
	Media  ESPNMedia         `json:"media"`
	Lang   string            `json:"lang"`
	Region string            `json:"region"`
}

type ESPNBroadcastType struct {
	ID        string `json:"id"`
	ShortName string `json:"shortName"`
}

type ESPNMarket struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type ESPNMedia struct {
	ShortName string `json:"shortName"`
}

type ESPNCompetitionStatus struct {
	Clock        float64        `json:"clock"`
	DisplayClock string         `json:"displayClock"`
	Period       int            `json:"period"`
	Type         ESPNStatusType `json:"type"`
	IsTBDFlex    bool           `json:"isTBDFlex"`
}

type ESPNStatusType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	State       string `json:"state"`
	Completed   bool   `json:"completed"`
	Description string `json:"description"`
	Detail      string `json:"detail"`
	ShortDetail string `json:"shortDetail"`
}

type ESPNRequestedSeason struct {
	Year        int    `json:"year"`
	Type        int    `json:"type"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}
