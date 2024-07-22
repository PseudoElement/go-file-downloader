package validators_module

import "reflect"

/* Struct has "" or nil field - returns true and empty field name, else returns false and "" */
func HasEmptyField[T any](obj T) (bool, string) {
	v := reflect.ValueOf(obj)

	for i := 0; i < v.NumField(); i++ {
		key := v.Type().Field(i).Name
		value := v.Field(i).Interface()
		if value == nil || value == "" {
			return true, key
		}
	}

	return false, ""
}
