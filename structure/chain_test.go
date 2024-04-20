package structure_test

import (
	"testing"

	"github.com/hjwalt/runway/structure"
	"github.com/stretchr/testify/assert"
)

func TestChainBasic(t *testing.T) {
	assert := assert.New(t)

	chain := structure.NewChain[int64, string]()

	chain.PushFirst(1, "one")
	chain.PushFirst(2, "two")
	chain.PushFirst(3, "three")
	chain.PushLast(4, "four")
	chain.PushLast(5, "five")
	chain.PushFirst(6, "six")
	chain.PushLast(7, "seven")

	assert.Equal(int64(7), chain.Size())

	assert.Equal("seven", chain.PeekLast().Value())
	assert.Equal("seven", chain.PopLast().Value())
	assert.False(chain.PeekFirst().Prev().Present())
	assert.False(chain.PeekLast().Next().Present())

	assert.Equal("six", chain.PeekFirst().Value())
	assert.Equal("five", chain.PeekLast().Value())
	assert.False(chain.PeekFirst().Prev().Present())
	assert.False(chain.PeekLast().Next().Present())

	assert.Equal(int64(6), chain.Size())

	assert.Equal("six", chain.PopFirst().Value())
	assert.False(chain.PeekFirst().Prev().Present())
	assert.False(chain.PeekLast().Next().Present())

	assert.Equal("five", chain.PopLast().Value())
	assert.False(chain.PeekFirst().Prev().Present())
	assert.False(chain.PeekLast().Next().Present())

	assert.Equal(int64(4), chain.Size())

	assert.Equal("three", chain.PeekFirst().Value())
	assert.Equal("four", chain.PeekLast().Value())

	assert.Equal("three", chain.PopFirst().Value())
	assert.False(chain.PeekFirst().Prev().Present())
	assert.False(chain.PeekLast().Next().Present())

	assert.Equal("two", chain.PopFirst().Value())
	assert.False(chain.PeekFirst().Prev().Present())
	assert.False(chain.PeekLast().Next().Present())

	assert.Equal("one", chain.PopFirst().Value())
	assert.False(chain.PeekFirst().Prev().Present())
	assert.False(chain.PeekLast().Next().Present())

	assert.Equal("four", chain.PopFirst().Value())

	assert.Equal(int64(0), chain.Size())

	assert.False(chain.PeekFirst().Present())
	assert.False(chain.PopFirst().Present())
	assert.False(chain.PeekLast().Present())
	assert.False(chain.PopLast().Present())
}

func TestChainDetachFirst(t *testing.T) {
	assert := assert.New(t)

	chain := structure.NewChain[int64, string]()

	chain.PushFirst(1, "one")
	chain.PushFirst(2, "two")
	chain.PushFirst(3, "three")
	chain.PushLast(7, "seven")
	chain.PushLast(4, "four")
	chain.PushLast(5, "five")
	chain.PushFirst(6, "six")
	chain.PushLast(7, "seven")

	assert.Equal(int64(8), chain.Size())

	detached := chain.DetachFirst(1)

	assert.Equal("one", detached.Value())
	assert.Equal(int64(7), chain.Size())
	assert.Equal("six", chain.PeekFirst().Value())
	assert.Equal("seven", chain.PeekLast().Value())
	assert.False(chain.PeekFirst().Prev().Present())
	assert.False(chain.PeekLast().Next().Present())

	detached = chain.DetachFirst(6)
	assert.Equal("six", detached.Value())
	assert.Equal(int64(6), chain.Size())
	assert.Equal("three", chain.PeekFirst().Value())
	assert.Equal("seven", chain.PeekLast().Value())
	assert.False(chain.PeekFirst().Prev().Present())
	assert.False(chain.PeekLast().Next().Present())

	detached = chain.DetachFirst(7)
	assert.Equal("seven", detached.Value())
	assert.Equal(int64(5), chain.Size())
	assert.Equal("three", chain.PeekFirst().Value())
	assert.Equal("seven", chain.PeekLast().Value())
	assert.False(chain.PeekFirst().Prev().Present())
	assert.False(chain.PeekLast().Next().Present())

	detached = chain.DetachFirst(7)
	assert.Equal("seven", detached.Value())
	assert.Equal(int64(4), chain.Size())
	assert.Equal("three", chain.PeekFirst().Value())
	assert.Equal("five", chain.PeekLast().Value())
	assert.False(chain.PeekFirst().Prev().Present())
	assert.False(chain.PeekLast().Next().Present())
}

func TestChainDetachLast(t *testing.T) {
	assert := assert.New(t)

	chain := structure.NewChain[int64, string]()

	chain.PushFirst(1, "one")
	chain.PushFirst(2, "two")
	chain.PushFirst(3, "three")
	chain.PushLast(4, "four")
	chain.PushLast(5, "five")
	chain.PushFirst(6, "six")
	chain.PushFirst(1, "one")
	chain.PushLast(7, "seven")

	assert.Equal(int64(8), chain.Size())

	detached := chain.DetachLast(1)

	assert.Equal("one", detached.Value())
	assert.Equal(int64(7), chain.Size())
	assert.Equal("one", chain.PeekFirst().Value())
	assert.Equal("seven", chain.PeekLast().Value())
	assert.False(chain.PeekFirst().Prev().Present())
	assert.False(chain.PeekLast().Next().Present())

	detached = chain.DetachLast(1)

	assert.Equal("one", detached.Value())
	assert.Equal(int64(6), chain.Size())
	assert.Equal("six", chain.PeekFirst().Value())
	assert.Equal("seven", chain.PeekLast().Value())
	assert.False(chain.PeekFirst().Prev().Present())
	assert.False(chain.PeekLast().Next().Present())
}
