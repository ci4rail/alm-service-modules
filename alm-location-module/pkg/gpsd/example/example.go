package main

import (
	"alm-location-module/pkg/gpsd"
	"fmt"
	"os"
)

func main() {
	c, err := gpsd.NewClient("localhost:2947")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c.Register(gpsd.Tpv, func(r interface{}) {
		data := r.(*gpsd.TpvObj)
		fmt.Printf("Lat: %f\n", data.Lat)
		fmt.Printf("Lon: %f\n", data.Lon)
	})
	done, err := c.Watch()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer c.Close()
	<-done
}
