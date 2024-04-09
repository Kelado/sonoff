package main

import (
	"log"
	"net/http"

	"github/Kelado/sonoff/internal/basicr3"
	"github/Kelado/sonoff/internal/controllers"
	"github/Kelado/sonoff/internal/logs"
	"github/Kelado/sonoff/internal/registry"
)

// Commands to find device's network information
// dns-sd -B _ewelink._tcp
// ping HOSTNAME.local

const (
	ID   = "" // Device ID from mDNS
	IP   = "" // Device IP
	PORT = "" // Device port
)

func main() {
	device := basicr3.NewSwitch(ID, "Bedroom", IP, PORT)
	log.Println(device)

	registry := registry.NewDeviceRegistry()
	registry.Register(device)

	logger := &logs.Logger{}

	controller := controllers.NewController(registry, logger)

	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServer(http.Dir(".")))

	mux.HandleFunc("GET /", controller.HomePage)

	mux.HandleFunc("POST /switch/{id}/on", controller.TurnOnDevice)
	mux.HandleFunc("POST /switch/{id}/off", controller.TurnOffDevice)

	mux.HandleFunc("GET /logs/last", controller.GetLastLogEntry)

	mux.HandleFunc("GET /schedule/form", controller.ScheduleForm)
	mux.HandleFunc("POST /schedule", controller.Schedule)

	s := &http.Server{
		Addr:    ":5001",
		Handler: mux,
	}

	log.Println("Server started listening at Port" + s.Addr)
	s.ListenAndServe()

}
