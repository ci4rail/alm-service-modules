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
	conf "alm-mqtt-module/internal/config"
	"alm-mqtt-module/internal/version"
	"alm-mqtt-module/pkg/avro"
	"alm-mqtt-module/pkg/client"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/eclipse/paho.golang/paho"

	"github.com/nats-io/nats.go"
)

const (
	connectTimeoutSeconds int = 30
)

var (
	config                  *conf.Config
	natsClient              *nats.Conn
	deviceID                string
	newConfigRegisterChan   chan string
	newConfigUnregisterChan chan string
	pubChan                 chan conf.MqttMessage
)

func mqttHandler(msg *paho.Publish) {
	fmt.Printf("New MQTT message for '%s'\n", msg.Topic)

	m := make(map[string]interface{})
	m["payload"] = msg.Payload
	m["acqTime"] = time.Now().Unix()
	m["device"] = deviceID
	avro, err := avro.Writer(m, client.DataCodec)
	if err != nil {
		fmt.Println(err)
	}

	config.MessageChannelsMutex.Lock()
	ch := config.GetChannelsForTopic(msg.Topic)
	for k := range ch {
		ch[k] <- avro
	}
	config.MessageChannelsMutex.Unlock()
}

func main() {
	log.Printf("alm-mqtt-module version: %s\n", version.Version)

	mqttServer := "mosquitto:1883"
	if env := os.Getenv("MQTT_SERVER"); len(env) > 0 {
		mqttServer = env
	}

	natsServer := "nats"
	if env := os.Getenv("NATS_SERVER"); len(env) > 0 {
		natsServer = env
	}

	deviceID = "null"
	if env := os.Getenv("IOTEDGE_DEVICEID"); len(env) > 0 {
		deviceID = env
	}

	// Connect Options
	opts := []nats.Option{nats.Name("alm-mqtt-module"), nats.Timeout(30 * time.Second)}
	opts = setupConnOptions(opts)
	natsClientChan := make(chan *nats.Conn)
	go func() {
		for i := 0; i < connectTimeoutSeconds; i++ {
			if natsClient, err := nats.Connect(natsServer, opts...); err != nil {
				log.Printf("Connect failed: %s\n", err)
				log.Printf("Reconnecting to '%s'\n", natsServer)
			} else {
				log.Printf("Connected to '%s'\n", natsServer)
				natsClientChan <- natsClient
				return
			}
			time.Sleep(time.Second)
		}

		log.Fatal("Cannot connect to NAS server.")
	}()

	natsClient = <-natsClientChan
	defer natsClient.Close()

	mqttClientChan := make(chan *paho.Client)
	go func() {
		conn, err := net.Dial("tcp", mqttServer)
		if err != nil {
			log.Fatalf("Failed to connect to %s: %s", mqttServer, err)
		}

		// From https://github.com/eclipse/paho.golang/blob/336f2adf08b8233199ac8132b8dd12cbb8c69eca/paho/client.go
		// client.Conn *MUST* be set to an already connected net.Conn before
		// Connect() is called.
		client := paho.NewClient(paho.ClientConfig{
			Conn:   conn,
			Router: paho.NewSingleHandlerRouter(mqttHandler),
		})

		for i := 0; i < connectTimeoutSeconds; i++ {
			// Connect Client to MQTT Broker
			res, err := client.Connect(context.Background(), &paho.Connect{})
			if err != nil {
				log.Printf("Failed to connect to %s: %s", mqttServer, err.Error())
			} else if res.ReasonCode != 0 {
				log.Printf("Failed to connect with reason: %d - %s", res.ReasonCode, res.Properties.ReasonString)
			} else {
				println("Connected to MQTT Broker successfully")
				mqttClientChan <- client
				return
			}
			time.Sleep(time.Second)
		}
		log.Fatal("Cannot connect to MQTT broker.")
	}()

	mqttClient := <-mqttClientChan

	// Channels to register and unregister for 100 simultations requests
	newConfigRegisterChan = make(chan string, 100)
	newConfigUnregisterChan = make(chan string, 100)
	pubChan = make(chan conf.MqttMessage, 100)

	config = conf.NewConfig("alm-mqtt-module", natsClient, newConfigRegisterChan, newConfigUnregisterChan, pubChan)

	go config.HandleConfigRequests()
	go config.HandlePublishRequests()

	for {
		select {
		case newMqttTopic := <-newConfigRegisterChan:
			fmt.Printf("Subscribing '%s'\n", newMqttTopic)

			if _, err := (*mqttClient).Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					newMqttTopic: {QoS: 1},
				},
			}); err != nil {
				log.Fatal(err)
			}

		case removeMqttTopic := <-newConfigUnregisterChan:
			fmt.Printf("Unsubscribing '%s'\n", removeMqttTopic)
			if _, err := (*mqttClient).Unsubscribe(context.Background(), &paho.Unsubscribe{
				Topics: []string{removeMqttTopic},
			}); err != nil {
				log.Fatal(err)
			}

		case mqttMessage := <-pubChan:
			fmt.Printf("Publish message to topic '%s'\n", mqttMessage.Topic)

			if _, err := (*mqttClient).Publish(context.Background(), &paho.Publish{
				QoS:     1,
				Topic:   mqttMessage.Topic,
				Payload: mqttMessage.Payload,
			}); err != nil {
				log.Fatal(err)
			}

		default:
			time.Sleep(time.Second)
		}
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
