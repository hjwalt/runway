package eip

import (
	"net/url"
	"strings"

	"github.com/hjwalt/runway/format"
)

type ChannelSyncType int

const (
	ChannelSync  ChannelSyncType = 0 // i.e. HTTP, AMQP Request-Reply
	ChannelAsync ChannelSyncType = 1 // i.e. Pub/Sub, Kafka
)

// Metadata on EIP channel, not to be confused with golang channel
type Channel[V any] interface {
	// Channel utilities
	GetFormat() format.Format[V]
	SetFormat(format.Format[V])
	GetSync() ChannelSyncType
	SetSync(ChannelSyncType)
	// URL op
	SetUrl(*url.URL)
	GetScheme() string
	GetHost() string // host or host:port
	GetHostname() string
	GetPort() string // returns "" if no port specified
	// Param op
	PutParam(ChannelParam)
	AllParam() ChannelParam
	SetParam(string, string)
	GetParam(string) string
	DelParam(string)
}

func ChannelFromUrl[V any](urlstr string) (Channel[V], error) {
	newchannel := &channel[V]{
		param: ChannelParam{},
	}

	parsed, parseErr := url.Parse(urlstr)
	if parseErr != nil {
		return nil, parseErr
	}

	newchannel.SetUrl(parsed)

	return newchannel, nil
}

type channel[V any] struct {
	format format.Format[V]
	sync   ChannelSyncType
	url    *url.URL
	param  ChannelParam
}

func (c *channel[V]) GetFormat() format.Format[V] {
	return c.format
}

func (c *channel[V]) SetFormat(f format.Format[V]) {
	c.format = f
}

func (c *channel[V]) GetSync() ChannelSyncType {
	return c.sync
}

func (c *channel[V]) SetSync(s ChannelSyncType) {
	c.sync = s
}

func (c *channel[V]) SetUrl(u *url.URL) {
	c.url = u
	for k, vs := range u.Query() {
		for _, v := range vs {
			c.SetParam(k, v)
		}
	}
}

func (c *channel[V]) GetScheme() string {
	return c.url.Scheme
}

func (c *channel[V]) GetHost() string {
	return c.url.Host
}

func (c *channel[V]) GetHostname() string {
	return c.url.Hostname()
}

func (c *channel[V]) GetPort() string {
	return c.url.Port()
}

func (c *channel[V]) PutParam(p ChannelParam) {
	c.param = p.Clone()
}

func (c *channel[V]) AllParam() ChannelParam {
	return c.param
}

func (c *channel[V]) SetParam(k string, v string) {
	c.param.Set(k, v)
}

func (c *channel[V]) GetParam(k string) string {
	return c.param.Get(k)
}

func (c *channel[V]) DelParam(k string) {
	c.param.Del(k)
}

type ChannelParam map[string]string

func (h ChannelParam) Clone() ChannelParam {
	newh := ChannelParam{}
	for k, v := range h {
		newh[k] = v
	}
	return newh
}

func (h ChannelParam) Set(key, value string) {
	k := h.Key(key)
	h[k] = value
}

func (h ChannelParam) Get(key string) string {
	k := h.Key(key)
	if v, ok := h[k]; ok {
		return v
	}
	return ""
}

func (h ChannelParam) Del(key string) {
	k := h.Key(key)
	delete(h, k)
}

func (h ChannelParam) Key(key string) string {
	return strings.ToLower(key)
}
