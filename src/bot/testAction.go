package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func handleCommandTest(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println("test action")
	// createClockChannel(s, i.GuildID)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource, // Immediate response
		Data: &discordgo.InteractionResponseData{
			Content: "time counting", // The actual message
		}})
}

// func createClockChannel(s *discordgo.Session, guildID string) string {
// 	channel, err := s.GuildChannelCreate(guildID, "Clock Channel", discordgo.ChannelTypeGuildText)
// 	if err != nil {
// 		log.Println("Error creating channel:", err)
// 		return ""
// 	}
// 	s.ChannelMessageSend(guildID, fmt.Sprintf("⏳ Clock started in <#%s>", channel.ID))
// 	return channel.ID
// }

// func updateClock(s *discordgo.Session, channelId string) {
// 	ticker := time.NewTicker(time.Minute * 5)
// 	defer ticker.Stop()

// 	startTime := time.Now()
// 	for range ticker.C {
// 		elapsed := time.Since(startTime)
// 		minutes := int(elapsed.Minutes())
// 		seconds := int(elapsed.Seconds()) % 60

// 		newName := fmt.Sprintf("⏳ %dm %ds", minutes, seconds)
// 		_, err := s.ChannelEditComplex(channelId, &discordgo.ChannelEdit{Name: newName})
// 		if err != nil {
// 			log.Println("Error updating channel name:", err)
// 		}
// 	}
// }

// func deleteClock(s *discordgo.Session, channelId string) string {
// 	channel, err := s.ChannelDelete(channelId)
// 	if err != nil {
// 		log.Println("Error creating channel:", err)
// 		return ""
// 	}
// 	return channel.ID
// }
