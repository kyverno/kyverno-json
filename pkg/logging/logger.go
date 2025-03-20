package logging

import (
	"fmt"
	"io"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogLevel represents the severity level of logs
type LogLevel int

const (
	// ErrorLevel logs only errors
	ErrorLevel LogLevel = iota
	// WarnLevel logs warnings and errors
	WarnLevel
	// InfoLevel logs info messages and above
	InfoLevel
	// DebugLevel logs debug messages and above
	DebugLevel
)

var logLevelNames = map[LogLevel]string{
	ErrorLevel: "ERROR",
	WarnLevel:  "WARN",
	InfoLevel:  "INFO",
	DebugLevel: "DEBUG",
}

// Logger is an interface for logging operations
type Logger interface {
	// Error logs an error message
	Error(msg string, keysAndValues ...interface{})
	// Warn logs a warning message
	Warn(msg string, keysAndValues ...interface{})
	// Info logs an info message
	Info(msg string, keysAndValues ...interface{})
	// Debug logs a debug message
	Debug(msg string, keysAndValues ...interface{})
	// WithValues adds key-value pairs to the logger
	WithValues(keysAndValues ...interface{}) Logger
}

// Options configures the logger
type Options struct {
	// Level is the log level
	Level LogLevel
	// Out is the output writer
	Out io.Writer
}

// zapLogger implements the Logger interface using zap
type zapLogger struct {
	logger *zap.Logger
}

var (
	// defaultLogger is the default logger
	defaultLogger Logger
	// defaultOptions are the default logging options
	defaultOptions = Options{
		Level: InfoLevel,
		Out:   os.Stdout,
	}
	// mutex ensures that logger initialization is thread-safe
	mutex = sync.Mutex{}
)

// Initialize initializes the logger with the given options
func Initialize(opts Options) {
	mutex.Lock()
	defer mutex.Unlock()

	zapLevel := zapcore.InfoLevel
	switch opts.Level {
	case ErrorLevel:
		zapLevel = zapcore.ErrorLevel
	case WarnLevel:
		zapLevel = zapcore.WarnLevel
	case InfoLevel:
		zapLevel = zapcore.InfoLevel
	case DebugLevel:
		zapLevel = zapcore.DebugLevel
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(config)

	writeSyncer := zapcore.AddSync(opts.Out)
	core := zapcore.NewCore(encoder, writeSyncer, zapLevel)

	logger := zap.New(core)
	defaultLogger = &zapLogger{
		logger: logger,
	}
}

// Default returns the default logger
func Default() Logger {
	mutex.Lock()
	defer mutex.Unlock()

	if defaultLogger == nil {
		Initialize(defaultOptions)
	}
	return defaultLogger
}

// Error logs an error message
func (l *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.logger.Error(msg, toZapFields(keysAndValues...)...)
}

// Warn logs a warning message
func (l *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warn(msg, toZapFields(keysAndValues...)...)
}

// Info logs an info message
func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Info(msg, toZapFields(keysAndValues...)...)
}

// Debug logs a debug message
func (l *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Debug(msg, toZapFields(keysAndValues...)...)
}

// WithValues adds key-value pairs to the logger
func (l *zapLogger) WithValues(keysAndValues ...interface{}) Logger {
	return &zapLogger{
		logger: l.logger.With(toZapFields(keysAndValues...)...),
	}
}

// toZapFields converts a list of key-value pairs to zap fields
func toZapFields(keysAndValues ...interface{}) []zap.Field {
	if len(keysAndValues) == 0 {
		return nil
	}
	if len(keysAndValues)%2 != 0 {
		keysAndValues = append(keysAndValues, "MISSING")
	}
	fields := make([]zap.Field, 0, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			key = fmt.Sprintf("%v", keysAndValues[i])
		}
		fields = append(fields, zap.Any(key, keysAndValues[i+1]))
	}
	return fields
}

// ParseLevel parses a log level string into a LogLevel
func ParseLevel(level string) (LogLevel, error) {
	switch level {
	case "error", "ERROR":
		return ErrorLevel, nil
	case "warn", "WARN":
		return WarnLevel, nil
	case "info", "INFO":
		return InfoLevel, nil
	case "debug", "DEBUG":
		return DebugLevel, nil
	default:
		return InfoLevel, fmt.Errorf("unknown log level: %s", level)
	}
}

// LogError logs an error message using the default logger
func LogError(msg string, keysAndValues ...interface{}) {
	Default().Error(msg, keysAndValues...)
}

// LogWarn logs a warning message using the default logger
func LogWarn(msg string, keysAndValues ...interface{}) {
	Default().Warn(msg, keysAndValues...)
}

// LogInfo logs an info message using the default logger
func LogInfo(msg string, keysAndValues ...interface{}) {
	Default().Info(msg, keysAndValues...)
}

// LogDebug logs a debug message using the default logger
func LogDebug(msg string, keysAndValues ...interface{}) {
	Default().Debug(msg, keysAndValues...)
}
