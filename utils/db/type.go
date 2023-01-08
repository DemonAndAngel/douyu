package db

import (
	"database/sql/driver"
	"douyu/utils/helpers"
	"encoding/json"
	"errors"
	"fmt"
)

type Uint64ArrayStruct []uint64

func (array Uint64ArrayStruct) Value() (driver.Value, error) {
	str := helpers.JsonMarshal(array)
	if str == "" {
		str = "[]"
	}
	return []byte(str), nil
}
func (array *Uint64ArrayStruct) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return json.Unmarshal(b, array)
}

type StringStringArrayStruct [][]string

func (array StringStringArrayStruct) Value() (driver.Value, error) {
	str := helpers.JsonMarshal(array)
	if str == "" {
		str = "[]"
	}
	return []byte(str), nil
}
func (array *StringStringArrayStruct) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return json.Unmarshal(b, array)
}

type StringArrayStruct []string

func (array StringArrayStruct) Value() (driver.Value, error) {
	str := helpers.JsonMarshal(array)
	if str == "" {
		str = "[]"
	}
	return []byte(str), nil
}
func (array *StringArrayStruct) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return json.Unmarshal(b, array)
}

type ArrayMapStruct []map[string]interface{}

func (data ArrayMapStruct) Value() (driver.Value, error) {
	str := helpers.JsonMarshal(data)
	if str == "" {
		str = "{}"
	}
	return []byte(str), nil
}

func (data *ArrayMapStruct) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return json.Unmarshal(b, data)
}

func GetArrayMapStructWithStr(str string) ArrayMapStruct {
	data := ArrayMapStruct{}
	json.Unmarshal([]byte(str), &data)
	return data
}

type MapStruct map[string]interface{}

func (data MapStruct) Value() (driver.Value, error) {
	str := helpers.JsonMarshal(data)
	if str == "" {
		str = "{}"
	}
	return []byte(str), nil
}

func (data *MapStruct) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return json.Unmarshal(b, data)
}
func GetMapStructWithStr(str string) MapStruct {
	data := MapStruct{}
	_ = json.Unmarshal([]byte(str), &data)
	return data
}
func GetStringArrayStructWithStr(str string) StringArrayStruct {
	data := StringArrayStruct{}
	_ = json.Unmarshal([]byte(str), &data)
	return data
}
func GetStringStringArrayStructWithStr(str string) StringStringArrayStruct {
	data := StringStringArrayStruct{}
	_ = json.Unmarshal([]byte(str), &data)
	return data
}
