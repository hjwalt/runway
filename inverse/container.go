package inverse

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/hjwalt/runway/logger"
)

func NewContainer() Container {
	return &container{
		resolvers: map[string][]Resolver{},
		lastCache: map[string]any{},
		allCache:  map[string][]any{},
	}
}

type Resolver func(context.Context, Container) (any, error)

func Qualifier(name string) qualifier {
	return qualifier{Name: name}
}

type qualifier struct {
	Name string
}

type Container interface {
	Reset()
	Clear()
	Invalid(error)
	Valid() bool
	Error() error
	Add(q string, r Resolver)
	AddVal(q string, v any)
	Get(ctx context.Context, q string) (any, error)
	GetAll(ctx context.Context, q string) ([]any, error)
}

type container struct {
	mutex     sync.Mutex
	resolvers map[string][]Resolver
	allCache  map[string][]any
	lastCache map[string]any
	err       error
}

func (c *container) Reset() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.resolvers = map[string][]Resolver{}
	c.allCache = map[string][]any{}
	c.lastCache = map[string]any{}
	c.err = nil
}

func (c *container) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.allCache = map[string][]any{}
	c.lastCache = map[string]any{}
}

func (c *container) Valid() bool {
	return c.err == nil
}

func (c *container) Invalid(err error) {
	c.err = err
}

func (c *container) Error() error {
	return c.err
}

func (c *container) Add(q string, r Resolver) {
	if !c.Valid() {
		logger.Error("invalid container")
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if resolvers, qualifierExist := c.resolvers[q]; qualifierExist {
		c.resolvers[q] = append(resolvers, r)
	} else {
		c.resolvers[q] = []Resolver{r}
	}
}

func (c *container) AddVal(q string, v any) {
	if !c.Valid() {
		logger.Error("invalid container")
		return
	}

	c.Add(q, func(ctx context.Context, container Container) (any, error) { return v, nil })
}

func (c *container) Get(ctx context.Context, q string) (any, error) {
	if !c.Valid() {
		return nil, getError(q, c.Error())
	}

	// check for cached data
	cached, cachedExist := c.lastCache[q]
	if cachedExist {
		return cached, nil
	}

	// check for resolver
	resolvers, getResolverErr := c.getResolver(ctx, q)
	if getResolverErr != nil {
		return nil, getResolverErr
	}

	lastResolver := resolvers[len(resolvers)-1]
	resolvedVal, resolveErr := c.resolve(ctx, q, lastResolver)
	if resolveErr != nil {
		return nil, resolveErr
	}

	// cache
	c.lastCache[q] = resolvedVal

	return resolvedVal, nil
}

func (c *container) GetAll(ctx context.Context, q string) ([]any, error) {
	if !c.Valid() {
		return []any{}, getError(q, c.Error())
	}

	// check for cached data
	cached, cachedExist := c.allCache[q]
	if cachedExist {
		return cached, nil
	}

	// check for resolver
	resolvers, getResolverErr := c.getResolver(ctx, q)
	if getResolverErr != nil {
		return []any{}, getResolverErr
	}

	resultList := make([]any, len(resolvers))
	for ri, resolver := range resolvers {
		resolvedVal, resolveErr := c.resolve(ctx, q, resolver)
		if resolveErr != nil {
			return []any{}, resolveErr
		}
		resultList[ri] = resolvedVal
	}

	// cache
	c.allCache[q] = resultList

	return resultList, nil
}

func (c *container) getResolver(ctx context.Context, q string) ([]Resolver, error) {
	// check for context and loop
	if ctx == nil {
		return []Resolver{}, getError(q, ErrInverseNilContext)
	}
	if ctx.Value(Qualifier(q)) != nil {
		c.Invalid(ErrInverseResolveLoop)
		return []Resolver{}, getError(q, ErrInverseResolveLoop)
	}

	// resolve
	c.mutex.Lock()
	resolvers, qualifierExist := c.resolvers[q]
	c.mutex.Unlock()

	if !qualifierExist {
		return []Resolver{}, getError(q, ErrInverseResolverMissing)
	}

	return resolvers, nil
}

func (c *container) resolve(ctx context.Context, q string, resolver Resolver) (any, error) {
	injectedVal, injectedErr := resolver(context.WithValue(ctx, Qualifier(q), "already-resolved"), c)
	if injectedErr != nil {
		return nil, errors.Join(injectedErr, ErrInverseResolveError, fmt.Errorf("qualifier being resolved is %s", q))
	}
	return injectedVal, nil
}

func getError(q string, err error) error {
	return errors.Join(err, fmt.Errorf("qualifier being resolved is %s", q))
}

var (
	ErrInverseNilContext      = errors.New("container received nil context")
	ErrInverseResolveLoop     = errors.New("container dependency loop")
	ErrInverseResolverMissing = errors.New("container does not contain qualifier")
	ErrInverseResolveError    = errors.New("container failed to resolve qualifier")
	ErrInverseCastError       = errors.New("container value cast error")
)
