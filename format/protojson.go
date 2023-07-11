package format

import (
	"github.com/hjwalt/runway/reflect"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type ProtojsonFormat[T proto.Message] struct{}

func (helper ProtojsonFormat[T]) Default() T {
	return reflect.Construct[T]()
}

func (helper ProtojsonFormat[T]) Marshal(value T) ([]byte, error) {
	marshaller := protojson.MarshalOptions{
		Multiline: false,
		Indent:    "",
	}
	return marshaller.Marshal(value)
}

func (helper ProtojsonFormat[T]) Unmarshal(value []byte) (T, error) {
	if len(value) == 0 {
		return helper.Default(), nil
	}
	protoMessage := helper.Default()
	unmarshaller := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
	err := unmarshaller.Unmarshal(value, protoMessage)
	if err != nil {
		return helper.Default(), err
	}
	return protoMessage, nil
}

func Protojson[T proto.Message]() Format[T] {
	return ProtojsonFormat[T]{}
}
