package format

import (
	"errors"
)

// assuming byte compatibility, i.e. bytes <-> proto, string <-> json
func Convert[V1 any, V2 any](
	v V1,
	v1 Format[V1],
	v2 Format[V2],
) (V2, error) {

	// serialise value
	valueBytes, marshalErr := v1.Marshal(v)
	if marshalErr != nil {
		return v2.Default(), errors.Join(ErrFormatConversionMarshal, marshalErr)
	}

	// deserialise value
	nextValue, unmarshallErr := v2.Unmarshal(valueBytes)
	if unmarshallErr != nil {
		return v2.Default(), errors.Join(ErrFormatConversionUnmarshal, unmarshallErr)
	}

	return nextValue, nil
}

var (
	ErrFormatConversionMarshal   = errors.New("format conversion marshal error")
	ErrFormatConversionUnmarshal = errors.New("format conversion unmarshal error")
)
