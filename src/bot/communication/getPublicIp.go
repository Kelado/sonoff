package communication

func SendPublicIp(channelId string, ip string) {
	SendMessage(channelId, "> Public IP: **"+ip+"**")
}
