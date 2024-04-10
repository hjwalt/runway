package managed

import (
	"context"

	"github.com/hjwalt/runway/inverse"
)

const (
	QualifierLifecycle   = "lifecycle"
	QualifierRetry       = "retry"
	QualifierHttpHandler = "http-handler"
	QualifierHealth      = "health"
)

// lifecycle: register -> resolve -> clean
type Component interface {
	Name() string
	Register(context.Context, inverse.Container) error
	Resolve(context.Context, inverse.Container) error
	Clean() error
}

// lifecycle: register -> resolve -> start -> stop -> clean
type Service interface {
	Component

	Start() error
	Stop() error
}

// special archetypes

type Configuration interface {
	Component

	Has(key string) bool
	Get() map[string]string
	GetString(key string, defaultValue string) string
	GetBool(key string, defaultValue bool) bool
	GetInt32(key string, defaultValue int32) int32
	GetInt64(key string, defaultValue int64) int64
	GetUint32(key string, defaultValue uint32) uint32
	GetUint64(key string, defaultValue uint64) uint64
}

type Health interface {
	Component

	GetString() map[string]string
	SetString(component string, key string, value string)
	GetBool() map[string]bool
	SetBool(component string, key string, value bool)
	GetInt() map[string]int64
	SetInt(component string, key string, value int64)
	IncInt(component string, key string, value int64)
	DecInt(component string, key string, value int64)
}

type Retry interface {
	Component

	Do(fnToDo func(int64) error) error
}

type Lifecycle interface {
	Service

	Running() bool
	Error(error)
}

// Running functions

type Runnable interface {
	Run() error
}

type Channel[T any] interface {
	Channel() (<-chan T, error)
	Loop(T) error
}

type Loop interface {
	Loop() error
}
