package communication

func WelcomeMessage(channelId string) {
	communicator.s.ChannelMessageSend(channelId, "Let's build a DIY smart home together!")
}
