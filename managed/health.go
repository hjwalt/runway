package managed

import (
	"context"
	"strings"
	"sync"

	"github.com/hjwalt/runway/inverse"
)

func AddHealth(ic inverse.Container) {
	AddComponent(ic, &health{})
}

func GetHealth(container inverse.Container, ctx context.Context) (Health, error) {
	return inverse.GenericGet[Health](container, ctx, QualifierHealth)
}

type health struct {
	mutex   sync.Mutex
	strings map[string]string
	ints    map[string]int64
	bools   map[string]bool
}

func (c *health) Name() string {
	return QualifierHealth
}

func (c *health) Register(ctx context.Context, ic inverse.Container) error {
	if c.strings == nil {
		c.strings = map[string]string{}
	}
	ic.AddVal(c.Name(), c)
	return nil
}

func (c *health) Resolve(context.Context, inverse.Container) error {
	return nil
}

func (c *health) Clean() error {
	return nil
}

func (c *health) GetString() map[string]string {
	return c.strings
}

func (c *health) SetString(component string, key string, value string) {
	c.strings[healthKey(component, key)] = value
}

func (c *health) GetBool() map[string]bool {
	return c.bools
}

func (c *health) SetBool(component string, key string, value bool) {
	c.bools[healthKey(component, key)] = value
}

func (c *health) GetInt() map[string]int64 {
	return c.ints
}

func (c *health) SetInt(component string, key string, value int64) {
	c.ints[healthKey(component, key)] = value
}

func (c *health) IncInt(component string, key string, value int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	curr, exists := c.ints[healthKey(component, key)]
	if !exists {
		curr = 0
	}
	c.ints[healthKey(component, key)] = curr + value
}

func (c *health) DecInt(component string, key string, value int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	curr, exists := c.ints[healthKey(component, key)]
	if !exists {
		curr = 0
	}
	c.ints[healthKey(component, key)] = curr - value
}
func healthKey(component string, key string) string {
	return strings.Join([]string{component, key}, "_")
}
