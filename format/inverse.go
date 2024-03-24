package format

import (
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/reflect"
)

func RegisterFormat[T any](c inverse.Container, f Format[T]) {
	inverse.GenericAddVal(c, "format-"+reflect.TypeName(f.Default()), f)
}

func RegisterMask(c inverse.Container, name string, m Mask) {
	inverse.GenericAddVal(c, "format-"+name, m)
}
