package main

import (
	"alm-mqtt-module/pkg/avro"
	"alm-mqtt-module/pkg/client"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	natsServer := "nats"

	natsClient, err := nats.Connect(natsServer)
	if err != nil {
		log.Fatal(err)
	}

	subject, err := client.RegisterMqttTopic("alm-mqtt-module", "simulation/temperature", natsClient)
	if err != nil {
		log.Fatal(err)
	}
	cleanup := func() {
		err := client.UnregisterNatsSubject("alm-mqtt-module", subject, natsClient)
		if err != nil {
			log.Fatal(err)
		}
		natsClient.Close()
	}

	fmt.Printf("Subscribing to nats subject '%s'\n", subject)
	_, err = natsClient.Subscribe(subject, func(msg *nats.Msg) {
		avro, err := avro.NewAvroReader(msg.Data)
		if err != nil {
			log.Fatal(err)
		}
		j, err := avro.AvroToByteString()
		if err != nil {
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
		if i >= 10 {
			cleanup()
			break
		}
	}
}
