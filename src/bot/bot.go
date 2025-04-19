package bot

import (
	"fmt"
	"github/Kelado/sonoff/src/bot/communication"
	"github/Kelado/sonoff/src/services"
	"os"

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
	Session *discordgo.Session
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

	b.Session, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Failed initializing Discord Session:", err)
		return
	}

	u, err := b.Session.User("@me")
	if err != nil {
		fmt.Println("Failed getting current User:", err)
		return
	}

	b.id = u.ID
	b.Session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	err = b.Session.Open()
	if err != nil {
		fmt.Println("Failed opening connection to Discord:", err)
		return
	}

	b.Session.AddHandler(baseInteractionDispacher)
	registerCommands(b.Session)

	communication.InitCommunicator(b.Session)
	communication.WelcomeMessage(GeneralChannelID)

	fmt.Println("Bot is now connected!")
}

func (b *Bot) Stop() {
	b.Session.Close()
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
