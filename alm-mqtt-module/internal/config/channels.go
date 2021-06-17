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
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// Channels type to store all channels
type Channels struct {
	// subChannels are channels other application stream data into this application
	subChannels map[string][]string
	// basename is the name of this application to prefix on created channels
	basename string

	mutex sync.Mutex
}

// NewChannels function to create a new channel managment
func NewChannels(basename string) Channels {
	return Channels{
		subChannels: make(map[string][]string),
		basename:    basename,
	}
}

// RegisterSub function to register to a MQTT topic to get a mapped nats subject
func (c *Channels) RegisterSub(topic string) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	newChannel := fmt.Sprintf("%s.%s", c.basename, uuid.New().String())
	c.subChannels[topic] = append(c.subChannels[topic], newChannel)
	return newChannel, nil
}

// Get function to get all registered nats subjects for a specific MQTT topic
func (c *Channels) Get(topic string) []string {
	return c.subChannels[topic]
}

// GetTopic function to get the corresponding MQTT topic to a nats subscription
func (c *Channels) GetTopic(subject string) (string, error) {
	for t := range c.subChannels {
		for _, s := range c.subChannels[t] {
			if s == subject {
				return t, nil
			}
		}
	}
	return "", fmt.Errorf("no topic found for subject '%s'", subject)
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
		c.subChannels[topic] = removeFromStringSlice(c.subChannels[topic], index)
		return topic, nil
	}
	return "", fmt.Errorf("mapped subject was not registered at all")
}

func removeFromStringSlice(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
