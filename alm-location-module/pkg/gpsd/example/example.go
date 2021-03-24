package main

import (
	"alm-location-module/pkg/gpsd"
	"fmt"
	"os"

	"github.com/relvacode/iso8601"
)

func main() {
	c, err := gpsd.NewClient("localhost:2947")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c.RegisterTpv(func(r interface{}) {
		data := r.(*gpsd.Tpv)
		fmt.Printf("Lat: %f\n", data.Lat)
		fmt.Printf("Lon: %f\n", data.Lon)
		t, _ := iso8601.Parse([]byte(data.Time))
		fmt.Printf("Timestamp: %d\n", t.Unix())
	})

	done, err := c.Watch()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer c.Close()
	<-done
}
