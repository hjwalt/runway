package reflect_test

import (
	"encoding/binary"
	"testing"

	"github.com/hjwalt/runway/reflect"
	"github.com/stretchr/testify/assert"
)

func TestEndian(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(binary.LittleEndian, reflect.Endian())
}
