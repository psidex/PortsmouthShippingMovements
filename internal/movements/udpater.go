package movements

import (
	"log"
)

// UpdateMovements updates the movement lists in the given MovementStorage whilst logging anything that happens.
func UpdateMovements(movementStore *MovementStorage) {
	log.Println("Updating movement store")
	err := movementStore.UpdateAllMovements()
	if err != nil {
		log.Printf("Error when updating movement store: %s", err)
	}
}
