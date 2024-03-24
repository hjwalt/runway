package reflect

import "reflect"

func TypeName(v any) string {
	return reflect.TypeOf(v).String()
}
