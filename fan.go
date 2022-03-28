/*
Toggles a LED on physical pin 35 (GPIO pin 19)
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"pwm/drv"
	"strconv"
	"syscall"
	"time"
	//"github.com/stianeikeland/go-rpio"
)

func main() {

	/*
		err := rpio.Open()
		if err != nil {
			os.Exit(1)
		}
		defer rpio.Close()

		pin := rpio.Pin(19)
		pin.Mode(rpio.Pwm)
		pin.Freq(10000)
		//pin.Freq(64)
		// the LED will be blinking at 2000Hz
		pin.DutyCycle(100, 100)
		// (source frequency divided by cycle length => 64000/32 = 2000)

		// five times smoothly fade in and out
		fmt.Println("Starting")

		go func() {
			for i := 0; i < 50; i++ {
				fmt.Printf("LoopID: %d", i)
				for i := uint32(0); i < 32; i++ { // increasing brightness
					fmt.Printf("Internal LoopID: %d", i)
					//pin.DutyCycle(i, 32)
					// time.Sleep(time.Second/32)
					time.Sleep(time.Second)
				}
				time.Sleep(time.Second * 3)
				for i := uint32(32); i > 0; i-- { // decreasing brightness
					fmt.Printf("External LoopID: %d", i)
					//pin.DutyCycle(i, 32)
					time.Sleep(time.Second)
				}
			}

		}()

	*/

	err := drv.GpioOpen()
	if err != nil {
		fmt.Printf("Erorr opening gpio: %s\n", err)
	}

	drv.LightEnable(true)
	time.Sleep(time.Second * 5)
	drv.LightEnable(false)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("enter fan speed 0 - 100 => ")
	go func() {
		for scanner.Scan() {
			var u64 uint64
			var u32 uint32
			u64, err = strconv.ParseUint(scanner.Text(), 10, 32)
			if err != nil {
				fmt.Println(err)
			}
			u32 = uint32(u64)
			drv.SetFanSpeed(u32)
			fmt.Printf("setting fan speed to %d", u32)
		}
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	s := <-sigc
	drv.GpioClose()
	// stopGPIO(&pin)
	fmt.Printf("leaving main %v", s)

}

/*
func stopGPIO(pin *rpio.Pin) {
	fmt.Println("Stopping GPIO")
	rpio.StopPwm()
	rpio.Close()
}
*/
