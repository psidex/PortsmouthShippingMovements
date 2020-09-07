package movements

var abbreviationMap = map[string]string{
	"SRJ":     "SOUTH RAILWAY JETTY",
	"SRJ(S)":  "SOUTH RAILWAY JETTY (SOUTH)",
	"SRJ(C)":  "SOUTH RAILWAY JETTY (CENTRE)",
	"SRJ(N)":  "SOUTH RAILWAY JETTY (NORTH)",
	"SJ":      "SHEER JETTY",
	"VJ":      "VICTORY JETTY",
	"PRJ":     "PRINCESS ROYAL JETTY",
	"NCJ":     "NORTH CORNER JETTY",
	"NCJ(W)":  "NORTH CORNER JETTY (WEST)",
	"NCJ(E)":  "NORTH CORNER JETTY (EAST)",
	"SWW":     "SOUTH WEST WALL",
	"SW":      "SOUTH WALL",
	"SW(W)":   "SOUTH WALL (WEST)",
	"SW(E)":   "SOUTH WALL (EAST)",
	"NW":      "NORTH WALL",
	"NW(E)":   "NORTH WALL (EAST)",
	"NW(W)":   "NORTH WALL (WEST)",
	"NWW":     "NORTH WEST WALL",
	"NWW(S)":  "NORTH WEST WALL (SOUTH)",
	"NWW(N)":  "NORTH WEST WALL (NORTH)",
	"FLJ":     "FOUNTAIN LAKE JETTY",
	"FLJ1":    "FOUNTAIN LAKE JETTY 1",
	"FLJ2":    "FOUNTAIN LAKE JETTY 2",
	"FLJ3":    "FOUNTAIN LAKE JETTY 3",
	"FLJ4":    "FOUNTAIN LAKE JETTY 4",
	"FLJ5":    "FOUNTAIN LAKE JETTY 5",
	"OFJ":     "OIL FUEL JETTY",
	"OFJ(N)":  "OIL FUEL JETTY (NORTH)",
	"BII":     "BASIN 2",
	"BIII":    "BASIN 3",
	"O/B":     "OUTBOARD",
	"OSB":     "OUTER SPIT BUOY",
	"HBR":     "HARBOUR",
	"UHAF":    "UPPER HARBOUR AMMUNITIONING FACILITY",
	"Z M’RGS": "Z MOORINGS",
	"BP":      "BEDENHAM PIER",
	"HC":      "HASLAR CREEK",
	"PC":      "PORTCHESTER CREEK",
	"PP":      "PETROL PIER",
	"SH":      "SPIT HEAD",
	"PIP":     "PORTSMOUTH INTERNATIONAL PORT",
	"WLM":     "WIGHTLINK MOORINGS",
	"RCY":     "ROYAL CLARENCE YARD",
	"BT/TX":   "BOAT TRANSFER",
	"RAAON":   "REMAIN AT ANCHOR OVERNIGHT",
	"TCL":     "TANK CLEANER",
	"HORB":    "HOLD OFF RE-BERTH",
	"WIND":    "WIND SHIP (COLD MOVE USING TUGS TO TURN SHIP AND RE-BERTH)",
	"LNTM":    "LOCAL NOTICE TO MARINERS",
}

// locationFromAbbreviation returns a Location struct for a given abbreviation.
// If no location name can be found, the name is also set to the abbreviation.
func locationFromAbbreviation(abbreviation string) Location {
	name := abbreviation
	if locationName, ok := abbreviationMap[abbreviation]; ok {
		name = locationName
	}
	return Location{
		Abbreviation: abbreviation,
		Name:         name,
	}
}
