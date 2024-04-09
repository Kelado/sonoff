package scheduler

import (
	"github/Kelado/sonoff/internal/basicr3"
	"log"
	"time"
)

type Scheduler struct {
}

type Job struct {
	Device *basicr3.Switch
	Action string
	At     time.Time
}

func (j *Job) Run() {
	go func() {
		log.Println("Job: started")
		timeUntilRun := j.At.Sub(time.Now())
		if timeUntilRun > 0 {
			time.Sleep(timeUntilRun)
		}

		log.Println("Job: execute action")
		switch j.Action {
		case "on":
			j.Device.SetOn()
			break
		case "off":
			j.Device.SetOff()
			break
		default:
			log.Println("Uknown action ", j.Action)
		}
	}()
}

func (s *Scheduler) Schedule(job *Job) {
	job.Run()
}
