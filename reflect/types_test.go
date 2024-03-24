package reflect_test

import (
	"testing"

	"github.com/hjwalt/runway/reflect"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTypeof(t *testing.T) {
	assert := assert.New(t)

	// Test with an integer
	intType := reflect.TypeName(42)
	assert.Equal("int", intType)

	// Test with a string
	stringType := reflect.TypeName("hello")
	assert.Equal("string", stringType)

	// Test with a custom struct
	type Person struct {
		Name string
		Age  int
	}
	personType := reflect.TypeName(Person{})
	assert.Equal("reflect_test.Person", personType)

	// Test with a pointer
	pointerType := reflect.TypeName(&Person{})
	assert.Equal("*reflect_test.Person", pointerType)

	// Test with a pointer var
	var personPointer *Person
	personPointerType := reflect.TypeName(personPointer)
	assert.Equal("*reflect_test.Person", personPointerType)

	// Test with pointer reflect construct
	reflectType := reflect.TypeName(reflect.Construct[*timestamppb.Timestamp]())
	assert.Equal("*timestamppb.Timestamp", reflectType)

	// Test with a slice
	sliceType := reflect.TypeName([]int{})
	assert.Equal("[]int", sliceType)

	// Test with a map
	mapType := reflect.TypeName(map[string]int{})
	assert.Equal("map[string]int", mapType)

	// Test with a function
	funcType := reflect.TypeName(func() {})
	assert.Equal("func()", funcType)
}
