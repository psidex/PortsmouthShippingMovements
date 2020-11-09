package main

import (
	"github.com/gorilla/mux"
	"github.com/psidex/PortsmouthShippingMovements/internal/api"
	"github.com/psidex/PortsmouthShippingMovements/internal/bing"
	"github.com/psidex/PortsmouthShippingMovements/internal/config"
	"github.com/psidex/PortsmouthShippingMovements/internal/images"
	"github.com/psidex/PortsmouthShippingMovements/internal/movements"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"time"
)

// check checks if the error is not nil, if it is, it calls log.Fatalf.
func check(err error) {
	if err != nil {
		log.Fatalf("Fatal Error: %s", err)
	}
}

func main() {
	c, err := config.LoadConfig()
	check(err)

	bingApiHttpClient := &http.Client{Timeout: time.Second * 10}
	webScraperHttpClient := &http.Client{Timeout: time.Second * 10}

	imageSearchApi := bing.NewImageSearchApi(bingApiHttpClient, c.BingImageSearchApiKey)
	imageUrlMan, err := images.NewUrlManager(imageSearchApi, c.ImageStoragePath)
	check(err)

	movementScraper := movements.NewScraper(webScraperHttpClient, c.ContactEmail)
	movementMan := movements.NewManager(imageUrlMan, movementScraper)

	// Load initial data.
	movements.UpdateMovements(movementMan)

	// Start a cron to run the update function using given config string.
	cr := cron.New()
	_, err = cr.AddFunc(c.UpdateCronString, func() { movements.UpdateMovements(movementMan) })
	check(err)
	cr.Start()

	apiRoute := api.NewMovementApi(movementMan)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/movements", apiRoute.GetShippingMovements)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.Println("Serving http on all available interfaces @ port 8080")
	check(http.ListenAndServe(":8080", router))
}
