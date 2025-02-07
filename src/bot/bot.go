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

	b.session.AddHandler(baseInteractionDispacher)
	registerCommands(b.session)

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
