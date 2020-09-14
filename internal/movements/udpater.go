package movements

import (
	"log"
	"time"
)

// UpdateMovementsPeriodically periodically updates the 2 movement lists in the given MovementStorage.
func UpdateMovementsPeriodically(movementStore *MovementStorage, scraper MovementScraper) {
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

		time.Sleep(time.Hour * 8)
	}
}
