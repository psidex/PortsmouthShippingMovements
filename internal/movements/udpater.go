package movements

import (
	"time"
)

// UpdateMovementsPeriodically takes a MovementHandler and updates the 2 movement lists every sleepDuration.
func UpdateMovementsPeriodically(handler *MovementHandler, sleepDuration time.Duration) {
	// TODO: Error handling.
	todayMovements, _ := GetTodayMovements()
	handler.SetTodayMovements(todayMovements)

	tomorrowMovements, _ := GetTomorrowMovements()
	handler.SetTomorrowMovements(tomorrowMovements)

	time.Sleep(sleepDuration)
}
