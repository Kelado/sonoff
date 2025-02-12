package basicr3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Switch struct {
	ID    string
	Name  string
	Ip    string
	Port  string
	State State
}

type State struct {
	State      string
	Pulse      string
	PulseWidth int
}

func NewSwitch(id, name, ip, port string) *Switch {
	s := Switch{
		ID:   id,
		Name: name,
		Ip:   ip,
		Port: port,
	}

	// s.SetPulseOff()

	// s.Sync()

	return &s
}

func (s *Switch) SetOn() {
	url := s.getBaseUrl() + "/zeroconf/switch"

	data := map[string]interface{}{
		"deviceid": s.ID,
		"data": map[string]interface{}{
			"switch": "on",
		},
	}

	sendCommand(url, data)
}

func (s *Switch) SetOff() {
	url := s.getBaseUrl() + "/zeroconf/switch"

	data := map[string]interface{}{
		"deviceid": s.ID,
		"data": map[string]interface{}{
			"switch": "off",
		},
	}

	sendCommand(url, data)
}

func (s *Switch) SetPulseOn(pulseWidth int) {
	url := s.getBaseUrl() + "/zeroconf/pulse"

	data := map[string]interface{}{
		"deviceid": s.ID,
		"data": map[string]interface{}{
			"pulse":      "on",
			"pulseWidth": pulseWidth,
		},
	}

	sendCommand(url, data)
}

func (s *Switch) SetPulseOff() {
	url := s.getBaseUrl() + "/zeroconf/pulse"

	data := map[string]interface{}{
		"deviceid": s.ID,
		"data": map[string]interface{}{
			"pulse":      "off",
			"pulseWidth": 0,
		},
	}

	sendCommand(url, data)
}

func (s *Switch) GetInfo() InfoResp {
	url := s.getBaseUrl() + "/zeroconf/info"

	data := map[string]interface{}{
		"deviceid": s.ID,
		"data":     map[string]interface{}{},
	}

	respEncoded, _ := sendCommand(url, data)

	var resp InfoResp
	err := json.Unmarshal(respEncoded, &resp)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}

func (s *Switch) Sync() {
	resp := s.GetInfo()

	state := State{
		State:      resp.Data.Switch,
		Pulse:      resp.Data.Pulse,
		PulseWidth: resp.Data.PulseWidth,
	}
	s.State = state
}

func (s *Switch) getBaseUrl() string {
	return "http://" + s.Ip + ":" + s.Port
}

func (s *Switch) String() string {
	return fmt.Sprintf("Device: %s, Name: %s, State: %s, Pulse: %s",
		s.ID, s.Name, s.State.State, s.State.Pulse)
}

func (s *Switch) VerboseString() string {
	url := s.getBaseUrl() + "/zeroconf/info"

	data := map[string]interface{}{
		"deviceid": s.ID,
		"data":     map[string]interface{}{},
	}

	respEncoded, _ := sendCommand(url, data)
	return string(respEncoded)
}

func encode(data map[string]interface{}) *bytes.Buffer {
	jsonValue, _ := json.Marshal(data)
	return bytes.NewBuffer(jsonValue)
}

func sendCommand(url string, data map[string]interface{}) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, encode(data))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	return body, nil
}
