package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// Debug level information
func Debug(msg string, fields ...Field) {
	getLogger().Debug(msg, fields...)
}

// Info level information
func Info(msg string, fields ...Field) {
	getLogger().Info(msg, fields...)
}

// Warn level information
func Warn(msg string, fields ...Field) {
	getLogger().Warn(msg, fields...)
}

// Error level information
func Error(msg string, fields ...Field) {
	getLogger().Error(msg, fields...)
}

// Panic level information
func Panic(msg string, fields ...Field) {
	getLogger().Panic(msg, fields...)
}

// Fatal level information
//func Fatal(msg string, fields ...Field) {
//	getLogger().Fatal(msg, fields...)
//}

// Debugf format level information
func Debugf(format string, a ...interface{}) {
	getLogger().Debug(fmt.Sprintf(format, a...))
}

// Infof format level information
func Infof(format string, a ...interface{}) {
	getLogger().Info(fmt.Sprintf(format, a...))
}

// Warnf format level information
func Warnf(format string, a ...interface{}) {
	getLogger().Warn(fmt.Sprintf(format, a...))
}

// Errorf format level information
func Errorf(format string, a ...interface{}) {
	getLogger().Error(fmt.Sprintf(format, a...))
}

// Fatalf format level information
//func Fatalf(format string, a ...interface{}) {
//	getLogger().Fatal(fmt.Sprintf(format, a...))
//}

// Sync flushing any buffered log entries, applications should take care to call Sync before exiting.
func Sync() error {
	err := getLogger().Sync()
	if err != nil && !strings.Contains(err.Error(), "/dev/stdout") {
		return err
	}
	return nil
}

// WithFields carrying field information
func WithFields(fields ...Field) *zap.Logger {
	return getLogger().With(fields...)
}
