package reflect

import (
	"encoding/binary"

	"github.com/hjwalt/runway/trusted"
)

func Endian() binary.ByteOrder {
	return trusted.Endian()
}
