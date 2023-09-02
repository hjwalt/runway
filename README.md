Runway to make a solid and simple golang software.

## Using

```
go get github.com/hjwalt/runway
```

```
// environment get
environment.GetString("ENV_NAME_TO_SOME_URL_TO_STORE", "default url")

// logging
logger.Debug("debug")
logger.Info("info", zap.String("some key", "some value"))

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

// reflection utility functions
val := reflect.GetBool("false")
```

## Developing

Makefile is heavily used.

```
make test
make update
make tidy
make htmlcov
```

Last coverage: 98.5%

## To Do

1. Go doc
2. Unit test actions

## Why

Golang is a simple language, where anyone can do pretty much everything with the standard language.
However, as I worked with Golang, there are pieces of code that continues to get replicated in many projects.

Some of it includes:

1. Dependency inversion and configuration list parameter management
2. Reflection utility functions
3. Logging utility functions
4. Environment variable retrieval with sane defaults

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

These principles are not object oriented construct, its just easier to achieve with object oriented language in comparison to some other languages.
The principles can be achieved with idiomatic golang.

## Reasoning

### Environment

Environment utility function facades `os.Getenv` with additional logic to perform simple type conversion.

### Logger

Logging functions in golang are fairly basic, but for a production ready software, usually some of the following are needed:

1. JSON based logging for indexing
2. Level changes
3. File based logging instead of stdout

The default does exactly that, with the option to change to however you want to configure the providers.
At the moment, only `zap` is added.

### Reflect

I do not and I think I should not use reflection often in Golang.
However, some type conversions to specific basic types are fairly common, and the utility helps with exactly that.

It reduces the clutter of having multiple utility functions scattered all over the place.
However, the way type conversions are done are very very specific, so do read the codes.

### Inverse

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

Example:

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

### Runtime

Runtimes are independent context that needs to be maintained from the start until the end of the program execution.
Runtimes usually also needs to be instantiated and cleaned up.

Examples of runtimes:

1. HTTP server
2. Kafka consumer
3. Kafka producer
4. Database connection

Why is it important to standardise?

1. Golang is too barebone. This kind of runtime management needs to keep getting rewritten.
2. Resource management is not a difficult problem but the cost of mistakes is high. Your program can hang, you can have connections interrupted without clean disconnect, and many more. Reducing the potential scope of failure is in general a good idea.