package logger

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"base/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New(cfg config.LogConfig) (*zap.Logger, error) {
	level := parseLevel(cfg.Level)
	encoder := newEncoder(cfg.Format)
	cores, err := buildCores(cfg, encoder, level)
	if err != nil {
		return nil, err
	}

	return zap.New(zapcore.NewTee(cores...), zap.AddCaller()), nil
}

func buildCores(cfg config.LogConfig, encoder zapcore.Encoder, level zapcore.Level) ([]zapcore.Core, error) {
	var cores []zapcore.Core

	switch strings.ToLower(cfg.Output) {
	case "console":
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
	case "file":
		fileCore, err := newFileCore(cfg, encoder, level)
		if err != nil {
			return nil, err
		}
		cores = append(cores, fileCore)
	case "both":
		fileCore, err := newFileCore(cfg, encoder, level)
		if err != nil {
			return nil, err
		}
		cores = append(cores, fileCore)
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
	default:
		return nil, os.ErrInvalid
	}

	return cores, nil
}

func newFileCore(cfg config.LogConfig, encoder zapcore.Encoder, level zapcore.Level) (zapcore.Core, error) {
	if err := os.MkdirAll(cfg.Directory, 0o755); err != nil {
		return nil, err
	}

	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(cfg.Directory, cfg.Filename),
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		Compress:   cfg.Compress,
	})
	return zapcore.NewCore(encoder, writer, level), nil
}

func newEncoder(format string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.MessageKey = "msg"
	encoderConfig.CallerKey = "caller"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	if strings.EqualFold(format, "json") {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

type Field = zap.Field

func String(key, val string) Field { return zap.String(key, val) }
func Error(err error) Field        { return zap.Error(err) }
