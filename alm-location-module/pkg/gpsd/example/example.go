package main

import (
	"alm-location-module/pkg/gpsd"
	"fmt"
	"os"
	"time"

	"github.com/relvacode/iso8601"
)

var (
	gpsdHost = "localhost:2947"
)

type position struct {
	lat       float64
	lon       float64
	timestamp time.Time
}

func main() {
	gpsdHostEnv := os.Getenv("GPSD_HOST")
	if len(gpsdHostEnv) > 0 {
		gpsdHost = gpsdHostEnv
	}

	c, err := gpsd.NewClient(gpsdHost)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	newPositionChan := make(chan position)

	c.RegisterTpv(func(r interface{}) {
		tpv := r.(*gpsd.Tpv)
		t, err := iso8601.Parse([]byte(tpv.Time))
		if err != nil {
			fmt.Println(err)
		}
		pos := position{
			lat:       tpv.Lat,
			lon:       tpv.Lon,
			timestamp: t,
		}
		fmt.Println(pos)
		newPositionChan <- pos
	})

	_, err = c.Watch()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		newPos := <-newPositionChan
		fmt.Println(newPos)
		fmt.Printf("Sending value: %f,%f,%d\n", newPos.lat, newPos.lon, newPos.timestamp.Unix())
	}
}
