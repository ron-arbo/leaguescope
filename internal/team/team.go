package team

import (
	"slices"
	"time"

	"golang.org/x/exp/rand"
)

type Team struct {
	Name       string `json:"name"` // Could split this into struct containing city + team, but let's keep it simple for now
	Conference string `json:"conference"`
	Division   string `json:"division"`
}

const (
	// Conferences
	AFC = "American Football Conference"
	NFC = "National Football Conference"

	// Divisions
	AFCEast  = "AFC East"
	AFCNorth = "AFC North"
	AFCSouth = "AFC South"
	AFCWest  = "AFC West"
	NFCEast  = "NFC East"
	NFCNorth = "NFC North"
	NFCSouth = "NFC South"
	NFCWest  = "NFC West"
)

// Definitions for all teams, divisons, and conferences
var (
	NewEnglandPatriots = Team{
		Name:       "New England Patriots",
		Conference: AFC,
		Division:   AFCEast,
	}

	NewYorkJets = Team{
		Name:       "New York Jets",
		Conference: AFC,
		Division:   AFCEast,
	}

	BuffaloBills = Team{
		Name:       "Buffalo Bills",
		Conference: AFC,
		Division:   AFCEast,
	}

	MiamiDolphins = Team{
		Name:       "Miami Dolphins",
		Conference: AFC,
		Division:   AFCEast,
	}

	PittsburghSteelers = Team{
		Name:       "Pittsburgh Steelers",
		Conference: AFC,
		Division:   AFCNorth,
	}

	BaltimoreRavens = Team{
		Name:       "Baltimore Ravens",
		Conference: AFC,
		Division:   AFCNorth,
	}

	ClevelandBrowns = Team{
		Name:       "Cleveland Browns",
		Conference: AFC,
		Division:   AFCNorth,
	}

	CincinnatiBengals = Team{
		Name:       "Cincinnati Bengals",
		Conference: AFC,
		Division:   AFCNorth,
	}

	TennesseeTitans = Team{
		Name:       "Tennessee Titans",
		Conference: AFC,
		Division:   AFCSouth,
	}

	IndianapolisColts = Team{
		Name:       "Indianapolis Colts",
		Conference: AFC,
		Division:   AFCSouth,
	}

	JacksonvilleJaguars = Team{
		Name:       "Jacksonville Jaguars",
		Conference: AFC,
		Division:   AFCSouth,
	}

	HoustonTexans = Team{
		Name:       "Houston Texans",
		Conference: AFC,
		Division:   AFCSouth,
	}

	KansasCityChiefs = Team{
		Name:       "Kansas City Chiefs",
		Conference: AFC,
		Division:   AFCWest,
	}

	LasVegasRaiders = Team{
		Name:       "Las Vegas Raiders",
		Conference: AFC,
		Division:   AFCWest,
	}

	LosAngelesChargers = Team{
		Name:       "Los Angeles Chargers",
		Conference: AFC,
		Division:   AFCWest,
	}

	DenverBroncos = Team{
		Name:       "Denver Broncos",
		Conference: AFC,
		Division:   AFCWest,
	}

	DallasCowboys = Team{
		Name:       "Dallas Cowboys",
		Conference: NFC,
		Division:   NFCEast,
	}

	NewYorkGiants = Team{
		Name:       "New York Giants",
		Conference: NFC,
		Division:   NFCEast,
	}

	PhiladelphiaEagles = Team{
		Name:       "Philadelphia Eagles",
		Conference: NFC,
		Division:   NFCEast,
	}

	WashingtonCommanders = Team{
		Name:       "Washington Commanders",
		Conference: NFC,
		Division:   NFCEast,
	}

	GreenBayPackers = Team{
		Name:       "Green Bay Packers",
		Conference: NFC,
		Division:   NFCNorth,
	}

	MinnesotaVikings = Team{
		Name:       "Minnesota Vikings",
		Conference: NFC,
		Division:   NFCNorth,
	}

	ChicagoBears = Team{
		Name:       "Chicago Bears",
		Conference: NFC,
		Division:   NFCNorth,
	}

	DetroitLions = Team{
		Name:       "Detroit Lions",
		Conference: NFC,
		Division:   NFCNorth,
	}

	TampaBayBuccaneers = Team{
		Name:       "Tampa Bay Buccaneers",
		Conference: NFC,
		Division:   NFCSouth,
	}

	NewOrleansSaints = Team{
		Name:       "New Orleans Saints",
		Conference: NFC,
		Division:   NFCSouth,
	}

	CarolinaPanthers = Team{
		Name:       "Carolina Panthers",
		Conference: NFC,
		Division:   NFCSouth,
	}

	AtlantaFalcons = Team{
		Name:       "Atlanta Falcons",
		Conference: NFC,
		Division:   NFCSouth,
	}

	LosAngelesRams = Team{
		Name:       "Los Angeles Rams",
		Conference: NFC,
		Division:   NFCWest,
	}

	SanFrancisco49ers = Team{
		Name:       "San Francisco 49ers",
		Conference: NFC,
		Division:   NFCWest,
	}

	SeattleSeahawks = Team{
		Name:       "Seattle Seahawks",
		Conference: NFC,
		Division:   NFCWest,
	}

	ArizonaCardinals = Team{
		Name:       "Arizona Cardinals",
		Conference: NFC,
		Division:   NFCWest,
	}

	AFCEastTeams  = []Team{NewEnglandPatriots, NewYorkJets, BuffaloBills, MiamiDolphins}
	AFCNorthTeams = []Team{PittsburghSteelers, BaltimoreRavens, ClevelandBrowns, CincinnatiBengals}
	AFCSouthTeams = []Team{TennesseeTitans, IndianapolisColts, JacksonvilleJaguars, HoustonTexans}
	AFCWestTeams  = []Team{KansasCityChiefs, LasVegasRaiders, LosAngelesChargers, DenverBroncos}

	NFCEastTeams  = []Team{DallasCowboys, NewYorkGiants, PhiladelphiaEagles, WashingtonCommanders}
	NFCNorthTeams = []Team{GreenBayPackers, MinnesotaVikings, ChicagoBears, DetroitLions}
	NFCSouthTeams = []Team{TampaBayBuccaneers, NewOrleansSaints, CarolinaPanthers, AtlantaFalcons}
	NFCWestTeams  = []Team{LosAngelesRams, SanFrancisco49ers, SeattleSeahawks, ArizonaCardinals}

	AFCTeams = slices.Concat(AFCEastTeams, AFCNorthTeams, AFCSouthTeams, AFCWestTeams)
	NFCTeams = slices.Concat(NFCEastTeams, NFCNorthTeams, NFCSouthTeams, NFCWestTeams)

	NFLTeams = slices.Concat(AFCTeams, NFCTeams)

	Conferences = []string{AFC, NFC}
	Divisions   = []string{AFCEast, AFCNorth, AFCSouth, AFCWest, NFCEast, NFCNorth, NFCSouth, NFCWest}
)

