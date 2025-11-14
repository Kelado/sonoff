package bot

import (
	"fmt"
	"github/Kelado/sonoff/src/bot/communication"
	"log"

	"github.com/bwmarrin/discordgo"
)

func handleCommandGetPublicIp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ip, err := ipService.GetPublicIp()
	if err != nil {
		log.Println("Failed to get public IP:", err)
		return
	}
	fmt.Println("Public IP:", ip.String())
	communication.SendPublicIp(i.ChannelID, ip.String())
}
