package bot_actions

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func ClearChannel(session *discordgo.Session, channelId string) {
	for {
		messages, err := session.ChannelMessages(channelId, 100, "", "", "")
		if err != nil {
			log.Fatalf("Error fetching messages: %v", err)
		}
		if len(messages) == 0 {
			fmt.Println("Channel is clean!")
			break
		}

		// Bulk delete if possible
		if len(messages) > 1 {
			var messageIDs []string
			for _, msg := range messages {
				messageIDs = append(messageIDs, msg.ID)
			}
			err = session.ChannelMessagesBulkDelete(channelId, messageIDs)
			if err != nil {
				log.Printf("Bulk delete failed, deleting individually: %v", err)
				// If bulk delete fails, delete one by one
				for _, msg := range messages {
					session.ChannelMessageDelete(channelId, msg.ID)
					time.Sleep(500 * time.Millisecond) // Avoid rate limits
				}
			}
		} else {
			// Delete single message
			err := session.ChannelMessageDelete(channelId, messages[0].ID)
			if err != nil {
				log.Printf("Failed to delete message: %v", err)
			}
		}

		time.Sleep(200 * time.Millisecond) // Rate limit handling
	}
}
