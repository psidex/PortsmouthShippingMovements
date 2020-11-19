package main

import (
	"github.com/gorilla/mux"
	"github.com/psidex/PortsmouthShippingMovements/internal/config"
	"github.com/psidex/PortsmouthShippingMovements/internal/qhm"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
)

// check checks if the error is not nil, if it is, it calls log.Fatalf.
func check(err error) {
	if err != nil {
		log.Fatalf("Fatal Error: %s", err)
	}
}

func main() {
	c, err := config.Get()
	check(err)

	// Create movement manager from config and load initial data.
	movementMan := qhm.NewMovementManager(c)
	movementMan.Update()

	// Start a cron to run the update function using given config string.
	cr := cron.New()
	_, err = cr.AddFunc(c.UpdateCronString, func() { movementMan.Update() })
	check(err)
	cr.Start()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/movements", movementMan.SendCurrentMovements)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.Println("Serving http on all available interfaces @ port 8080")
	check(http.ListenAndServe(":8080", router))
}
