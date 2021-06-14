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

package config

import (
	"alm-mqtt-module/pkg/avro"
	schema "alm-mqtt-module/pkg/schema"
	"encoding/json"
	"testing"

	"github.com/linkedin/goavro"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func createNatsMessage(assert *assert.Assertions, msg map[string]interface{}, schema string) nats.Msg {
	pubRequestCodec, err := goavro.NewCodec(schema)
	assert.Nil(err)

	bytes, err := avro.Writer(msg, pubRequestCodec)
	assert.Nil(err)

	return nats.Msg{
		Subject: "",
		Reply:   "",
		Data:    bytes,
		Sub:     &nats.Subscription{},
	}

}

func TestParsePublishRequestWithPublishSchemaAndJsonPayload(t *testing.T) {
	assert := assert.New(t)

	topic := "abc123/.-"

	// Create json payload
	type message struct {
		Counter     int    `json:"counter"`
		Temperature int    `json:"temperature"`
		UnicornName string `json:"unicornName"`
	}
	payload := message{
		Counter:     42,
		Temperature: 36,
		UnicornName: "Diamond Butter",
	}
	bytePayload, _ := json.Marshal(payload)

	// Prepare nats message contents
	msg := make(map[string]interface{})
	msg["topic"] = topic
	msg["payload"] = bytePayload

	natsMsg := createNatsMessage(assert, msg, schema.PubRequest)

	req := parsePublishRequest(&natsMsg)
	assert.Equal(req.Payload, bytePayload)
	assert.Equal(req.Topic, topic)
}

func TestParsePublishRequestWithPublishSchemaAndStringPayload(t *testing.T) {
	assert := assert.New(t)

	msg := make(map[string]interface{})
	topic := "abc123/.-"
	msg["topic"] = topic

	// Create string payload
	payload := "payload/.-!?"
	bytePayload := []byte(payload)
	msg["payload"] = bytePayload

	natsMsg := createNatsMessage(assert, msg, schema.PubRequest)

	req := parsePublishRequest(&natsMsg)
	assert.Equal(req.Payload, bytePayload)
	assert.Equal(string(req.Payload), payload)
	assert.Equal(req.Topic, topic)
}
