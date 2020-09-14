package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/psidex/PortsmouthShippingMovements/internal/api"
	"github.com/psidex/PortsmouthShippingMovements/internal/bing"
	"github.com/psidex/PortsmouthShippingMovements/internal/config"
	"github.com/psidex/PortsmouthShippingMovements/internal/images"
	"github.com/psidex/PortsmouthShippingMovements/internal/movements"
	"log"
	"net/http"
	"os"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	accessLogFile, err := os.OpenFile(c.AccessLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer accessLogFile.Close()

	imageSearchApi := bing.NewImageSearchApi(c.BingImageSearchApiKey)
	imageStore, err := images.NewShipImageUrlStorage(imageSearchApi, c.ImageStoragePath)
	if err != nil {
		log.Fatal(err)
	}

	scraper := movements.NewMovementScraper(c.ContactEmail)
	movementStore := movements.NewMovementStorage(imageStore)
	go movements.UpdateMovementsPeriodically(movementStore, scraper)

	apiRoute := api.MovementApi{MovementStore: movementStore}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/movements", apiRoute.GetShippingMovements)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	loggedRouter := handlers.LoggingHandler(accessLogFile, router)
	err = http.ListenAndServe(":8080", loggedRouter)
	if err != nil {
		log.Fatal(err)
	}
}
