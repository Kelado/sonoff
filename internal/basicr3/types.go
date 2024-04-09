package basicr3

type InfoResp struct {
	Seq   int          `json:"seq"`
	Error int          `json:"error"`
	Data  InfoDataResp `json:"data"`
}

type InfoDataResp struct {
	Switch         string `json:"switch"`
	Startup        string `json:"startup"`
	Pulse          string `json:"pulse"`
	PulseWidth     int    `json:"pulseWidth"`
	SSID           string `json:"ssid"`
	OtaUnlock      bool   `json:"otaUnlock"`
	FwVersion      string `json:"fwVersion"`
	Deviceid       string `json:"deviceid"`
	BSSID          string `json:"bssid"`
	SignalStrength int    `json:"signalStrength"`
}
