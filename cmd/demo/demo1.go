package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	aio "github.com/adafruit/io-client-go"
	"github.com/dfense/demo1/drv"
	"github.com/dfense/demo1/net"
	"github.com/dfense/demo1/settings"
)

func render(label string, f *aio.Feed) {
	sfeed, _ := json.MarshalIndent(f, "", "  ")
	fmt.Printf("--- %v\n", label)
	fmt.Println(string(sfeed))
}

func ShowAll(client *aio.Client) {
	// Get the list of all available feeds
	feeds, _, err := client.Feed.All()
	if err != nil {
		fmt.Println("UNEXPECTED ERROR!", err)
		panic(err)
	}

	for _, feed := range feeds {
		render(feed.Name, feed)
	}
}

func main() {
	// client := aio.NewClient(os.Getenv("ADAFRUIT_IO_KEY"))
	// ShowAll(client)
	configFile := flag.String("configfile", "/boot/settings.yaml", "full path and file name to config file")
	log.Println(configFile)
	flag.Parse()
	config, err := settings.GetConfigFromFile(*configFile)
	if err != nil {
		log.Fatalf("config file read error : %s", err)
	}

	drv.GpioOpen()

	dataChannel := make(chan interface{})
	// bme680 := drv.NewBME680(drv.PRIMARY)
	// bme680.Start(dataChannel)
	// Make a new channel which will be used to ensure we get all output

	bme680S := drv.NewBME680(drv.SECONDARY)
	bme680S.Start(dataChannel)

	log.Println("bme started, now connecting to adafruit")

	// network mqtt agent connection to cloud
	client := net.Client{}

	// never returns if no connection.. !!
	client.Connect(*config)
	go client.Listen(dataChannel)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	s := <-sigc
	fmt.Println(s)

	bme680S.Stop()
	drv.GpioClose()

	log.Println("Stopping main")

}
