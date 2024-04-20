package structure

import (
	"sync"
)

func NewChain[I comparable, V any]() Chain[I, V] {
	return &chain[I, V]{
		mutex: sync.Mutex{},
		first: emptyChainNode[I, V](),
		last:  emptyChainNode[I, V](),
	}
}

func EmptyChainNode[I comparable, V any]() ChainNode[I, V] {
	return emptyChainNode[I, V]()
}

// to be used for queues, stacks, and linked list.
// thread / goroutine safe.
// requires explicit comparable identifier for more straightforward operations.
type Chain[I comparable, V any] interface {
	PushFirst(I, V)
	PeekFirst() ChainNode[I, V]
	PopFirst() ChainNode[I, V]
	DetachFirst(I) ChainNode[I, V]

	PushLast(I, V)
	PeekLast() ChainNode[I, V]
	PopLast() ChainNode[I, V]
	DetachLast(I) ChainNode[I, V]

	Size() int64
}

type chain[I comparable, V any] struct {
	mutex sync.Mutex
	first *chainNode[I, V]
	last  *chainNode[I, V]
	size  int64
}

func (c *chain[I, V]) PushFirst(identity I, value V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	node := newChainNode[I, V](identity, value)

	if c.first != nil && c.first.Present() {
		// chain -> node -> prevfirst -> ... -> last
		prevfirst := c.first

		c.first = node
		prevfirst.prev = node
		node.next = prevfirst
	} else {
		// chain -> first -> node <- last
		c.first = node
		c.last = node
	}

	c.size += 1
}

func (c *chain[I, V]) PeekFirst() ChainNode[I, V] {
	return c.first
}

func (c *chain[I, V]) PopFirst() ChainNode[I, V] {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.detachFirst()
}

func (c *chain[I, V]) DetachFirst(identity I) ChainNode[I, V] {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	node := c.first

	for node != nil && node.Present() {
		if node.Identity() == identity {
			break
		}

		node = node.next
	}

	return c.detach(node)
}

func (c *chain[I, V]) PushLast(identity I, value V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	node := newChainNode[I, V](identity, value)

	if c.last != nil && c.last.Present() {
		// chain -> first -> ... -> prevlast -> node <- last
		prevlast := c.last

		c.last = node
		prevlast.next = node
		node.prev = prevlast
	} else {
		// chain -> first -> node <- last
		c.first = node
		c.last = node
	}

	c.size += 1
}

func (c *chain[I, V]) PeekLast() ChainNode[I, V] {
	return c.last
}

func (c *chain[I, V]) PopLast() ChainNode[I, V] {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.detachLast()
}

func (c *chain[I, V]) DetachLast(identity I) ChainNode[I, V] {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	node := c.last

	for node != nil && node.Present() {
		if node.Identity() == identity {
			break
		}

		node = node.prev
	}

	return c.detach(node)
}

func (c *chain[I, V]) detach(node *chainNode[I, V]) *chainNode[I, V] {
	if node == c.first {
		return c.detachFirst()
	}

	if node == c.last {
		return c.detachLast()
	}

	if node != nil && node.Present() {
		// chain -> first -> ... -> prevnode -> node -> nextnode -> ... -> last
		prevnode := node.prev
		nextnode := node.next

		// chain -> first -> ... -> prevnode -> nextnode -> ... -> last
		prevnode.next = nextnode
		nextnode.prev = prevnode

		node.next = emptyChainNode[I, V]()
		node.prev = emptyChainNode[I, V]()

		c.size -= 1

		return node
	}

	return emptyChainNode[I, V]()
}

func (c *chain[I, V]) detachFirst() *chainNode[I, V] {
	if c.first != nil && c.first.Present() {
		// chain -> first -> currfirst -> nextfirst -> ... -> last
		currfirst := c.first
		nextfirst := c.first.next

		if nextfirst != nil && nextfirst.Present() {
			// chain -> first -> nextfirst -> ... -> last
			c.first = nextfirst
			c.first.prev = emptyChainNode[I, V]()
		} else {
			// chain -> first -> empty <- last
			c.first = emptyChainNode[I, V]()
			c.last = emptyChainNode[I, V]()
		}

		c.size -= 1

		currfirst.next = emptyChainNode[I, V]()
		currfirst.prev = emptyChainNode[I, V]()

		return currfirst
	} else {
		return emptyChainNode[I, V]()
	}
}

func (c *chain[I, V]) detachLast() *chainNode[I, V] {
	if c.last != nil && c.last.Present() {
		// chain -> first -> ... -> prevlast -> currlast <- last
		currlast := c.last
		prevlast := c.last.prev

		if prevlast != nil && prevlast.Present() {
			// chain -> first -> ... -> prevlast <- last
			c.last = prevlast
			c.last.next = emptyChainNode[I, V]()
		} else {
			// chain -> first -> empty <- last
			c.first = emptyChainNode[I, V]()
			c.last = emptyChainNode[I, V]()
		}

		c.size -= 1

		currlast.next = emptyChainNode[I, V]()
		currlast.prev = emptyChainNode[I, V]()

		return currlast
	} else {
		return emptyChainNode[I, V]()
	}
}

func (c *chain[I, V]) Size() int64 {
	return c.size
}

type ChainNode[I comparable, V any] interface {
	Present() bool
	Identity() I
	Value() V
	Next() ChainNode[I, V]
	Prev() ChainNode[I, V]
}

func emptyChainNode[I comparable, V any]() *chainNode[I, V] {
	return &chainNode[I, V]{
		present: false,
	}
}

func newChainNode[I comparable, V any](identity I, value V) *chainNode[I, V] {
	return &chainNode[I, V]{
		present: true,
		ident:   identity,
		value:   value,
		next:    emptyChainNode[I, V](),
		prev:    emptyChainNode[I, V](),
	}
}

type chainNode[I comparable, V any] struct {
	next    *chainNode[I, V]
	prev    *chainNode[I, V]
	ident   I
	value   V
	present bool
}

func (n *chainNode[I, V]) Present() bool {
	return n.present
}

func (n *chainNode[I, V]) Identity() I {
	return n.ident
}

func (n *chainNode[I, V]) Value() V {
	return n.value
}

func (n *chainNode[I, V]) Next() ChainNode[I, V] {
	if !n.present {
		return emptyChainNode[I, V]()
	}
	return n.next
}

func (n *chainNode[I, V]) Prev() ChainNode[I, V] {
	if !n.present {
		return emptyChainNode[I, V]()
	}
	return n.prev
}
