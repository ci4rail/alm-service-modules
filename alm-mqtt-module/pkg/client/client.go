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

package client

import (
	"alm-mqtt-module/pkg/avro"
	"alm-mqtt-module/pkg/schema"

	// no need for a name
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/linkedin/goavro"
	"github.com/nats-io/nats.go"
)

var (
	//go:embed avro_schemas/dataSchema.avsc
	dataSchema string
	// DataCodec is the parsed dataSchema.avsc
	DataCodec *goavro.Codec = avro.CreateSchema(dataSchema)
)

// Client is a struct containing client relevant data
type Client struct {
	nats   *nats.Conn
	target string
}

// NewClient creates a new client for talking to `alm-mqtt-module`
func NewClient(target string, nats *nats.Conn) *Client {
	return &Client{
		nats:   nats,
		target: target,
	}
}

// RegisterMqttTopic is used to let a client register to a specific MQTT topic.
// This functions returns a nats subject the client can subscribe to in order to read the
// forwarded message.
func (c *Client) RegisterMqttTopic(topic string) (schema.RegisterSubResponseType, error) {
	msg := make(map[string]interface{})
	msg["topic"] = topic
	registerSubRequestCodec, err := goavro.NewCodec(schema.RegisterSubRequest)
	if err != nil {
		return schema.RegisterSubResponseType{}, err
	}
	bytes, err := avro.Writer(msg, registerSubRequestCodec)
	if err != nil {
		return schema.RegisterSubResponseType{}, err
	}

	response, err := c.nats.Request(fmt.Sprintf("%s.config.register", c.target), bytes, 2*time.Second)
	if err != nil {
		if c.nats.LastError() != nil {
			return schema.RegisterSubResponseType{}, fmt.Errorf("%v for request", c.nats.LastError())
		}
		return schema.RegisterSubResponseType{}, err
	}

	avro, _ := avro.NewReader(response.Data)
	j, _ := avro.ByteString()

	res := schema.RegisterSubResponseType{}
	err = json.Unmarshal(j, &res)
	if err != nil {
		return schema.RegisterSubResponseType{}, err
	}
	if len(res.Error) > 0 {
		return schema.RegisterSubResponseType{}, fmt.Errorf("%s", res.Error)
	}
	return res, nil
}

// UnregisterNatsSubject is used to unregister a client from a specific nats subject
// previously registered using 'RegisterMqttTopic'.
func (c *Client) UnregisterNatsSubject(subject string) error {
	msg := make(map[string]interface{})
	msg["subject"] = subject
	unregisterSubRequestCodec, err := goavro.NewCodec(schema.UnregisterSubRequest)
	if err != nil {
		return err
	}
	bytes, err := avro.Writer(msg, unregisterSubRequestCodec)
	if err != nil {
		return err
	}

	response, err := c.nats.Request(fmt.Sprintf("%s.config.unregister", c.target), bytes, 2*time.Second)
	if err != nil {
		if c.nats.LastError() != nil {
			return fmt.Errorf("%v for request", c.nats.LastError())
		}
		return err
	}

	avro, _ := avro.NewReader(response.Data)
	j, _ := avro.ByteString()

	res := schema.UnregisterSubResponseType{}
	err = json.Unmarshal(j, &res)
	if err != nil {
		return err
	}
	if res.Error != "" {
		return fmt.Errorf("%s", res.Error)
	}
	return nil
}

// PublishOnMqttTopic is used to to send to a specific MQTT topic.
func (c *Client) PublishOnMqttTopic(topic string, payload []byte) error {
	msg := make(map[string]interface{})
	msg["topic"] = topic
	msg["payload"] = payload
	pubRequestCodec, err := goavro.NewCodec(schema.PubRequest)
	if err != nil {
		return err
	}

	bytes, err := avro.Writer(msg, pubRequestCodec)
	if err != nil {
		return err
	}
	response, err := c.nats.Request(fmt.Sprintf("%s.publish", c.target), bytes, 2*time.Second)
	if err != nil {
		if c.nats.LastError() != nil {
			return fmt.Errorf("%v for request", c.nats.LastError())
		}
		return err
	}

	avro, _ := avro.NewReader(response.Data)
	j, _ := avro.ByteString()

	res := schema.PubResponseType{}
	err = json.Unmarshal(j, &res)
	if err != nil {
		return err
	}
	if res.Error != "" {
		return fmt.Errorf("%s", res.Error)
	}
	return nil
}

// RequestReply is used to to send to a specific MQTT topic.
func (c *Client) RequestReply(topic string, payload []byte, timeout int32) ([]byte, error) {
	msg := make(map[string]interface{})
	msg["topic"] = topic
	msg["payload"] = payload
	msg["timeout"] = timeout
	codec, err := goavro.NewCodec(schema.ReqResRequest)
	if err != nil {
		return []byte{}, err
	}

	bytes, err := avro.Writer(msg, codec)
	if err != nil {
		return []byte{}, err
	}
	response, err := c.nats.Request(fmt.Sprintf("%s.request-response", c.target), bytes, (time.Duration(timeout)*time.Millisecond)+(5*time.Second))
	if err != nil {
		if c.nats.LastError() != nil {
			return []byte{}, fmt.Errorf("%v for request", c.nats.LastError())
		}
		return []byte{}, err
	}

	avro, err := avro.NewReader(response.Data)
	if err != nil {
		return []byte{}, err
	}

	m, err := avro.Map()
	if err != nil {
		return []byte{}, err
	}

	res := schema.ReqResResponsetType{
		Payload: m["payload"].([]byte),
		Error:   m["error"].(string),
	}

	if res.Error != "" {
		return []byte{}, fmt.Errorf(res.Error)
	}

	return res.Payload, nil
}
