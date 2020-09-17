package movements

import (
	"github.com/psidex/PortsmouthShippingMovements/internal/images"
	"github.com/psidex/PortsmouthShippingMovements/internal/shipinfo"
	"sync"
)

// MovementStorage is a goroutine-safe controller for storing and updating movement lists for today and tomorrow.
type MovementStorage struct {
	mu                *sync.Mutex                // The mutex for this handler. Being a pointer ensures that it's never copied.
	imageUrlStore     images.ShipImageUrlStorage // The storage manager for ship image urls.
	scraper           MovementScraper            // The web scraper for Movement data.
	todayMovements    []Movement                 // A slice containing Movement structs for today.
	tomorrowMovements []Movement                 // A slice containing Movement structs for tomorrow.
}

// NewMovementStorage creates a new MovementStorage.
// Any MovementStorage should be passed as a pointer as the setters reassign the movement slice fields.
func NewMovementStorage(imageUrlStore images.ShipImageUrlStorage, scraper MovementScraper) *MovementStorage {
	return &MovementStorage{
		mu:            &sync.Mutex{},
		imageUrlStore: imageUrlStore,
		scraper:       scraper,
	}
}

// GetMovements returns the 2 Movement slices stored by the MovementStorage.
func (m MovementStorage) GetMovements() (todayMovements []Movement, tomorrowMovements []Movement) {
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
func (m MovementStorage) postProcessMovements(movementSlice []Movement) {
	// We have to iterate this way so we can change the values in the slice. Using range would make a copy of the
	// elements which means we wouldn't be able to change the actual values inside the slice.
	// This could be sped up using goroutines but the API has a req/sec limit that I don't want to hit.
	for i := 0; i < len(movementSlice); i++ {
		if movementSlice[i].Type == Move {
			url := m.imageUrlStore.GetUrlForShip(movementSlice[i].Name)
			movementSlice[i].ImageUrl = url
			movementSlice[i].VesselFinderUrl = shipinfo.GetShipVesselFinderUrl(movementSlice[i].Name)
		}
	}
}

// UpdateAllMovements updates the list of movements for both days.
func (m *MovementStorage) UpdateAllMovements() error {
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
