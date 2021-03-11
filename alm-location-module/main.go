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
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/linkedin/goavro/v2"
	"github.com/nats-io/nats.go"
)

const (
	defaultUpdateIntervalMs int = 1000
	delta                       = 0.00016
)

func main() {
	log.Printf("alm-location-module version: %s\n", version.Version)
	updateIntervalMs := defaultUpdateIntervalMs
	if i := os.Getenv("UPDATE_INTERVAL_MS"); i != "" {
		interval, err := strconv.Atoi(i)
		updateIntervalMs = interval
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Info: using update interval in %d milliseconds\n", updateIntervalMs)
	} else {
		log.Printf("Info: env UPDATE_INTERVAL_MS. Using default %d\n", updateIntervalMs)
	}

	natsServer := "nats"
	if env := os.Getenv("NATS_SERVER"); len(env) > 0 {
		natsServer = env
	}

	deviceID := "null"
	if env := os.Getenv("IOTEDGE_DEVICEID"); len(env) > 0 {
		deviceID = env
	}

	// Connect Options.
	opts := []nats.Option{nats.Name("ads-node-module")}
	opts = setupConnOptions(opts)

	// Connect to NATS
	nc, err := nats.Connect(natsServer, opts...)
	if err != nil {
		log.Fatal(err)
	}
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
	lat := 49.445273
	lon := 11.082713

	for {
		// Define avro message content
		msg["device"] = deviceID
		msg["acqTime"] = time.Now().Unix()
		msg["lat"] = lat
		msg["lon"] = lon

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
		time.Sleep(time.Duration(updateIntervalMs) * time.Millisecond)
		rand.Seed(time.Now().UnixNano())
		randLat := -delta/2 + rand.Float64()*(delta)
		randLon := -delta/2 + rand.Float64()*(delta)
		lat += randLat
		lon += randLon
		fmt.Printf("Sending value: %f,%f\n", lat, lon)
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
