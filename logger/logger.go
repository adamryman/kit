package logger

import (
	"fmt"
	"os"

	stdlog "log"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/go-stack/stack"
)

type LogFormat int32

const (
	FormatLogfmt LogFormat = iota
	FormatJSON
	FormatNop
	FormatDefault = FormatJSON
)

var (
	defaultLevels Leveler
	defaultFormat LogFormat = FormatJSON
)

// nominalStackDepth is how far up we need to go in the stack to get the function that called "Log()"
// This is a result of go-kit's log.Valuer pattern.
const nominalStackDepth = 3

func init() {
	stdlog.SetFlags(stdlog.Llongfile)
	Init(defaultFormat)
}

// Init sets the default logger to the desired format and initializes a few
// logging values.  Note that all keyvals added by calling AddDefaultKeyvals
// before Init will be removed.
func Init(format LogFormat) {
	var l log.Logger
	switch format {
	case FormatJSON:
		l = log.NewJSONLogger(os.Stdout)
	case FormatLogfmt:
		l = log.NewLogfmtLogger(os.Stdout)
	case FormatNop:
		l = log.NewNopLogger()
	default:
		panic(fmt.Errorf("invalid log format: %v", format))
	}
	defaultFormat = format
	defaultLevels = levels{l}

	AddDefaultKeyvals(
		"ts", log.DefaultTimestampUTC,
		"caller", caller(nominalStackDepth),
		"function", function(nominalStackDepth),
	)
}

// AddDefaultKeyvals adds values to the defaultLevels, having them be on all logs
func AddDefaultKeyvals(keyvals ...interface{}) {
	defaultLevels = defaultLevels.With(keyvals...)

	stdlibAdapter := log.NewStdlibAdapter(
		defaultLevels.With(
			// Function does not work with the stdlib redirector
			"function", "stdlibLoggerRedirect",
		).Info(),
	)

	stdlog.SetOutput(stdlibAdapter)
}

// Info returns an info level logger
func Info() log.Logger {
	return defaultLevels.Info()
}

// Debug returns an debug level logger
func Debug() log.Logger {
	return defaultLevels.Debug()
}

// Warn returns an warn level logger
func Warn() log.Logger {
	return defaultLevels.Warn()
}

// Error returns an errror level logger.
func Error() log.Logger {
	return defaultLevels.Error()
}

// Crit returns an crit level logger.
func Crit() log.Logger {
	return defaultLevels.Crit()
}

// With adds the key value pairs to the Leveler.
func With(keyvals ...interface{}) Leveler {
	return defaultLevels.With(keyvals...)
}

// WithCustomDepth provides the basic logging context, but evaluates
// caller and function at a different stack depth. This allows the
// creation of convenience functions that wrap log operations without
// obscuring the intended caller and function values.
func WithCustomDepth(depth int, keyvals ...interface{}) Leveler {
	return defaultLevels.WithCustomDepth(depth, keyvals...)
}

// LogError is a convenience function for logging an error with the default
// context and a consistent "err" key for the error value. Other key-values
// can be provided after the error value.
func LogError(err error, keyvals ...interface{}) error {
	// NOTE: Adam: LogError used to return `defaultLevels.LogError(err, keyvals...)`,
	// though that resulting in the key "function" being populated with
	// "LogError" rather than the function calling `LogError`.
	return defaultLevels.WithCustomDepth(nominalStackDepth+1, keyvals...).Error().Log("err", err)
}

// Leveler provides an interface for modifying a log.Context with extra key
// value pairs and enforces the choosing of a level (Info, Debug, Warn, Error,
// Crit) before returning a logger
type Leveler interface {
	// Info returns an info level logger.
	Info() log.Logger
	// Debug returns an debug level logger.
	Debug() log.Logger
	// Warn returns an warn level logger.
	Warn() log.Logger
	// Error returns an error level logger.
	Error() log.Logger
	// Crit returns an crit level logger.
	Crit() log.Logger

	// With adds the key value pairs to the Leveler
	With(keyval ...interface{}) Leveler
	// WithCustomDepth provides the basic logging context, but evaluates
	// caller and function at a different stack depth. This allows the
	// creation of convenience functions that wrap log operations without
	// obscuring the intended caller and function values.
	WithCustomDepth(depth int, keyval ...interface{}) Leveler
	// LogError is a convenience function for logging an error with the default
	// context and a consistent "err" key for the error value. Other key-values
	// can be provided after the error value.
	LogError(err error, keyvals ...interface{}) error
}

// levels implements Leveler and stores an internal logger
// that has all previous key value pairs added.
type levels struct {
	internalLogger log.Logger
}

func (l levels) Info() log.Logger {
	return level.Info(l.internalLogger)
}
func (l levels) Debug() log.Logger {
	return level.Debug(l.internalLogger)
}
func (l levels) Warn() log.Logger {
	return level.Warn(l.internalLogger)
}
func (l levels) Error() log.Logger {
	return level.Error(l.internalLogger)
}
func (l levels) Crit() log.Logger {
	return log.WithPrefix(l.internalLogger, level.Key(), "crit")
}

func (l levels) With(keyvals ...interface{}) Leveler {
	return levels{log.With(l.internalLogger, keyvals...)}
}
func (l levels) WithCustomDepth(depth int, keyvals ...interface{}) Leveler {
	return l.With(
		"caller", caller(depth),
		"function", function(depth),
	).With(keyvals...)
}
func (l levels) LogError(err error, keyvals ...interface{}) error {
	return l.WithCustomDepth(nominalStackDepth+2, keyvals...).Error().Log("err", err)
}

// Helpers

// caller returns the file and line number of the caller as:
// caller returns the path of source file relative to the compile time GOPATH
// and line number of the calling line.
// mv-universe-v2/lib/logger.go:12
func caller(depth int) log.Valuer {
	return func() interface{} {
		return fmt.Sprintf("%+v", stack.Caller(depth))
	}
}

// function returns the function name of the calling function.
func function(depth int) log.Valuer {
	return func() interface{} {
		return fmt.Sprintf("%n", stack.Caller(depth))
	}
}
