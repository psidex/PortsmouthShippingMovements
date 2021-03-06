package qhm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

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

// abbrvRx contains regex and replacement text for things in an abbreviation that need replacing this way.
var abbrvRx = map[*regexp.Regexp]string{
	regexp.MustCompile(`\((| )N(| )\)`): "(North)",
	regexp.MustCompile(`\((| )E(| )\)`): "(East)",
	regexp.MustCompile(`\((| )S(| )\)`): "(South)",
	regexp.MustCompile(`\((| )W(| )\)`): "(West)",
	regexp.MustCompile(`\((| )C(| )\)`): "(Centre)",
}

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

	// Potentially a more efficient way of doing this, but benchmarking this I don't think it matters too much.
	for r, text := range abbrvRx {
		abbrv = r.ReplaceAllString(abbrv, text)
	}

	parsed := ""
	for _, w := range strings.Split(abbrv, " ") {
		runed := []rune(w)
		if first := runed[0]; unicode.IsNumber(first) {
			berth, _ := strconv.Atoi(string(first))
			parsed += fmt.Sprintf("Berth %d ", berth)
			w = string(runed[1:])
		}

		// With the above berth number logic and also splitting on spaces, trim(w) might be "".
		if strings.TrimSpace(w) != "" {
			if full, ok := qhmAbbreviations[w]; ok {
				parsed += full
			} else {
				parsed += w
			}
			parsed += " "
		}
	}
	return strings.TrimSpace(parsed)
}
