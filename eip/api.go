package eip

import (
	"github.com/hjwalt/runway/managed"
	"github.com/hjwalt/runway/optional"
)

type SubscriberFactory[V any] interface {
	managed.Component

	Create(Channel[V], Handler[V]) (Subscriber[V], error)
}

type Subscriber[V any] interface {
	managed.Service

	Channel() Channel[V]
}

type PublisherFactory[V any] interface {
	managed.Component

	Create(Channel[V]) (Publisher[V], error)
}

type Publisher[V any] interface {
	managed.Service

	Channel() Channel[V]
	Publish([]Message[V]) error
}

type Handler[V any] interface {
	managed.Service

	Handle([]Message[V]) error
}

type Procedure[V any] interface {
	managed.Service

	Channel() Channel[V]
	Call(Message[V]) (Message[V], error)
}

type Batch[V1 any, V2 any] func([]Message[V1]) ([]Message[V2], error)

type Pipe[V1 any, V2 any] func(Message[V1]) (Message[V2], error)

type Filter[V any] func(Message[V]) (optional.Optional[Message[V]], error)

type Splitter[V1 any, V2 any] func(Message[V1]) ([]Message[V2], error)
