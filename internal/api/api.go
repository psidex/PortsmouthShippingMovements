package api

import (
	"encoding/json"
	"github.com/psidex/PortsmouthShippingMovements/internal/movements"
	"log"
	"net/http"
)

// MovementApi contains data and methods for serving movement related API functions.
type MovementApi struct {
	MovementStore *movements.MovementStorage
}

// GetShippingMovements is an endpoint for shipping movement data.
func (a MovementApi) GetShippingMovements(w http.ResponseWriter, r *http.Request) {
	todayMovements, tomorrowMovements := a.MovementStore.GetMovements()

	currentMovements := map[string][]movements.Movement{
		"today":    todayMovements,
		"tomorrow": tomorrowMovements,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(currentMovements)
	if err != nil {
		log.Print("GetShippingMovements: encoding API response failed: %s", err)
	}
}
