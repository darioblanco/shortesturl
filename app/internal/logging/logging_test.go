package logging

import (
	"testing"

	"github.com/darioblanco/shortesturl/app/internal/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestNewLogger_Development(t *testing.T) {
	logger, err := NewLogger(&config.Values{Environment: "dev"})
	assert.NoError(t, err)
	assert.Equal(t, "development", logger.(*builtinLogger).Preset)
}

func TestNewLogger_Production(t *testing.T) {
	logger, err := NewLogger(&config.Values{Environment: "prod"})
	assert.NoError(t, err)
	assert.Equal(t, "production", logger.(*builtinLogger).Preset)
}

func TestNewLoggerWithCore(t *testing.T) {
	zapLogger := zap.NewNop()
	logger := NewLoggerWithCore(zapLogger)
	assert.Equal(t, zapLogger, logger.coreLogger)
	assert.Equal(t, zapLogger.Sugar(), logger.logger)
}

func TestLoggerNamed(t *testing.T) {
	testLogger := zaptest.NewLogger(t)
	logger := &builtinLogger{logger: testLogger.Sugar()}
	newLogger := logger.Named("test")
	assert.NotEqual(t, newLogger, logger)
}

func TestLoggerWith(t *testing.T) {
	testLogger := zaptest.NewLogger(t)
	logger := &builtinLogger{logger: testLogger.Sugar()}
	newLogger := logger.With("hello", "world")
	assert.Equal(t, newLogger, logger)
}

func TestLoggerGetCoreLogger(t *testing.T) {
	testLogger := zaptest.NewLogger(t)
	logger := &builtinLogger{coreLogger: testLogger, logger: testLogger.Sugar()}
	assert.Equal(t, logger.coreLogger.Sugar(), logger.logger)
}

func createZapLogger(t *testing.T, expectedLevel zapcore.Level) *zap.Logger {
	return zaptest.NewLogger(t, zaptest.WrapOptions(zap.Hooks(func(e zapcore.Entry) error {
		if e.Level != expectedLevel {
			t.Fatalf("it should log in %s", expectedLevel)
		}
		return nil
	})))
}

func TestLoggerDebug(t *testing.T) {
	testLogger := createZapLogger(t, zap.DebugLevel)
	logger := &builtinLogger{logger: testLogger.Sugar()}
	logger.Debug("My message", "key", "value")
}

func TestLoggerInfo(t *testing.T) {
	testLogger := createZapLogger(t, zap.InfoLevel)
	logger := &builtinLogger{logger: testLogger.Sugar()}
	logger.Info("My message", "key", "value")
}

func TestLoggerError(t *testing.T) {
	testLogger := createZapLogger(t, zap.DPanicLevel)
	logger := &builtinLogger{logger: testLogger.Sugar()}
	logger.Error("My message", "key", "value")
}

func TestLoggerWarn(t *testing.T) {
	testLogger := createZapLogger(t, zap.WarnLevel)
	logger := &builtinLogger{logger: testLogger.Sugar()}
	logger.Warn("My message", "key", "value")
}
