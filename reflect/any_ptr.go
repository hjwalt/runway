package reflect

import "reflect"

// returns indirected value (or the same value if it is not a pointer), value is not default value (nil)
func GetValue(input any) (any, bool) {
	if input == nil {
		return input, false
	}

	if !IsPointer(input) {
		return input, true
	}

	reflectedValue := reflect.Indirect(reflect.ValueOf(input))
	if !reflectedValue.IsValid() {
		return input, false
	}

	return reflectedValue.Interface(), true
}
