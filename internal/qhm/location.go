package qhm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	nRegex, _ = regexp.Compile(`\((| )N(| )\)`)
	eRegex, _ = regexp.Compile(`\((| )E(| )\)`)
	sRegex, _ = regexp.Compile(`\((| )S(| )\)`)
	wRegex, _ = regexp.Compile(`\((| )W(| )\)`)
	cRegex, _ = regexp.Compile(`\((| )C(| )\)`)
)

// Location holds the names for a single location, to be used in the Movement type.
type Location struct {
	Abbreviation string `json:"abbreviation"` // The abbreviation of the location.
	Name         string `json:"name"`         // The full name of the location.
}

// newLocation creates a Location from an abbreviation.
func newLocation(abbrv string) Location {
	return Location{
		Abbreviation: abbrv,
		Name:         parseAbbreviation(abbrv),
	}
}

func parseAbbreviation(abbrv string) string {
	if len(abbrv) <= 0 {
		return ""
	}

	parsed := ""

	// Potentially a more efficient way of doing this, but benchmarking this I don't think it matters too much.
	abbrv = nRegex.ReplaceAllString(abbrv, "(North)")
	abbrv = eRegex.ReplaceAllString(abbrv, "(East)")
	abbrv = sRegex.ReplaceAllString(abbrv, "(South)")
	abbrv = wRegex.ReplaceAllString(abbrv, "(West)")
	abbrv = cRegex.ReplaceAllString(abbrv, "(Centre)")

	for _, w := range strings.Split(abbrv, " ") {
		berth := 0
		runed := []rune(w)

		if first := runed[0]; unicode.IsNumber(first) {
			berth, _ = strconv.Atoi(string(first))
			w = string(runed[1:])
		}

		if full, ok := qhmAbbreviations[w]; ok {
			parsed += full
		} else {
			parsed += w
		}

		if berth > 0 {
			parsed += fmt.Sprintf("Berth %d", berth)
		}

		parsed += " "
	}

	return strings.TrimSpace(parsed)
}

// qhmAbbreviations is a string:string map of QHM abbreviation to full text.
var qhmAbbreviations = map[string]string{
	"NAB":     "Nab Tower",
	"SRJ":     "South Railway Jetty",
	"SJ":      "Sheer Jetty",
	"VJ":      "Victory Jetty",
	"PRJ":     "Princess Royal Jetty",
	"NCJ":     "North Corner Jetty",
	"SWW":     "South West Wall",
	"SW":      "South Wall",
	"NW":      "North Wall",
	"NWW":     "North West Wall",
	"FLJ":     "Fountain Lake Jetty",
	"OFJ":     "Oil Fuel Jetty",
	"BII":     "Basin 2",
	"BIII":    "Basin 3",
	"O/B":     "Outboard",
	"OSB":     "Outer Spit Buoy",
	"HBR":     "Harbour",
	"UHAF":    "Upper Harbour Ammunitioning Facility",
	"Z M’RGS": "Z Moorings",
	"BP":      "Bedenham Pier",
	"HC":      "Haslar Creek",
	"PC":      "Portchester Creek",
	"PP":      "Petrol Pier",
	"SH":      "Spit Head",
	"PIP":     "Portsmouth International Port",
	"WLM":     "Wightlink Moorings",
	"RCY":     "Royal Clarence Yard",
	"BT/TX":   "Boat Transfer",
	"RAAON":   "Remain At Anchor Overnight",
	"TCL":     "Tank Cleaner",
	"HORB":    "Hold Off Re-Berth",
	"WIND":    "Wind Ship (Cold Move Using Tugs To Turn Ship And Re-Berth)",
}
