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

package avro

import (
	"bytes"

	"github.com/linkedin/goavro"
)

// Reader is the struct to handle the avro reader
type Reader struct {
	ocfr   *goavro.OCFReader
	codec  *goavro.Codec
	schema string
	data   map[string]interface{}
}

// NewReader Creates a new avro reader that takes []byte and returns a Reader object.
func NewReader(data []byte) (*Reader, error) {
	o, err := goavro.NewOCFReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	readCodec, err := goavro.NewCodec(o.Codec().Schema())
	if err != nil {
		return nil, err
	}
	return &Reader{
		ocfr:   o,
		codec:  readCodec,
		schema: o.Codec().Schema(),
	}, nil
}

// JSON returns a string that contains a JSON of the read data
func (a *Reader) JSON() (string, error) {
	bytes, err := a.ByteString()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ByteString returns a slice of bytes containing the JSON string
func (a *Reader) ByteString() ([]byte, error) {
	if a.data == nil {
		var err error
		a.data, err = a.Map()
		if err != nil {
			return nil, err
		}
	}
	jbytes, err := a.codec.TextualFromNative(nil, a.data)
	if err != nil {
		return nil, err
	}
	return jbytes, nil
}

// Map returns a map containing all key value pairs
func (a *Reader) Map() (map[string]interface{}, error) {
	if a.data == nil {
		for a.ocfr.Scan() {
			datum, err := a.ocfr.Read()
			if err != nil {
				return nil, err
			}
			a.data = datum.(map[string]interface{})
		}
	}
	return a.data, nil
}
