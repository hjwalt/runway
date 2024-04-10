package managed

import (
	"context"
	"errors"
	"strings"

	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/reflect"
)

const (
	QualifierConfigPrefix = "config"
)

var ErrConfigMissingTarget = errors.New("configuration target missing")

func EmptyConfig(component string) Configuration {
	return &config{
		component: component,
	}
}

func NewConfig(component string, configuration map[string]string) Configuration {
	return &config{
		component:     component,
		configuration: configuration,
	}
}

func ResolveConfig(ctx context.Context, container inverse.Container, component string) (Configuration, error) {
	return inverse.GenericGet[Configuration](container, ctx, configKey(component))
}

type config struct {
	component     string
	configuration map[string]string
}

func (c *config) Name() string {
	return configKey(c.component)
}

func (c *config) Register(ctx context.Context, ic inverse.Container) error {
	if len(c.component) == 0 {
		return ErrConfigMissingTarget
	}
	if c.configuration == nil {
		c.configuration = map[string]string{}
	}
	ic.AddVal(c.Name(), c)
	return nil
}

func (c *config) Resolve(context.Context, inverse.Container) error {
	return nil
}

func (c *config) Clean() error {
	return nil
}

func (c *config) Has(key string) bool {
	_, exists := c.configuration[key]
	return exists
}

func (c *config) Get() map[string]string {
	return c.configuration
}

func (c *config) GetString(key string, defaultValue string) string {
	if value, exists := c.configuration[key]; exists {
		return value
	} else {
		return defaultValue
	}
}

func (c *config) GetBool(key string, defaultValue bool) bool {
	if value, exists := c.configuration[key]; exists {
		return reflect.GetBool(value)
	} else {
		return defaultValue
	}
}

func (c *config) GetInt32(key string, defaultValue int32) int32 {
	if value, exists := c.configuration[key]; exists {
		return reflect.GetInt32(value)
	} else {
		return defaultValue
	}
}

func (c *config) GetInt64(key string, defaultValue int64) int64 {
	if value, exists := c.configuration[key]; exists {
		return reflect.GetInt64(value)
	} else {
		return defaultValue
	}
}

func (c *config) GetUint32(key string, defaultValue uint32) uint32 {
	if value, exists := c.configuration[key]; exists {
		return reflect.GetUint32(value)
	} else {
		return defaultValue
	}
}

func (c *config) GetUint64(key string, defaultValue uint64) uint64 {
	if value, exists := c.configuration[key]; exists {
		return reflect.GetUint64(value)
	} else {
		return defaultValue
	}
}

func configKey(target string) string {
	return strings.Join([]string{QualifierConfigPrefix, target}, "_")
}
