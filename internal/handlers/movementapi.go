package handlers

import (
	"encoding/json"
	"github.com/psidex/PortsmouthShippingMovements/internal/qhm"
	"log"
	"net/http"
)

// MovementApi contains data and methods for serving movement related API functions.
type MovementApi struct {
	movementMan *qhm.MovementManager
}

// NewMovementApi creates a new MovementApi.
func NewMovementApi(MovementManager *qhm.MovementManager) MovementApi {
	return MovementApi{movementMan: MovementManager}
}

// GetShippingMovements is an endpoint for shipping movement data.
func (a MovementApi) GetShippingMovements(w http.ResponseWriter, r *http.Request) {
	todayMovements, tomorrowMovements := a.movementMan.Movements()

	currentMovements := map[string][]qhm.Movement{
		"today":    todayMovements,
		"tomorrow": tomorrowMovements,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(currentMovements)
	if err != nil {
		log.Printf("GetShippingMovements: encoding API response failed: %s", err)
	}
}
