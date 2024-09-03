//go:build !solution

package reversemap

import "reflect"

func ReverseMap(forward interface{}) interface{} {
	forwardMap := reflect.ValueOf(forward)
	reverseMap := reflect.MakeMap(reflect.MapOf(forwardMap.Type().Elem(), forwardMap.Type().Key()))

	for _, key := range forwardMap.MapKeys() {
		value := forwardMap.MapIndex(key)
		reverseMap.SetMapIndex(value, key)
	}

	return reverseMap.Interface()
}
