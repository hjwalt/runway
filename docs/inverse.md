---
layout: page
title: Inverse
---

I have looked into dig, fx, and wire, if there are any more I should look at please add in the issues or discussions.

- Wire is impractical to use in multi integration complex set up (I used Kafka, Postgres, HTTP, and some other MQ integrations).
- Dig is not recommended within Golang community due to excessive use of reflection.
- Fx is too complicated (maybe just for me).

This is a compromise that:

1. Does not use reflection (excessively, only interface casting and default value construction).
2. Uses simple injector method.
3. Allows multiple constructors for the same qualifier, and get all of them if needed.
4. Does loop detection (assuming context are passed properly).
5. Saves constructed instances for reuse

Any resolution errors will be runtime detected not compile time detected, fortunately or unfortunately.
I try to use the injected instances as early as possible (as if it is Java beans construction).

## Using

```
// inverse registration
inverse.Register(
	"test-1-qualifier", 
	func(ctx context.Context) (any, error) { 
		return "test-1-value", nil 
	},
)

// inverse resolution
val, err := inverse.GetLast[string](context.Background(), "test-1-qualifier")
val, err := inverse.GetAll[string](context.Background(), "test-1-qualifier")
```

## Examples

Multi level injector:

```
func MultiLevelInjector(ctx context.Context) (any, error) {
	depA, err := inverse.GetLast[DepA](ctx, "depAQualifier")
	if err != nil {
		return nil, err
	}
	depB, err := inverse.GetLast[DepB](ctx, "depBQualifier")
	if err != nil {
		return nil, err
	}
	return ThisStruct {
		DepA: depA,
		DepB: depB,
	}, nil
}
```
