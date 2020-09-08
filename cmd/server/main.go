package main

import (
	"github.com/gorilla/mux"
	"github.com/psidex/PortsmouthShippingMovements/internal/api"
	"github.com/psidex/PortsmouthShippingMovements/internal/movements"
	"log"
	"net/http"
	"time"
)

func main() {
	movementHandler := movements.NewMovementHandler()
	go movements.UpdateMovementsPeriodically(movementHandler, time.Hour*12)
	apiRoute := api.MovementApi{MovementHandler: movementHandler}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/movements", apiRoute.GetShippingMovements)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// TODO: is log.Fatal like this a bad idea?
	log.Fatal(http.ListenAndServe(":8080", router))
}
