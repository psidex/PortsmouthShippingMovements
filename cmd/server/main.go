package main

import (
	"github.com/psidex/PortsmouthShippingMovements/internal/movements"
	"log"
)

func main() {
	shipMovements, err := movements.GetTomorrowMovements()
	if err != nil {
		log.Fatalf("Fatal error getting movements: %v", err)
	}
	for _, m := range shipMovements {
		if m.Type == movements.Notice {
			log.Printf("Notice: %s at %s", m.Name, m.Time)
		} else {
			log.Printf("%s is moving from %s to %s at %s\n", m.Name, m.From.Name, m.To.Name, m.Time)
		}
	}
}
