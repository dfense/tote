package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	aio "github.com/adafruit/io-client-go"
	"github.com/dfense/demo1/drv"
	"github.com/dfense/demo1/net"
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

	bme680 := drv.NewBME680(drv.PRIMARY)
	// Make a new channel which will be used to ensure we get all output
	dataChannel := make(chan interface{})
	bme680.Start(dataChannel)

	bme680S := drv.NewBME680(drv.SECONDARY)
	bme680S.Start(dataChannel)

	log.Println("bme started, now connecting to adafruit")

	client := net.Client{}
	client.Connect()
	go client.Listen(dataChannel)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	//go func() {
	s := <-sigc
	fmt.Println(s)

	bme680.Stop()
	// ... do something ...

	//}()

	log.Println("Stopping main")

}
