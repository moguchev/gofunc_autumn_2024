package logger

import (
	rkentry "github.com/rookie-ninja/rk-entry/v2/entry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapEntry - logger entry name (you can override it if your logger entry has another name)
var ZapEntry = "zap"

func getZapLogerrOrDefault() *zap.Logger {
	log := rkentry.GlobalAppCtx.GetLoggerEntry(ZapEntry)
	if log != nil {
		return log.Logger
	}

	log = rkentry.GlobalAppCtx.GetLoggerEntryDefault()
	if log != nil {
		return log.Logger
	}

	return zap.New(zapcore.NewNopCore())
}

// Sugar wraps the Logger to provide a more ergonomic, but slightly slower,
// API. Sugaring a Logger is quite inexpensive, so it's reasonable for a
// single application to use both Loggers and SugaredLoggers, converting
// between them on the boundaries of performance-sensitive code.
func Sugar() *zap.SugaredLogger {
	log := getZapLogerrOrDefault()
	return log.Sugar()
}

// Named adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func Named(s string) *zap.Logger {
	log := getZapLogerrOrDefault()
	return log.Named(s)
}

// WithOptions clones the current Logger, applies the supplied Options, and
// returns the resulting Logger. It's safe to use concurrently.
func WithOptions(opts ...zap.Option) *zap.Logger {
	log := getZapLogerrOrDefault()
	return log.WithOptions(opts...)
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func With(fields ...zap.Field) *zap.Logger {
	log := getZapLogerrOrDefault()
	return log.With(fields...)
}

// Level reports the minimum enabled level for this logger.
//
// For NopLoggers, this is [zapcore.InvalidLevel].
func Level() zapcore.Level {
	log := getZapLogerrOrDefault()
	return log.Level()
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...zap.Field) {
	log := getZapLogerrOrDefault()
	log.Debug(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg string, fields ...zap.Field) {
	log := getZapLogerrOrDefault()
	log.Info(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(msg string, fields ...zap.Field) {
	log := getZapLogerrOrDefault()
	log.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(msg string, fields ...zap.Field) {
	log := getZapLogerrOrDefault()
	log.Error(msg, fields...)
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func DPanic(msg string, fields ...zap.Field) {
	log := getZapLogerrOrDefault()
	log.DPanic(msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, fields ...zap.Field) {
	log := getZapLogerrOrDefault()
	log.Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, fields ...zap.Field) {
	log := getZapLogerrOrDefault()
	log.Fatal(msg, fields...)
}
