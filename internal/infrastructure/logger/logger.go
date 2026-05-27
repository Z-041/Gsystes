package logger

import (
	"os"
	"time"

	"github.com/gsystes/backend/internal/infrastructure/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var globalLogger *zap.Logger

func InitLogger(cfg config.LogConfig) error {
	writeSyncer := getLogWriter(cfg)
	encoder := getEncoder()

	var level zapcore.Level
	if err := level.Set(cfg.Level); err != nil {
		level = zapcore.InfoLevel
	}

	core := zapcore.NewCore(encoder, writeSyncer, level)
	globalLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zap.ReplaceGlobals(globalLogger)
	return nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(cfg config.LogConfig) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	ws := zapcore.AddSync(lumberJackLogger)
	if cfg.Level == "debug" {
		ws = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), ws)
	}
	return ws
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

func With(fields ...zap.Field) *zap.Logger {
	return globalLogger.With(fields...)
}

func Sync() error {
	return globalLogger.Sync()
}

func GetLogger() *zap.Logger {
	return globalLogger
}

func StringField(key, value string) zap.Field {
	return zap.String(key, value)
}

func IntField(key string, value int) zap.Field {
	return zap.Int(key, value)
}

func AnyField(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

func ErrorField(err error) zap.Field {
	return zap.Error(err)
}

func DurationField(key string, value time.Duration) zap.Field {
	return zap.Duration(key, value)
}

func UintField(key string, value uint) zap.Field {
	return zap.Uint(key, value)
}
