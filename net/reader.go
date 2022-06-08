package net

import (
	"log"

	"github.com/dfense/demo1/drv"
)

func ReadForever(inbound <-chan interface{}) {

	for {

		data := <-inbound
		switch v := data.(type) {
		case *drv.BME680:
			log.Printf("BME: %+v\n", v)
		case string:
			log.Printf("string: %s\n", v)
		default:
			log.Printf("channel no match  %T\n", data)

		}
	}
}
