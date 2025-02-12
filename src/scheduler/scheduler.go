package scheduler

import (
	"github/Kelado/sonoff/src/basicr3"
	"log"
	"time"
)

var JobManager Scheduler

type Scheduler struct{}

type Job struct {
	Device *basicr3.Switch
	Action string
	At     time.Time
}

func (j *Job) run() {
	go func() {
		log.Println("Job: started")

		timeUntilRun := time.Until(j.At)
		if timeUntilRun > 0 {
			time.Sleep(timeUntilRun)
		}

		log.Println("Job: execute action")
		switch j.Action {
		case "on":
			j.Device.SetOn()
		case "off":
			j.Device.SetOff()
		default:
			log.Println("Uknown action ", j.Action)
		}
	}()
}

func (s *Scheduler) Schedule(job *Job) {
	job.run()
}

func (s *Scheduler) ScheduleAction(at time.Time, action func()) {
	go func() {
		timeUntilRun := time.Until(at)
		if timeUntilRun > 0 {
			time.Sleep(timeUntilRun)
		}

		action()
	}()
}

