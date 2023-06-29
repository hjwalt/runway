package stringy_test

import (
	"testing"

	"github.com/hjwalt/runway/stringy"
	"github.com/stretchr/testify/assert"
)

func TestSnake(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("snake_case_this", stringy.ToLowerSnakeCase("snakeCaseThis"))
	assert.Equal("snake_case_this", stringy.ToLowerSnakeCase("SnakeCaseThis"))

	assert.Equal("SNAKE_CASE_THIS", stringy.ToUpperSnakeCase("snakeCaseThis"))
	assert.Equal("SNAKE_CASE_THIS", stringy.ToUpperSnakeCase("SnakeCaseThis"))

	// odd cases

	assert.Equal("s_snake_case_this", stringy.ToLowerSnakeCase("SSnakeCaseThis"))
}
