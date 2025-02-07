package bot

import (
	"fmt"
	services "github/Kelado/sonoff/src/services"
	"time"

	"github.com/bwmarrin/discordgo"
)

const ()

// Interactions
const (
	TH_TurnOn          = "th_turn_on"
	TH_TurnOff         = "th_turn_off"
	TH_Schedule        = "th_schedule"
	TH_HourSelection   = "th_hour_selection"
	TH_MinuteSelection = "th_minute_selection"
)

type InteractionResponseData discordgo.InteractionResponseData
type ActionHandler func(*discordgo.Session, *discordgo.InteractionCreate)

var th_handler = map[string]ActionHandler{
	TH_TurnOn:          turnOn,
	TH_TurnOff:         turnOff,
	TH_Schedule:        schedule,
	TH_HourSelection:   handleHourSelection,
	TH_MinuteSelection: handleMinutesSelection,
}

var scheduleDate = struct {
	hour   string
	minute string
}{}

func handleCommandAction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "🌡️ Thermostat Controller",
		Description: "It was about time to take a bath 🦨💨👃🤮🤢",
		Color:       0x00ff00,
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Turn ON",
					Style:    discordgo.PrimaryButton,
					CustomID: TH_TurnOn,
					Emoji:    &discordgo.ComponentEmoji{Name: "🔥"},
				},
				discordgo.Button{
					Label:    "Turn OFF",
					Style:    discordgo.DangerButton,
					CustomID: TH_TurnOff,
					Emoji:    &discordgo.ComponentEmoji{Name: "❄️"},
				},
				discordgo.Button{
					Label:    "Schedule your bath",
					Style:    discordgo.SecondaryButton,
					CustomID: TH_Schedule,
					Emoji:    &discordgo.ComponentEmoji{Name: "📅"},
				},
			},
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})

}

func actionInteractionDispatcher(s *discordgo.Session, i *discordgo.InteractionCreate) {
	t := i.MessageComponentData().CustomID
	if _, ok := th_handler[t]; !ok {
		respond(InteractionResponseData{Content: "Careful there little one!"}, s, i)
	}
	th_handler[t](s, i)
}

func turnOff(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switchService.TurnOffByName(services.Thermostat)
	response := InteractionResponseData{
		Content: "❄️ Thermostat is now **OFF**!",
	}
	respond(response, s, i)
}

func turnOn(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switchService.TurnOnByName(services.Thermostat)
	response := InteractionResponseData{
		Content: "🔥 Thermostat is now **ON**!",
	}
	respond(response, s, i)
}

func schedule(s *discordgo.Session, i *discordgo.InteractionCreate) {
	response := InteractionResponseData{
		Content: fmt.Sprintf("⏰ Select a time :"),
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    TH_HourSelection,
						Placeholder: "Select hour...",
						Options:     generateHourOptions(),
					},
				},
			},
		},
	}
	respond(response, s, i)
}

func handleHourSelection(s *discordgo.Session, i *discordgo.InteractionCreate) {
	scheduleDate.hour = i.MessageComponentData().Values[0]
	response := InteractionResponseData{
		Content: fmt.Sprintf("⏰ Select a time:"),
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    TH_MinuteSelection,
						Placeholder: "Select minute...",
						Options:     generateMinuteOptions(),
					},
				},
			},
		},
	}
	respond(response, s, i)
}

func handleMinutesSelection(s *discordgo.Session, i *discordgo.InteractionCreate) {
	scheduleDate.minute = i.MessageComponentData().Values[0]
	targetTimeStr := scheduleDate.hour + ":" + scheduleDate.minute

	response := InteractionResponseData{
		Content: "You want to go for a bath at " + scheduleDate.hour + ":" + scheduleDate.minute,
	}
	location := time.Local
	currentTime := time.Now()
	targetTime, err := time.ParseInLocation("15:04", targetTimeStr, location)
	if err != nil {
		fmt.Println("Error parsing target time:", err)
	}
	targetTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), targetTime.Hour(), targetTime.Minute(), 0, 0, location)
	if targetTime.Before(currentTime) {
		fmt.Errorf("WTF are you doing man")
	}
	durationUntilTarget := targetTime.Sub(currentTime)
	fmt.Printf("Time until %s: %v\n", targetTime.Format("15:04"), durationUntilTarget)
	go func() {
		time.Sleep(durationUntilTarget)
		fmt.Println("Time to perform the scheduled action!")
		switchService.TurnOnByName(services.Thermostat)
	}()

	respond(response, s, i)
}

func respond(response InteractionResponseData, s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: (*discordgo.InteractionResponseData)(&response),
	})
}

func generateHourOptions() []discordgo.SelectMenuOption {
	var options []discordgo.SelectMenuOption
	currentHour := time.Now().Hour()
	for h := currentHour; h < 24; h++ {
		options = append(options, discordgo.SelectMenuOption{
			Label: fmt.Sprintf("%02d", h),
			Value: fmt.Sprintf("%02d", h),
		})
	}
	return options
}

func generateMinuteOptions() []discordgo.SelectMenuOption {
	var options []discordgo.SelectMenuOption
	currentMinute := time.Now().Minute()

	startMinute := ((currentMinute / 10) + 1) * 10
	if startMinute == 60 {
		startMinute = 0
	}

	for m := startMinute; m < 60; m += 10 {
		options = append(options, discordgo.SelectMenuOption{
			Label: fmt.Sprintf("%02d", m),
			Value: fmt.Sprintf("%02d", m),
		})
	}
	return options
}
