package shipinfo

import "regexp"

const titleRegexString = "^(MV )"

var (
	titleRegex = regexp.MustCompile(titleRegexString)
)

// shipNameFromTitle returns a ships name from its title. The title can be slightly different, e.g. "MV" at the start
// for motor vessel. This is not part of the actual name.
func shipNameFromTitle(title string) string {
	return titleRegex.ReplaceAllString(title, "")
}

// GetShipVesselFinderUrl returns a URL to a vessel finder search for the given ship.
func GetShipVesselFinderUrl(shipTitle string) string {
	return "https://www.vesselfinder.com/vessels?name=" + shipNameFromTitle(shipTitle)
}
