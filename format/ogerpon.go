package format

import (
	"errors"
	"strings"
)

// ogerpon the teal mask pokemon. for unit testing broken format masking
type OgerponFormat struct {
}

func (helper OgerponFormat) Mask(value []byte) ([]byte, error) {
	if strings.ToLower(string(value)) == "wellspring" {
		return []byte{}, ErrWellspringMask
	}
	if strings.ToLower(string(value)) == "error" {
		return []byte{}, ErrBasic
	}
	return reverse(value), nil
}

func (helper OgerponFormat) Unmask(value []byte) ([]byte, error) {
	reversed := reverse(value)
	if strings.ToLower(string(reversed)) == "hearthflame" {
		return []byte{}, ErrHearthflameMask
	}
	if strings.ToLower(string(reversed)) == "error" {
		return []byte{}, ErrBasic
	}
	return reversed, nil
}

func Ogerpon() Mask {
	return OgerponFormat{}
}

var (
	ErrWellspringMask  = errors.New("wellspring")
	ErrHearthflameMask = errors.New("hearthflame")
)

func reverse(s []byte) []byte {
	reversed := make([]byte, len(s))
	slen := len(s)
	for i, val := range s {
		reversed[slen-i-1] = val
	}
	return reversed
}
