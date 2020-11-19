package qhm

import (
	"encoding/json"
	"github.com/peterbourgon/diskv/v3"
	"github.com/psidex/PortsmouthShippingMovements/internal/bing"
	"github.com/psidex/PortsmouthShippingMovements/internal/config"
	"github.com/psidex/PortsmouthShippingMovements/internal/vesselfinder"
	"log"
	"net/http"
	"strings"
	"time"
)

// MovementManager is a goroutine-safe manager for storing, updating, and serving movement lists for today and tomorrow.
type MovementManager struct {
	imageSearch       bing.ImageSearchApi // The API for searching for image urls.
	searchLimiter     <-chan time.Time    // Basic rate limiter for image searches.
	urlCache          *diskv.Diskv        // Cache for image urls.
	scraper           scraper             // The web scraper for QHM data.
	todayMovements    []Movement          // A slice containing Movement structs for today.
	tomorrowMovements []Movement          // A slice containing Movement structs for tomorrow.
}

// NewMovementManager creates a new MovementManager.
func NewMovementManager(c config.Config) *MovementManager {
	d := diskv.New(diskv.Options{
		BasePath:     c.StoragePath,
		CacheSizeMax: 1024 * 1024, // 1MB
	})

	searchHttpClient := &http.Client{Timeout: time.Second * 10}
	scrapeHttpClient := &http.Client{Timeout: time.Second * 10}

	imageSearch := bing.NewImageSearchApi(searchHttpClient, c.BingImageSearchApiKey)
	scraper := newScraper(scrapeHttpClient, c.ContactEmail)

	// Allow a max of 1 req a second, bit strong but guarantees no rate limit errors.
	limiter := time.Tick(1 * time.Second)

	// Any MovementManager should be passed as a pointer as the setters reassign the movement slice fields.
	return &MovementManager{
		imageSearch:   imageSearch,
		searchLimiter: limiter,
		urlCache:      d,
		scraper:       scraper,
	}
}

// insertMetadata fetches and inserts metadata for each movement. Values are edited in place so nothing is returned.
func (m MovementManager) insertMetadata(movementSlice []Movement) {
	// We iterate this way so we can change the values in place.
	for i := 0; i < len(movementSlice); i++ {
		if movementSlice[i].Type != Move {
			continue
		}

		shipTitle := movementSlice[i].Name

		// If there are multiple ships referenced in one movement, just get image for first one.
		splitters := []string{",", "&", " AND "}
		for _, splitter := range splitters {
			if strings.Contains(shipTitle, splitter) {
				shipTitle = strings.Split(shipTitle, splitter)[0]
			}
		}

		shipTitle = strings.TrimSpace(shipTitle)
		var imageUrl string

		if m.urlCache.Has(shipTitle) {
			iuBytes, err := m.urlCache.Read(shipTitle)
			imageUrl = string(iuBytes)
			if err != nil {
				log.Printf("Error reading %s from urlCache: %s\n", shipTitle, err)
			}
		} else {
			log.Printf("Searching for image of: \"Portsmouth %s\"", shipTitle)
			<-m.searchLimiter
			var err error

			// Prepend "Portsmouth " to image search so that a generic name will still find a relevant image.
			imageUrl, err = m.imageSearch.SearchForImageUrl("Portsmouth " + shipTitle)
			if err != nil {
				log.Printf("Error searching for image of \"Portsmouth %s\": %s\n", shipTitle, err)
			} else {
				// If write fails we'll just get the url again next time.
				_ = m.urlCache.Write(shipTitle, []byte(imageUrl))
			}
		}

		movementSlice[i].ImageUrl = imageUrl
		movementSlice[i].VesselFinderUrl = vesselfinder.GetSearchUrl(shipTitle)
	}
}

// Update updates the list of movements for both days.
func (m *MovementManager) Update() {
	log.Println("Updating movement store")

	tToday := time.Now()
	tTomorrow := tToday.AddDate(0, 0, 1)

	todayMovements, err := m.scraper.getMovements(tToday)
	if err != nil {
		log.Printf("Error scraping movements for today: %s", err)
		return
	}
	m.insertMetadata(todayMovements)
	m.todayMovements = todayMovements

	tomorrowMovements, err := m.scraper.getMovements(tTomorrow)
	if err != nil {
		log.Printf("Error scraping movements for tomorrow: %s", err)
		return
	}
	m.insertMetadata(tomorrowMovements)
	m.tomorrowMovements = tomorrowMovements
}

// SendCurrentMovements is an endpoint for handling a request for "current" movement data.
func (m *MovementManager) SendCurrentMovements(w http.ResponseWriter, _ *http.Request) {
	currentMovements := map[string][]Movement{
		"today":    m.todayMovements,
		"tomorrow": m.tomorrowMovements,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(currentMovements)
	if err != nil {
		log.Printf("SendCurrentMovements: encoding API response failed: %s", err)
	}
}
