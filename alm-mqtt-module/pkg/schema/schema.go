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
