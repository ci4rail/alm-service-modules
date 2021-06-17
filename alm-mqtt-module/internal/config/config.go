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

package config

import (
	"alm-mqtt-module/pkg/avro"
	schema "alm-mqtt-module/pkg/schema"
	"sync"
	"time"

	"encoding/json"
	"fmt"
	"log"

	"github.com/linkedin/goavro"
	"github.com/nats-io/nats.go"
)

const (
	// Timeout in seconds. When this timeout exceeds the corresponding channels will be removed.
	timeout = 5
)

// MqttMessage store MQTT message information required for transmission
type MqttMessage struct {
	Topic   string
	Payload []byte
}

type subjectChannelMapping struct {
	channel chan []byte
	subject string
}

// Config type to store configuration
type Config struct {
	registerSubRequestCodec    *goavro.Codec
	registerSubResponseCodec   *goavro.Codec
	unregisterSubRequestCodec  *goavro.Codec
	unregisterSubResponseCodec *goavro.Codec
	pubRequestCodec            *goavro.Codec
	pubResponseCodec           *goavro.Codec
	nats                       *nats.Conn
	basename                   string
	channels                   Channels
	newConfigRegisterChan      chan string
	newConfigUnregisterChan    chan string
	newConfigSendMQTTChan      chan MqttMessage
	MessageChannelsMutex       sync.Mutex
	MessageChannels            map[string][]subjectChannelMapping
	subscribed                 map[string]bool
}

// NewConfig creates a new config containing all channel definitions
func NewConfig(basename string, natsConn *nats.Conn, newConfigRegisterChan, newConfigUnregisterChan chan string, newConfigSendMQTTChan chan MqttMessage) *Config {
	return &Config{
		registerSubRequestCodec:    schema.RegisterSubRequestCodec,
		registerSubResponseCodec:   schema.RegisterSubResponseCodec,
		unregisterSubRequestCodec:  schema.UnregisterSubRequestCodec,
		unregisterSubResponseCodec: schema.UnregisterSubResponseCodec,
		pubRequestCodec:            schema.PubRequestCodec,
		pubResponseCodec:           schema.PubResponseCodec,
		nats:                       natsConn,
		basename:                   basename,
		channels:                   NewChannels(basename),
		newConfigRegisterChan:      newConfigRegisterChan,
		newConfigUnregisterChan:    newConfigUnregisterChan,
		newConfigSendMQTTChan:      newConfigSendMQTTChan,
		MessageChannels:            make(map[string][]subjectChannelMapping),
		subscribed:                 make(map[string]bool),
	}
}

func (c *Config) configHandlerRegister(msg *nats.Msg) {
	req := parseConfigRegisterRequest(msg)
	fmt.Printf("Register for '%s'\n", req.Topic)
	subject, err := c.channels.RegisterSub(req.Topic)

	var errText string = ""
	if err != nil {
		fmt.Println(err)
		errText = err.Error()
	}
	subjectChannelMapping := subjectChannelMapping{
		channel: make(chan []byte, 20),
		subject: subject,
	}
	c.MessageChannelsMutex.Lock()
	c.MessageChannels[req.Topic] = append(c.MessageChannels[req.Topic], subjectChannelMapping)
	c.MessageChannelsMutex.Unlock()
	c.subscribed[subject] = true
	go func(channel chan []byte, subject string) {
		for {
			avro, ok := <-channel
			if !ok {
				break
			}

			fmt.Printf("\t-> nats '%s'\n", subject)

			if c.subscribed[subject] {
				_, err = c.nats.Request(subject, avro, time.Duration(timeout)*time.Second)
				if err != nil {
					fmt.Printf("Subject '%s' timed out. Unregistering.\n", subject)
					_, err := c.cleanupSubject(subject)

					if err != nil {
						log.Fatal(err)
					}
					break
				}
			}
		}
	}(subjectChannelMapping.channel, subject)

	res := schema.RegisterSubResponseType{
		Subject: subject,
		Error:   errText,
	}
	r, err := c.createConfigRegisterResponse(res)
	if err != nil {
		log.Fatal(err)
	}
	err = msg.Respond(r)
	if err != nil {
		log.Fatal(err)
	}
	c.newConfigRegisterChan <- req.Topic
}

func (c *Config) configHandlerUnregister(msg *nats.Msg) {
	var errText string = ""
	req := parseConfigUnregisterRequest(msg)
	fmt.Printf("Unregister for '%s'\n", req.Subject)
	_, err := c.cleanupSubject(req.Subject)
	if err != nil {
		errText = err.Error()
	}

	res := schema.UnregisterSubResponseType{
		Error: errText,
	}
	r, err := c.createConfigUnregisterResponse(res)
	if err != nil {
		log.Fatal(err)
	}
	err = msg.Respond(r)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Config) handlerPublish(msg *nats.Msg) {
	var errText string = ""
	req := parsePublishRequest(msg)
	fmt.Printf("Received Publish Request for '%s'\n", req.Topic)

	if req.Topic == "" {
		errText = "Empty topic received"
	} else {
		mqttMessage := MqttMessage{
			Topic:   req.Topic,
			Payload: req.Payload,
		}
		c.newConfigSendMQTTChan <- mqttMessage
	}

	res := schema.PubResponseType{
		Error: errText,
	}
	r, err := c.createPublishResponse(res)
	if err != nil {
		log.Fatal(err)
	}
	err = msg.Respond(r)
	if err != nil {
		log.Fatal(err)
	}
}

