package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	assert := assert.New(t)
	channels := NewChannels("module1")
	ch1_1, err := channels.RegisterSub("topic1")
	assert.Nil(err)
	ch1_2, err := channels.RegisterSub("topic1")
	assert.Nil(err)
	ch2_1, err := channels.RegisterSub("topic2")
	assert.Nil(err)
	registered := channels.Get("topic1")
	assert.Contains(registered, ch1_1)
	assert.Contains(registered, ch1_2)
	registered = channels.Get("topic2")
	assert.Contains(registered, ch2_1)
}

func TestUnRegister(t *testing.T) {
	assert := assert.New(t)
	channels := NewChannels("module1")
	ch1_1, err := channels.RegisterSub("topic1")
	assert.Nil(err)
	ch1_2, err := channels.RegisterSub("topic1")
	assert.Nil(err)
	ch2_1, err := channels.RegisterSub("topic2")
	assert.Nil(err)
	registered := channels.Get("topic1")
	assert.Contains(registered, ch1_1)
	assert.Contains(registered, ch1_2)
	registered = channels.Get("topic2")
	assert.Contains(registered, ch2_1)

	topic, err := channels.UnregisterSub(ch1_1)
	assert.Nil(err)
	assert.Equal(topic, "topic1")
	topic, err = channels.UnregisterSub(ch1_1)
	assert.Equal(topic, "")
	assert.NotNil(err)
	registered = channels.Get("topic1")
	assert.NotContains(registered, ch1_1)
	assert.Contains(registered, ch1_2)
}
