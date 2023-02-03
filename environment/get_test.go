package environment_test

import (
	"os"
	"testing"

	"github.com/hjwalt/runway/environment"
	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	os.Setenv("ENVIRONMENT_TEST_STRING", "test")
	os.Setenv("ENVIRONMENT_TEST_STRING_WHITESPACE", "   asdf    ")

	assert := assert.New(t)

	assert.Equal("test", environment.GetString("ENVIRONMENT_TEST_STRING", "default"))
	assert.Equal("asdf", environment.GetString("ENVIRONMENT_TEST_STRING_WHITESPACE", "default"))
	assert.Equal("default", environment.GetString("ENVIRONMENT_TEST_STRING_NOT_FOUND", "default"))
}

func TestGetInt64(t *testing.T) {
	os.Setenv("ENVIRONMENT_TEST_INT", "1")
	os.Setenv("ENVIRONMENT_TEST_INT_WHITESPACE", "           2              ")
	os.Setenv("ENVIRONMENT_TEST_INT_NOT_PARSEABLE", "2asdfasdfasdf2              ")

	assert := assert.New(t)

	assert.Equal(int64(1), environment.GetInt64("ENVIRONMENT_TEST_INT", 0))
	assert.Equal(int64(2), environment.GetInt64("ENVIRONMENT_TEST_INT_WHITESPACE", 0))
	assert.Equal(int64(0), environment.GetInt64("ENVIRONMENT_TEST_INT_NOT_FOUND", 0))
	assert.Equal(int64(100), environment.GetInt64("ENVIRONMENT_TEST_INT_NOT_PARSEABLE", 100))
}

func TestGetBool(t *testing.T) {
	os.Setenv("ENVIRONMENT_TEST_BOOL", "true")
	os.Setenv("ENVIRONMENT_TEST_BOOL_WHITESPACE", "          true                  ")
	os.Setenv("ENVIRONMENT_TEST_BOOL_NOT_PARSEABLE", "          asdfasdfdasf                  ")

	assert := assert.New(t)

	assert.Equal(true, environment.GetBool("ENVIRONMENT_TEST_BOOL", false))
	assert.Equal(true, environment.GetBool("ENVIRONMENT_TEST_BOOL_WHITESPACE", false))
	assert.Equal(false, environment.GetBool("ENVIRONMENT_TEST_BOOL_NOT_FOUND", false))
	assert.Equal(false, environment.GetBool("ENVIRONMENT_TEST_BOOL_NOT_PARSEABLE", false))
}
