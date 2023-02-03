package reflect_test

import (
	"testing"
	"time"

	"github.com/hjwalt/runway/reflect"
	"github.com/stretchr/testify/assert"
)

func TestGetField(t *testing.T) {
	assert := assert.New(t)

	strPtr := "test2"

	testStruct := TargetTest{
		Str:    "test",
		StrPtr: &strPtr,
	}

	res, err := reflect.GetField(testStruct, "Str")
	assert.NoError(err)
	assert.Equal("test", res)

	res, err = reflect.GetField(&testStruct, "Str")
	assert.NoError(err)
	assert.Equal("test", res)

	res, err = reflect.GetField(testStruct, "StrPtr")
	assert.NoError(err)
	assert.Equal(&strPtr, res)

	res, err = reflect.GetField(testStruct, "NumPtr")
	assert.NoError(err)
	assert.Nil(res)

	res, err = reflect.GetField(testStruct, "Time")
	assert.NoError(err)
	assert.Equal(time.Time{}, res)

	_, err = reflect.GetField(testStruct, "UnknownField")
	assert.Error(err)

	_, err = reflect.GetField("test", "Time")
	assert.Error(err)

	var testNilStruct *TargetTest

	res, err = reflect.GetField(testNilStruct, "Time")
	assert.NoError(err)
	assert.Nil(res)
}
