package eip

import "strings"

const (
	HeaderId             = "id"
	HeaderKey            = "metadata_key"
	HeaderMetadataPrefix = "metadata_" // denotes headers that should not be published into brokers
)

type Message[V any] interface {
	// Header op
	PutHeader(MessageHeader)
	AllHeader() MessageHeader
	AddHeader(string, string)
	SetHeader(string, string)
	GetHeader(string) string
	ValHeader(string) []string
	DelHeader(string)

	// Body op
	GetBody() V
	SetBody(V)
}

func NewMessage[V any]() Message[V] {
	return &message[V]{
		header: MessageHeader{},
	}
}

func CopyMessage[V any](m Message[V]) Message[V] {
	newm := NewMessage[V]()
	newm.PutHeader(m.AllHeader())
	newm.SetBody(m.GetBody())
	return newm
}

type message[V any] struct {
	header MessageHeader
	body   V
}

func (m *message[V]) PutHeader(h MessageHeader) {
	m.header = h.Clone()
}

func (m *message[V]) AllHeader() MessageHeader {
	return m.header
}

func (m *message[V]) AddHeader(h string, v string) {
	m.header.Add(h, v)
}

func (m *message[V]) SetHeader(h string, v string) {
	m.header.Set(h, v)
}

func (m *message[V]) GetHeader(h string) string {
	return m.header.Get(h)
}

func (m *message[V]) ValHeader(h string) []string {
	return m.header.Val(h)
}

func (m *message[V]) DelHeader(h string) {
	m.header.Del(h)
}

func (m *message[V]) GetBody() V {
	return m.body
}

func (m *message[V]) SetBody(v V) {
	m.body = v
}

type MessageHeader map[string][]string

func (h MessageHeader) Clone() MessageHeader {
	newh := MessageHeader{}
	for k, v := range h {
		newh[k] = append([]string{}, v...)
	}
	return newh
}

func (h MessageHeader) Add(key, value string) {
	k := h.Key(key)
	h[k] = append(h[k], value)
}

func (h MessageHeader) Set(key, value string) {
	k := h.Key(key)
	h[k] = []string{value}
}

func (h MessageHeader) Get(key string) string {
	k := h.Key(key)
	v := h[k]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

func (h MessageHeader) Val(key string) []string {
	k := h.Key(key)
	if v, ok := h[k]; ok {
		return v
	}
	return []string{}
}

func (h MessageHeader) Del(key string) {
	k := h.Key(key)
	delete(h, k)
}

func (h MessageHeader) Key(key string) string {
	return strings.ToLower(key)
}
