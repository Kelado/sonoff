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

	bot_actions "github/Kelado/sonoff/src/bot/actions"
	"github/Kelado/sonoff/src/bot/communication"

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
	namedayService := services.NewNamedayService()

	discordBot := bot.New(switchService)
	discordBot.Start()

	// Register all cron jobs
	actionManager := services.NewActionService()
	actionManager.AddRepeatedEveryDay(&services.Action{
		Name:   "close-forgotten-thermostat",
		Hour:   1,
		Minute: 0,
		Action: func() {
			switchService.TurnOffByName(services.Thermostat)
		},
	})

	actionManager.AddRepeatedEveryDay(&services.Action{
		Name:   "clear-thermostat-channel",
		Hour:   0,
		Minute: 0,
		Action: func() {
			bot_actions.ClearChannel(discordBot.Session, bot.ThermostatChannelId)
		},
	})
	actionManager.AddRepeatedEveryDay(&services.Action{
		Name:   "post-celebrating-names-for-today",
		Hour:   0,
		Minute: 0,
		Action: func() {
			communication.SendCelebratingNamesForToday(bot.GeneralChannelID, namedayService.GetCelebratingNamesForToday())
		},
	})

	// Wait for termination
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	discordBot.Stop()
}
