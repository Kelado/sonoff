package services

import (
	"log"
	"time"
)

type Action struct {
	Name   string
	Hour   int
	Minute int
	Action func()
}

type ActionService struct{}

func NewActionService() *ActionService {
	return &ActionService{}
}

func (s *ActionService) AddRepeatedEveryDay(a *Action) {
	go func() {
		for {
			now := time.Now()
			nextRun := time.Date(now.Year(), now.Month(), now.Day(), a.Hour, a.Minute, now.Second(), 0, now.Location())

			// If nextRun is in the past, schedule for tomorrow
			if nextRun.Before(now) {
				nextRun = nextRun.Add(24 * time.Hour)
			}

			log.Println("Action " + a.Name + " will run again at " + nextRun.String())
			timeUntilRun := time.Until(nextRun)
			log.Println("Action " + a.Name + " will wait for " + timeUntilRun.String())
			time.Sleep(timeUntilRun)

			// Execute the action
			log.Println("Executing action " + a.Name)
			a.Action()

			// Wait for the triggerAt time has passed
			time.Sleep(time.Minute)
		}
	}()

}
