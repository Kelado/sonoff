package bot

import (
	"fmt"
	"github/Kelado/sonoff/src/services"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	TurnOn  = "turn_on"
	TurnOff = "turn_off"

	ScheduleStart   = "schedule_start"
	HourSelection   = "hour_selection"
	MinuteSelection = "minute_selection"
)

var scheduleDate = struct {
	hour   string
	minute string
}{}

func (b *Bot) sendActionControlls(s *discordgo.Session, channelID string) {
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
					CustomID: TurnOn,
					Emoji:    &discordgo.ComponentEmoji{Name: "🔥"},
				},
				discordgo.Button{
					Label:    "Turn OFF",
					Style:    discordgo.DangerButton,
					CustomID: TurnOff,
					Emoji:    &discordgo.ComponentEmoji{Name: "❄️"},
				},
				discordgo.Button{
					Label:    "Schedule your bath",
					Style:    discordgo.SecondaryButton,
					CustomID: ScheduleStart,
					Emoji:    &discordgo.ComponentEmoji{Name: "📅"},
				},
			},
		},
	}

	s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embed:      embed,
		Components: components,
	})
}

func (b *Bot) actionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response string
	switch i.MessageComponentData().CustomID {
	case TurnOn:
		b.switchService.TurnOnByName(services.Thermostat)
		response = "✅ Lamp is now **ON**! 💡"
	case TurnOff:
		b.switchService.TurnOffByName(services.Thermostat)
		response = "✅ Lamp is now **OFF**! 🕶️"
	case ScheduleStart:
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("⏰ Select a time :"),
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.SelectMenu{
								CustomID:    HourSelection,
								Placeholder: "Select hour...",
								Options:     generateHourOptions(),
							},
						},
					},
				},
			},
		})
	case HourSelection:
		hour := i.MessageComponentData().Values[0]
		scheduleDate.hour = hour
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("⏰ Select a time : (hour: " + hour + ")"),
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.SelectMenu{
								CustomID:    MinuteSelection,
								Placeholder: "Select minute...",
								Options:     generateMinuteOptions(),
							},
						},
					},
				},
			},
		})
	case MinuteSelection:
		minute := i.MessageComponentData().Values[0]
		scheduleDate.minute = minute
		response = "You want to go for a bath at " + scheduleDate.hour + ":" + minute

		targetTimeStr := scheduleDate.hour + ":" + scheduleDate.minute
		location := time.Local
		currentTime := time.Now()
		targetTime, err := time.ParseInLocation("15:04", targetTimeStr, location)
		if err != nil {
			fmt.Println("Error parsing target time:", err)
			return
		}
		targetTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), targetTime.Hour(), targetTime.Minute(), 0, 0, location)
		if targetTime.Before(currentTime) {
			fmt.Errorf("WTF are you doing man")
			return
		}
		durationUntilTarget := targetTime.Sub(currentTime)
		fmt.Printf("Time until %s: %v\n", targetTime.Format("15:04"), durationUntilTarget)
		go func() {
			time.Sleep(durationUntilTarget)
			fmt.Println("Time to perform the scheduled action!")
			b.switchService.TurnOnByName(services.Thermostat)
		}()

	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

// Generate hour options (0-23)
func generateHourOptions() []discordgo.SelectMenuOption {
	var options []discordgo.SelectMenuOption
	for h := 0; h < 24; h++ {
		options = append(options, discordgo.SelectMenuOption{
			Label: fmt.Sprintf("%02d", h),
			Value: fmt.Sprintf("%02d", h),
		})
	}
	return options
}

// Generate minute options (0, 30)
func generateMinuteOptions() []discordgo.SelectMenuOption {
	var options []discordgo.SelectMenuOption
	for m := 0; m < 60; m += 10 {
		options = append(options, discordgo.SelectMenuOption{
			Label: fmt.Sprintf("%02d", m),
			Value: fmt.Sprintf("%02d", m),
		})
	}
	return options
}
