package main

import "encoding/json"

// BytesToStruct unmarshals JSON bytes into a struct
func BytesToStruct(jsonBytes []byte, structData interface{}) error {
	return json.Unmarshal(jsonBytes, structData)
}

// MapToStruct converts a map to a struct using JSON marshaling
func MapToStruct(data map[string]interface{}, obj interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, obj)
}

// StructToJSONString converts a struct to a JSON string
func StructToJSONString(structData interface{}) (string, error) {
	jsonBytes, err := json.Marshal(structData)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// InterfaceToStruct converts any interface to a struct using JSON marshaling
func InterfaceToStruct(source interface{}, des interface{}) error {
	jsonBytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, des)
}
