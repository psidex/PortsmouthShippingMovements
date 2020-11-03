package movements

import (
	"log"
)

// UpdateMovements updates the movement lists in the given Manager whilst logging anything that happens.
func UpdateMovements(movementManager *Manager) {
	log.Println("Updating movement store")
	err := movementManager.UpdateAllMovements()
	if err != nil {
		log.Printf("Error when updating movement store: %s", err)
	}
}
