package reflect_test

import (
	"testing"
	"time"

	"github.com/hjwalt/runway/reflect"
	"github.com/hjwalt/runway/structure"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TargetTest struct {
	Str        string
	StrPtr     *string
	Num        int
	NumPtr     *int
	Num8       int8
	Num8Ptr    *int8
	Num16      int16
	Num32      int32
	Num64      int64
	Unum       uint
	UnumPtr    *uint
	Unum8      uint8
	Unum16     uint16
	Unum32     uint32
	Unum64     uint64
	Float32    float32
	Float64    float64
	Float64Ptr *float64
	Bool       bool
	BoolPtr    *bool
	TimeProto  *timestamppb.Timestamp
	TimePtr    *time.Time
	Time       time.Time
	Existing   int64
	Slice      []int64
	hidden     string
}

func TestMapperNil(t *testing.T) {
	assert := assert.New(t)

	mapper := reflect.NewMapper()

	source := map[string]interface{}{
		"str":        nil,
		"str_ptr":    nil,
		"num":        nil,
		"num8":       nil,
		"num16":      nil,
		"num32":      nil,
		"num64":      nil,
		"unum":       nil,
		"unum8":      nil,
		"unum16":     nil,
		"unum32":     nil,
		"unum64":     nil,
		"float32":    nil,
		"float64":    nil,
		"bool":       nil,
		"time_proto": nil,
		"time_ptr":   nil,
		"time":       nil,
		"slice":      nil,
	}

	target := &TargetTest{
		Existing: 100,
	}

	mapper.Set(target, source)

	assert.Equal("", target.Str)
	assert.Nil(target.StrPtr)
	assert.Equal(int(0), target.Num)
	assert.Nil(target.NumPtr)
	assert.Equal(int8(0), target.Num8)
	assert.Equal(int16(0), target.Num16)
	assert.Equal(int32(0), target.Num32)
	assert.Equal(int64(0), target.Num64)
	assert.Equal(uint(0), target.Unum)
	assert.Nil(target.UnumPtr)
	assert.Equal(uint8(0), target.Unum8)
	assert.Equal(uint16(0), target.Unum16)
	assert.Equal(uint32(0), target.Unum32)
	assert.Equal(uint64(0), target.Unum64)
	assert.Equal(float32(0), target.Float32)
	assert.Equal(float64(0), target.Float64)
	assert.Nil(target.Float64Ptr)
	assert.Equal(false, target.Bool)
	assert.Nil(target.BoolPtr)
	assert.Equal(time.Time{}, target.Time)
	assert.Nil(target.TimePtr)
	assert.Nil(target.TimeProto)
	assert.Equal(0, len(target.Slice))
	assert.Equal(int64(0), target.Existing) // TODO: existing value is getting replaced without any value attached
}

func TestMapperMappingPointer(t *testing.T) {
	assert := assert.New(t)

	mapper := reflect.NewMapper()

	strPtr := "str_ptr"
	numPtr := int(0)
	unumPtr := uint(0)
	floatPtr := float64(64)
	boolPtr := true
	timePtr := time.UnixMilli(1648888888888).UTC()

	source := map[string]interface{}{
		"str":        "str",
		"str_ptr":    "str_ptr",
		"num":        0,
		"num_ptr":    0,
		"num8":       int8(8),
		"num16":      int16(16),
		"num32":      int32(32),
		"num64":      int64(64),
		"unum":       uint(0),
		"unum_ptr":   uint(0),
		"unum8":      uint8(8),
		"unum16":     uint16(16),
		"unum32":     uint32(32),
		"unum64":     uint64(64),
		"float32":    float32(32),
		"float64":    float64(64),
		"float64ptr": float64(64),
		"time_proto": "2022-04-02T08:41:28.888Z",
		"time_ptr":   "2022-04-02T08:41:28.888Z",
		"time":       "2022-04-02T08:41:28.888Z",
		"bool":       true,
		"bool_ptr":   true,
		"slice":      []int64{1},
		"hidden":     "test",
	}

	target := &TargetTest{
		Existing: 100,
	}

	mapper.Set(target, source)

	assert.Equal("str", target.Str)
	assert.Equal(&strPtr, target.StrPtr)
	assert.Equal(int(0), target.Num)
	assert.Equal(&numPtr, target.NumPtr)
	assert.Equal(int8(8), target.Num8)
	assert.Equal(int16(16), target.Num16)
	assert.Equal(int32(32), target.Num32)
	assert.Equal(int64(64), target.Num64)
	assert.Equal(uint(0), target.Unum)
	assert.Equal(&unumPtr, target.UnumPtr)
	assert.Equal(uint8(8), target.Unum8)
	assert.Equal(uint16(16), target.Unum16)
	assert.Equal(uint32(32), target.Unum32)
	assert.Equal(uint64(64), target.Unum64)
	assert.Equal(float32(32), target.Float32)
	assert.Equal(float64(64), target.Float64)
	assert.Equal(&floatPtr, target.Float64Ptr)
	assert.Equal(true, target.Bool)
	assert.Equal(&boolPtr, target.BoolPtr)
	assert.Equal(time.UnixMilli(1648888888888).UTC(), target.Time)
	assert.Equal(&timePtr, target.TimePtr)
	assert.Equal(&timestamppb.Timestamp{Seconds: 1648888888, Nanos: 888000000}, target.TimeProto)
	assert.Equal(1, len(target.Slice))
	assert.Equal([]int64{1}, target.Slice)
	assert.Equal(int64(0), target.Existing) // TODO: existing value is getting replaced without any value attached
	assert.Equal("", target.hidden)
}

func TestMapperMapping(t *testing.T) {
	assert := assert.New(t)

	mapper := reflect.NewMapper()

	source := map[string]interface{}{
		"str": "str",
	}

	target := TargetTest{}

	target = mapper.Set(target, source).(TargetTest)

	assert.Equal("str", target.Str)
}

func TestMapperNilTarget(t *testing.T) {
	assert := assert.New(t)

	mapper := reflect.NewMapper()

	source := map[string]interface{}{
		"str":        nil,
		"str_ptr":    nil,
		"num":        nil,
		"num8":       nil,
		"num16":      nil,
		"num32":      nil,
		"num64":      nil,
		"unum":       nil,
		"unum8":      nil,
		"unum16":     nil,
		"unum32":     nil,
		"unum64":     nil,
		"float32":    nil,
		"float64":    nil,
		"bool":       nil,
		"time_proto": nil,
		"time":       nil,
		"slice":      nil,
	}

	res := mapper.Set(nil, source)
	assert.Nil(res)

	var target *TargetTest

	res2 := mapper.Set(target, source)
	assert.Nil(res2)
}

func TestMapperMappingCustomFieldSearch(t *testing.T) {
	assert := assert.New(t)

	mapper := reflect.NewMapper(
		reflect.WithMapperFieldSearch(func(fieldName string) []string {
			defaults := reflect.DefaultStringSearch(fieldName)
			if fieldName == "Str" {
				defaults = append(defaults, "str_test")
			}
			return defaults
		}),
	)

	source := map[string]interface{}{
		"str_test": "str",
	}

	target := TargetTest{}

	target = mapper.Set(target, source).(TargetTest)

	assert.Equal("str", target.Str)
}

func TestMapperMappingVoid(t *testing.T) {
	assert := assert.New(t)

	mapper := reflect.NewMapper()

	source := map[string]interface{}{
		"str": "str",
	}

	var target structure.Void

	newval := mapper.Set(target, source)

	assert.Nil(newval)
	assert.Nil(target)
}
