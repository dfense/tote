package net

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dfense/demo1/drv"
	"github.com/dfense/demo1/settings"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	c mqtt.Client
}

// secure in encrpypted file
const (
	passwd = "4f73eae9093ad766953dba916ad97d111a4dfa70"
)

var (
	config settings.Config
)

func msgRcvd(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("Received message on topic: %s\n Message: %s\n", message.Topic(), message.Payload())
	if message.Topic() == config.FeedPrefix+config.Totename+"-fanspeed" {
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

// subscribe to command channel
// create message handlers
// connect to server and wait for commands
// create a publish function
func (c *Client) Connect(conf settings.Config) {

	config = conf
	opts := mqtt.NewClientOptions().AddBroker(config.IOURL).SetClientID(config.Secureid)
	// TODO read in account info
	opts.SetUsername(config.Username)
	// TODO read in secret
	opts.SetPassword(passwd)

	log.Printf("Connecting to : %s using %s username: %s secureID: %s \n", config.IOURL, opts.ClientID, config.Username, config.Secureid)
	//set OnConnect handler as anonymous function
	//after connected, subscribe to topic
	opts.OnConnect = func(c mqtt.Client) {

		if token := c.Subscribe(config.FeedPrefix+config.Totename+"-commands", 0, msgRcvd); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
		}

		if token := c.Subscribe(config.FeedPrefix+config.Totename+"-fanspeed", 0, msgRcvd); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
		}

	}

	c.c = mqtt.NewClient(opts)
	// try forever until a connection is made
	for {
		if token := c.c.Connect(); token.Wait() && token.Error() != nil {
			log.Println(token.Error())
			time.Sleep(time.Millisecond * 500)
		} else {
			break
		}
	}
}

// channel listener for any other application that needs
// to send data to cloud, just write message to the inbound channel
// wire it up outside of here in a controller/config layer
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
				c.c.Publish(config.FeedPrefix+config.Totename+"-barometric"+strconv.Itoa(sensorNumber), 1, false, fmt.Sprintf("{\"value\": %f}", v.Pressure))
				c.c.Publish(config.FeedPrefix+config.Totename+"-temp"+strconv.Itoa(sensorNumber), 1, false, fmt.Sprintf("{\"value\": %f}", v.RawTempC*1.8+32))
				c.c.Publish(config.FeedPrefix+config.Totename+"-RH"+strconv.Itoa(sensorNumber), 1, false, fmt.Sprintf("{\"value\": %f}", v.RawRH))
				c.c.Publish(config.FeedPrefix+config.Totename+"-CO2-"+strconv.Itoa(sensorNumber), 1, false, fmt.Sprintf("{\"value\": %f}", v.Co2PPM))
				log.Printf("BME: %+v\n", v)
			}
		case string:
			log.Printf("string: %s\n", v)
		default:
			log.Printf("channel no match  %T\n", data)

		}
	}
}
