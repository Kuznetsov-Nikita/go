//go:build !solution

package jsonlist

import (
	"encoding/json"
	"io"
	"reflect"
)

func Marshal(w io.Writer, slice interface{}) error {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		return &json.UnsupportedTypeError{Type: sliceValue.Type()}
	}

	encoder := json.NewEncoder(w)
	for i := 0; i < sliceValue.Len(); i++ {
		if err := encoder.Encode(sliceValue.Index(i).Interface()); err != nil {
			return err
		}
	}

	return nil
}

func Unmarshal(r io.Reader, slice interface{}) error {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Ptr || sliceValue.Elem().Kind() != reflect.Slice {
		return &json.UnsupportedTypeError{Type: sliceValue.Type()}
	}

	decoder := json.NewDecoder(r)
	for {
		elem := reflect.New(sliceValue.Elem().Type().Elem()).Interface()
		if err := decoder.Decode(elem); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		sliceValue.Elem().Set(reflect.Append(sliceValue.Elem(), reflect.ValueOf(elem).Elem()))
	}

	return nil
}
