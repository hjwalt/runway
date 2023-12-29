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

func (helper Base64Format) Mask(value []byte) ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(value)), nil
}

func (helper Base64Format) Unmask(value []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(value))
}

func Base64() Format[string] {
	return Base64Format{}
}

func Base64Mask() Mask {
	return Base64Format{}
}
