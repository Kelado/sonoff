package bot

import (
	"fmt"
	"github/Kelado/sonoff/src/scheduler"
	"github/Kelado/sonoff/src/services"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	DefaultOnDuration = 40 * time.Minute
)

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
	log.Println("Selected hour: " + scheduleDate.hour)
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
	log.Println("Selected minutes: " + scheduleDate.minute)
	targetTimeStr := scheduleDate.hour + ":" + scheduleDate.minute

	response := InteractionResponseData{
		Content: "Thermostat will be turned 🔥ON at " + scheduleDate.hour + ":" + scheduleDate.minute,
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
	fmt.Printf("Time until %s, %v\n", targetTime.Format(time.Kitchen), durationUntilTarget)
	fmt.Printf("Thermostat will open at: %v\n", targetTime.Format(time.Kitchen))
	scheduler.JobManager.ScheduleAction(targetTime, func() {
		switchService.TurnOnByName(services.Thermostat)
		communicator.sendMessage(ThermostatChannelId, fmt.Sprintf("Thermostat just turned 🔥ON: %v", time.Now().Format(time.Kitchen)))
		channelId := communicator.createClockChannel()
		go communicator.updateClock(channelId)

		closeAt := time.Now().Add(DefaultOnDuration)
		fmt.Printf("Thermostat will close at: %v\n", closeAt.Format(time.Kitchen))
		scheduler.JobManager.ScheduleAction(closeAt, func() {
			switchService.TurnOffByName(services.Thermostat)
			communicator.sendMessage(ThermostatChannelId, fmt.Sprintf("Thermostat just turned ❄️OFF: %v", time.Now().Format(time.Kitchen)))
			communicator.deleteClock(channelId)
		})
	})

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

	// For some reason I can not start from minute < 10
	startMinute := 10
	minuteNow := time.Now().Minute()
	if minuteNow > startMinute {
		if minuteNow%2 == 0 {
			startMinute = minuteNow + 2
		} else {
			startMinute = minuteNow + 3
		}
	}

	for m := startMinute; m < 60; m += 2 {
		options = append(options, discordgo.SelectMenuOption{
			Label: fmt.Sprintf("%02d", int(m)),
			Value: fmt.Sprintf("%02d", int(m)),
		})
	}
	return options
}

// func createClockChannel(s *discordgo.Session, guildID string) string {
// 	channel, err := s.GuildChannelCreate(guildID, "Clock Channel", discordgo.ChannelTypeGuildText)
// 	if err != nil {
// 		log.Println("Error creating channel:", err)
// 		return ""
// 	}
// 	s.ChannelMessageSend(guildID, fmt.Sprintf("⏳ Clock started in <#%s>", channel.ID))
// 	return channel.ID
// }

// func updateClock(s *discordgo.Session, channelId string) {
// 	ticker := time.NewTicker(time.Minute * 5)
// 	defer ticker.Stop()

// 	startTime := time.Now()
// 	for range ticker.C {
// 		elapsed := time.Since(startTime)
// 		minutes := int(elapsed.Minutes())
// 		seconds := int(elapsed.Seconds()) % 60

// 		newName := fmt.Sprintf("⏳ %dm %ds", minutes, seconds)
// 		_, err := s.ChannelEditComplex(channelId, &discordgo.ChannelEdit{Name: newName})
// 		if err != nil {
// 			log.Println("Error updating channel name:", err)
// 		}
// 	}
// }

// func deleteClock(s *discordgo.Session, channelId string) string {
// 	channel, err := s.ChannelDelete(channelId)
// 	if err != nil {
// 		log.Println("Error creating channel:", err)
// 		return ""
// 	}
// 	return channel.ID
// }
