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

// dailyMovementHtmlToStruct takes the body from a request to dailyMovementUrl and extracts the movements.
func dailyMovementHtmlToStruct(body io.ReadCloser) ([]Movement, error) {
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

		movements = append(movements, thisMovement)
	})

	return movements, nil
}

// getMovements returns a slice of Movement structs containing the data for the given date.
func getMovements(dt time.Time) ([]Movement, error) {
	query := dailyMovementUrl + dt.Format("02/01/2006") // dd/mm/yyyy

	resp, err := http.Get(query)
	if err != nil {
		return []Movement{}, err
	}
	defer resp.Body.Close()

	movements, err := dailyMovementHtmlToStruct(resp.Body)
	if err != nil {
		return []Movement{}, err
	}

	return movements, nil
}

// GetTodayMovements returns a slice of Movement structs containing the data for today.
func GetTodayMovements() ([]Movement, error) {
	dt := time.Now()
	return getMovements(dt)
}

// GetTomorrowMovements returns a slice of Movement structs containing the data for tomorrow.
func GetTomorrowMovements() ([]Movement, error) {
	dt := time.Now()
	tomorrow := dt.AddDate(0, 0, 1)
	return getMovements(tomorrow)
}
