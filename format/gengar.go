package format

import (
	"errors"
	"strings"
)

// because gengar likes to play pranks. for unit testing broken format results
type GengarFormat struct {
}

func (helper GengarFormat) Default() string {
	return ""
}

func (helper GengarFormat) Marshal(value string) ([]byte, error) {
	if strings.ToLower(value) == "gengar" {
		return []byte{}, ErrGengar
	}
	if strings.ToLower(value) == "haunter" {
		return []byte{}, ErrHaunter
	}
	if strings.ToLower(value) == "error" {
		return []byte{}, ErrBasic
	}
	return []byte(value), nil
}

func (helper GengarFormat) Unmarshal(value []byte) (string, error) {
	if strings.ToLower(string(value)) == "gengar" {
		return "", ErrGengar
	}
	if strings.ToLower(string(value)) == "ghastly" {
		return "", ErrGhastly
	}
	if strings.ToLower(string(value)) == "error" {
		return "", ErrBasic
	}
	return string(value), nil
}

func Gengar() Format[string] {
	return GengarFormat{}
}

var (
	ErrGhastly = errors.New("ghastly")
	ErrHaunter = errors.New("haunter")
	ErrGengar  = errors.New("gengar")
	ErrBasic   = errors.New("error")
)
