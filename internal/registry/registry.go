package registry

import (
	"github/Kelado/sonoff/internal/basicr3"
)

var INCR_INT = 0

type DeviceRegistry struct {
	devices map[string]basicr3.Switch
}

func NewDeviceRegistry() *DeviceRegistry {
	return &DeviceRegistry{
		devices: make(map[string]basicr3.Switch),
	}
}

func (dr *DeviceRegistry) Register(device *basicr3.Switch) {
	dr.devices[device.ID] = *device
}

func (dr *DeviceRegistry) Get(id string) basicr3.Switch {
	return dr.devices[id]
}

func (dr *DeviceRegistry) GetFirst() basicr3.Switch {
	for _, v := range dr.devices {
		return v
	}
	return basicr3.Switch{}
}

func (dr *DeviceRegistry) List() []basicr3.Switch {
	list := make([]basicr3.Switch, 0, len(dr.devices))
	for _, device := range dr.devices {
		list = append(list, device)
	}
	return list
}
