package reflect

import "reflect"

func IsPointer(input any) bool {
	if input == nil {
		return true
	}
	if reflect.TypeOf(input).Kind() == reflect.Pointer {
		return true
	}
	return false
}
