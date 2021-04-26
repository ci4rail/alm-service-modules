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
	//go:embed avro_schemas/dataSchema.avro
	dataSchema string
	// DataCodec can be used to decode avro data
	DataCodec = avro.CreateSchema(dataSchema)
)

// RegisterMqttTopic is used to let a client register to a specific MQTT topic.
// This functions returns a nats subject the client can subscribe to in order to read the
// forwarded message.
func RegisterMqttTopic(target string, topic string, client *nats.Conn) (string, error) {
	msg := make(map[string]interface{})
	msg["topic"] = topic
	registerSubRequestCodec, err := goavro.NewCodec(schema.RegisterSubRequest)
	if err != nil {
		return "", err
	}
	bytes, err := avro.NewAvroWriter(msg, registerSubRequestCodec)
	if err != nil {
		return "", err
	}

	response, err := client.Request(fmt.Sprintf("%s.config.register", target), bytes, 2*time.Second)
	if err != nil {
		if client.LastError() != nil {
			return "", fmt.Errorf("%v for request", client.LastError())
		}
		return "", err
	}

	avro, _ := avro.NewAvroReader(response.Data)
	j, _ := avro.AvroToByteString()

	res := schema.RegisterSubResponseType{}
	err = json.Unmarshal(j, &res)
	if err != nil {
		return "", err
	}
	if res.Error != "" {
		return "", fmt.Errorf("%s", res.Error)
	}
	return res.Subject, nil
}

// UnregisterNatsSubject is used to unregister a client from a specific nats subject
// previously registered using 'RegisterMqttTopic'.
func UnregisterNatsSubject(target string, subject string, client *nats.Conn) error {
	msg := make(map[string]interface{})
	msg["subject"] = subject
	unregisterSubRequestCodec, err := goavro.NewCodec(schema.UnregisterSubRequest)
	if err != nil {
		return err
	}
	bytes, err := avro.NewAvroWriter(msg, unregisterSubRequestCodec)
	if err != nil {
		return err
	}

	response, err := client.Request(fmt.Sprintf("%s.config.unregister", target), bytes, 2*time.Second)
	if err != nil {
		if client.LastError() != nil {
			return fmt.Errorf("%v for request", client.LastError())
		}
		return err
	}

	avro, _ := avro.NewAvroReader(response.Data)
	j, _ := avro.AvroToByteString()

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
