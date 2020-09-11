package movements

import (
	"github.com/psidex/PortsmouthShippingMovements/internal/images"
	"github.com/psidex/PortsmouthShippingMovements/internal/shipinfo"
	"sync"
)

// MovementStorage is a goroutine-safe controller for storing movement lists for today and tomorrow.
type MovementStorage struct {
	mu                *sync.Mutex                // The mutex for this handler. It being a pointer ensures that it's never copied.
	imageUrlStore     images.ShipImageUrlStorage // The storage manager for ship image urls.
	todayMovements    []Movement                 // A slice containing Movement structs for today.
	tomorrowMovements []Movement                 // A slice containing Movement structs for tomorrow.
}

// NewMovementStorage creates a new MovementStorage.
// Any MovementStorage should be passed as a pointer as the setters reassign the movement slice fields.
func NewMovementStorage(imageUrlStore images.ShipImageUrlStorage) *MovementStorage {
	return &MovementStorage{
		mu:            &sync.Mutex{},
		imageUrlStore: imageUrlStore,
	}
}

// TodayMovements returns the stored list of movements for today.
func (m MovementStorage) TodayMovements() []Movement {
	// This Lock() doesn't need to be at the top of the function but putting it here helps prevent bugs in the future.
	m.mu.Lock()
	defer m.mu.Unlock()
	// For each of these getters and setters, a copy is made of the slice so that it can't be edited outside of these
	// functions (as a slice is a reference type).
	tmp := make([]Movement, len(m.todayMovements))
	copy(tmp, m.todayMovements)
	return tmp
}

// TomorrowMovements returns the stored list of movements for tomorrow.
func (m MovementStorage) TomorrowMovements() []Movement {
	m.mu.Lock()
	defer m.mu.Unlock()
	tmp := make([]Movement, len(m.tomorrowMovements))
	copy(tmp, m.tomorrowMovements)
	return tmp
}

// postProcessMovements does post processing for the scraped movement data, such as setting an image.
func (m MovementStorage) postProcessMovements(movementSlice []Movement) {
	// We have to iterate this way so we can change the values in the slice. Using range would make a copy of the
	// elements which means we wouldn't be able to change the actual values inside the slice.
	for i := 0; i < len(movementSlice); i++ {
		if movementSlice[i].Type == Move {
			url := m.imageUrlStore.GetUrlForShip(movementSlice[i].Name)
			movementSlice[i].ImageUrl = url
			movementSlice[i].VesselFinderUrl = shipinfo.GetShipVesselFinderUrl(movementSlice[i].Name)
		}
	}
}

// SetTodayMovements sets the list of movements for today. Will set image urls for given Movements.
func (m *MovementStorage) SetTodayMovements(movementSlice []Movement) {
	m.mu.Lock()
	defer m.mu.Unlock()
	tmp := make([]Movement, len(movementSlice))
	copy(tmp, movementSlice)
	m.postProcessMovements(tmp)
	m.todayMovements = tmp
}

// SetTomorrowMovements sets the list of movements for tomorrow. Will set image urls for given Movements.
func (m *MovementStorage) SetTomorrowMovements(movementSlice []Movement) {
	m.mu.Lock()
	defer m.mu.Unlock()
	tmp := make([]Movement, len(movementSlice))
	copy(tmp, movementSlice)
	m.postProcessMovements(tmp)
	m.tomorrowMovements = tmp
}
