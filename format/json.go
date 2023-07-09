package format

import (
	"encoding/json"

	reflect "github.com/hjwalt/runway/reflect"
)

type JsonFormat[T any] struct{}

func (helper JsonFormat[T]) Default() T {
	return reflect.Construct[T]()
}

func (helper JsonFormat[T]) Marshal(value T) ([]byte, error) {
	if reflect.IsNil(value) {
		return []byte{}, nil
	}
	if reflect.IsPointer(value) {
		return json.Marshal(value)
	} else {
		return json.Marshal(&value)
	}
}

func (helper JsonFormat[T]) Unmarshal(value []byte) (T, error) {
	if len(value) == 0 {
		return helper.Default(), nil
	}
	jsonMessage := helper.Default()
	if reflect.IsPointer(jsonMessage) {
		err := json.Unmarshal(value, jsonMessage)
		return jsonMessage, err
	} else {
		err := json.Unmarshal(value, &jsonMessage)
		return jsonMessage, err
	}
}

func Json[T any]() Format[T] {
	return JsonFormat[T]{}
}
