package qhm

// Location holds the names for a single location, to be used in the Movement type.
type Location struct {
	Abbreviation string `json:"abbreviation"` // The abbreviation of the location.
	Name         string `json:"name"`         // The full name of the location.
}

// newLocation creates a Location from an abbreviation.
func newLocation(abbreviation string) Location {
	name := abbreviation
	if locationName, ok := locationAbbreviations[abbreviation]; ok {
		name = locationName
	}
	return Location{
		Abbreviation: abbreviation,
		Name:         name,
	}
}

// locationAbbreviations is a string:string map of location abbreviation to full name.
var locationAbbreviations = map[string]string{
	"NAB":       "Nab Tower",
	"SRJ":       "South Railway Jetty",
	"SRJ ( S )": "South Railway Jetty (South)",
	"SRJ ( C )": "South Railway Jetty (Centre)",
	"SRJ ( N )": "South Railway Jetty (North)",
	"SJ":        "Sheer Jetty",
	"VJ":        "Victory Jetty",
	"PRJ":       "Princess Royal Jetty",
	"NCJ":       "North Corner Jetty",
	"NCJ ( W )": "North Corner Jetty (West)",
	"NCJ ( E )": "North Corner Jetty (East)",
	"SWW":       "South West Wall",
	"SW":        "South Wall",
	"SW ( W )":  "South Wall (West)",
	"SW ( E )":  "South Wall (East)",
	"NW":        "North Wall",
	"NW ( E )":  "North Wall (East)",
	"NW ( W )":  "North Wall (West)",
	"NWW":       "North West Wall",
	"NWW ( S )": "North West Wall (South)",
	"NWW ( N )": "North West Wall (North)",
	"FLJ":       "Fountain Lake Jetty",
	"FLJ 1":     "Fountain Lake Jetty 1",
	"FLJ 2":     "Fountain Lake Jetty 2",
	"FLJ 3":     "Fountain Lake Jetty 3",
	"FLJ 4":     "Fountain Lake Jetty 4",
	"FLJ 5":     "Fountain Lake Jetty 5",
	"OFJ":       "Oil Fuel Jetty",
	"OFJ ( N )": "Oil Fuel Jetty (North)",
	"BII":       "Basin 2",
	"BIII":      "Basin 3",
	"O/B":       "Outboard",
	"OSB":       "Outer Spit Buoy",
	"HBR":       "Harbour",
	"UHAF":      "Upper Harbour Ammunitioning Facility",
	"Z M’RGS":   "Z Moorings",
	"BP":        "Bedenham Pier",
	"HC":        "Haslar Creek",
	"PC":        "Portchester Creek",
	"PP":        "Petrol Pier",
	"SH":        "Spit Head",
	"PIP":       "Portsmouth International Port",
	"WLM":       "Wightlink Moorings",
	"RCY":       "Royal Clarence Yard",
	"BT/TX":     "Boat Transfer",
	"RAAON":     "Remain At Anchor Overnight",
	"TCL":       "Tank Cleaner",
	"HORB":      "Hold Off Re-Berth",
	"WIND":      "Wind Ship (Cold Move Using Tugs To Turn Ship And Re-Berth)",
}