func removeFromSubjectChannelMappingSlice(s []subjectChannelMapping, i int) []subjectChannelMapping {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (c *Config) createConfigRegisterResponse(res schema.RegisterSubResponseType) ([]byte, error) {
	msg := make(map[string]interface{})
	msg["subject"] = res.Subject
	msg["error"] = res.Error
	return avro.Writer(msg, c.registerSubResponseCodec)
}

func (c *Config) createConfigUnregisterResponse(res schema.UnregisterSubResponseType) ([]byte, error) {
	msg := make(map[string]interface{})
	msg["error"] = res.Error
	return avro.Writer(msg, c.unregisterSubResponseCodec)
}

func (c *Config) createPublishResponse(res schema.PubResponseType) ([]byte, error) {
	msg := make(map[string]interface{})
	msg["error"] = res.Error
	return avro.Writer(msg, c.pubResponseCodec)
}

func parseConfigRegisterRequest(msg *nats.Msg) schema.RegisterSubRequestType {
	avro, err := avro.NewReader(msg.Data)
	if err != nil {
		log.Fatal(err)
	}
	j, err := avro.ByteString()
	if err != nil {
		log.Fatal(err)
	}
	res := schema.RegisterSubRequestType{}
	err = json.Unmarshal(j, &res)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func parseConfigUnregisterRequest(msg *nats.Msg) schema.UnregisterSubRequestType {
	avro, err := avro.NewReader(msg.Data)
	if err != nil {
		log.Fatal(err)
	}
	j, err := avro.ByteString()
	if err != nil {
		log.Fatal(err)
	}
	res := schema.UnregisterSubRequestType{}
	err = json.Unmarshal(j, &res)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func parsePublishRequest(msg *nats.Msg) schema.PubRequestType {
	avro, err := avro.NewReader(msg.Data)
	if err != nil {
		log.Fatal(err)
	}

	m, err := avro.Map()
	if err != nil {
		log.Fatal(err)
	}

	return schema.PubRequestType{
		Topic:   m["topic"].(string),
		Payload: m["payload"].([]byte),
	}
}

// HandleConfigRequests registeres for configuration requests on the nats server
func (c *Config) HandleConfigRequests() {
	if _, err := c.nats.Subscribe(fmt.Sprintf("%s.config.register", c.basename), c.configHandlerRegister); err != nil {
		log.Fatal(err)
	}
	if _, err := c.nats.Subscribe(fmt.Sprintf("%s.config.unregister", c.basename), c.configHandlerUnregister); err != nil {
		log.Fatal(err)
	}
}

// HandlePublishRequests register handler for publish requests on the nats server
func (c *Config) HandlePublishRequests() {
	if _, err := c.nats.Subscribe(fmt.Sprintf("%s.publish", c.basename), c.handlerPublish); err != nil {
		log.Fatal(err)
	}
}

// GetRegistrations gets all client registrations for a specific topic
func (c *Config) GetRegistrations(topic string) []string {
	return c.channels.Get(topic)
}

// GetChannelsForTopic returns all go channels that feed the handling routines of all nats subscriptions for a given topic
func (c *Config) GetChannelsForTopic(topic string) map[string]chan []byte {
	ret := make(map[string]chan []byte)
	for _, mapping := range c.MessageChannels[topic] {
		ret[mapping.subject] = mapping.channel
	}
	return ret
}

func (c *Config) cleanupSubject(subject string) (string, error) {
	var err error
	delete(c.subscribed, subject)
	topic, err := c.channels.GetTopic(subject)
	if err != nil {
		fmt.Println(err)
	}
	_, err = c.channels.UnregisterSub(subject)
	if err != nil {
		fmt.Println(err)
	}

	c.MessageChannelsMutex.Lock()
	for i, chMapp := range c.MessageChannels[topic] {
		if chMapp.subject == subject {
			close(chMapp.channel)
			c.MessageChannels[topic] = removeFromSubjectChannelMappingSlice(c.MessageChannels[topic], i)
		}
	}
	// remove topic from message channels if no further subjects / channels contained
	if len(c.MessageChannels[topic]) <= 0 {
		delete(c.MessageChannels, topic)
	}
	c.MessageChannelsMutex.Unlock()

	if len(c.channels.Get(topic)) == 0 {
		c.newConfigUnregisterChan <- topic
	}

	return topic, err
}
