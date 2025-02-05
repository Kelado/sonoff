package services

import (
	"github/Kelado/sonoff/src/basicr3"
	"github/Kelado/sonoff/src/registry"
)

type SwitchType string

const (
	Thermostat SwitchType = "THERMOSTAT"
)

type SwitchService struct {
	registry *registry.DeviceRegistry
}

func NewSwitchService(registry *registry.DeviceRegistry) *SwitchService {
	return &SwitchService{registry: registry}
}

func (s *SwitchService) GetByName(name SwitchType) basicr3.Switch {
	return s.registry.GetByName(string(name))
}

func (s *SwitchService) IsActiveByName(name SwitchType) bool {
	device := s.registry.GetByName(string(name))
	return device.State.State == "on"
}

func (s *SwitchService) TurnOffByName(name SwitchType) {
	device := s.registry.GetByName(string(name))
	device.SetOff()
}

func (s *SwitchService) TurnOnByName(name SwitchType) {
	device := s.registry.GetByName(string(name))
	device.SetOn()
}
