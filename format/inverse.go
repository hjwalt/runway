package format

import (
	"context"

	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/reflect"
)

func RegisterFormat[T any](c inverse.Container, f Format[T]) {
	inverse.GenericAddVal(c, "format-"+reflect.TypeName(reflect.Construct[T]()), f)
}

func RetrieveFormat[T any](c inverse.Container, ctx context.Context) (Format[T], error) {
	return inverse.GenericGet[Format[T]](c, ctx, "format-"+reflect.TypeName(reflect.Construct[T]()))
}

func RegisterMask(c inverse.Container, name string, m Mask) {
	inverse.GenericAddVal(c, "format-"+name, m)
}

func RetrieveMask(c inverse.Container, name string, ctx context.Context) (Mask, error) {
	return inverse.GenericGet[Mask](c, ctx, "format-"+name)
}
