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
	"alm-mqtt-module/pkg/client"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	natsServer := "nats"
	opts := []nats.Option{nats.Name("alm-mqtt-module-publish"), nats.Timeout(3 * time.Second)}
	natsClient, err := nats.Connect(natsServer, opts...)
	if err != nil {
		log.Fatal(err)
	}
	client := client.NewClient("alm-mqtt-module", natsClient)

	counter := 0
	for {
		time.Sleep(time.Second)
		counter++
		type message struct {
			Counter int `json:"counter"`
		}
		msg := message{
			Counter: counter,
		}

		b, _ := json.Marshal(msg)
		fmt.Printf("Sending message: %s\n", b)
		res, err := client.RequestReply("/pis/cmd/state", b, int32(5000))
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Received message: %s\n", res)
		}

		if counter >= 20 {
			return
		}
	}
}
