package movements

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/psidex/PortsmouthShippingMovements/internal/images"
	"io"
	"net/http"
	"strings"
	"time"
)

// This file has a couple of "FIXME" comments that are to do with debugging, to be fixed later.

// UAString is the custom User Agent string for web requests made by this program.
// TODO: Currently hardcoded to my email address, probably change this to read from a config file in the future?
const UAString = "PortsmouthShippingMovements/0.1 (simjenner3@gmail.com)"

// dailyMovementUrl is the base URL for daily movements.
//FIXME: const dailyMovementUrl = "https://www.royalnavy.mod.uk/qhm/portsmouth/shipping-movements/daily-movements?date="
const dailyMovementUrl = "http://127.0.0.1:8000/qhm.html"

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
		} else {
			// If it's not a notice it's a ship so we need an image.
			thisMovement.ImageUrl = images.GetImageForShip(thisMovement.Name)
		}

		movements = append(movements, thisMovement)
	})

	return movements, nil
}

// getMovements returns a slice of Movement structs containing the data for the given date.
func getMovements(dt time.Time) ([]Movement, error) {
	//FIXME: query := dailyMovementUrl + dt.Format("02/01/2006") // dd/mm/yyyy
	query := dailyMovementUrl

	client := &http.Client{}

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return []Movement{}, err
	}

	req.Header.Set("User-Agent", UAString)

	resp, err := client.Do(req)
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
