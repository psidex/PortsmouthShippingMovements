package movements

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
	"time"
)

// dailyMovementUrl is the base URL for daily movements.
const dailyMovementUrl = "https://www.royalnavy.mod.uk/qhm/portsmouth/shipping-movements/daily-movements?date="

// MovementScraper is for dealing with requesting and parsing movements from the QHM.
type MovementScraper struct {
	client   *http.Client // The http Client.
	uaString string       // uaString is the custom User Agent string for web requests made by this program.
}

// NewMovementScraper creates a new MovementScraper.
func NewMovementScraper(contactEmail string) MovementScraper {
	return MovementScraper{
		uaString: "PortsmouthShippingMovements/0.1 (" + contactEmail + ")",
		client:   &http.Client{Timeout: time.Second * 10},
	}
}

// dailyMovementHtmlToStruct takes the body from a request to dailyMovementUrl and extracts the movements.
func (m MovementScraper) dailyMovementHtmlToStruct(body io.ReadCloser) ([]Movement, error) {
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
func (m MovementScraper) getMovements(dt time.Time) ([]Movement, error) {
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
func (m MovementScraper) GetTodayMovements() ([]Movement, error) {
	dt := time.Now()
	return m.getMovements(dt)
}

// GetTomorrowMovements returns a slice of Movement structs containing the data for tomorrow.
func (m MovementScraper) GetTomorrowMovements() ([]Movement, error) {
	dt := time.Now()
	tomorrow := dt.AddDate(0, 0, 1)
	return m.getMovements(tomorrow)
}
