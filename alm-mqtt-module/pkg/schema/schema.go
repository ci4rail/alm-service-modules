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

// PubRequestType is the struct for an Publish request
type PubRequestType struct {
	Topic   string `json:"topic"`
	Payload []byte `json:"payload"`
}

// PubResponseType is the struct for an Publish response
type PubResponseType struct {
	Error string `json:"error"`
}

// ReqResRequestType is the struct for an `request respsonse` request
type ReqResRequestType struct {
	Topic   string `json:"topic"`
	Payload []byte `json:"payload"`
	Timeout int32  `json:"timeout"`
}

// ReqResResponsetType is the struct for an `request respsonse` response
type ReqResResponsetType struct {
	Payload []byte `json:"payload"`
	Error   string `json:"error"`
}

// RegisterSubRequest is the text file loaded schema for RegisterSubRequests
//go:embed avro_schemas/registerSubRequest.avsc
var RegisterSubRequest string

// RegisterSubRequestCodec is the prepared avro codec for RegisterSubRequests
var RegisterSubRequestCodec = avro.CreateSchema(RegisterSubRequest)

// RegisterSubResponse is the text file loaded schema for RegisterSubResponses
//go:embed avro_schemas/registerSubResponse.avsc
var RegisterSubResponse string

// RegisterSubResponseCodec is the prepared avro codec for RegisterSubResponses
var RegisterSubResponseCodec = avro.CreateSchema(RegisterSubResponse)

// UnregisterSubRequest is the text file loaded schema for UnregisterSubRequests
//go:embed avro_schemas/unregisterSubRequest.avsc
var UnregisterSubRequest string

// UnregisterSubRequestCodec is the prepared avro codec for UnregisterSubRequests
var UnregisterSubRequestCodec = avro.CreateSchema(UnregisterSubRequest)

// UnregisterSubResponse is the text file loaded schema for UnregisterSubResponses
//go:embed avro_schemas/unregisterSubResponse.avsc
var UnregisterSubResponse string

// UnregisterSubResponseCodec is the prepared avro codec for UnregisterSubResponses
var UnregisterSubResponseCodec = avro.CreateSchema(UnregisterSubResponse)

// PubRequest is the text file loaded schema for PublishRequests
//go:embed avro_schemas/pubRequest.avsc
var PubRequest string

// PubRequestCodec is the prepared avro codec for PubRequests
var PubRequestCodec = avro.CreateSchema(PubRequest)

// PubResponse is the text file loaded schema for PubResponses
//go:embed avro_schemas/pubResponse.avsc
var PubResponse string

// PubResponseCodec is the prepared avro codec for PubResponses
var PubResponseCodec = avro.CreateSchema(PubResponse)

// ReqResRequest is the text file loaded schema for Request Response Requests
//go:embed avro_schemas/reqResRequest.avsc
var ReqResRequest string

// ReqResRequestCodec is the prepared avro codec for Request Response Requests
var ReqResRequestCodec = avro.CreateSchema(ReqResRequest)

// ReqResResponse is the text file loaded schema for Request Response Resposes
//go:embed avro_schemas/reqResResponse.avsc
var ReqResResponse string

// ReqResResponseCodec is the prepared avro codec for Request Response Resposes
var ReqResResponseCodec = avro.CreateSchema(ReqResResponse)
