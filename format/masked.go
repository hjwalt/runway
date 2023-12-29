package format

import (
	"errors"
)

type MaskedFormat[T any] struct {
	mask   Mask
	actual Format[T]
}

func (helper MaskedFormat[T]) Default() T {
	return helper.actual.Default()
}

func (helper MaskedFormat[T]) Marshal(value T) ([]byte, error) {
	byteToMask, marshalErr := helper.actual.Marshal(value)
	if marshalErr != nil {
		return []byte{}, errors.Join(ErrMaskActualMarshal, marshalErr)
	}

	maskedBytes, maskingErr := helper.mask.Mask(byteToMask)
	if maskingErr != nil {
		return []byte{}, errors.Join(ErrMaskMarshal, maskingErr)
	}

	return maskedBytes, nil
}

func (helper MaskedFormat[T]) Unmarshal(value []byte) (T, error) {
	bytesToUnmarshal, unmaskingErr := helper.mask.Unmask(value)
	if unmaskingErr != nil {
		return helper.Default(), errors.Join(ErrMaskUnmarshal, unmaskingErr)
	}

	valueToReturn, unmarshalErr := helper.actual.Unmarshal(bytesToUnmarshal)
	if unmarshalErr != nil {
		return helper.Default(), errors.Join(ErrMaskActualUnmarshal, unmarshalErr)
	}

	return valueToReturn, nil
}

func Masked[T any](mask Mask, actual Format[T]) Format[T] {
	return MaskedFormat[T]{mask: mask, actual: actual}
}

var (
	ErrMaskActualMarshal   = errors.New("masking error while converting type to bytes")
	ErrMaskMarshal         = errors.New("masking error while masking bytes")
	ErrMaskActualUnmarshal = errors.New("masking error while converting bytes to type")
	ErrMaskUnmarshal       = errors.New("masking error while unmasking bytes")
)
