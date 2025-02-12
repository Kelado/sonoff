package bot

import (
	"fmt"
	"github/Kelado/sonoff/src/services"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

var config Config

var (
	GuildId             string
	ThermostatChannelId string
	GeneralChannelID    string
)

// Services
var (
	switchService *services.SwitchService
)

type Bot struct {
	id      string
	session *discordgo.Session
}

func New(switchSvc *services.SwitchService) *Bot {
	switchService = switchSvc
	return &Bot{}
}

func (b *Bot) Start() {
	GuildId = os.Getenv("GUILD_ID")
	ThermostatChannelId = os.Getenv("THERMOSTAT_CHANNEL_ID")
	GeneralChannelID = os.Getenv("GENERAL_CHANNEL_ID")

	config, err := ReadConfig()
	if err != nil {
		fmt.Println("Failed reading configuration:", err)
		return
	}

	b.session, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Failed initializing Discord Session:", err)
		return
	}

	u, err := b.session.User("@me")
	if err != nil {
		fmt.Println("Failed getting current User:", err)
		return
	}

	b.id = u.ID
	b.session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	err = b.session.Open()
	if err != nil {
		fmt.Println("Failed opening connection to Discord:", err)
		return
	}

	b.clearChannel(ThermostatChannelId)

	b.session.AddHandler(baseInteractionDispacher)
	registerCommands(b.session)

	initCommunicator(b.session)
	// communicator.WelcomeMessage()

	fmt.Println("Bot is now connected!")
}

func (b *Bot) Stop() {
	b.session.Close()
	fmt.Println("Bot is now stopped!")
}

// This is gonna be a little messy, but it leaves everything else to be clean and modular
func baseInteractionDispacher(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		commandHandler(s, i)
	case discordgo.InteractionMessageComponent:
		actionInteractionDispatcher(s, i)
	}
}

func (b *Bot) clearChannel(channelId string) {
	for {
		messages, err := b.session.ChannelMessages(channelId, 100, "", "", "")
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
			err = b.session.ChannelMessagesBulkDelete(channelId, messageIDs)
			if err != nil {
				log.Printf("Bulk delete failed, deleting individually: %v", err)
				// If bulk delete fails, delete one by one
				for _, msg := range messages {
					b.session.ChannelMessageDelete(channelId, msg.ID)
					time.Sleep(500 * time.Millisecond) // Avoid rate limits
				}
			}
		} else {
			// Delete single message
			err := b.session.ChannelMessageDelete(channelId, messages[0].ID)
			if err != nil {
				log.Printf("Failed to delete message: %v", err)
			}
		}

		time.Sleep(200 * time.Millisecond) // Rate limit handling
	}
}
