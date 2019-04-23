package core

import (
	"bytes"
	"encoding/gob"
)

func serializeType(value interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func deserializeType(data []byte, val interface{}) error {
	buf := bytes.NewReader(data)
	encoder := gob.NewDecoder(buf)
	return encoder.Decode(val)
}
