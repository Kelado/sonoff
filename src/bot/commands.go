package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

const (
	CommandStatus   = "status"
	CommandAction   = "action"
	CommandTest     = "test"
	CommandPublicIP = "ip"
)

var classTrackerCommands = []*discordgo.ApplicationCommand{
	{
		Name:        CommandPublicIP,
		Description: "Get the public IP address of the server",
	},
}

var thermostatCommands = []*discordgo.ApplicationCommand{
	{
		Name:        CommandStatus,
		Description: "Check the status of thermostat",
	},
	{
		Name:        CommandAction,
		Description: "Open thermostat controller",
	},
	{
		Name:        CommandTest,
		Description: "This is a test action command",
	},
}
var commands = append(classTrackerCommands, thermostatCommands...)

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	CommandStatus:   handleCommandStatus,
	CommandAction:   handleCommandAction,
	CommandTest:     handleCommandTest,
	CommandPublicIP: handleCommandGetPublicIp,
}

func commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var actionPermissions = map[string][]*discordgo.ApplicationCommand{
		ThermostatChannelId:   thermostatCommands,
		ClassTrackerChannelID: classTrackerCommands,
	}

	availableCommands, ok := actionPermissions[i.ChannelID]
	if !ok {
		return
	}

	cmd := i.ApplicationCommandData().Name
	if _, ok := commandHandlers[cmd]; !ok {
		return
	}

	// Check if the command is available in this channel
	commandAllowed := false
	for _, availableCmd := range availableCommands {
		if availableCmd.Name == cmd {
			commandAllowed = true
			break
		}
	}
	if !commandAllowed {
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
