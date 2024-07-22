package slice_utils

import (
	"errors"
	"reflect"
)

type conditionWithIndex[T any] func(val T, ind int) bool

func Contains[T any](arr []T, val any) bool {
	for _, el := range arr {
		reflected_el := reflect.ValueOf(el)
		reflected_val := reflect.ValueOf(val)
		if reflected_el.Interface() == reflected_val.Interface() {
			return true
		}
	}
	return false
}

func Find[T any](arr []T, cond func(T) bool) (interface{}, error) {
	for i := 0; i < len(arr); i++ {
		el := arr[i]
		if cond(el) {
			return el, nil
		}
	}
	val := new(T)
	return *val, errors.New("Element not found!")
}

func Filter[T any](arr []T, cond conditionWithIndex[T]) []T {
	var filtered []T
	for i, el := range arr {
		needPush := cond(el, i)
		if needPush {
			filtered = append(filtered, el)
		}
	}
	return filtered
}
