/*
Copyright © 2021 edgefarm.io

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
	"alm-mqtt-module/pkg/avro"
	"alm-mqtt-module/pkg/client"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	natsServer := "nats"
	opts := []nats.Option{nats.Name("alm-mqtt-module-example"), nats.Timeout(3 * time.Second)}
	natsClient, err := nats.Connect(natsServer, opts...)
	if err != nil {
		log.Fatal(err)
	}
	client := client.NewClient("alm-mqtt-module", natsClient)
	res, err := client.RegisterMqttTopic("simulation/temperature")
	if err != nil {
		log.Fatal(err)
	}

	cleanup := func() {
		err := client.UnregisterNatsSubject(res.Subject)
		if err != nil {
			log.Fatal(err)
		}
		natsClient.Close()
	}

	fmt.Printf("Subscribing to nats subject '%s'\n", res.Subject)
	_, err = natsClient.Subscribe(res.Subject, func(msg *nats.Msg) {
		avro, err := avro.NewReader(msg.Data)
		if err != nil {
			log.Fatal(err)
		}
		j, err := avro.ByteString()
		if err != nil {
			log.Fatal(err)
		}
		if err := msg.Respond([]byte{}); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", string(j))
	})
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	for {
		time.Sleep(time.Second)
		i++
		type message struct {
			Counter int `json:"counter"`
		}
		msg := message{
			Counter: i,
		}

		b, _ := json.Marshal(msg)
		err = client.PublishOnMqttTopic("example/app", b)
		if err != nil {
			fmt.Println("Error:", err)
			cleanup()
		}
		if i >= 20 {
			cleanup()
			return
		}
	}
}
