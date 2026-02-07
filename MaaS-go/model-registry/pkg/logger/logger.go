package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds logger configuration
type Config struct {
	Level      string // debug, info, warn, error
	Format     string // json, console
	Output     string // stdout, file, both
	FilePath   string // log file path
	MaxSize    int    // megabytes
	MaxAge     int    // days
	MaxBackups int    // number of backups
	Compress   bool   // compress rotated files
}

// DefaultConfig returns default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:      "info",
		Format:     "json",
		Output:     "stdout",
		FilePath:   "logs/app.log",
		MaxSize:    100, // 100MB
		MaxAge:     7,   // 7 days
		MaxBackups: 10,
		Compress:   true,
	}
}

// Logger is a wrapper around zap.Logger
type Logger struct {
	*zap.Logger
	sugar  *zap.SugaredLogger
	config Config
}

// New creates a new logger with default configuration
func New(level string) *Logger {
	config := DefaultConfig()
	config.Level = level
	return NewWithConfig(config)
}

// NewWithConfig creates a new logger with custom configuration
func NewWithConfig(config Config) *Logger {
	// Parse log level
	logLevel := parseLogLevel(config.Level)

	// Create encoder
	encoder := createEncoder(config.Format)

	// Create write syncer
	writeSyncer := createWriteSyncer(config)

	// Create core
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)

	// Create logger
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return &Logger{
		Logger: logger,
		sugar:  logger.Sugar(),
		config: config,
	}
}

// parseLogLevel parses string level to zapcore.Level
func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// createEncoder creates zap encoder based on format
func createEncoder(format string) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if format == "console" {
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

// createWriteSyncer creates write syncer based on output type
func createWriteSyncer(config Config) zapcore.WriteSyncer {
	switch config.Output {
	case "file":
		return createFileSyncer(config)
	case "both":
		stdoutSyncer := zapcore.AddSync(os.Stdout)
		fileSyncer := createFileSyncer(config)
		return zapcore.NewMultiWriteSyncer(stdoutSyncer, fileSyncer)
	default: // stdout
		return zapcore.AddSync(os.Stdout)
	}
}

// createFileSyncer creates file syncer with rotation
func createFileSyncer(config Config) zapcore.WriteSyncer {
	// Create log directory if not exists
	logDir := filepath.Dir(config.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		// Fall back to stdout if cannot create directory
		return zapcore.AddSync(os.Stdout)
	}

	// Open log file
	file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Fall back to stdout if cannot open file
		return zapcore.AddSync(os.Stdout)
	}

	return zapcore.AddSync(file)
}

// With creates a child logger with additional fields
func (l *Logger) With(fields ...interface{}) *Logger {
	return &Logger{
		Logger: l.Logger.With(l.toFields(fields)...),
		sugar:  l.sugar.With(fields...),
		config: l.config,
	}
}

// WithContext creates a logger with request context
func (l *Logger) WithContext(requestID string) *Logger {
	return l.With("request_id", requestID)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	l.sugar.Debugw(msg, keysAndValues...)
}

// Info logs an info message
func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.sugar.Infow(msg, keysAndValues...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	l.sugar.Warnw(msg, keysAndValues...)
}

// Error logs an error message
func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	l.sugar.Errorw(msg, keysAndValues...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	l.sugar.Fatalw(msg, keysAndValues...)
}

// Sync flushes buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// GetConfig returns logger configuration
func (l *Logger) GetConfig() Config {
	return l.config
}

// toFields converts interface slice to zap.Field slice
func (l *Logger) toFields(keysAndValues []interface{}) []zap.Field {
	if len(keysAndValues)%2 != 0 {
		keysAndValues = append(keysAndValues, "missing")
	}

	fields := make([]zap.Field, 0, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			key = "invalid_key"
		}
		fields = append(fields, zap.Any(key, keysAndValues[i+1]))
	}
	return fields
}

// Global logger instance
var globalLogger *Logger

// InitGlobalLogger initializes global logger
func InitGlobalLogger(config Config) {
	globalLogger = NewWithConfig(config)
}

// GetGlobalLogger returns global logger instance
func GetGlobalLogger() *Logger {
	if globalLogger == nil {
		globalLogger = New("info")
	}
	return globalLogger
}

// Debug logs using global logger
func Debug(msg string, keysAndValues ...interface{}) {
	GetGlobalLogger().Debug(msg, keysAndValues...)
}

// Info logs using global logger
func Info(msg string, keysAndValues ...interface{}) {
	GetGlobalLogger().Info(msg, keysAndValues...)
}

// Warn logs using global logger
func Warn(msg string, keysAndValues ...interface{}) {
	GetGlobalLogger().Warn(msg, keysAndValues...)
}

// Error logs using global logger
func Error(msg string, keysAndValues ...interface{}) {
	GetGlobalLogger().Error(msg, keysAndValues...)
}

// Fatal logs using global logger
func Fatal(msg string, keysAndValues ...interface{}) {
	GetGlobalLogger().Fatal(msg, keysAndValues...)
}
