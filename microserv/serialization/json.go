package serialization

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
)

func SerializeToJsonReader(model interface{}) (io.Reader, error) {
	var jsonBytes, err = json.Marshal(model)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(jsonBytes), nil
}

func SerializeToJsonWriter(model interface{}, writer io.Writer) error {
	return json.NewEncoder(writer).Encode(model)
}

func DeserializeFromJsonReader(reader io.Reader, valueType reflect.Type) (interface{}, error) {
	var modelRef = reflect.New(valueType).Elem().Interface()
	err := json.NewDecoder(reader).Decode(&modelRef)
	if err != nil {
		if err == io.EOF {
			err = nil
		} else {
			return modelRef, err
		}
	}

	return modelRef, err
}
