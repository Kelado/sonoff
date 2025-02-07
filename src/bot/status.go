package bot

import (
	services "github/Kelado/sonoff/src/services"

	"github.com/bwmarrin/discordgo"
)

func handleCommandStatus(s *discordgo.Session, i *discordgo.InteractionCreate) {
	isActive := switchService.IsActiveByName(services.Thermostat)
	if isActive {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource, // Immediate response
			Data: &discordgo.InteractionResponseData{
				Content: "thermostat is ON right now", // The actual message
			}})
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource, // Immediate response
			Data: &discordgo.InteractionResponseData{
				Content: "thermostat is OFF right now", // The actual message
			}})
	}
}
