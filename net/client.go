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

// subscribe to command channel

// create message handlers

// connect to server and wait for commands

// create a publish function

func (c *Client) Connect() {

	opts := mqtt.NewClientOptions().AddBroker("ssl://io.adafruit.com:8883").SetClientID("dfensenetAmerica")
	opts.SetUsername("dfense")
	opts.SetPassword("4f73eae9093ad766953dba916ad97d111a4dfa70")
	c.c = mqtt.NewClient(opts)
	if token := c.c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	msgRcvd := func(client mqtt.Client, message mqtt.Message) {
		fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
	}

	if token := c.c.Subscribe("dfense/feeds/test", 0, msgRcvd); token.Wait() && token.Error() != nil {
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
			log.Printf("BME: %s\n", v)
			throttle++
			if throttle > 3 {
				throttle = 0
			}
			sensorNumber := 1
			if v.I2c == drv.SECONDARY {
				sensorNumber = 2
			}
			// slow it down by half , put in hysterisis buffer soon...
			if throttle == 0 {
				c.c.Publish("dfense/feeds/t1000-barometric"+strconv.Itoa(sensorNumber), 1, false, fmt.Sprintf("{\"value\": %f}", v.Pressure))
				c.c.Publish("dfense/feeds/t1000-temp"+strconv.Itoa(sensorNumber), 1, false, fmt.Sprintf("{\"value\": %f}", v.RawTempC))
				c.c.Publish("dfense/feeds/t1000-RH"+strconv.Itoa(sensorNumber), 1, false, fmt.Sprintf("{\"value\": %f}", v.RawRH))
				c.c.Publish("dfense/feeds/t1000-CO2-"+strconv.Itoa(sensorNumber), 1, false, fmt.Sprintf("{\"value\": %f}", v.Co2PPM))
				log.Printf("BME: %f\n", v.RawRH)
			}
		case string:
			log.Printf("string: \n", v)
		default:
			log.Printf("channel no match  %T\n", data)

		}
	}
}
