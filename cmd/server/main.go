package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/psidex/PortsmouthShippingMovements/internal/api"
	"github.com/psidex/PortsmouthShippingMovements/internal/bing"
	"github.com/psidex/PortsmouthShippingMovements/internal/config"
	"github.com/psidex/PortsmouthShippingMovements/internal/images"
	"github.com/psidex/PortsmouthShippingMovements/internal/movements"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
	"time"
)

// check checks if the error is not nil, if it is, log and exit with 1.
func check(err error) {
	if err != nil {
		log.Fatalf("Error during initialization: %s", err)
	}
}

func main() {
	// Read and parse config.
	c, err := config.LoadConfig()
	check(err)

	// Set up webserver access log file.
	accessLogFile, err := os.OpenFile(c.AccessLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	check(err)
	defer accessLogFile.Close()

	// Custom http client for all web requests (Bing API and scraper for QHM page).
	// The API and scraper can share a client as the API will only be requested after a scraping has happened, so there
	// won't be any blocking.
	httpClient := &http.Client{Timeout: time.Second * 10}

	// Set up image storage, web scraper and then movement storage.
	imageSearchApi := bing.NewImageSearchApi(httpClient, c.BingImageSearchApiKey)
	imageUrlStore, err := images.NewShipImageUrlStorage(imageSearchApi, c.ImageStoragePath)
	check(err)

	scraper := movements.NewMovementScraper(httpClient, c.ContactEmail)
	movementStore := movements.NewMovementStorage(imageUrlStore, scraper)

	// Load initial data.
	movements.UpdateMovements(movementStore)

	// Start a cron to run the update function at midnight, 8am, and 4pm.
	cr := cron.New()
	_, err = cr.AddFunc(c.UpdateCronString, func() { movements.UpdateMovements(movementStore) })
	check(err)
	cr.Start()

	// Set up all the web server stuff.
	apiRoute := api.MovementApi{MovementStore: movementStore}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/movements", apiRoute.GetShippingMovements)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	loggedRouter := handlers.LoggingHandler(accessLogFile, router)

	// Start serving.
	log.Println("Listening on http://127.0.0.1:8080")
	err = http.ListenAndServe(":8080", loggedRouter)
	check(err)
}
