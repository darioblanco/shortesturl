package logging

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// NewTest returns a Logger instance with a dummy core log provided by the core log library
func NewTest(t testing.TB) *builtinLogger {
	testLogger := zaptest.NewLogger(t, zaptest.Level(zap.ErrorLevel))
	return NewLoggerWithCore(testLogger)
}

// Mocked allows mocking log methods
type Mocked struct {
	mock.Mock
}

func (l *Mocked) Named(name string) *builtinLogger {
	args := l.Called(name)
	return args.Get(0).(*builtinLogger)
}

func (l *Mocked) With(args ...interface{}) *builtinLogger {
	calledArgs := l.Called(args)
	return calledArgs.Get(0).(*builtinLogger)
}

func (l *Mocked) Debug(msg string, keysAndValues ...interface{}) {
	l.Called(msg, keysAndValues)
}

func (l *Mocked) Info(msg string, keysAndValues ...interface{}) {
	l.Called(msg, keysAndValues)
}

func (l *Mocked) Error(msg string, keysAndValues ...interface{}) {
	l.Called(msg, keysAndValues)
}

func (l *Mocked) Warn(msg string, keysAndValues ...interface{}) {
	l.Called(msg, keysAndValues)
}
