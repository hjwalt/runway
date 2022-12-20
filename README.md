# Runway

Runway to make a solid and simple golang software.

## Using

```
go get github.com/hjwalt/runway
```

```
logger.Debug("debug")
logger.Info("info", zap.String("some key", "some value"))

inverse.Register("test-1-qualifier", func(ctx context.Context) (any, error) { return "test-1-value", nil })
val, err := inverse.GetLast[string](context.Background(), "test-1-qualifier")
val, err := inverse.GetAll[string](context.Background(), "test-1-qualifier")

val := reflect.GetBool("false")
```

For more functions do read the codebase.

## To Do

1. Go doc
2. More tests

## Why

Golang is a simple language, where anyone can do pretty much everything with the standard language.
However, as I worked with Golang, there are pieces of code that continues to get replicated in many projects.

Some of it includes:

1. Dependency inversion management
2. Reflection utility functions
3. Logging utility functions

This project is a combination of those.

## Principle

### KISS

Keep it simple and stupid. 
The catch phrase of all competitive programmers.
This should also apply for complex softwares.

- Sane defaults
- Popular libraries
- Constant improvements
- Easy to use utility functions

### SOLID

The principles of SOLID (refer to Clean Architecture by Robert C. Martin for proper wordings) applies:

1. Single responsibility: each software components has one and only one reason to change
2. Open closed: behaviour of systems should be changed by adding new code, rather than changing existing code
3. Liskov substitution: software components must adhere to contracts (interfaces) to be substitutable
4. Interface segregation: avoid depending on things not used
5. Dependency inversion: depend on policies (interfaces) rather than details (implementation)

These principles are NOT object oriented construct, its just easier to achieve with object oriented language in comparison to languages prior.
In this case I am talking about ancient versions of C / C++. 
Modern C and C++ can achieve SOLID.

Meaning, these can be achieved with idiomatic golang, without ugly complex implementations.

## Reasoning

### Logger

Logging functions in golang are fairly basic, but for a production ready software, usually some of the following are needed:

1. JSON based logging for indexing
2. Level changes
3. File based logging instead of stdout

The default does exactly that, with the option to change to however you want to configure the providers.
At the moment, only `zap` is added.

### Reflect

I do not nor I think I should use reflection often in Golang.
However, some type conversions to specific basic types are fairly common, and the utility helps with exactly that.

It reduces the clutter of having multiple utility functions scattered all over the place.
However, the way type conversions are done is very very specific, so do read the codes.

### Inverse

I have looked into dig, fx, and wire, if there are any more I should look at please add in the issues.

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

Example:

```
func MultiLevelInjector(ctx context.Context) (any, error) {
	depA, err := inject.GetLast[DepA](ctx, "depAQualifier")
	if err != nil {
		return nil, err
	}
	depB, err := inject.GetLast[DepB](ctx, "depBQualifier")
	if err != nil {
		return nil, err
	}
	return ThisStruct {
		DepA: depA,
		DepB: depB,
	}, nil
}
```
