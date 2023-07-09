package format

import (
	"github.com/hjwalt/runway/logger"
)

// assuming byte compatibility, i.e. bytes <-> proto, string <-> json
func Convert[V1 any, V2 any](
	v V1,
	v1 Format[V1],
	v2 Format[V2],
) (V2, error) {

	// serialise value
	valueBytes, err := v1.Marshal(v)
	if err != nil {
		logger.ErrorErr("conversion value serialisation failure", err)
		return v2.Default(), err
	}

	// deserialise value
	return v2.Unmarshal(valueBytes)
}
