package team

// Team IDs from the ESPN API (http://site.api.espn.com/apis/site/v2/sports/football/nfl/teams)
// Hardcode for now, but we could get these using the API if needed
const (
	idArizonaCardinals     = "22"
	idAtlantaFalcons       = "1"
	idBaltimoreRavens      = "33"
	idBuffaloBills         = "2"
	idCarolinaPanthers     = "29"
	idChicagoBears         = "3"
	idCincinnatiBengals    = "4"
	idClevelandBrowns      = "5"
	idDallasCowboys        = "6"
	idDenverBroncos        = "7"
	idDetroitLions         = "8"
	idGreenBayPackers      = "9"
	idHoustonTexans        = "34"
	idIndianapolisColts    = "11"
	idJacksonvilleJaguars  = "30"
	idKansasCityChiefs     = "12"
	idLasVegasRaiders      = "13"
	idLosAngelesChargers   = "24"
	idLosAngelesRams       = "14"
	idMiamiDolphins        = "15"
	idMinnesotaVikings     = "16"
	idNewEnglandPatriots   = "17"
	idNewOrleansSaints     = "18"
	idNewYorkGiants        = "19"
	idNewYorkJets          = "20"
	idPhiladelphiaEagles   = "21"
	idPittsburghSteelers   = "23"
	idSanFrancisco49ers    = "25"
	idSeattleSeahawks      = "26"
	idTampaBayBuccaneers   = "27"
	idTennesseeTitans      = "10"
	idWashingtonCommanders = "28"
)

var TeamIDs = []string{
	idNewEnglandPatriots, idNewYorkJets, idBuffaloBills, idMiamiDolphins, // AFC East
	idPittsburghSteelers, idBaltimoreRavens, idClevelandBrowns, idCincinnatiBengals, // AFC North
	idTennesseeTitans, idIndianapolisColts, idJacksonvilleJaguars, idHoustonTexans, // AFC South
	idKansasCityChiefs, idLasVegasRaiders, idLosAngelesChargers, idDenverBroncos, // AFC West

	idDallasCowboys, idNewYorkGiants, idPhiladelphiaEagles, idWashingtonCommanders, // NFC East
	idGreenBayPackers, idMinnesotaVikings, idChicagoBears, idDetroitLions, // NFC North
	idTampaBayBuccaneers, idNewOrleansSaints, idCarolinaPanthers, idAtlantaFalcons, // NFC South
	idLosAngelesRams, idSanFrancisco49ers, idSeattleSeahawks, idArizonaCardinals, // NFC West
}
