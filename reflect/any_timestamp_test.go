package reflect_test

import (
	"testing"
	"time"

	"github.com/hjwalt/runway/reflect"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTimeToTime(t *testing.T) {
	assert := assert.New(t)

	timeTotest := time.UnixMilli(1648888888888).UTC()
	timeProto := &timestamppb.Timestamp{Seconds: 1648888888, Nanos: 888000000}

	timeTotestWithoutMilli := time.UnixMilli(1648888888000).UTC()

	assert.Equal(&timeTotest, reflect.GetTimestampMs(timeTotest))
	assert.Equal(&timeTotest, reflect.GetTimestampMs(&timeTotest))
	assert.Equal(&timeTotest, reflect.GetTimestampMs(timeProto))
	assert.Equal(&timeTotest, reflect.GetTimestampMs(1648888888888))
	assert.Equal(&timeTotest, reflect.GetTimestampMs("2022-04-02T08:41:28.888Z"))
	assert.Equal(&timeTotest, reflect.GetTimestampMs("2022-04-02T15:41:28.888+07:00"))
	assert.Equal(&timeTotestWithoutMilli, reflect.GetTimestampMs("2022-04-02T08:41:28Z"))

	var nilTime *time.Time

	assert.Nil(reflect.GetTimestampMs("2022-04-02T15:41:28.888+07.00"))
	assert.Nil(reflect.GetTimestampMs(StructForTest{Message: "test"}))
	assert.Nil(reflect.GetTimestampMs(nilTime))
	assert.Nil(reflect.GetTimestampMs(nil))
}

func TestTimeToProtobuf(t *testing.T) {
	assert := assert.New(t)

	timeTotest := time.UnixMilli(1648888888888).UTC()
	timeProto := &timestamppb.Timestamp{Seconds: 1648888888, Nanos: 888000000}

	assert.Equal(timeProto, reflect.GetProtoTimestampMs(timeTotest))
	assert.Equal(timeProto, reflect.GetProtoTimestampMs(timeProto))
	assert.Nil(reflect.GetProtoTimestampMs(nil))
}
