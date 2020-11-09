package api

import (
	"encoding/json"
	"github.com/psidex/PortsmouthShippingMovements/internal/movements"
	"log"
	"net/http"
)

// MovementApi contains data and methods for serving movement related API functions.
type MovementApi struct {
	movementManager *movements.Manager
}

// NewMovementApi creates a new MovementApi.
func NewMovementApi(MovementManager *movements.Manager) MovementApi {
	return MovementApi{movementManager: MovementManager}
}

// GetShippingMovements is an endpoint for shipping movement data.
func (a MovementApi) GetShippingMovements(w http.ResponseWriter, r *http.Request) {
	todayMovements, tomorrowMovements := a.movementManager.GetMovements()

	currentMovements := map[string][]movements.Movement{
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
