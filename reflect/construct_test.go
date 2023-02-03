package reflect_test

import (
	"context"
	"testing"

	"github.com/hjwalt/runway/reflect"
	"github.com/stretchr/testify/assert"
)

type StructForTest struct {
	Message string
}

type InterfaceForTest interface {
	Test()
}

type FuncForTest func(ctx context.Context)

func TestConstruct(t *testing.T) {
	assert := assert.New(t)

	ptrFormat := reflect.Construct[*StructForTest]()
	assert.Equal(&StructForTest{}, ptrFormat)

	nonPtrFormat := reflect.Construct[StructForTest]()
	assert.Equal(StructForTest{}, nonPtrFormat)

	interfaceType := reflect.Construct[InterfaceForTest]()
	assert.Nil(interfaceType)

	funcType := reflect.Construct[FuncForTest]()
	assert.Nil(funcType)
}
