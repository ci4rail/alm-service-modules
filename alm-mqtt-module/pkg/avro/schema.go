package avro

import (
	"log"

	"github.com/linkedin/goavro"
)

// CreateSchema creates an avro codec from an avro schema text file
func CreateSchema(schema string) *goavro.Codec {
	codec, err := goavro.NewCodec(schema)
	if err != nil {
		log.Fatal(err)
	}
	return codec
}
