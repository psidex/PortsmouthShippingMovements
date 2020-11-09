package movements

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/psidex/PortsmouthShippingMovements/internal/shipinfo"
	"io"
	"net/http"
	"strings"
	"time"
)

// dailyMovementUrl is the base URL for daily movements.
//const dailyMovementUrl = "http://127.0.0.1:8000/qhm.html?q="
const dailyMovementUrl = "https://www.royalnavy.mod.uk/qhm/portsmouth/shipping-movements/daily-movements?date="

// Scraper is for dealing with requesting and parsing movements from the QHM.
type Scraper struct {
	client   *http.Client // The http Client.
	uaString string       // uaString is the custom User Agent string for web requests made by this program.
}

// NewScraper creates a new Scraper for scraping movement data from the QHM page.
func NewScraper(client *http.Client, contactEmail string) Scraper {
	return Scraper{
		uaString: "PortsmouthShippingMovements/0.1 (" + contactEmail + ")",
		client:   client,
	}
}

// dailyMovementHtmlToStruct takes the body from a request to dailyMovementUrl and extracts the movements.
func (m Scraper) dailyMovementHtmlToStruct(body io.ReadCloser) ([]Movement, error) {
	var movements []Movement

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return []Movement{}, err
	}

	doc.Find(".qhm-shipping-movements>tbody>tr").Each(func(i int, s *goquery.Selection) {
		thisMovement := Movement{
			Type: Move,
		}

		s.Find("td").Each(func(i int, s *goquery.Selection) {
			dataName, exists := s.Attr("data-th")
			if exists {
				// Name names and other things have long whitespace in-between words sometimes, not sure why.
				tdTextSplit := strings.Fields(s.Text())
				tdText := strings.Join(tdTextSplit, " ")

				switch dataName {
				case "Ser":
					thisMovement.Position = tdText
				case "Time":
					thisMovement.Time = tdText
				case "Ship":
					thisMovement.Name = tdText
				case "From":
					thisMovement.From = locationFromAbbreviation(tdText)
				case "To":
					thisMovement.To = locationFromAbbreviation(tdText)
				case "Methods":
					thisMovement.Method = tdText
				case "Tug":
					thisMovement.Remarks = tdText
				}

			}
		})

		if thisMovement.From.Name == "" && thisMovement.To.Name == "" {
			thisMovement.Type = Notice
		}

		// Movement images will be set
		movements = append(movements, thisMovement)
	})

	return movements, nil
}

// getMovements returns a slice of Movement structs containing the data for the given date.
func (m Scraper) getMovements(dt time.Time) ([]Movement, error) {
	query := dailyMovementUrl + dt.Format("02/01/2006") // dd/mm/yyyy

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return []Movement{}, err
	}

	req.Header.Set("User-Agent", m.uaString)

	resp, err := m.client.Do(req)
	if err != nil {
		return []Movement{}, err
	}
	defer resp.Body.Close()

	movements, err := m.dailyMovementHtmlToStruct(resp.Body)
	if err != nil {
		return []Movement{}, err
	}

	return movements, nil
}

// GetTodayMovements returns a slice of Movement structs containing the data for today.
func (m Scraper) GetTodayMovements() ([]Movement, error) {
	dt := time.Now()
	return m.getMovements(dt)
}

// GetTomorrowMovements returns a slice of Movement structs containing the data for tomorrow.
func (m Scraper) GetTomorrowMovements() ([]Movement, error) {
	dt := time.Now()
	tomorrow := dt.AddDate(0, 0, 1)
	return m.getMovements(tomorrow)
}

// locationFromAbbreviation returns a Location struct for a given abbreviation.
// If no location name can be found, the name is also set to the abbreviation.
func locationFromAbbreviation(abbreviation string) Location {
	name := abbreviation
	if locationName, ok := shipinfo.LocationAbbreviations[abbreviation]; ok {
		name = locationName
	}
	return Location{
		Abbreviation: abbreviation,
		Name:         name,
	}
}