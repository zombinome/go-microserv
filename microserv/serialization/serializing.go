package serialization

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
)

var ErrUnsupportedFieldType = errors.New("unsupported field type")

func WriteToQueryParameters(model interface{}) (url.Values, error) {
	if model == nil {
		return nil, nil
	}

	return serializeToParams(reflect.ValueOf(model), reflect.TypeOf(model), tagParam)
}

func WriteFormParameters(model interface{}) (map[string][]string, error) {
	if model == nil {
		return nil, nil
	}

	return serializeToParams(reflect.ValueOf(model), reflect.TypeOf(model), tagForm)
}

func serializeToParams(modelRef reflect.Value, modelType reflect.Type, tag string) (map[string][]string, error) {
	var modelVal = modelRef
	if modelRef.Kind() == reflect.Ptr || modelRef.Kind() == reflect.Interface {
		modelVal = modelRef.Elem() // Ptr -> Interface
	}
	var fieldCount = modelVal.NumField()
	var result = make(map[string][]string)
	for i := 0; i < fieldCount; i++ {

		var paramName = modelType.Field(i).Tag.Get(tag)
		if paramName == "" {
			continue
		}

		var fieldRef = modelVal.Field(i)
		result[paramName] = convertFieldValueToQueryParamValue(fieldRef)
	}
	return result, nil
}

func convertFieldValueToQueryParamValue(fieldRef reflect.Value) []string {
	if fieldRef.Kind() == reflect.Array {
		var length = fieldRef.Len()
		var result = make([]string, length)
		for i := 0; i < length; i++ {
			result[i], _ = convertFieldValueToQueryString(fieldRef.Index(i))

		}
		return result
	}

	var strValue, err = convertFieldValueToQueryString(fieldRef)
	if err != nil {
		// TODO: Add logging warning here

	}
	return []string{strValue}
}

func convertFieldValueToQueryString(fieldRef reflect.Value) (string, error) {
	if fieldRef.IsNil() {
		return "", nil
	}

	switch fieldRef.Kind() {
	case reflect.Bool:
		if fieldRef.Bool() {
			return "true", nil
		} else {
			return "false", nil
		}

	case reflect.String:
		return fieldRef.String(), nil

	case reflect.Float32:
	case reflect.Float64:
		return strconv.FormatFloat(fieldRef.Float(), 'f', -1, 64), nil

	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
		return strconv.FormatInt(fieldRef.Int(), 10), nil

	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
		return strconv.FormatUint(fieldRef.Uint(), 10), nil

	case reflect.Struct:
		return fieldRef.String(), nil // TODO: Add serialization of date and time
	}

	return "", ErrUnsupportedFieldType
}
