package team

import (
	"fmt"
	"slices"
)

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
		ID:         idNewEnglandPatriots,
		Name:       TeamName{"New England", "Patriots"},
		Conference: AFC,
		Division:   AFCEast,
	}

	NewYorkJets = Team{
		ID:         idNewYorkJets,
		Name:       TeamName{"New York", "Jets"},
		Conference: AFC,
		Division:   AFCEast,
	}

	BuffaloBills = Team{
		ID:         idBuffaloBills,
		Name:       TeamName{"Buffalo", "Bills"},
		Conference: AFC,
		Division:   AFCEast,
	}

	MiamiDolphins = Team{
		ID:         idMiamiDolphins,
		Name:       TeamName{"Miami", "Dolphins"},
		Conference: AFC,
		Division:   AFCEast,
	}

	PittsburghSteelers = Team{
		ID:         idPittsburghSteelers,
		Name:       TeamName{"Pittsburgh", "Steelers"},
		Conference: AFC,
		Division:   AFCNorth,
	}

	BaltimoreRavens = Team{
		ID:         idBaltimoreRavens,
		Name:       TeamName{"Baltimore", "Ravens"},
		Conference: AFC,
		Division:   AFCNorth,
	}

	ClevelandBrowns = Team{
		ID:         idClevelandBrowns,
		Name:       TeamName{"Cleveland", "Browns"},
		Conference: AFC,
		Division:   AFCNorth,
	}

	CincinnatiBengals = Team{
		ID:         idCincinnatiBengals,
		Name:       TeamName{"Cincinnati", "Bengals"},
		Conference: AFC,
		Division:   AFCNorth,
	}

	TennesseeTitans = Team{
		ID:         idTennesseeTitans,
		Name:       TeamName{"Tennessee", "Titans"},
		Conference: AFC,
		Division:   AFCSouth,
	}

	IndianapolisColts = Team{
		ID:         idIndianapolisColts,
		Name:       TeamName{"Indianapolis", "Colts"},
		Conference: AFC,
		Division:   AFCSouth,
	}

	JacksonvilleJaguars = Team{
		ID:         idJacksonvilleJaguars,
		Name:       TeamName{"Jacksonville", "Jaguars"},
		Conference: AFC,
		Division:   AFCSouth,
	}

	HoustonTexans = Team{
		ID:         idHoustonTexans,
		Name:       TeamName{"Houston", "Texans"},
		Conference: AFC,
		Division:   AFCSouth,
	}

	KansasCityChiefs = Team{
		ID:         idKansasCityChiefs,
		Name:       TeamName{"Kansas City", "Chiefs"},
		Conference: AFC,
		Division:   AFCWest,
	}

	LasVegasRaiders = Team{
		ID:         idLasVegasRaiders,
		Name:       TeamName{"Las Vegas", "Raiders"},
		Conference: AFC,
		Division:   AFCWest,
	}

	LosAngelesChargers = Team{
		ID:         idLosAngelesChargers,
		Name:       TeamName{"Los Angeles", "Chargers"},
		Conference: AFC,
		Division:   AFCWest,
	}

	DenverBroncos = Team{
		ID:         idDenverBroncos,
		Name:       TeamName{"Denver", "Broncos"},
		Conference: AFC,
		Division:   AFCWest,
	}

	DallasCowboys = Team{
		ID:         idDallasCowboys,
		Name:       TeamName{"Dallas", "Cowboys"},
		Conference: NFC,
		Division:   NFCEast,
	}

	NewYorkGiants = Team{
		ID:         idNewYorkGiants,
		Name:       TeamName{"New York", "Giants"},
		Conference: NFC,
		Division:   NFCEast,
	}

	PhiladelphiaEagles = Team{
		ID:         idPhiladelphiaEagles,
		Name:       TeamName{"Philadelphia", "Eagles"},
		Conference: NFC,
		Division:   NFCEast,
	}

	WashingtonCommanders = Team{
		ID:         idWashingtonCommanders,
		Name:       TeamName{"Washington", "Commanders"},
		Conference: NFC,
		Division:   NFCEast,
	}

	GreenBayPackers = Team{
		ID:         idGreenBayPackers,
		Name:       TeamName{"Green Bay", "Packers"},
		Conference: NFC,
		Division:   NFCNorth,
	}

	MinnesotaVikings = Team{
		ID:         idMinnesotaVikings,
		Name:       TeamName{"Minnesota", "Vikings"},
		Conference: NFC,
		Division:   NFCNorth,
	}

	ChicagoBears = Team{
		ID:         idChicagoBears,
		Name:       TeamName{"Chicago", "Bears"},
		Conference: NFC,
		Division:   NFCNorth,
	}

	DetroitLions = Team{
		ID:         idDetroitLions,
		Name:       TeamName{"Detroit", "Lions"},
		Conference: NFC,
		Division:   NFCNorth,
	}

	TampaBayBuccaneers = Team{
		ID:         idTampaBayBuccaneers,
		Name:       TeamName{"Tampa Bay", "Buccaneers"},
		Conference: NFC,
		Division:   NFCSouth,
	}

	NewOrleansSaints = Team{
		ID:         idNewOrleansSaints,
		Name:       TeamName{"New Orleans", "Saints"},
		Conference: NFC,
		Division:   NFCSouth,
	}

	CarolinaPanthers = Team{
		ID:         idCarolinaPanthers,
		Name:       TeamName{"Carolina", "Panthers"},
		Conference: NFC,
		Division:   NFCSouth,
	}

	AtlantaFalcons = Team{
		ID:         idAtlantaFalcons,
		Name:       TeamName{"Atlanta", "Falcons"},
		Conference: NFC,
		Division:   NFCSouth,
	}

	LosAngelesRams = Team{
		ID:         idLosAngelesRams,
		Name:       TeamName{"Los Angeles", "Rams"},
		Conference: NFC,
		Division:   NFCWest,
	}

	SanFrancisco49ers = Team{
		ID:         idSanFrancisco49ers,
		Name:       TeamName{"San Francisco", "49ers"},
		Conference: NFC,
		Division:   NFCWest,
	}

	SeattleSeahawks = Team{
		ID:         idSeattleSeahawks,
		Name:       TeamName{"Seattle", "Seahawks"},
		Conference: NFC,
		Division:   NFCWest,
	}

	ArizonaCardinals = Team{
		ID:         idArizonaCardinals,
		Name:       TeamName{"Arizona", "Cardinals"},
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

func (tn TeamName) String() string {
	return fmt.Sprintf("%s %s", tn.Location, tn.Name)
}
