package bot

import (
	"fmt"
	"github/Kelado/sonoff/src/services"
	"os"

	"github.com/bwmarrin/discordgo"
)

var config Config

var (
	ThermostatChannelId string
	GeneralChannelID    string
)

type Bot struct {
	id      string
	session *discordgo.Session

	switchService *services.SwitchService
}

func New(switchService *services.SwitchService) *Bot {
	return &Bot{
		switchService: switchService,
	}
}

func (b *Bot) Start() {
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

	b.session.AddHandler(b.actionHandler)
	b.session.AddHandler(b.messageHandler)

	err = b.session.Open()
	if err != nil {
		fmt.Println("Failed opening connection to Discord:", err)
		return
	}
	fmt.Println("Bot is now connected!")
}

func (b *Bot) Stop() {
	b.session.Close()
	fmt.Println("Bot is now stopped!")
}

func (b *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == b.id {
		return
	}

	if m.ChannelID == ThermostatChannelId {
		cmd := m.Content
		switch cmd {
		case "status":
			isActive := b.switchService.IsActiveByName(services.Thermostat)
			if isActive {
				s.ChannelMessageSend(ThermostatChannelId, "thermostat is ON right now")
			} else {
				s.ChannelMessageSend(ThermostatChannelId, "thermostat is OFF right now")
			}
		default:
			b.sendActionControlls(s, m.ChannelID)
		}
	}
}
