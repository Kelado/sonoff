package communication

import "strings"

func SendCelebratingNamesForToday(channelId string, names []string) {
	SendMessage(channelId, "Celebrating names for today: "+strings.Join(names, ", ")+" ")
}
