package qhm

import (
	"github.com/psidex/PortsmouthShippingMovements/internal/images"
	"github.com/psidex/PortsmouthShippingMovements/internal/vesselfinder"
	"log"
	"strings"
	"sync"
	"time"
)

// MovementManager is a goroutine-safe manager for storing and updating movement lists for today and tomorrow.
type MovementManager struct {
	mu                *sync.Mutex       // The mutex for this handler. Being a pointer ensures that it's never copied.
	imageUrlMan       images.UrlManager // The storage manager for ship image urls.
	scraper           Scraper           // The web scraper for QHM data.
	todayMovements    []Movement        // A slice containing Movement structs for today.
	tomorrowMovements []Movement        // A slice containing Movement structs for tomorrow.
}

// NewMovementManager creates a new MovementManager.
func NewMovementManager(imageUrlMan images.UrlManager, scraper Scraper) *MovementManager {
	// Any MovementManager should be passed as a pointer as the setters reassign the movement slice fields.
	return &MovementManager{
		mu:          &sync.Mutex{},
		imageUrlMan: imageUrlMan,
		scraper:     scraper,
	}
}

// Movements returns the 2 Movement slices stored by the MovementManager.
func (m MovementManager) Movements() (todayMovements []Movement, tomorrowMovements []Movement) {
	// This Lock() doesn't need to be at the top of the function but putting it here helps prevent bugs in the future.
	m.mu.Lock()
	defer m.mu.Unlock()

	todayMovements = make([]Movement, len(m.todayMovements))
	tomorrowMovements = make([]Movement, len(m.tomorrowMovements))

	// A copy is made of the slices so that they can't be edited outside of this func (as a slice is a reference type).
	copy(todayMovements, m.todayMovements)
	copy(tomorrowMovements, m.tomorrowMovements)

	return todayMovements, tomorrowMovements
}

// insertMetadata fetches and inserts metadata for each movement. Values are edited in place so nothing is returned.
func (m MovementManager) insertMetadata(movementSlice []Movement) {
	// We iterate this way so we can change the values in place.
	for i := 0; i < len(movementSlice); i++ {
		if movementSlice[i].Type == Move {
			shipTitle := movementSlice[i].Name

			// If there are multiple ships referenced in one movement, just get image for first one.
			splitters := []string{",", "&", " AND "}
			for _, splitter := range splitters {
				if strings.Contains(shipTitle, splitter) {
					shipTitle = strings.Split(shipTitle, splitter)[0]
				}
			}

			shipTitle = strings.TrimSpace(shipTitle)

			// Prepend "Portsmouth " to image search so that a generic name like "TUG" will still show a relevant image.
			movementSlice[i].ImageUrl = m.imageUrlMan.GetUrl("Portsmouth " + shipTitle)
			movementSlice[i].VesselFinderUrl = vesselfinder.GetSearchUrl(shipTitle)
		}
	}
}

// Update updates the list of movements for both days.
func (m *MovementManager) Update() {
	log.Println("Updating movement store")
	m.mu.Lock()
	defer m.mu.Unlock()

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
