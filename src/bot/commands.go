package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

const (
	CommandStatus = "status"
	CommandAction = "action"
)

// To add new command
//  1. Declare it below
//  2. Register its handlerin commandHandlers
var commands = []*discordgo.ApplicationCommand{
	{
		Name:        CommandStatus,
		Description: "Check the status of thermostat",
	},
	{
		Name:        CommandAction,
		Description: "Open thermostat controller",
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	CommandStatus: handleCommandStatus,
	CommandAction: handleCommandAction,
}

func commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ChannelID != ThermostatChannelId {
		return
	}

	cmd := i.ApplicationCommandData().Name
	if _, ok := commandHandlers[cmd]; !ok {
		return
	}
	commandHandlers[cmd](s, i)
}

func registerCommands(s *discordgo.Session) {
	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			log.Fatalf("Cannot create command: %v", err)
		}
	}
}
