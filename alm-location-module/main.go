/*
Copyright © 2021 Ci4Rail GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"alm-location-module/internal/version"
	"alm-location-module/pkg/gpsd"
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/linkedin/goavro/v2"
	"github.com/nats-io/nats.go"
	iso8601 "github.com/relvacode/iso8601"
)

const (
	connectTimeoutSeconds int = 30
)

var (
	// GpsdHost is defined by running this container with `--add-host=host.docker.internal:host-gateway`
	gpsdHost    = "host.docker.internal:2947"
	invalidSent = false
	noFixSent   = false
)

type position struct {
	lat       float64
	lon       float64
	mode      int32
	timestamp time.Time
}

func main() {
	gpsdHostEnv := os.Getenv("GPSD_HOST")
	if len(gpsdHostEnv) > 0 {
		gpsdHost = gpsdHostEnv
	}

	log.Printf("alm-location-module version: %s\n", version.Version)

	natsServer := "nats"
	if env := os.Getenv("NATS_SERVER"); len(env) > 0 {
		natsServer = env
	}

	deviceID := "null"
	if env := os.Getenv("IOTEDGE_DEVICEID"); len(env) > 0 {
		deviceID = env
	}

	// Connect Options
	opts := []nats.Option{nats.Name("ads-node-module"), nats.Timeout(30 * time.Second)}
	opts = setupConnOptions(opts)
	ncChan := make(chan *nats.Conn)
	go func() {
		for i := 0; i < connectTimeoutSeconds; i++ {
			if nc, err := nats.Connect(natsServer, opts...); err != nil {
				log.Printf("Connect failed: %s\n", err)
				log.Printf("Reconnecting to '%s'\n", natsServer)
			} else {
				log.Printf("Connected to '%s'\n", natsServer)
				ncChan <- nc
				return
			}
			time.Sleep(time.Second)
		}
	}()

	nc := <-ncChan
	defer nc.Close()

	// avro schema defintion
	codec, err := goavro.NewCodec(`
	{
		"type": "record",
		"name": "service.location.gps",
		"doc": "device location by lat and long",
		"fields" : [
		{
			"name": "device",
			"type": "string"
		},
		{
			"name": "acqTime",
			"type": {
				"type": "int",
				"doc": "unix timestamp",
				"logicalType": "timestamp"
			}
		},
		{
			"name": "lat",
			"type": "double"
		},
		{
			"name": "lon",
			"type": "double"
		}
		]
	}`)

	if err != nil {
		fmt.Println(err)
	}

	msg := make(map[string]interface{})

	newPositionChan := make(chan position)

	gpsChan := make(chan *gpsd.Connection)
	go func() {
		for i := 0; i < connectTimeoutSeconds; i++ {
			if gpsdClient, err := gpsd.NewClient(gpsdHost); err != nil {
				log.Printf("Connect failed: %s\n", err)
				log.Printf("Reconnecting to '%s'\n", gpsdHost)
			} else {
				log.Printf("Connected to '%s'\n", gpsdHost)
				gpsChan <- gpsdClient
				return
			}
			time.Sleep(time.Second)
		}
	}()

	gpsClient := <-gpsChan
	gpsClient.RegisterTpv(func(r interface{}) {
		tpv := r.(*gpsd.Tpv)
		t, err := iso8601.Parse([]byte(tpv.Time))
		if err != nil {
			fmt.Println(err)
		}
		pos := position{
			lat:       tpv.Lat,
			lon:       tpv.Lon,
			mode:      int32(tpv.Mode),
			timestamp: t,
		}

		fmt.Printf("Received value: %f,%f,%d,%d\n", pos.lat, pos.lon, pos.mode, pos.timestamp.Unix())
		if pos.mode >= 2 {
			invalidSent = false
			noFixSent = false
			newPositionChan <- pos
			return
		}

		// This state means: the gps modem did not have any fix at all. Therefore
		// a single notifaction is sent with the system timestamp.
		if pos.lat == 0.0 && pos.lon == 0.0 && pos.mode == 0 {
			if !invalidSent {
				invalidSent = true
				noFixSent = false
				pos.timestamp = time.Now()
				fmt.Printf("Invalid data. Informing only once with system timestamp: %d", pos.timestamp.Unix())
				newPositionChan <- pos
			}
			return
		}

		if pos.mode == 1 {
			if !noFixSent {
				if pos.lat == 0.0 && pos.lon == 0.0 {
					pos.timestamp = time.Now()
					fmt.Printf("Invalid data. Informing only once with system timestamp: %d", pos.timestamp.Unix())
				} else {
					fmt.Printf("Lost GPS Fix. Informing only once with data timestamp: %d", pos.timestamp.Unix())
				}
				noFixSent = true
				invalidSent = false
				newPositionChan <- pos
			}
			return
		}
	})

	_, err = gpsClient.Watch()
	if err != nil {
		fmt.Println(err)
	}

	for {
		newPos := <-newPositionChan
		fmt.Println(newPos)
		// Define avro message content
		msg["device"] = deviceID
		msg["acqTime"] = newPos.timestamp.Unix()
		msg["lat"] = newPos.lat
		msg["lon"] = newPos.lon

		bin := new(bytes.Buffer)
		if err != nil {
			log.Fatalf("Failed to create event.avro file: %v", err)
		}

		ocfw, err := goavro.NewOCFWriter(goavro.OCFConfig{
			W:     bin,
			Codec: codec,
		})
		if err != nil {
			log.Fatalf("Failed to create the OCF Writer: %v", err)
		}

		err = ocfw.Append([]interface{}{msg})
		if err != nil {
			log.Fatalf("Failed to append to bin: %v", err)
		}

		err = nc.Publish("service.location", bin.Bytes())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Sending value: %f,%f,%d,%d\n", newPos.lat, newPos.lon, newPos.mode, newPos.timestamp.Unix())
	}
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}
