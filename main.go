package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github/Kelado/sonoff/src/basicr3"
	"github/Kelado/sonoff/src/bot"
	"github/Kelado/sonoff/src/registry"
	"github/Kelado/sonoff/src/services"

	"github.com/joho/godotenv"
)

// Commands to find device's network information
// dns-sd -B _ewelink._tcp
// ping HOSTNAME.local
var (
	ID   string // Device ID from mDNS
	IP   string // Device IP
	PORT string // Device port
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ID = os.Getenv("THERMOSTAT_SWITCH_ID")     // Device ID from mDNS
	IP = os.Getenv("THERMOSTAT_SWITCH_IP")     // Device IP
	PORT = os.Getenv("THERMOSTAT_SWITCH_PORT") // Device port

	device := basicr3.NewSwitch(ID, string(services.Thermostat), IP, PORT)
	log.Println(device)

	registry := registry.NewDeviceRegistry()
	registry.Register(device)

	switchService := services.NewSwitchService(registry)
	bot := bot.New(switchService)
	bot.Start()

	// Wait for termination
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	bot.Stop()
}
