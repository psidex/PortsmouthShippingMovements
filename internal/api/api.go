package api

import (
	"encoding/json"
	"github.com/psidex/PortsmouthShippingMovements/internal/movements"
	"net/http"
)

// MovementApi contains data and methods for serving movement related API functions.
type MovementApi struct {
	MovementHandler *movements.MovementHandler
}

// GetShippingMovements is an endpoint for shipping movement data.
func (a MovementApi) GetShippingMovements(w http.ResponseWriter, r *http.Request) {
	currentMovements := map[string][]movements.Movement{
		"today":    a.MovementHandler.TodayMovements(),
		"tomorrow": a.MovementHandler.TomorrowMovements(),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// TODO: Handle err!
	_ = json.NewEncoder(w).Encode(currentMovements)
}
