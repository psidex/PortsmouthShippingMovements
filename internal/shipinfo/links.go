package shipinfo

import (
	"net/url"
	"regexp"
)

const titlePrefixRegexString = "^(MV )"

var (
	titlePrefixRegex = regexp.MustCompile(titlePrefixRegexString)
)

// shipNameFromTitle returns a ships name from its title. The title can be slightly different, e.g. "MV" at the start
// for motor vessel. This is not part of the actual name.
func shipNameFromTitle(title string) string {
	return titlePrefixRegex.ReplaceAllString(title, "")
}

// GetShipVesselFinderUrl returns a URL to a vessel finder search for the given ship.
func GetShipVesselFinderUrl(shipTitle string) string {
	return "https://www.vesselfinder.com/vessels?name=" + url.QueryEscape(shipNameFromTitle(shipTitle))
}
