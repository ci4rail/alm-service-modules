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
	registered = channels.Get("topic1")
	assert.NotContains(registered, ch1_1)
	assert.Contains(registered, ch1_2)
	assert.Nil(err)
	assert.Equal(topic, "topic1")
	topic, err = channels.UnregisterSub(ch1_1)
	assert.Equal(topic, "")
	assert.NotNil(err)
	registered = channels.Get("topic1")
	assert.NotContains(registered, ch1_1)
	assert.Contains(registered, ch1_2)
}
