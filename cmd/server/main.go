package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/psidex/PortsmouthShippingMovements/internal/api"
	"github.com/psidex/PortsmouthShippingMovements/internal/config"
	"github.com/psidex/PortsmouthShippingMovements/internal/images"
	"github.com/psidex/PortsmouthShippingMovements/internal/movements"
	"log"
	"net/http"
	"os"
	"time"
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

	imageStore, err := images.NewShipImageUrlStorage(c.BingImageSearchApiKey, c.ImageStoragePath)
	if err != nil {
		log.Fatal(err)
	}

	movementStore := movements.NewMovementStorage(imageStore)
	go movements.UpdateMovementsPeriodically(movementStore, time.Hour*8)

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
