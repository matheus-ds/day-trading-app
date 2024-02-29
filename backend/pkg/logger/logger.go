package logger

import (
	"log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLog *zap.Logger

type Field = zapcore.Field

func init() {
	var err error
	config := zap.NewProductionConfig()
	zapLog, err = config.Build(zap.AddCallerSkip(1))

	if err != nil {
		log.Printf("can't initialize zap logger: %v", err)
	}
	defer zapLog.Sync()
}

// wrapper for converting types to zapcore fields

func String(k, v string) Field {
	return zap.String(k, v)
}

func Duration(k string, d time.Duration) Field {
	return zap.Duration(k, d)
}

func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}

func Int(k string, i int) Field {
	return zap.Int(k, i)
}

func Array(key string, val zapcore.ArrayMarshaler) Field {
	return zap.Array(key, val)
}

func Int64(k string, i int64) Field {
	return zap.Int64(k, i)
}

func ErrorType(v error) Field {
	return zap.Error(v)
}

// wrapper for zap logging methods

func Info(message string, fields ...zap.Field) {
	zapLog.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	zapLog.Debug(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	zapLog.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	zapLog.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	zapLog.Fatal(message, fields...)
}
