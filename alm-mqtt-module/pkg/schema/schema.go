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

package schema

import (
	"alm-mqtt-module/pkg/avro"
	// no need for a name here
	_ "embed"
)

// RegisterSubRequestType is the struct used for a Register Subscription request
type RegisterSubRequestType struct {
	Topic string `json:"topic"`
}

// RegisterSubResponseType is the struct for a Register Subscription response
type RegisterSubResponseType struct {
	Subject string `json:"subject"`
	Error   string `json:"error"`
}

// UnregisterSubRequestType is the struct for an Unregister Subscription request
type UnregisterSubRequestType struct {
	Subject string `json:"subject"`
}

// UnregisterSubResponseType is the struct for an Unregister Subscription response
type UnregisterSubResponseType struct {
	Error string `json:"error"`
}

// SendMQTTRequestType is the struct for an Unregister Subscription request
type SendMQTTRequestType struct {
	Topic   string `json:"topic"`
	Payload byte   `json:"payload"`
}

// SendMQTTResponseType is the struct for an Unregister Subscription response
type SendMQTTResponseType struct {
	Error string `json:"error"`
}

// RegisterSubRequest is the text file loaded schema for RegisterSubRequests
//go:embed avro_schemas/registerSubRequest.avro
var RegisterSubRequest string

// RegisterSubRequestCodec is the prepared avro codec for RegisterSubRequests
var RegisterSubRequestCodec = avro.CreateSchema(RegisterSubRequest)

// RegisterSubResponse is the text file loaded schema for RegisterSubResponses
//go:embed avro_schemas/registerSubResponse.avro
var RegisterSubResponse string

// RegisterSubResponseCodec is the prepared avro codec for RegisterSubResponses
var RegisterSubResponseCodec = avro.CreateSchema(RegisterSubResponse)

// UnregisterSubRequest is the text file loaded schema for UnregisterSubRequests
//go:embed avro_schemas/unregisterSubRequest.avro
var UnregisterSubRequest string

// UnregisterSubRequestCodec is the prepared avro codec for UnregisterSubRequests
var UnregisterSubRequestCodec = avro.CreateSchema(UnregisterSubRequest)

// UnregisterSubResponse is the text file loaded schema for UnregisterSubResponses
//go:embed avro_schemas/unregisterSubResponse.avro
var UnregisterSubResponse string

// UnregisterSubResponseCodec is the prepared avro codec for UnregisterSubResponses
var UnregisterSubResponseCodec = avro.CreateSchema(UnregisterSubResponse)

// SendMQTTRequest is the text file loaded schema for SendMQTTRequests
//go:embed avro_schemas/sendMQTTRequest.avro
var SendMQTTRequest string

// SendMQTTRequestCodec is the prepared avro codec for SendMQTTRequests
var SendMQTTRequestCodec = avro.CreateSchema(SendMQTTRequest)

// SendMQTTResponse is the text file loaded schema for SendMQTTResponses
//go:embed avro_schemas/sendMQTTResponse.avro
var SendMQTTResponse string

// SendMQTTResponseCodec is the prepared avro codec for SendMQTTResponses
var SendMQTTResponseCodec = avro.CreateSchema(SendMQTTResponse)
