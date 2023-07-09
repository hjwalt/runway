package format

import (
	"github.com/hjwalt/runway/reflect"
	"google.golang.org/protobuf/proto"
)

type ProtobufFormat[T proto.Message] struct{}

func (helper ProtobufFormat[T]) Default() T {
	return reflect.Construct[T]()
}

func (helper ProtobufFormat[T]) Marshal(value T) ([]byte, error) {
	return proto.Marshal(value)
}

func (helper ProtobufFormat[T]) Unmarshal(value []byte) (T, error) {
	protoMessage := helper.Default()
	err := proto.Unmarshal(value, protoMessage)
	if err != nil {
		return helper.Default(), err
	}
	return protoMessage, nil
}

func Protobuf[T proto.Message]() Format[T] {
	return ProtobufFormat[T]{}
}
