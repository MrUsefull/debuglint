# DebugLint

## Overview

DebugLint is a Go static analysis tool that identifies performance-impacting debug log statements. It detects expensive function calls in debug log arguments that execute even when debug logging is disabled, helping prevent production performance issues.

**Key Features:**

- **Multi-library support** - Works with zap and slog out of the box, extensible to other logging libraries
- **Smart function detection** - Automatically discovers debug wrapper functions and logging helpers
- **Comprehensive analysis** - Detects both direct violations and nested expensive calls

## Architecture

DebugLint uses a two-phase analysis approach built on Go's `golang.org/x/tools/go/analysis` framework:

### Phase 1: LogIdentify Analyzer

- **Purpose**: Discovers all debug-related functions in the codebase
- **Detection**: Identifies three types of functions:
  - **Known debug functions** (e.g., `zap.Logger.Debug`)
  - **Normal debug wrappers** (functions containing only debug calls)
  - **Guard-required wrappers** (debug functions with control flow like loops)

### Phase 2: DebugLint Analyzer  

- **Purpose**: Flags unguarded expensive function calls in debug statements
- **Analysis**: Uses results from LogIdentify to determine what constitutes a debug function
- **Reporting**: Issues `DL-1` violations for performance-impacting patterns

### Package Structure

```text
├── cmd/                    # CLI entry point
├── pkg/
│   ├── debuglint/         # Main linting analyzer (Phase 2)
│   ├── logidentify/       # Debug function discovery (Phase 1)  
│   └── linterrs/          # Common error types
└── internal/
    ├── configs/           # Library-specific configurations
    ├── funcs/             # Function description and categorization
    └── tools/             # Build utilities
```

## Installation

TODO: Run and install instructions
TODO: Publish for golangci-lint?

## Examples

### Function arguments are evaluated prior to function calls

That means expensive calls can cause debug logging to negatively impact performance.

```go
// bad - expensiveFunctionHere is always called
logging.Debug("my message", expensiveFunctionHere())

// fixed - expensiveFunctionHere is only called when debug logs are enabled
if logging.IsDebugEnabled() {
    logging.Debug("my message", expensiveFunctionHere())
}
```

#### zap logger debug enabled checks

There are multiple ways to check if debug logging is enabled with zap. The best place
to find up to date documentation is the [zap](https://pkg.go.dev/go.uber.org/zap#example-Logger.Check)
go packages page.

Example check methodology, lifted directly from the zap documentation

```go
 if ce := logger.Check(zap.DebugLevel, "debugging"); ce != nil {
  // If debug-level log output isn't enabled or if zap's sampling would have
  // dropped this log entry, we don't allocate the slice that holds these
  // fields.
  ce.Write(
   zap.String("foo", "bar"),
   zap.String("baz", "quux"),
  )
 }
```

Example logger.Level methodology lifted directly from zap's package documentation

```go
if logger.Level()  <= zapcore.DebugLevel {
    logger.Debug("message here", zap.String("myExpensiveStr", myExpensiveStr()))
}
```

#### Custom check functions

### Variables used only in debug scopes can be similarly expensive

Linter does not currently detect these issues.

```go
// bad - myExpensiveVar is still computed when debug logging is disabled
var myExpensiveVar = expensiveFunctionHere()
if logging.IsDebugEnabled() {
    logging.Debug("my message", myExpensiveVar)
}

// fixed - limit variable to debug scope
if logging.IsDebugEnabled() {
    var myExpensiveVar = expensiveFunctionHere()
    logging.Debug("my message", myExpensiveVar)
}
```

### Debug logging wrappers or aliases can also be expensive

```go
func myWeirdDebugFn(input []string) {
    logging.Debug("Some Message")
    for str := range input {
        logging.Debug("This is an example string", str)
    }
}

func MyFunction() {
    // bad
    myWeirdDebugFn(expensiveInputComputation())

    // fixed
    if logging.IsDebugEnabled() {
        myWeirdDebugFn(expensiveInputComputation())
    }
}
```

## Support

Create an issue in this repository.

## Roadmap

- [ ] Basic config parsing
- [ ] Document config format
- [ ] golangci-lint integration
- [ ] Detect var usage

## Contributing

Open a Pull Request.

## License

MIT

## Project status

In development
