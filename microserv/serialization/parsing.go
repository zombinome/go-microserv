package serialization

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strconv"
)

func ReadRequest(requestType reflect.Type, request *http.Request) (interface{}, error) {
	if requestType == nil {
		return nil, nil
	}

	var requestModelRef = reflect.New(requestType)

	var err error = nil

	// Mapping query params
	var queryParams = request.URL.Query()
	if len(queryParams) > 0 {
		fillModelFromMap(requestModelRef, queryParams)
	}

	// Mapping form values from multipart form data or from regular form data
	if request.MultipartForm != nil && len(request.MultipartForm.Value) > 0 {
		fillModelFromMap(requestModelRef, request.MultipartForm.Value)
	} else if request.Form != nil && len(request.Form) > 0 {
		fillModelFromMap(requestModelRef, request.Form)
	}

	// JSON lib quirk
	// You can't pass pointer to struct wrapped in interface{}, as in this case JSON would be
	// deserialized as map[string]interface{}
	// But you can pass pointer to struct wrapped in interface{}, which will work
	err = json.NewDecoder(request.Body).Decode(requestModelRef.Interface())
	var requestModel = requestModelRef.Elem().Interface()
	if err != nil {
		if err == io.EOF {
			err = nil
		} else {
			return requestModel, err
		}
	}

	return requestModel, err
}

func fillModelFromMap(modelRef reflect.Value, queryParams map[string][]string) error {
	var modelVal = modelRef.Elem() // Ptr -> Interface
	var modelType reflect.Type = modelVal.Type()

	var fieldCount = modelVal.NumField()
	for i := 0; i < fieldCount; i++ {

		var paramName = modelType.Field(i).Tag.Get("param")
		if paramName == "" {
			continue
		}

		if paramValues, paramPresent := queryParams[paramName]; paramPresent {
			var fieldRef = modelVal.Field(i)
			if fieldRef.CanSet() {
				if fieldRef.Kind() == reflect.Array {
					setFieldValueFromArray(&fieldRef, paramValues)
				} else if len(paramValues) > 0 {
					setFieldValueFromString(&fieldRef, fieldRef.Kind(), paramValues[0])
				}
			}
		}
	}
	return nil
}

func setFieldValueFromString(fieldRef *reflect.Value, kind reflect.Kind, value string) {
	switch kind {
	case reflect.String:
		fieldRef.SetString(value)

	case reflect.Bool:
		var boolVal, _ = strconv.ParseBool(value)
		fieldRef.SetBool(boolVal)

	case reflect.Int:
		var intVal, _ = strconv.ParseInt(value, 10, 64)
		fieldRef.SetInt(intVal)

	case reflect.Uint:
		var uintVal, _ = strconv.ParseUint(value, 10, 64)
		fieldRef.SetUint(uintVal)

	case reflect.Float32:
		var floatVal, _ = strconv.ParseFloat(value, 32)
		fieldRef.SetFloat(floatVal)

	case reflect.Float64:
		var floatVal, _ = strconv.ParseFloat(value, 64)
		fieldRef.SetFloat(floatVal)

	case reflect.Int8:
		var intVal, _ = strconv.ParseInt(value, 10, 8)
		fieldRef.SetInt(intVal)

	case reflect.Int16:
		var intVal, _ = strconv.ParseInt(value, 10, 16)
		fieldRef.SetInt(intVal)

	case reflect.Int32:
		var intVal, _ = strconv.ParseInt(value, 10, 32)
		fieldRef.SetInt(intVal)

	case reflect.Int64:
		var intVal, _ = strconv.ParseInt(value, 10, 64)
		fieldRef.SetInt(intVal)

	case reflect.Uint8:
		var intVal, _ = strconv.ParseUint(value, 10, 8)
		fieldRef.SetUint(intVal)

	case reflect.Uint16:
		var intVal, _ = strconv.ParseUint(value, 10, 16)
		fieldRef.SetUint(intVal)

	case reflect.Uint32:
		var intVal, _ = strconv.ParseUint(value, 10, 32)
		fieldRef.SetUint(intVal)

	case reflect.Uint64:
		var intVal, _ = strconv.ParseUint(value, 10, 64)
		fieldRef.SetUint(intVal)
	}
}

func setFieldValueFromArray(fieldRef *reflect.Value, values []string) {
	var itemKind = fieldRef.Elem().Kind()

	result := reflect.MakeSlice(reflect.SliceOf(fieldRef.Elem().Type()), len(values), cap(values))
	for i := 0; i < len(values); i++ {
		resultItemRef := result.Index(i)
		setFieldValueFromString(&resultItemRef, itemKind, values[i])
	}
}
