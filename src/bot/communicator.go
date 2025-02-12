package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

var communicator Communicator

const (
	TimeFormatHourMinutes = "15:15"
)

type Communicator struct {
	s *discordgo.Session
}

func initCommunicator(session *discordgo.Session) {
	communicator = Communicator{session}
}

func (c *Communicator) WelcomeMessage() {
	c.s.ChannelMessageSend(GeneralChannelID, "Let's build a DIY smart home together!")
}

func (c *Communicator) createClockChannel() string {
	channel, err := c.s.GuildChannelCreate(GuildId, "⏰ Clock", discordgo.ChannelTypeGuildText)
	if err != nil {
		log.Println("Error creating channel:", err)
		return ""
	}
	c.s.ChannelMessageSend(GuildId, fmt.Sprintf("⏳ Clock started in <#%s>", channel.ID))
	return channel.ID
}

func (c *Communicator) sendMessage(channelId string, content string) {
	c.s.ChannelMessageSend(channelId, content)
}

func (c *Communicator) updateClock(channelId string) {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	startTime := time.Now()
	for range ticker.C {
		elapsed := time.Since(startTime)
		minutes := int(elapsed.Minutes())
		seconds := int(elapsed.Seconds()) % 60

		newName := fmt.Sprintf("⏳|%dm:%ds", minutes, seconds)
		_, err := c.s.ChannelEditComplex(channelId, &discordgo.ChannelEdit{Name: newName})
		if err != nil {
			log.Println("Clock Channel closed")
			return
		}
	}
}

func (c *Communicator) deleteClock(channelId string) string {
	channel, err := c.s.ChannelDelete(channelId)
	if err != nil {
		log.Println("Error creating channel:", err)
		return ""
	}
	return channel.ID
}
