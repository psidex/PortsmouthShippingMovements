package qhm

import (
	"github.com/psidex/PortsmouthShippingMovements/internal/images"
	"github.com/psidex/PortsmouthShippingMovements/internal/vesselfinder"
	"log"
	"strings"
	"sync"
	"time"
)

// MovementManager is a goroutine-safe controller for storing and updating movement lists for today and tomorrow.
type MovementManager struct {
	mu                *sync.Mutex       // The mutex for this handler. Being a pointer ensures that it's never copied.
	imageUrlMan       images.UrlManager // The storage manager for ship image urls.
	scraper           Scraper           // The web scraper for QHM data.
	todayMovements    []Movement        // A slice containing Movement structs for today.
	tomorrowMovements []Movement        // A slice containing Movement structs for tomorrow.
}

// NewMovementManager creates a new MovementManager.
// Any MovementManager should be passed as a pointer as the setters reassign the movement slice fields.
func NewMovementManager(imageUrlMan images.UrlManager, scraper Scraper) *MovementManager {
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

// postProcess does post processing for the scraped movement data, such as setting an image.
// It directly changes the values stored by the slice so it does not return anything.
func (m MovementManager) postProcess(movementSlice []Movement) {
	// We have to iterate this way so we can change the values in the slice. Using range would make a copy of the
	// elements which means we wouldn't be able to change the actual values inside the slice.
	for i := 0; i < len(movementSlice); i++ {
		if movementSlice[i].Type == Move {
			// Create a query to search for an image of the ship.
			query := movementSlice[i].Name

			// If there are multiple ships referenced in one movement, just get image for first one.
			if strings.Contains(query, ",") {
				query = strings.Split(query, ",")[0]
			} else if strings.Contains(query, "&") {
				query = strings.Split(query, "&")[0]
			}

			query = strings.TrimSpace(query)

			// Prepend "Portsmouth " to image search so that a generic name like "TUG" will still show a relevant image.
			query = "Portsmouth " + query

			movementSlice[i].ImageUrl = m.imageUrlMan.GetUrl(query)
			movementSlice[i].VesselFinderUrl = vesselfinder.GetSearchUrl(movementSlice[i].Name)
		}
	}
}

// Update updates the list of movements for both days.
func (m *MovementManager) Update() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	tToday := time.Now()
	tTomorrow := tToday.AddDate(0, 0, 1)

	todayMovements, err := m.scraper.getMovements(tToday)
	if err != nil {
		return err
	}
	m.postProcess(todayMovements)
	m.todayMovements = todayMovements

	tomorrowMovements, err := m.scraper.getMovements(tTomorrow)
	if err != nil {
		return err
	}
	m.postProcess(tomorrowMovements)
	m.tomorrowMovements = tomorrowMovements

	return nil
}

// PrettyUpdate runs m.Update with logging and error handling.
func (m *MovementManager) PrettyUpdate() {
	log.Println("Updating movement store")
	err := m.Update()
	if err != nil {
		log.Printf("Error when updating movement store: %s", err)
	}
}
