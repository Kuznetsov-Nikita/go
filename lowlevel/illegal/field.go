//go:build !solution

package illegal

import (
	"reflect"
	"unsafe"
)

func SetPrivateField(obj interface{}, name string, value interface{}) {
	v := reflect.ValueOf(obj).Elem()
	field := v.FieldByName(name)
	if field.IsValid() {
		field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
		field.Set(reflect.ValueOf(value))
	}
}
