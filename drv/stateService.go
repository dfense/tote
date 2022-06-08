package drv

type NetworkState int64
type SensorFansBulbState int64

const (
	BootState NetworkState = iota
	AcquiringNetwork
	HotSpotEnabled
	NetworkAware
	DisconnectedW_Hotspot
	DisconnectedNoNetwork
	ShuttingDown
)

const (
	FansRunning SensorFansBulbState = iota
	LightOn
	FansAndLightRunning
	SensorReadTransmitting
)

func (s NetworkState) String() string {
	switch s {
	case BootState:
		return "Booting"
	case AcquiringNetwork:
		return "AcquringNetwork"
	case HotSpotEnabled:
		return "HotSpotEnabled"
	case NetworkAware:
		return "NetworkAware"
	case ShuttingDown:
		return "ShuttingDown"
	}
	return "unknown"
}

func (s SensorFansBulbState) String() string {
	switch s {
	case FansRunning:
		return "FansRunning"
	case LightOn:
		return "LightOn"
	case FansAndLightRunning:
		return "FansAndLightRunning"
	case SensorReadTransmitting:
		return "SensorReadTransmitting"
	}
	return "unknown"
}

func LedChange(ledNumber int, LedMessage string) {
	//resp, err := http.Get("http://webcode.me")
}

func FansOn() {

}

func LightTurnedOn() {

}

func setState() {

}
