package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func UnmarshalOmitEmpty(b []byte, v any) error {
	err := json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	reflected := reflect.Indirect(reflect.ValueOf(v))
	for i := 0; i < reflected.NumField(); i++ {
		field := reflected.Field(i)
		if field.IsZero() {
			return fmt.Errorf("\"%s\" can't be empty", reflected.Type().Field(i).Name)
		}
	}
	return nil
}
