package logging

import (
	"github.com/darioblanco/shortesturl/app/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger logs application events
type Logger interface {
	Named(name string) Logger
	With(args ...interface{}) Logger
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
}

type builtinLogger struct {
	coreLogger *zap.Logger
	logger     *zap.SugaredLogger
	Preset     string
}

// NewLogger creates a logging instance
func NewLogger(conf *config.Values) (Logger, error) {
	var (
		zapLogger *zap.Logger
		err       error
		preset    string
	)
	if conf.Environment == "prod" ||
		conf.Environment == "stage" ||
		conf.Environment == "test" {
		zapLogger, err = zap.NewProduction()
		zapLogger = zapLogger.With(
			zap.Field{
				Key:    "environment",
				Type:   zapcore.StringType,
				String: conf.Environment,
			},
			zap.Field{
				Key:    "version",
				Type:   zapcore.StringType,
				String: conf.Version,
			},
		)
		preset = "production"
	} else {
		zapLogger, err = zap.NewDevelopment()
		preset = "development"
	}
	defer zapLogger.Sync()
	return &builtinLogger{
		coreLogger: zapLogger,
		logger:     zapLogger.Sugar(),
		Preset:     preset,
	}, err
}

// NewLoggerWithCore creates an abstracted logger
func NewLoggerWithCore(coreLogger *zap.Logger) *builtinLogger {
	return &builtinLogger{
		coreLogger: coreLogger,
		logger:     coreLogger.Sugar(),
	}
}

// Named adds a sub-scope to the logger's name. Scopes are joined by
// periods. By default, Loggers have no scope.
func (l *builtinLogger) Named(name string) Logger {
	return &builtinLogger{
		coreLogger: l.coreLogger,
		logger:     l.logger.Named(name),
		Preset:     l.Preset,
	}
}

// With adds a variadic number of fields to the logging context. It accepts a
// mix of strongly-typed Field objects and loosely-typed key-value pairs. When
// processing pairs, the first element of the pair is used as the field key
// and the second as the field value.
//
// For example,
//   logger.With(
//     "hello", "world",
//     "failure", errors.New("oh no"),
//     Stack(),
//     "count", 42,
//     "user", User{Name: "alice"},
//  )
// is the equivalent of Zap's unsugared log:
//   unsugared.With(
//     String("hello", "world"),
//     String("failure", "oh no"),
//     Stack(),
//     Int("count", 42),
//     Object("user", User{Name: "alice"}),
//   )
//
// Note that the keys in key-value pairs should be strings. In development,
// passing a non-string key panics. In production, the logger is more
// forgiving: a separate error is logged, but the key-value pair is skipped
// and execution continues. Passing an orphaned key triggers similar behavior:
// panics in development and errors in production.
func (l *builtinLogger) With(args ...interface{}) Logger {
	return &builtinLogger{
		coreLogger: l.coreLogger,
		logger:     l.logger.With(args),
		Preset:     l.Preset,
	}
}

// Debug logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg)
func (l *builtinLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

// Info logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *builtinLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

// Error logs a message with some additional context. In development, the
// logger then panics, otherwise, it logs at ErrorLevel. This behavior
// makes it easier to catch errors that are theoretically possible,
// but shouldn't actually happen, without crashing in production.
// The variadic key-value pairs are treated as they are in With.
func (l *builtinLogger) Error(msg string, keysAndValues ...interface{}) {
	l.logger.DPanicw(msg, keysAndValues...)
}

// Warn logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *builtinLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}
