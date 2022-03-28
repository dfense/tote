package net

import (
	"fmt"
	"log"
	"strconv"

	"github.com/dfense/demo1/drv"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	c mqtt.Client
}

// secure in encrpypted file
const (
	iourl      = "ssl://io.adafruit.com:8883"
	secureid   = "dfensenetAmerica"
	passwd     = "4f73eae9093ad766953dba916ad97d111a4dfa70"
	feedPrefix = "dfense/feeds/"
	nodeName   = "t1000" //read in from hostname or file?
)

// subscribe to command channel
// create message handlers
// connect to server and wait for commands
// create a publish function

func (c *Client) Connect() {

	opts := mqtt.NewClientOptions().AddBroker(iourl).SetClientID(secureid)
	// TODO read in account info
	opts.SetUsername("dfense")
	// TODO read in secret
	opts.SetPassword(passwd)

	c.c = mqtt.NewClient(opts)
	for {
		if token := c.c.Connect(); token.Wait() && token.Error() != nil {
			log.Println(token.Error())
		} else {
			break
		}
	}

	msgRcvd := func(client mqtt.Client, message mqtt.Message) {
		fmt.Printf("Received message on topic: %s\n Message: %s\n", message.Topic(), message.Payload())
		if message.Topic() == feedPrefix+nodeName+"-fanspeed" {
			speed := string(message.Payload())
			var u64 uint64
			var u32 uint32
			u64, err := strconv.ParseUint(speed, 10, 32)
			if err != nil {
				fmt.Println(err)
			}
			u32 = uint32(u64)
			drv.SetFanSpeed(u32)
			fmt.Printf("setting fan speed to %d", u32)
		} else if string(message.Payload()) == "uvlight_on" {
			drv.LightEnable(true)
		} else if string(message.Payload()) == "uvlight_off" {
			drv.LightEnable(false)
		} else if string(message.Payload()) == "2000" {
			log.Println("RESET")
		}
	}

	if token := c.c.Subscribe(feedPrefix+nodeName+"-commands", 0, msgRcvd); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	if token := c.c.Subscribe(feedPrefix+nodeName+"-fanspeed", 0, msgRcvd); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	// if token := c.c.Subscribe("example/topic", 0, msgRcvd); token.Wait() && token.Error() != nil {
	// 	fmt.Println(token.Error())
	// }

	//c.Publish("dfense/feeds/test", 1, false, `{"value": {"sensor-1":22.587,"sensor-2":13.182}}`)
	// if token := c.c.Publish("dfense/feeds/test", 1, false, `{"value": {"command": true, "score": 99, "message":"verbose"}}`); token.Wait() && token.Error() != nil {
	// 	fmt.Println(token.Error())
	// }
	// if token := c.c.Publish("dfense/feeds/t1barometric", 1, false, "4000"); token.Wait() && token.Error() != nil {
	// 	fmt.Println(token.Error())
	// }

}

func (c Client) Listen(inbound <-chan interface{}) {

	throttle := 0
	for {

		data := <-inbound
		switch v := data.(type) {
		case *drv.BME680:
			log.Printf("BME: %+v\n", v)
			throttle++
			if throttle > 3 {
				throttle = 0
			}
			sensorNumber := 1
			if v.I2c == drv.SECONDARY {
				sensorNumber = 2
			}
			// slow it down by half , put in hysterisis buffer soon...
			// alternate approach is send on payload with all values in it
			if throttle == 0 {
				c.c.Publish(feedPrefix+nodeName+"-barometric"+strconv.Itoa(sensorNumber), 1, false, fmt.Sprintf("{\"value\": %f}", v.Pressure))
				c.c.Publish(feedPrefix+nodeName+"-temp"+strconv.Itoa(sensorNumber), 1, false, fmt.Sprintf("{\"value\": %f}", v.RawTempC))
				c.c.Publish(feedPrefix+nodeName+"-RH"+strconv.Itoa(sensorNumber), 1, false, fmt.Sprintf("{\"value\": %f}", v.RawRH))
				c.c.Publish(feedPrefix+nodeName+"-CO2-"+strconv.Itoa(sensorNumber), 1, false, fmt.Sprintf("{\"value\": %f}", v.Co2PPM))
				log.Printf("BME: %+v\n", v)
			}
		case string:
			log.Printf("string: %s\n", v)
		default:
			log.Printf("channel no match  %T\n", data)

		}
	}
}
