package config

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// Channels type to store all channels
type Channels struct {
	// // pubChannels are channels that contain data other applications are interested in
	// pubChannels map[string][]string
	// subChannels are channels other application stream data into this application
	subChannels map[string][]string
	// basename is the name of this application to prefix on created channels
	basename string

	mutex sync.Mutex
}

// NewChannels function to create a new channel managment
func NewChannels(basename string) Channels {
	return Channels{
		// pubChannels: make(map[string][]string),
		subChannels: make(map[string][]string),
		basename:    basename,
	}
}

func (c *Channels) checkRegistration(topic string) error {
	// TODO: check if channel registration is allowed
	// Possible errors are: the application does not provide the wanted subject
	// For now, no error is output as a fix value which indicates that the operation is allowed
	return nil
}

// RegisterSub function to register to a MQTT topic to get a mapped nats topic
func (c *Channels) RegisterSub(topic string) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	newChannel := fmt.Sprintf("%s.%s", c.basename, uuid.New().String())
	err := c.checkRegistration(topic)
	if err != nil {
		return "", err
	}
	c.subChannels[topic] = append(c.subChannels[topic], newChannel)
	return newChannel, nil
}

// Get function to get all registered nats subjects for a specific MQTT topic
func (c *Channels) Get(topic string) []string {
	return c.subChannels[topic]
}

// UnregisterSub function to unregister a specific nats subject
func (c *Channels) UnregisterSub(subject string) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	index := -1
	topic := ""
	for t := range c.subChannels {
		for i := range c.subChannels[t] {
			if c.subChannels[t][i] == subject {
				index = i
				topic = t
				break
			}
		}
	}
	if index != -1 {
		c.subChannels[topic] = removeFromSlice(c.subChannels[topic], index)
		return topic, nil
	}
	return "", fmt.Errorf("mapped subject was not registered at all")
}

func removeFromSlice(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
