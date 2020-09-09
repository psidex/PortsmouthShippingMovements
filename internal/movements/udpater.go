package movements

import (
	"log"
	"time"
)

// UpdateMovementsPeriodically takes a MovementStorage and an images.ShipImageUrlStorage and updates the 2 movement lists
// in the MovementStorage every sleepDuration.
func UpdateMovementsPeriodically(movementStore *MovementStorage, sleepDuration time.Duration) {
	log.Println("Updating movementStore todayMovements")
	todayMovements, err := GetTodayMovements()
	if err != nil {
		log.Printf("Error when calling GetTodayMovements: %v", err)
	} else {
		movementStore.SetTodayMovements(todayMovements)
	}

	log.Println("Updating movementStore tomorrowMovements")
	tomorrowMovements, err := GetTomorrowMovements()
	if err != nil {
		log.Printf("Error when calling GetTomorrowMovements: %v", err)
	} else {
		movementStore.SetTomorrowMovements(tomorrowMovements)
	}

	time.Sleep(sleepDuration)
}
