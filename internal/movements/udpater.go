package movements

import (
	"log"
	"time"
)

// UpdateMovementsPeriodically updates the 2 movement lists in the MovementStorage every sleepDuration.
func UpdateMovementsPeriodically(movementStore *MovementStorage, scraper MovementScraper, sleepDuration time.Duration) {
	for {
		log.Println("Updating movementStore todayMovements")
		todayMovements, err := scraper.GetTodayMovements()
		if err != nil {
			log.Printf("Error when calling GetTodayMovements: %v", err)
		} else {
			movementStore.SetTodayMovements(todayMovements)
		}

		log.Println("Updating movementStore tomorrowMovements")
		tomorrowMovements, err := scraper.GetTomorrowMovements()
		if err != nil {
			log.Printf("Error when calling GetTomorrowMovements: %v", err)
		} else {
			movementStore.SetTomorrowMovements(tomorrowMovements)
		}

		time.Sleep(sleepDuration)
	}
}
