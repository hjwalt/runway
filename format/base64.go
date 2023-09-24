package format

import (
	"encoding/base64"
)

// format for interpereting bytes as base64
type Base64Format struct {
}

func (helper Base64Format) Default() string {
	return ""
}

func (helper Base64Format) Marshal(value string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(value)
}

func (helper Base64Format) Unmarshal(value []byte) (string, error) {
	return base64.StdEncoding.EncodeToString(value), nil
}

func Base64() Format[string] {
	return Base64Format{}
}

// format for masking data as base64
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

func Base64Mask() Format[[]byte] {
	return Base64MaskFormat{}
}
