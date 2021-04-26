package config

import (
	"alm-mqtt-module/pkg/avro"
	schema "alm-mqtt-module/pkg/schema"

	"encoding/json"
	"fmt"
	"log"

	"github.com/linkedin/goavro"
	"github.com/nats-io/nats.go"
)

// Config type to store configuration
type Config struct {
	registerSubRequestCodec    *goavro.Codec
	registerSubResponseCodec   *goavro.Codec
	unregisterSubRequestCodec  *goavro.Codec
	unregisterSubResponseCodec *goavro.Codec
	nats                       *nats.Conn
	basename                   string
	channels                   Channels
	newConfigRegisterChan      chan string
	newConfigUnregisterChan    chan string
}

// NewConfig creates a new config containing all channel definitions
func NewConfig(basename string, nats *nats.Conn, newConfigRegisterChan, newConfigUnregisterChan chan string) (*Config, error) {
	return &Config{
		registerSubRequestCodec:    schema.RegisterSubRequestCodec,
		registerSubResponseCodec:   schema.RegisterSubResponseCodec,
		unregisterSubRequestCodec:  schema.UnregisterSubRequestCodec,
		unregisterSubResponseCodec: schema.UnregisterSubResponseCodec,
		nats:                       nats,
		basename:                   basename,
		channels:                   NewChannels(basename),
		newConfigRegisterChan:      newConfigRegisterChan,
		newConfigUnregisterChan:    newConfigUnregisterChan,
	}, nil
}

func (c *Config) configHandlerRegister(msg *nats.Msg) {
	req := parseConfigRegisterRequest(msg)
	fmt.Printf("Register for '%s'\n", req.Topic)
	subject, err := c.channels.RegisterSub(req.Topic)
	var errText string = ""
	if err != nil {
		errText = err.Error()
	}
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

	topic, err := c.channels.UnregisterSub(req.Subject)
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

	if len(c.channels.Get(topic)) == 0 {
		c.newConfigUnregisterChan <- topic
	}
}

func (c *Config) createConfigRegisterResponse(res schema.RegisterSubResponseType) ([]byte, error) {
	msg := make(map[string]interface{})
	msg["subject"] = res.Subject
	msg["error"] = res.Error
	return avro.NewAvroWriter(msg, c.registerSubResponseCodec)
}

func (c *Config) createConfigUnregisterResponse(res schema.UnregisterSubResponseType) ([]byte, error) {
	msg := make(map[string]interface{})
	msg["error"] = res.Error
	return avro.NewAvroWriter(msg, c.unregisterSubResponseCodec)
}

func parseConfigRegisterRequest(msg *nats.Msg) schema.RegisterSubRequestType {
	avro, err := avro.NewAvroReader(msg.Data)
	if err != nil {
		log.Fatal(err)
	}
	j, err := avro.AvroToByteString()
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
	avro, err := avro.NewAvroReader(msg.Data)
	if err != nil {
		log.Fatal(err)
	}
	j, err := avro.AvroToByteString()
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

// HandleConfigRequests registeres for configuration requests on the nats server
func (c *Config) HandleConfigRequests() {
	if _, err := c.nats.Subscribe(fmt.Sprintf("%s.config.register", c.basename), c.configHandlerRegister); err != nil {
		log.Fatal(err)
	}
	if _, err := c.nats.Subscribe(fmt.Sprintf("%s.config.unregister", c.basename), c.configHandlerUnregister); err != nil {
		log.Fatal(err)
	}
}

// GetRegistrations gets all client registrations for a specific topic
func (c *Config) GetRegistrations(topic string) []string {
	return c.channels.Get(topic)
}
