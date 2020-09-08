package movements

import "sync"

// MovementHandler is a goroutine-safe controller for storing movement lists for today and tomorrow.
type MovementHandler struct {
	mu                *sync.Mutex // The mutex for this handler. It being a pointer ensures that it's never copied.
	todayMovements    []Movement  // A slice containing Movement structs for today.
	tomorrowMovements []Movement  // A slice containing Movement structs for tomorrow.
}

// NewMovementHandler creates a new MovementHandler.
// Any MovementHandler should be passed as a pointer as the setters reassign the movement slice fields.
func NewMovementHandler() *MovementHandler {
	return &MovementHandler{
		mu: &sync.Mutex{},
	}
}

// TodayMovements returns the stored list of movements for today.
func (m MovementHandler) TodayMovements() []Movement {
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
func (m MovementHandler) TomorrowMovements() []Movement {
	m.mu.Lock()
	defer m.mu.Unlock()
	tmp := make([]Movement, len(m.tomorrowMovements))
	copy(tmp, m.tomorrowMovements)
	return tmp
}

// SetTodayMovements sets the list of movements for today.
func (m *MovementHandler) SetTodayMovements(movementSlice []Movement) {
	m.mu.Lock()
	defer m.mu.Unlock()
	tmp := make([]Movement, len(movementSlice))
	copy(tmp, movementSlice)
	m.todayMovements = tmp
}

// SetTomorrowMovements sets the list of movements for tomorrow.
func (m *MovementHandler) SetTomorrowMovements(movementSlice []Movement) {
	m.mu.Lock()
	defer m.mu.Unlock()
	tmp := make([]Movement, len(movementSlice))
	copy(tmp, movementSlice)
	m.tomorrowMovements = tmp
}
