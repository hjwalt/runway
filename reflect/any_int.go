package reflect

import (
	"reflect"
	"strconv"

	"github.com/hjwalt/runway/logger"
	"go.uber.org/zap"
)

func GetIntBase(input any, bitSize int) int64 {
	if input == nil {
		return 0
	}
	// Get int from pointer
	if reflect.TypeOf(input).Kind() == reflect.Pointer {
		return GetIntBase(reflect.Indirect(reflect.ValueOf(input)).Interface(), bitSize)
	}
	var intValue int64
	switch input := input.(type) {
	case int:
		intValue = int64(input)
	case int8:
		intValue = int64(input)
	case int16:
		intValue = int64(input)
	case int32:
		intValue = int64(input)
	case int64:
		intValue = input
	case bool:
		if input {
			intValue = 1
		} else {
			intValue = 0
		}
	case string:
		if input == "" {
			input = "0"
		}
		var err error
		intValue, err = strconv.ParseInt(input, 10, bitSize)
		if err != nil {
			logger.WarnErr("string parse int failed", err)
			return 0
		}
	default:
		logger.Warn("conversion for int type failed", zap.Any("type", reflect.TypeOf(input)), zap.Any("value", input))
		return 0
	}

	return intValue
}

func GetInt(input any) int {
	return int(GetIntBase(input, 0))
}

func GetInt8(input any) int8 {
	return int8(GetIntBase(input, 8))
}

func GetInt16(input any) int16 {
	return int16(GetIntBase(input, 16))
}

func GetInt32(input any) int32 {
	return int32(GetIntBase(input, 32))
}

func GetInt64(input any) int64 {
	return int64(GetIntBase(input, 64))
}
