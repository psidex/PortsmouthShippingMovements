package api

import (
	"encoding/json"
	"github.com/psidex/PortsmouthShippingMovements/internal/movements"
	"net/http"
)

// MovementApi contains data and methods for serving movement related API functions.
type MovementApi struct {
	MovementStore *movements.MovementStorage
}

// GetShippingMovements is an endpoint for shipping movement data.
func (a MovementApi) GetShippingMovements(w http.ResponseWriter, r *http.Request) {
	currentMovements := map[string][]movements.Movement{
		"today":    a.MovementStore.TodayMovements(),
		"tomorrow": a.MovementStore.TomorrowMovements(),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// TODO: Error!
	_ = json.NewEncoder(w).Encode(currentMovements)
}
