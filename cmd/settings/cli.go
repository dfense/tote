package main

import (
	"flag"
	"log"

	"github.com/dfense/demo1/settings"
)

func main() {
	configFile := flag.String("configfile", "/boot/settings.yaml", "full path and file name to config file")
	flag.Parse()
	log.Printf("config-arg : %s\n", *configFile)

	config, err := settings.GetConfigFromFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Tote: %s\n", config.Totename)
}

/**
A file that looks like this should be loaded to test

totename: "Tote-102"
feedprefix: "dfense/feeds/"
iourl: "ssl://io.adafruit.com:8883"
secureid: "dfensenetAmerica"o

*/
