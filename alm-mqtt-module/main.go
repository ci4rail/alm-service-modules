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
	conf "alm-mqtt-module/internal/config"
	"alm-mqtt-module/internal/version"
	"alm-mqtt-module/pkg/avro"
	"alm-mqtt-module/pkg/client"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

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

func mqttHandler(c mqtt.Client, msg mqtt.Message) {
	fmt.Printf("New MQTT message for '%s'\n", msg.Topic())

	m := make(map[string]interface{})
	m["payload"] = msg.Payload()
	m["acqTime"] = time.Now().Unix()
	m["device"] = deviceID
	avro, err := avro.Writer(m, client.DataCodec)
	if err != nil {
		fmt.Println(err)
	}

	config.MessageChannelsMutex.Lock()
	ch := config.GetChannelsForTopic(msg.Topic())
	for k := range ch {
		// go func(k string) {
		ch[k] <- avro
		// }(k)
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
	}()

	natsClient = <-natsClientChan
	defer natsClient.Close()

	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.AddBroker(fmt.Sprintf("tcp://%s", mqttServer))

	mqttClientChan := make(chan *mqtt.Client)
	go func() {
		client := mqtt.NewClient(mqttOpts)
		for i := 0; i < connectTimeoutSeconds; i++ {
			if token := client.Connect(); token.Wait() && token.Error() != nil {
				log.Printf("Connect failed: %s\n", token.Error())
				log.Printf("Reconnecting to '%s'\n", mqttServer)
			} else {
				log.Printf("Connected to '%s'\n", mqttServer)
				mqttClientChan <- &client
				return
			}
			time.Sleep(time.Second)
		}
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
			if token := (*mqttClient).Subscribe(newMqttTopic, 1, mqttHandler); token.Wait() && token.Error() != nil {
				log.Fatal(token.Error())
			}

		case removeMqttTopic := <-newConfigUnregisterChan:
			fmt.Printf("Unsubscribing '%s'\n", removeMqttTopic)
			if token := (*mqttClient).Unsubscribe(removeMqttTopic); token.Wait() && token.Error() != nil {
				log.Fatal(token.Error())
			}

		case mqttMessage := <-pubChan:
			fmt.Printf("Sending message to topic '%s'\n", mqttMessage.Topic)

			if token := (*mqttClient).Publish(mqttMessage.Topic, 1, false, mqttMessage.Payload); token.Wait() && token.Error() != nil {
				log.Fatal(token.Error())
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
