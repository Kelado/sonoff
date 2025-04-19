package communication

import (
	"github.com/bwmarrin/discordgo"
)

var communicator Communicator

const (
	TimeFormatHourMinutes = "15:15"
)

type Communicator struct {
	s *discordgo.Session
}

func InitCommunicator(session *discordgo.Session) {
	communicator = Communicator{session}
}

func SendMessage(channelId string, content string) {
	communicator.s.ChannelMessageSend(channelId, content)
}
