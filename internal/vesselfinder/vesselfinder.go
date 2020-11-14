package vesselfinder

import (
	"net/url"
	"regexp"
)

const titlePrefixRegexString = "^(MV |SD )"

var (
	titlePrefixRegex = regexp.MustCompile(titlePrefixRegexString)
)

// shipNameFromTitle returns a ships name from its title. The title can be slightly different, e.g. "MV" at the start
// for motor vessel. This is not part of the actual name.
func shipNameFromTitle(title string) string {
	return titlePrefixRegex.ReplaceAllString(title, "")
}

// GetSearchUrl returns a URL to a vessel finder search for the given ship.
func GetSearchUrl(shipTitle string) string {
	return "https://www.vesselfinder.com/vessels?name=" + url.QueryEscape(shipNameFromTitle(shipTitle))
}
