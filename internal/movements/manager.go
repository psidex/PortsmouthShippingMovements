package movements

import (
	"github.com/psidex/PortsmouthShippingMovements/internal/images"
	"github.com/psidex/PortsmouthShippingMovements/internal/shipinfo"
	"strings"
	"sync"
)

// Manager is a goroutine-safe controller for storing and updating movement lists for today and tomorrow.
type Manager struct {
	mu                *sync.Mutex       // The mutex for this handler. Being a pointer ensures that it's never copied.
	imageUrlManager   images.UrlManager // The storage manager for ship image urls.
	scraper           Scraper           // The web scraper for Movement data.
	todayMovements    []Movement        // A slice containing Movement structs for today.
	tomorrowMovements []Movement        // A slice containing Movement structs for tomorrow.
}

// NewManager creates a new Manager.
// Any Manager should be passed as a pointer as the setters reassign the movement slice fields.
func NewManager(urlManager images.UrlManager, scraper Scraper) *Manager {
	return &Manager{
		mu:              &sync.Mutex{},
		imageUrlManager: urlManager,
		scraper:         scraper,
	}
}

// GetMovements returns the 2 Movement slices stored by the Manager.
func (m Manager) GetMovements() (todayMovements []Movement, tomorrowMovements []Movement) {
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

// postProcessMovements does post processing for the scraped movement data, such as setting an image.
// It directly changes the values stored by the slice so it is not required to return anything.
func (m Manager) postProcessMovements(movementSlice []Movement) {
	// We have to iterate this way so we can change the values in the slice. Using range would make a copy of the
	// elements which means we wouldn't be able to change the actual values inside the slice.
	// This could be sped up using goroutines but the API has a req/sec limit that I don't want to hit.
	for i := 0; i < len(movementSlice); i++ {
		if movementSlice[i].Type == Move {
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

			movementSlice[i].ImageUrl = m.imageUrlManager.GetUrl(query)
			movementSlice[i].VesselFinderUrl = shipinfo.GetVesselFinderUrl(movementSlice[i].Name)
		}
	}
}

// UpdateAllMovements updates the list of movements for both days.
func (m *Manager) UpdateAllMovements() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	todayMovements, err := m.scraper.GetTodayMovements()
	if err != nil {
		return err
	}
	m.postProcessMovements(todayMovements)
	m.todayMovements = todayMovements

	tomorrowMovements, err := m.scraper.GetTomorrowMovements()
	if err != nil {
		return err
	}
	m.postProcessMovements(tomorrowMovements)
	m.tomorrowMovements = tomorrowMovements

	return nil
}
