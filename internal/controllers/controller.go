package controllers

import (
	"github/Kelado/sonoff/internal/logs"
	"github/Kelado/sonoff/internal/registry"
	"github/Kelado/sonoff/internal/scheduler"
	view "github/Kelado/sonoff/view/sections"
	"log"
	"net/http"
	"time"
)

type Controller struct {
	registry *registry.DeviceRegistry
	logger   *logs.Logger
}

func NewController(registry *registry.DeviceRegistry, logger *logs.Logger) *Controller {
	return &Controller{
		registry: registry,
		logger:   logger,
	}
}

func (c *Controller) HomePage(w http.ResponseWriter, r *http.Request) {
	device := c.registry.GetFirst()
	view.Homepage(device).Render(r.Context(), w)
}

func (c *Controller) DevicePanel(w http.ResponseWriter, r *http.Request) {
	devices := c.registry.List()
	view.Devices(devices).Render(r.Context(), w)
}

func (c *Controller) TurnOnDevice(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s := c.registry.Get(id)

	s.SetOn()
	log.Printf("Turn ON device with id [%v] ", id)

	c.logger.AddNewEntry("Turn ON device")

	w.Header().Add("HX-Trigger", "newLogEntry")
	view.DeviceOpened(s).Render(r.Context(), w)
}

func (c *Controller) TurnOffDevice(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s := c.registry.Get(id)

	s.SetOff()
	log.Printf("Turn OFF device with id [%v] ", id)

	c.logger.AddNewEntry("Turn OFF device")

	w.Header().Add("HX-Trigger", "newLogEntry")
	view.DeviceClosed(s).Render(r.Context(), w)
}

func (c *Controller) GetLastLogEntry(w http.ResponseWriter, r *http.Request) {
	logEntry := c.logger.GetLastLog()
	view.LogEntry(logEntry).Render(r.Context(), w)
}

func (c *Controller) ScheduleForm(w http.ResponseWriter, r *http.Request) {
	layout := "2006-01-02T15:04"
	now := time.Now().Local().Format(layout)
	view.ScheduleForm(now).Render(r.Context(), w)
}

func (c *Controller) Schedule(w http.ResponseWriter, r *http.Request) {
	timeForm := r.FormValue("date-time")
	layout := "2006-01-02T15:04"
	loc, err := time.LoadLocation("Europe/Istanbul")
	if err != nil {
		log.Println("Error loading location:", err)
		return
	}

	parsedTime, err := time.ParseInLocation(layout, timeForm, loc)
	if err != nil {
		log.Println("Error parsing time:", err)
		return
	}

	action := r.FormValue("action")

	device := c.registry.GetFirst()

	s := scheduler.Scheduler{}
	job := scheduler.Job{
		Device: &device,
		Action: action,
		At:     parsedTime,
	}
	s.Schedule(&job)
	c.logger.AddNewEntry("Scheduled Job: turn " + action + " at: " + parsedTime.Format("2006-01-02T15:04"))
	w.Header().Add("HX-Trigger", "newLogEntry")

	view.ScheduleButton().Render(r.Context(), w)
}