// TODO: Could create map for this instead
func DisplayNameToTeam(displayName string) Team {
	switch displayName {
	case "New England Patriots":
		return NewEnglandPatriots
	case "New York Jets":
		return NewYorkJets
	case "Buffalo Bills":
		return BuffaloBills
	case "Miami Dolphins":
		return MiamiDolphins
	case "Pittsburgh Steelers":
		return PittsburghSteelers
	case "Baltimore Ravens":
		return BaltimoreRavens
	case "Cleveland Browns":
		return ClevelandBrowns
	case "Cincinnati Bengals":
		return CincinnatiBengals
	case "Tennessee Titans":
		return TennesseeTitans
	case "Indianapolis Colts":
		return IndianapolisColts
	case "Jacksonville Jaguars":
		return JacksonvilleJaguars
	case "Houston Texans":
		return HoustonTexans
	case "Kansas City Chiefs":
		return KansasCityChiefs
	case "Las Vegas Raiders":
		return LasVegasRaiders
	case "Los Angeles Chargers":
		return LosAngelesChargers
	case "Denver Broncos":
		return DenverBroncos
	case "Dallas Cowboys":
		return DallasCowboys
	case "New York Giants":
		return NewYorkGiants
	case "Philadelphia Eagles":
		return PhiladelphiaEagles
	case "Washington Commanders":
		return WashingtonCommanders
	case "Green Bay Packers":
		return GreenBayPackers
	case "Minnesota Vikings":
		return MinnesotaVikings
	case "Chicago Bears":
		return ChicagoBears
	case "Detroit Lions":
		return DetroitLions
	case "Tampa Bay Buccaneers":
		return TampaBayBuccaneers
	case "New Orleans Saints":
		return NewOrleansSaints
	case "Carolina Panthers":
		return CarolinaPanthers
	case "Atlanta Falcons":
		return AtlantaFalcons
	case "Los Angeles Rams":
		return LosAngelesRams
	case "San Francisco 49ers":
		return SanFrancisco49ers
	case "Seattle Seahawks":
		return SeattleSeahawks
	case "Arizona Cardinals":
		return ArizonaCardinals
	default:
		return Team{}
	}
}

func SameDivision(t1, t2 string) bool {
	return DisplayNameToTeam(t1).Division == DisplayNameToTeam(t2).Division
}

func SameConference(t1, t2 string) bool {
	return DisplayNameToTeam(t1).Conference == DisplayNameToTeam(t2).Conference
}

func Names(teams []Team) string {
	var out string
	for _, team := range teams {
		out = out + team.Name + " "
	}

	return out
}

func GetRandomTeams(count int, noTeams []Team, noDivs []string, noConf string) []Team {
	rand.Seed(uint64(time.Now().UnixNano()))
	out := make([]Team, 0, count)

	// TODO: Better to make these map[string]interface{}?
	seenTeams := make(map[string]bool)
	seenDivs := make(map[string]bool)

	for _, team := range noTeams {
		seenTeams[team.Name] = true
	}
	for _, div := range noDivs {
		seenDivs[div] = true
	}

	for len(out) < count {
		randomTeam := NFLTeams[rand.Intn(len(NFLTeams))]
		if !seenTeams[randomTeam.Name] && !seenDivs[randomTeam.Division] && randomTeam.Conference != noConf {
			out = append(out, randomTeam)
			// Update seenTeams so we don't have duplicates
			seenTeams[randomTeam.Name] = true
		}
	}

	return out
}
