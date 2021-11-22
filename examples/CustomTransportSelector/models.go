package main

import "reflect"

type TestRequestModel struct {
	Foo  string `json:"foo" param:"foo"`
	Bar  int    `json:"bar" param:"bar"`
	Zulu bool   `json:"zulu" param:"zulu"`
}

var TestRequestType = reflect.TypeOf(TestRequestModel{})

type TestResponseModel struct {
	Service string      `json:"service"`
	Result  interface{} `json:"result,omitempty"`
}

var TestResponseType = reflect.TypeOf(TestResponseModel{})
