package format

import (
	"encoding/base64"
)

type Base64MaskFormat struct {
}

func (helper Base64MaskFormat) Default() []byte {
	return []byte{}
}

func (helper Base64MaskFormat) Marshal(value []byte) ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(value)), nil
}

func (helper Base64MaskFormat) Unmarshal(value []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(value))
}

func Base64() Format[[]byte] {
	return Base64MaskFormat{}
}

func Base64String() Format[string] {
	return Masked(
		Base64(),
		String(),
	)
}
