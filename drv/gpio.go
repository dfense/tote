package drv

import (
	"fmt"

	"github.com/stianeikeland/go-rpio/v4"
)

type gpio struct {
	fan   rpio.Pin
	light rpio.Pin
}

var (
	rpiGpio gpio
	cycle   uint32 = 120
)

// must call to initialize pins
func GpioOpen() error {
	err := rpio.Open()
	if err != nil {
		return err
	}

	//set fan pwm and zero speed
	rpiGpio.fan = rpio.Pin(19)
	rpiGpio.fan.Mode(rpio.Pwm)
	rpiGpio.fan.Freq(64000)

	//set light
	rpiGpio.light = rpio.Pin(10)
	rpiGpio.light.Output()

	return nil
}

// set variable fan speed from 0 - 100
func SetFanSpeed(speed uint32) {
	if speed != 0 {
		speed += 20
	}
	rpiGpio.fan.DutyCycle(speed, cycle)
}

// turn uv bulb on or off
func LightEnable(onOff bool) {
	if onOff {
		fmt.Println("light turning on")
		rpiGpio.light.High()
	} else {
		fmt.Println("light turning off")
		rpiGpio.light.Low()
	}
}

// called when hotspot interface is alive and able to take conn
func HotSpot(enabled bool) {

}

// callled when internet service is detected
func Skynet(enabled bool) {

}

// called immediate on startup
func Booting() {

}

// turn fan off, turn light off
func GpioClose() {

	LightEnable(false)
	SetFanSpeed(uint32(0))

	rpio.StopPwm()
	rpio.Close()
}
