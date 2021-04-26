package avro

import (
	"bytes"

	"github.com/linkedin/goavro"
)

func NewAvroWriter(msg map[string]interface{}, codec *goavro.Codec) ([]byte, error) {
	bin := new(bytes.Buffer)

	ocfw, err := goavro.NewOCFWriter(goavro.OCFConfig{
		W:     bin,
		Codec: codec,
	})
	if err != nil {
		return nil, err
	}

	err = ocfw.Append([]interface{}{msg})
	if err != nil {
		return nil, err
	}
	return bin.Bytes(), nil
}
